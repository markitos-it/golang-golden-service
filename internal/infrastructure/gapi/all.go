package gapi

import (
	context "context"
	"markitos-it-svc-golden/internal/domain/services"

	status "google.golang.org/grpc/status"
)

func (s *Server) ListGoldens(ctx context.Context, in *ListGoldensRequest) (*ListGoldensResponse, error) {
	var service services.GoldenAllService = services.NewGoldenAllService(s.repository)
	response, err := service.Do()
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	return &ListGoldensResponse{
		Goldens: s.ToProtos(response.Data),
	}, nil
}
