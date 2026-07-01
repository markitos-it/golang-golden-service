package services

import "markitos-it-svc-golden/internal/domain/model"

type GoldenAllResponse struct {
	Data []*model.Golden `json:"data"`
}

type GoldenAllService struct {
	Repository model.GoldenRepository
}

func NewGoldenAllService(repository model.GoldenRepository) GoldenAllService {
	return GoldenAllService{
		Repository: repository,
	}
}

func (s GoldenAllService) Do() (*GoldenAllResponse, error) {
	goldens, err := s.Repository.All()
	if err != nil {
		return nil, err
	}

	return &GoldenAllResponse{
		Data: goldens,
	}, nil
}
