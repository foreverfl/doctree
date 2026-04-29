package config

import (
	"slices"
	"testing"
)

func TestPickEditor_PrecedenceVisualOverEditor(t *testing.T) {
	t.Setenv("VISUAL", "code --wait")
	t.Setenv("EDITOR", "nano")

	got := pickEditor()
	want := []string{"code", "--wait"}
	if !slices.Equal(got, want) {
		t.Errorf("pickEditor() = %v, want %v", got, want)
	}
}

func TestPickEditor_FallsBackToEditor(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "nano")

	got := pickEditor()
	want := []string{"nano"}
	if !slices.Equal(got, want) {
		t.Errorf("pickEditor() = %v, want %v", got, want)
	}
}

func TestPickEditor_FinalFallbackIsVi(t *testing.T) {
	t.Setenv("VISUAL", "")
	t.Setenv("EDITOR", "")

	got := pickEditor()
	want := []string{"vi"}
	if !slices.Equal(got, want) {
		t.Errorf("pickEditor() = %v, want %v", got, want)
	}
}
