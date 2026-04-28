package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nvl258/resume-screener/internal/model"
	"github.com/nvl258/resume-screener/internal/repository"
	"github.com/nvl258/resume-screener/pkg/ai"
	"github.com/redis/go-redis/v9"
)

type AnalysisService struct {
	resumeRepo *repository.ResumeRepository
	jobRepo    *repository.JobRepository
	resultRepo *repository.ResultRepository
	aiProvider ai.Provider
	rdb        *redis.Client
}

func NewAnalysisService(
	resumeRepo *repository.ResumeRepository,
	jobRepo *repository.JobRepository,
	resultRepo *repository.ResultRepository,
	aiProvider ai.Provider,
	rdb *redis.Client,
) *AnalysisService {
	return &AnalysisService{
		resumeRepo: resumeRepo,
		jobRepo:    jobRepo,
		resultRepo: resultRepo,
		aiProvider: aiProvider,
		rdb:        rdb,
	}
}

func (s *AnalysisService) Analyze(ctx context.Context, req *model.AnalyzeRequest) (*model.Result, error) {
	cacheKey := fmt.Sprintf("result:%s:%s", req.ResumeID, req.JobID)
	if cached, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil {
		var result model.Result
		if json.Unmarshal([]byte(cached), &result) == nil {
			return &result, nil
		}
	}

	resume, err := s.resumeRepo.FindByID(ctx, req.ResumeID)
	if err != nil {
		return nil, fmt.Errorf("resume not found")
	}
	job, err := s.jobRepo.FindByID(ctx, req.JobID)
	if err != nil {
		return nil, fmt.Errorf("job not found")
	}

	aiResult, err := s.aiProvider.AnalyzeResume(ctx, resume.ExtractedText, job.Description)
	if err != nil {
		return nil, fmt.Errorf("AI analysis failed: %w", err)
	}

	result := &model.Result{
		ResumeID:       req.ResumeID,
		JobID:          req.JobID,
		Score:          aiResult.Score,
		Strengths:      aiResult.Strengths,
		MissingSkills:  aiResult.MissingSkills,
		Recommendation: aiResult.Recommendation,
		RawAIResponse:  aiResult.RawResponse,
	}
	if err := s.resultRepo.Create(ctx, result); err != nil {
		return nil, fmt.Errorf("failed to save result: %w", err)
	}

	result.ResumeFileName = resume.FileName
	result.JobTitle = job.Title

	if data, err := json.Marshal(result); err == nil {
		s.rdb.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return result, nil
}

func (s *AnalysisService) GetResult(ctx context.Context, id uuid.UUID) (*model.Result, error) {
	return s.resultRepo.FindByID(ctx, id)
}

func (s *AnalysisService) GetRanking(ctx context.Context, jobID uuid.UUID) ([]*model.Result, error) {
	return s.resultRepo.ListByJob(ctx, jobID)
}
