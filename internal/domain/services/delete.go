package services

import (
	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/types"
)

type GoldenDeleteRequest struct {
	Id string `json:"id"`
}

type GoldenDeleteService struct {
	Repository model.GoldenRepository
}

func NewGoldenDeleteService(repository model.GoldenRepository) GoldenDeleteService {
	return GoldenDeleteService{
		Repository: repository,
	}
}

func (s GoldenDeleteService) Do(request GoldenDeleteRequest) error {
	securedId, err := types.NewGoldenId(request.Id)
	if err != nil {
		return err
	}

	if err := s.Repository.Delete(securedId); err != nil {
		return err
	}

	return nil
}
