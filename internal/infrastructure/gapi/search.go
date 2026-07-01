package gapi

import (
	context "context"
	"markitos-it-svc-golden/internal/domain/services"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *Server) SearchGoldens(ctx context.Context, in *SearchGoldensRequest) (*SearchGoldensResponse, error) {
	if in.PageNumber < 1 {
		return nil, status.Error(codes.InvalidArgument, "invalid page number")
	}

	if in.PageSize < 1 {
		return nil, status.Error(codes.InvalidArgument, "invalid page size")
	}

	var service services.GoldensearchService = services.NewGoldensearchService(s.repository)
	var request services.GoldensearchRequest = services.GoldensearchRequest{
		SearchTerm: in.SearchTerm,
		PageNumber: int(in.PageNumber),
		PageSize:   int(in.PageSize),
	}

	response, err := service.Do(request)
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())
	}

	domainGoldens := response.Data
	grpcCollection := s.ToProtos(domainGoldens)

	return &SearchGoldensResponse{
		Goldens: grpcCollection,
	}, nil
}
