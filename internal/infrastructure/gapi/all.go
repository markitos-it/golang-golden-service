package gapi

import (
	context "context"
	"markitos-it-svc-golden/internal/domain/services"
)

func (s *Server) ListGoldens(ctx context.Context, in *ListGoldensRequest) (*ListGoldensResponse, error) {
	var service services.GoldenAllService = services.NewGoldenAllService(s.repository)
	response, err := service.Do()
	if err != nil {
		return nil, s.ToStatusError(err)
	}

	return &ListGoldensResponse{
		Goldens: s.ToProtos(response.Data),
	}, nil
}
