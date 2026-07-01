package types

import "markitos-it-svc-golden/internal/domain/shared"

type GoldenId struct {
	value string
}

func NewGoldenId(value string) (*GoldenId, error) {
	if shared.IsUUIDv4(value) {
		return &GoldenId{value}, nil
	}

	return nil, shared.ErrGoldenBadRequest
}

func (b *GoldenId) Value() string {
	return b.value
}
