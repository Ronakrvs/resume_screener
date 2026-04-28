package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nvl258/resume-screener/internal/api/handlers"
	"github.com/nvl258/resume-screener/internal/api/middleware"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	auth     *handlers.AuthHandler
	resume   *handlers.ResumeHandler
	job      *handlers.JobHandler
	analysis *handlers.AnalysisHandler
	secret   string
	rdb      *redis.Client
}

func NewRouter(
	auth *handlers.AuthHandler,
	resume *handlers.ResumeHandler,
	job *handlers.JobHandler,
	analysis *handlers.AnalysisHandler,
	secret string,
	rdb *redis.Client,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
		}

		protected := api.Group("/")
		protected.Use(middleware.Auth(secret))
		protected.Use(middleware.RateLimit(rdb, 100, time.Minute))
		{
			protected.POST("/resume/upload", resume.Upload)
			protected.GET("/resume", resume.List)
			protected.GET("/resume/:id", resume.Get)

			protected.POST("/job", job.Create)
			protected.GET("/job", job.List)

			protected.POST("/analyze", analysis.Analyze)
			protected.GET("/results/:id", analysis.GetResult)
			protected.GET("/ranking/:jobId", analysis.GetRanking)
		}
	}

	return r
}
