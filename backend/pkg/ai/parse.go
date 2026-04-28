package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// extractJSON pulls the first complete JSON object out of raw model output,
// handling preamble prose, markdown fences, and trailing notes.
func extractJSON(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	// Strip markdown code fences
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	// Find the outermost { ... } block
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start == -1 || end == -1 || end <= start {
		return "", fmt.Errorf("no JSON object found in response")
	}

	return raw[start : end+1], nil
}

func parseAnalysisResult(raw string) (*AnalysisResult, error) {
	jsonStr, err := extractJSON(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to locate JSON in AI response: %w\nraw output: %s", err, raw)
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI JSON: %w\nextracted: %s", err, jsonStr)
	}
	result.RawResponse = jsonStr
	return &result, nil
}
