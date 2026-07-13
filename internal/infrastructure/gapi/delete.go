package gapi

import (
	context "context"
	"log"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteGolden(ctx context.Context, in *DeleteGoldenRequest) (*DeleteGoldenResponse, error) {
	if _, err := types.NewGoldenId(in.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, msgInvalidRequest)
	}

	request := services.GoldenDeleteRequest{Id: in.Id}

	var service services.GoldenDeleteService = services.NewGoldenDeleteService(s.repository)
	if err := service.Do(request); err != nil {
		log.Printf("[ERROR] DeleteGolden: %v\n", err)
		return nil, s.ToStatusError(err)
	}

	return &DeleteGoldenResponse{
		Deleted: request.Id,
	}, nil
}
