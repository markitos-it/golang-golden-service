package domain_test

import (
	"os"
	"path/filepath"
	"testing"

	"markitos-it-svc-golden/internal/domain/types"
)

func TestCanCreateValidGoldenPoster(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.jpg")

	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	poster, err := types.NewGoldenPoster(tempDir, "test.jpg")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if poster.Value() != "test.jpg" {
		t.Errorf("Expected 'test.jpg', got %v", poster.Value())
	}
}

func TestCannotCreateInvalidGoldenPoster(t *testing.T) {
	tempDir := t.TempDir()

	_, err := types.NewGoldenPoster(tempDir, "nonexistent.jpg")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestCanCreateEmptyGoldenPoster(t *testing.T) {
	tempDir := t.TempDir()

	poster, err := types.NewGoldenPoster(tempDir, "")
	if err != nil {
		t.Errorf("Expected no error for empty poster, got %v", err)
	}

	if poster.Value() != "" {
		t.Errorf("Expected empty string, got %v", poster.Value())
	}
}
