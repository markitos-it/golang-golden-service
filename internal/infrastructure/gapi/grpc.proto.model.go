package gapi

import (
	domain "markitos-it-svc-golden/internal/domain/model"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewGRPCGolden(in *domain.Golden) *Golden {
	return &Golden{
		Id:        in.Id,
		Name:      in.Name,
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt),
	}
}
