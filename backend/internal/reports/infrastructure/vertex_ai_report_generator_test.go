package infrastructure

import (
	"strings"
	"testing"
)

func TestNewVertexAIReportGeneratorRejectsMissingProjectID(t *testing.T) {
	generator, err := NewVertexAIReportGenerator("   ", "", "")
	if err == nil {
		t.Fatalf("expected project id required error")
	}
	if generator != nil {
		t.Fatalf("expected nil generator when project id is missing")
	}
	if err.Error() != "gcp project id is required" {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestBuildWeeklyReportPromptIncludesInstructionsAndInput(t *testing.T) {
	prompt := buildWeeklyReportPrompt(`{"workspace_name":"Algorithms","user_name":"Ana Gomez"}`)
	if !strings.Contains(prompt, "JSON:") {
		t.Fatalf("expected prompt to include JSON marker")
	}
	if !strings.Contains(prompt, `"workspace_name":"Algorithms"`) {
		t.Fatalf("expected prompt to include serialized input")
	}
	if !strings.Contains(prompt, "No uses Markdown") {
		t.Fatalf("expected prompt to contain instructions, got %q", prompt)
	}
	if !strings.Contains(prompt, "130 palabras") {
		t.Fatalf("expected prompt to include length rule, got %q", prompt)
	}
}

func TestBuildWeeklyReportPromptMentionsExpectedNaturalLanguageFields(t *testing.T) {
	prompt := buildWeeklyReportPrompt(`{"total_hours":6}`)
	expectedFragments := []string{
		"monitor o asistente",
		"workspace",
		"horas reportadas",
		"recomendaci",
	}

	for _, fragment := range expectedFragments {
		if !strings.Contains(strings.ToLower(prompt), fragment) {
			t.Fatalf("expected prompt to contain %q, got %q", fragment, prompt)
		}
	}
}
