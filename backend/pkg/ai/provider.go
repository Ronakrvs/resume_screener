package ai

import "context"

type AnalysisResult struct {
	Score          int      `json:"score"`
	Strengths      []string `json:"strengths"`
	MissingSkills  []string `json:"missing_skills"`
	Recommendation string   `json:"recommendation"`
	RawResponse    string   `json:"raw_response"`
}

type Provider interface {
	AnalyzeResume(ctx context.Context, resumeText, jobDescription string) (*AnalysisResult, error)
}
