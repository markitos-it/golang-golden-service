package services

import (
	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/domain/types"
)

type GoldenUpdateRequest struct {
	Id         string `json:"id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Content    string `json:"content"`
	PosterData string `json:"poster_data"`
}

type GoldenUpdateService struct {
	Repository model.GoldenRepository
	BaseDir    string
}

func NewGoldenUpdateService(repository model.GoldenRepository, baseDir string) GoldenUpdateService {
	return GoldenUpdateService{
		Repository: repository,
		BaseDir:    baseDir,
	}
}

func (s GoldenUpdateService) Do(request GoldenUpdateRequest) error {
	securedId, err := types.NewGoldenId(request.Id)
	if err != nil {
		return err
	}

	securedName, err := types.NewGoldenName(request.Name)
	if err != nil {
		return err
	}

	securedContent, err := types.NewGoldenContent(request.Content)
	if err != nil {
		return err
	}

	golden, err := s.Repository.One(securedId)
	if err != nil {
		return err
	}

	golden.Name = securedName.Value()
	golden.Content = securedContent.Value()

	if request.PosterData != "" {
		posterPath, err := shared.SaveBase64File(s.BaseDir, request.PosterData)
		if err != nil {
			return err
		}
		golden.Poster = posterPath
	}

	return s.Repository.Update(golden)
}
