package services

import (
	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
)

type GoldenCreateRequest struct {
	Name       string
	Content    string
	PosterData string
	/* ___CUSTOM_STRUCT_FIELDS___*/
}

type GoldenCreateResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Poster  string `json:"poster"`
	/* ___CUSTOM_STRUCT_FIELDS___*/
}

type GoldenCreateService struct {
	Repository model.GoldenRepository
	BaseDir    string
}

func NewGoldenCreateService(repository model.GoldenRepository, baseDir string) GoldenCreateService {
	return GoldenCreateService{
		Repository: repository,
		BaseDir:    baseDir,
	}
}

func (s GoldenCreateService) Do(request GoldenCreateRequest) (*GoldenCreateResponse, error) {
	posterPath := ""
	if request.PosterData != "" {
		var err error
		posterPath, err = shared.SaveBase64File(s.BaseDir, request.PosterData)
		if err != nil {
			return nil, err
		}
	}

	golden, err := model.NewGolden(shared.UUIDv4(), request.Name, request.Content, posterPath, s.BaseDir)
	if err != nil {
		return nil, err
	}

	if err := s.Repository.Create(golden); err != nil {
		return nil, err
	}

	return &GoldenCreateResponse{
		Id:      golden.Id,
		Name:    golden.Name,
		Content: golden.Content,
		Poster:  golden.Poster,
		/* ___CUSTOM_RESPONSE_FIELDS___*/
	}, nil
}
