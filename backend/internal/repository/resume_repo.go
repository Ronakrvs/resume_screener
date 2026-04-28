package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nvl258/resume-screener/internal/model"
)

type ResumeRepository struct {
	db *pgxpool.Pool
}

func NewResumeRepository(db *pgxpool.Pool) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) Create(ctx context.Context, resume *model.Resume) error {
	query := `INSERT INTO resumes (id, user_id, file_name, file_url, extracted_text)
	          VALUES ($1, $2, $3, $4, $5) RETURNING created_at`
	resume.ID = uuid.New()
	return r.db.QueryRow(ctx, query, resume.ID, resume.UserID, resume.FileName, resume.FileURL, resume.ExtractedText).
		Scan(&resume.CreatedAt)
}

func (r *ResumeRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Resume, error) {
	resume := &model.Resume{}
	query := `SELECT id, user_id, file_name, file_url, extracted_text, created_at FROM resumes WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&resume.ID, &resume.UserID, &resume.FileName, &resume.FileURL, &resume.ExtractedText, &resume.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("resume not found: %w", err)
	}
	return resume, nil
}

func (r *ResumeRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Resume, error) {
	query := `SELECT id, user_id, file_name, file_url, created_at FROM resumes WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resumes []*model.Resume
	for rows.Next() {
		resume := &model.Resume{}
		if err := rows.Scan(&resume.ID, &resume.UserID, &resume.FileName, &resume.FileURL, &resume.CreatedAt); err != nil {
			return nil, err
		}
		resumes = append(resumes, resume)
	}
	return resumes, nil
}
