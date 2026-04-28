package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nvl258/resume-screener/internal/model"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job *model.Job) error {
	query := `INSERT INTO jobs (id, user_id, title, description) VALUES ($1, $2, $3, $4) RETURNING created_at`
	job.ID = uuid.New()
	return r.db.QueryRow(ctx, query, job.ID, job.UserID, job.Title, job.Description).Scan(&job.CreatedAt)
}

func (r *JobRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	job := &model.Job{}
	query := `SELECT id, user_id, title, description, created_at FROM jobs WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&job.ID, &job.UserID, &job.Title, &job.Description, &job.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}
	return job, nil
}

func (r *JobRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*model.Job, error) {
	query := `SELECT id, user_id, title, description, created_at FROM jobs WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*model.Job
	for rows.Next() {
		job := &model.Job{}
		if err := rows.Scan(&job.ID, &job.UserID, &job.Title, &job.Description, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
