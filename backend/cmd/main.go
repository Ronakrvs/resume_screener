package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nvl258/resume-screener/configs"
	"github.com/nvl258/resume-screener/internal/api"
	"github.com/nvl258/resume-screener/internal/api/handlers"
	"github.com/nvl258/resume-screener/internal/repository"
	"github.com/nvl258/resume-screener/internal/service"
	pkgai "github.com/nvl258/resume-screener/pkg/ai"
	"github.com/nvl258/resume-screener/pkg/storage"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := configs.Load()
	ctx := context.Background()

	// Database
	db, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to database")

	// Run migrations
	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to init migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied")

	// Redis — parse URL or use as addr directly
	redisOpts, redisErr := redis.ParseURL(cfg.RedisURL)
	if redisErr != nil {
		redisOpts = &redis.Options{Addr: cfg.RedisURL}
	}
	rdb := redis.NewClient(redisOpts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis unavailable: %v — caching disabled", err)
	} else {
		log.Println("Connected to Redis")
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	resumeRepo := repository.NewResumeRepository(db)
	jobRepo := repository.NewJobRepository(db)
	resultRepo := repository.NewResultRepository(db)

	// AI Provider
	var aiProvider pkgai.Provider
	if cfg.OpenAIKey != "" {
		aiProvider = pkgai.NewOpenAIProvider(cfg.OpenAIKey)
		log.Println("Using OpenAI provider")
	} else {
		aiProvider = pkgai.NewOllamaProvider("", "")
		log.Println("Using Ollama provider (local)")
	}

	// Storage
	store := storage.NewS3Store(cfg.S3Region, cfg.S3Bucket, cfg.S3Endpoint, cfg.AWSKey, cfg.AWSSecret)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	resumeSvc := service.NewResumeService(resumeRepo, store)
	jobSvc := service.NewJobService(jobRepo)
	analysisSvc := service.NewAnalysisService(resumeRepo, jobRepo, resultRepo, aiProvider, rdb)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	resumeHandler := handlers.NewResumeHandler(resumeSvc)
	jobHandler := handlers.NewJobHandler(jobSvc)
	analysisHandler := handlers.NewAnalysisHandler(analysisSvc)

	// Router
	router := api.NewRouter(authHandler, resumeHandler, jobHandler, analysisHandler, cfg.JWTSecret, rdb)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}
