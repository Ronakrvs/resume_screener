package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Resume struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	FileName      string    `json:"file_name" db:"file_name"`
	FileURL       string    `json:"file_url" db:"file_url"`
	ExtractedText string    `json:"extracted_text,omitempty" db:"extracted_text"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Job struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Result struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ResumeID        uuid.UUID `json:"resume_id" db:"resume_id"`
	JobID           uuid.UUID `json:"job_id" db:"job_id"`
	Score           int       `json:"score" db:"score"`
	Strengths       []string  `json:"strengths" db:"strengths"`
	MissingSkills   []string  `json:"missing_skills" db:"missing_skills"`
	Recommendation  string    `json:"recommendation" db:"recommendation"`
	RawAIResponse   string    `json:"raw_ai_response,omitempty" db:"raw_ai_response"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	// Joined fields
	ResumeFileName  string    `json:"resume_file_name,omitempty" db:"resume_file_name"`
	JobTitle        string    `json:"job_title,omitempty" db:"job_title"`
}

type AnalyzeRequest struct {
	ResumeID uuid.UUID `json:"resume_id" binding:"required"`
	JobID    uuid.UUID `json:"job_id" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateJobRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
