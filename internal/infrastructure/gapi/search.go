package gapi

import (
	context "context"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/shared"
)

const maxPageSize int32 = 100

func (s *Server) SearchGoldens(ctx context.Context, in *SearchGoldensRequest) (*SearchGoldensResponse, error) {
	if in.PageNumber < 1 {
		return nil, s.ToStatusError(shared.ErrInvalidPageNumber)
	}

	if in.PageSize < 1 || in.PageSize > maxPageSize {
		return nil, s.ToStatusError(shared.ErrInvalidPageSize)
	}

	var service services.GoldensearchService = services.NewGoldensearchService(s.repository)
	var request services.GoldensearchRequest = services.GoldensearchRequest{
		SearchTerm: in.SearchTerm,
		PageNumber: int(in.PageNumber),
		PageSize:   int(in.PageSize),
	}

	response, err := service.Do(request)
	if err != nil {
		return nil, s.ToStatusError(err)
	}

	domainGoldens := response.Data
	grpcCollection := s.ToProtos(domainGoldens)

	return &SearchGoldensResponse{
		Goldens: grpcCollection,
	}, nil
}
