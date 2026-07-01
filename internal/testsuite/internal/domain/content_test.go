package domain_test

import (
	"encoding/base64"
	"testing"

	"markitos-it-svc-golden/internal/domain/types"
)

func TestCanCreateValidGoldenContent(t *testing.T) {
	content, err := types.NewGoldenContent("Hello World")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedEncoded := base64.StdEncoding.EncodeToString([]byte("Hello World"))
	if content.Value() != expectedEncoded {
		t.Errorf("Expected base64 encoded content, got %v", content.Value())
	}

	decoded, err := content.DecodedValue()
	if err != nil {
		t.Errorf("Expected no error decoding, got %v", err)
	}

	if decoded != "Hello World" {
		t.Errorf("Expected 'Hello World', got %v", decoded)
	}
}

func TestCanCreateEmptyGoldenContent(t *testing.T) {
	content, err := types.NewGoldenContent("")
	if err != nil {
		t.Errorf("Expected no error for empty content, got %v", err)
	}

	if content.Value() != base64.StdEncoding.EncodeToString([]byte("")) {
		t.Errorf("Expected base64 encoded empty string")
	}

	decoded, err := content.DecodedValue()
	if err != nil {
		t.Errorf("Expected no error decoding empty content, got %v", err)
	}

	if decoded != "" {
		t.Errorf("Expected empty string, got %v", decoded)
	}
}

func TestContentSanitization(t *testing.T) {
	maliciousContent := "<script>alert('xss')</script>Hello World"
	content, err := types.NewGoldenContent(maliciousContent)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	decoded, err := content.DecodedValue()
	if err != nil {
		t.Errorf("Expected no error decoding, got %v", err)
	}

	if decoded == maliciousContent {
		t.Error("Expected script tag to be sanitized")
	}

	if decoded != "Hello World" {
		t.Errorf("Expected sanitized content 'Hello World', got %v", decoded)
	}
}
