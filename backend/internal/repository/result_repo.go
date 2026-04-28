package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nvl258/resume-screener/internal/model"
)

type ResultRepository struct {
	db *pgxpool.Pool
}

func NewResultRepository(db *pgxpool.Pool) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r *ResultRepository) Create(ctx context.Context, result *model.Result) error {
	query := `INSERT INTO results (id, resume_id, job_id, score, strengths, missing_skills, recommendation, raw_ai_response)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at`
	result.ID = uuid.New()
	return r.db.QueryRow(ctx, query,
		result.ID, result.ResumeID, result.JobID, result.Score,
		result.Strengths, result.MissingSkills, result.Recommendation, result.RawAIResponse,
	).Scan(&result.CreatedAt)
}

func (r *ResultRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Result, error) {
	result := &model.Result{}
	query := `
		SELECT r.id, r.resume_id, r.job_id, r.score, r.strengths, r.missing_skills, r.recommendation, r.raw_ai_response, r.created_at,
		       res.file_name AS resume_file_name, j.title AS job_title
		FROM results r
		JOIN resumes res ON res.id = r.resume_id
		JOIN jobs j ON j.id = r.job_id
		WHERE r.id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&result.ID, &result.ResumeID, &result.JobID, &result.Score,
		&result.Strengths, &result.MissingSkills, &result.Recommendation, &result.RawAIResponse, &result.CreatedAt,
		&result.ResumeFileName, &result.JobTitle,
	)
	if err != nil {
		return nil, fmt.Errorf("result not found: %w", err)
	}
	return result, nil
}

func (r *ResultRepository) ListByJob(ctx context.Context, jobID uuid.UUID) ([]*model.Result, error) {
	query := `
		SELECT r.id, r.resume_id, r.job_id, r.score, r.strengths, r.missing_skills, r.recommendation, r.created_at,
		       res.file_name AS resume_file_name, j.title AS job_title
		FROM results r
		JOIN resumes res ON res.id = r.resume_id
		JOIN jobs j ON j.id = r.job_id
		WHERE r.job_id = $1
		ORDER BY r.score DESC`
	rows, err := r.db.Query(ctx, query, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*model.Result
	for rows.Next() {
		result := &model.Result{}
		if err := rows.Scan(
			&result.ID, &result.ResumeID, &result.JobID, &result.Score,
			&result.Strengths, &result.MissingSkills, &result.Recommendation, &result.CreatedAt,
			&result.ResumeFileName, &result.JobTitle,
		); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
