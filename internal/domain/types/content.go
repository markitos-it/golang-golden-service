package types

import (
	"encoding/base64"
	"markitos-it-svc-golden/internal/domain/shared"
	"regexp"
)

const (
	GOLDEN_CONTENT_ONE_MEGABYTE_PER_SECURITY = 1000000
)

type GoldenContent struct {
	value string
}

func NewGoldenContent(value string) (*GoldenContent, error) {
	if isValidGoldenContent(value) {
		sanitized := sanitizeGoldenContent(value)
		encoded := base64.StdEncoding.EncodeToString([]byte(sanitized))
		return &GoldenContent{encoded}, nil
	}

	return nil, shared.ErrGoldenBadRequest
}

func isValidGoldenContent(value string) bool {
	if value == "" {
		return true
	}

	sanitized := sanitizeGoldenContent(value)
	return len(sanitized) <= GOLDEN_CONTENT_ONE_MEGABYTE_PER_SECURITY
}

func sanitizeGoldenContent(value string) string {
	pattern := `(?i)<script[^>]*>.*?</script>`
	re := regexp.MustCompile(pattern)
	sanitized := re.ReplaceAllString(value, "")

	pattern = `(?i)<(script|iframe|object|embed|link|style)[^>]*>.*?</(script|iframe|object|embed|link|style)>`
	re = regexp.MustCompile(pattern)
	sanitized = re.ReplaceAllString(sanitized, "")

	return sanitized
}

func (c *GoldenContent) Value() string {
	return c.value
}

func (c *GoldenContent) DecodedValue() (string, error) {
	if c.value == "" {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(c.value)
	if err != nil {
		return "", shared.ErrGoldenBadRequest
	}

	return string(decoded), nil
}
