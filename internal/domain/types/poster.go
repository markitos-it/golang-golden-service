package types

import (
	"markitos-it-svc-golden/internal/domain/shared"
	"os"
	"path/filepath"
	"strings"
)

type GoldenPoster struct {
	value string
}

func NewGoldenPoster(baseDir, value string) (*GoldenPoster, error) {
	if isValidGoldenPoster(baseDir, value) {
		return &GoldenPoster{value}, nil
	}

	return nil, shared.ErrInvalidGoldenPoster
}

func isValidGoldenPoster(baseDir, value string) bool {
	if value == "" {
		return true
	}

	cleanBaseDir := strings.TrimSuffix(baseDir, "/")
	fullPath := filepath.Join(cleanBaseDir, value)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (p *GoldenPoster) Value() string {
	return p.value
}
