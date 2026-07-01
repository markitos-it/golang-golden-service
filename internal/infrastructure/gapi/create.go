package gapi

import (
	context "context"
	"log"
	"markitos-it-svc-golden/internal/domain/services"

	"google.golang.org/grpc/status"
)

func (s *Server) CreateGolden(ctx context.Context, req *CreateGoldenRequest) (*CreateGoldenResponse, error) {
	var request services.GoldenCreateRequest = services.GoldenCreateRequest{
		Name:       req.Name,
		Content:    req.Content,
		PosterData: req.PosterData,
		/* ___CUSTOM_FIELDS_TO_DOMAIN___*/
	}

	var service services.GoldenCreateService = services.NewGoldenCreateService(s.repository, s.config.BaseDir)
	entity, err := service.Do(request)
	if err != nil {
		log.Printf("❌ ERROR (CreateGolden): %v\n", err)
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &CreateGoldenResponse{
		Id:      entity.Id,
		Name:    entity.Name,
		Content: entity.Content,
		Poster:  entity.Poster,
		/* ___CUSTOM_FIELDS_TO_PROTO___*/
	}, nil
}
