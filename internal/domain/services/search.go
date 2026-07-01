package services

import (
	"markitos-it-svc-golden/internal/domain/model"
)

type GoldensearchResponse struct {
	Data []*model.Golden `json:"data"`
}

type GoldensearchRequest struct {
	SearchTerm string `json:"searchTerm"`
	PageNumber int    `json:"pageNumber" bindings:"min=1"`
	PageSize   int    `json:"pageSize" bindings:"min=10,max=100"`
}

type GoldensearchService struct {
	Repository model.GoldenRepository
}

func NewGoldensearchService(repository model.GoldenRepository) GoldensearchService {
	return GoldensearchService{
		Repository: repository,
	}
}

func (s GoldensearchService) Do(request GoldensearchRequest) (*GoldensearchResponse, error) {
	response, err := s.Repository.SearchAndPaginate(request.SearchTerm, request.PageNumber, request.PageSize)
	if err != nil {
		return nil, err
	}

	return &GoldensearchResponse{
		Data: response,
	}, nil
}
