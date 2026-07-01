package domain_test

import (
	"strings"
	"testing"

	"markitos-it-svc-golden/internal/domain/types"
)

func TestCanCreateValidGoldenName(t *testing.T) {
	validNames := []string{
		"ValidName",
		"AnotherValidName",
		"Valid Name With Spaces",
		"Short",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"}
	for _, name := range validNames {
		if _, err := types.NewGoldenName(name); err != nil {
			t.Errorf("Expected valid name, but got invalid: %s", name)
		}
	}

	invalidNames := []string{
		" InvalidName",
		"InvalidName ",
		"Invalid Name ",
		" Invalid Name",
		"Invalid@Name",
		"Invalid#Name",
		"Invalid123Name",
		"Invalid Name With Spaces ",
		" Invalid Name With Spaces",
		"Invalid Name With Spaces And Symbols!",
	}
	for _, name := range invalidNames {
		if _, err := types.NewGoldenName(name); err == nil {
			t.Errorf("Expected valid name, but got invalid: %s", name)
		}
	}

	invalidLengthNames := []string{
		strings.Repeat("a", types.GOLDEN_NAME_MAX_LENGTH+1),
		strings.Repeat("b", types.GOLDEN_NAME_MIN_LENGTH-1),
		"",
	}
	for _, name := range invalidLengthNames {
		if _, err := types.NewGoldenName(name); err == nil {
			t.Errorf("Expected invalid name, but got invalid: %s", name)
		}
	}

}
