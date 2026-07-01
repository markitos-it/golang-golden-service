package gapi

import (
	context "context"
	"log"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateGolden(ctx context.Context, in *UpdateGoldenRequest) (*UpdateGoldenResponse, error) {
	if _, err := types.NewGoldenId(in.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var service services.GoldenUpdateService = services.NewGoldenUpdateService(s.repository, s.config.BaseDir)
	var request services.GoldenUpdateRequest = services.GoldenUpdateRequest{
		Id:         in.Id,
		Name:       in.Name,
		Content:    in.Content,
		PosterData: in.PosterData,
	}
	if err := service.Do(request); err != nil {
		log.Printf("❌ ERROR (UpdateGolden): %v\n", err)
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &UpdateGoldenResponse{
		Updated: request.Id,
	}, nil
}
