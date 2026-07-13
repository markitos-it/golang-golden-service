package types

import (
	"markitos-it-svc-golden/internal/domain/shared"
	"regexp"
)

type GoldenName struct {
	value string
}

const GOLDEN_NAME_MAX_LENGTH = 100
const GOLDEN_NAME_MIN_LENGTH = 3

func NewGoldenName(value string) (*GoldenName, error) {
	if isValidGoldenName(value) {
		return &GoldenName{value}, nil
	}

	return nil, shared.ErrInvalidGoldenName
}

func isValidGoldenName(value string) bool {
	if len(value) > GOLDEN_NAME_MAX_LENGTH || len(value) < GOLDEN_NAME_MIN_LENGTH {
		return false
	}

	pattern := `^[a-zA-Z]{1}[a-zA-Z ]+[a-zA-Z]$|^[a-zA-Z]{1}$`
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return false
	}

	return matched
}

func (b *GoldenName) Value() string {
	return b.value
}
