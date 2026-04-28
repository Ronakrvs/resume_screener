package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		client: openai.NewClient(apiKey),
		model:  openai.GPT4oMini,
	}
}

const analysisPrompt = `You are an expert HR recruiter and resume analyst. Compare the resume below against the job description and return a JSON response.

RESUME:
%s

JOB DESCRIPTION:
%s

Return ONLY valid JSON in this exact format (no markdown, no explanation):
{
  "score": <integer 0-100 representing match percentage>,
  "strengths": [<list of specific strengths matching the job>],
  "missing_skills": [<list of skills/requirements in the job but missing from resume>],
  "recommendation": "<one paragraph recommendation for the recruiter>"
}

Scoring guide:
- 90-100: Exceptional match, highly recommended
- 70-89: Good match, worth interviewing
- 50-69: Partial match, some gaps
- 30-49: Weak match, significant gaps
- 0-29: Poor match, does not meet requirements`

func (p *OpenAIProvider) AnalyzeResume(ctx context.Context, resumeText, jobDescription string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(analysisPrompt, resumeText, jobDescription)

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: p.model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an expert HR recruiter. Always respond with valid JSON only."},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0.3,
		MaxTokens:   1000,
	})
	if err != nil {
		return nil, fmt.Errorf("openai request failed: %w", err)
	}

	return parseAnalysisResult(resp.Choices[0].Message.Content)
}
