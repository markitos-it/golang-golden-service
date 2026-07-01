package services

import (
	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/types"
)

type GoldenOneRequest struct {
	Id string `json:"id"`
}

type GoldenOneResponse struct {
	Data *model.Golden `json:"data"`
}

type GoldenOneService struct {
	Repository model.GoldenRepository
}

func NewGoldenOneService(repository model.GoldenRepository) GoldenOneService {
	return GoldenOneService{
		Repository: repository,
	}
}

func (s GoldenOneService) Do(request GoldenOneRequest) (*GoldenOneResponse, error) {
	securedId, err := types.NewGoldenId(request.Id)
	if err != nil {
		return nil, err
	}

	golden, err := s.Repository.One(securedId)
	if err != nil {
		return nil, err
	}

	return &GoldenOneResponse{golden}, nil
}
