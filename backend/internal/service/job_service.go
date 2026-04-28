package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nvl258/resume-screener/internal/model"
	"github.com/nvl258/resume-screener/internal/repository"
)

type JobService struct {
	jobRepo *repository.JobRepository
}

func NewJobService(jobRepo *repository.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

func (s *JobService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateJobRequest) (*model.Job, error) {
	job := &model.Job{UserID: userID, Title: req.Title, Description: req.Description}
	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *JobService) List(ctx context.Context, userID uuid.UUID) ([]*model.Job, error) {
	return s.jobRepo.ListByUser(ctx, userID)
}

func (s *JobService) Get(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	return s.jobRepo.FindByID(ctx, id)
}
