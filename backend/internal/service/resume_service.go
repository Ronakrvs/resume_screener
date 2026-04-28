package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nvl258/resume-screener/internal/model"
	"github.com/nvl258/resume-screener/internal/repository"
	"github.com/nvl258/resume-screener/pkg/parser"
	"github.com/nvl258/resume-screener/pkg/storage"
)

type ResumeService struct {
	resumeRepo *repository.ResumeRepository
	store      *storage.S3Store
}

func NewResumeService(resumeRepo *repository.ResumeRepository, store *storage.S3Store) *ResumeService {
	return &ResumeService{resumeRepo: resumeRepo, store: store}
}

func (s *ResumeService) Upload(ctx context.Context, userID uuid.UUID, filename string, data []byte, contentType string) (*model.Resume, error) {
	if len(data) > 10*1024*1024 {
		return nil, fmt.Errorf("file too large (max 10MB)")
	}

	text, err := parser.ExtractText(filename, data)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text: %w", err)
	}
	if len(strings.TrimSpace(text)) < 20 {
		return nil, fmt.Errorf("could not extract text from file — try a text-based PDF (not a scanned image)")
	}

	fileURL, err := s.store.Upload(ctx, filename, data, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	resume := &model.Resume{
		UserID:        userID,
		FileName:      filename,
		FileURL:       fileURL,
		ExtractedText: text,
	}
	if err := s.resumeRepo.Create(ctx, resume); err != nil {
		return nil, fmt.Errorf("failed to save resume: %w", err)
	}
	return resume, nil
}

func (s *ResumeService) List(ctx context.Context, userID uuid.UUID) ([]*model.Resume, error) {
	return s.resumeRepo.ListByUser(ctx, userID)
}

func (s *ResumeService) Get(ctx context.Context, id uuid.UUID) (*model.Resume, error) {
	return s.resumeRepo.FindByID(ctx, id)
}
