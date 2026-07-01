package model

import (
	"encoding/json"
	"log"
	"markitos-it-svc-golden/internal/domain/types"
	"time"
)

type Golden struct {
	Id        string    `json:"id" binding:"required,uuid"`
	Name      string    `json:"name" binding:"required"`
	Content   string    `json:"content"`
	Poster    string    `json:"poster"`
	CreatedAt time.Time `json:"created_at" binding:"required,datetime" default:"now"`
	UpdatedAt time.Time `json:"updated_at" binding:"required,datetime" default:"now"`
}

func NewGolden(id, name, content, poster, baseDir string) (*Golden, error) {
	secureId, err := types.NewGoldenId(id)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (NewGoldenId): %v\n", err)
		return nil, err
	}

	secureName, err := types.NewGoldenName(name)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (NewGoldenName): %v\n", err)
		return nil, err
	}

	secureContent, err := types.NewGoldenContent(content)
	if err != nil {
		log.Printf("❌ DEBUG ERROR (NewGoldenContent): %v\n", err)
		return nil, err
	}

	var securePoster *types.GoldenPoster
	if poster != "" {
		securePoster, err = types.NewGoldenPoster(baseDir, poster)
		if err != nil {
			log.Printf("❌ DEBUG ERROR (NewGoldenPoster): %v\n", err)
			return nil, err
		}
	}

	posterValue := ""
	if securePoster != nil {
		posterValue = securePoster.Value()
	}

	return &Golden{
		Id:        secureId.Value(),
		Name:      secureName.Value(),
		Content:   secureContent.Value(),
		Poster:    posterValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (h *Golden) GetEntityName() string {
	return "golden"
}

func (h *Golden) ToJSONString() (string, error) {
	bytes, err := json.Marshal(h)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *Golden) GetId() *types.GoldenId {
	result, _ := types.NewGoldenId(h.Id)
	return result
}

func (h *Golden) GetContent() *types.GoldenContent {
	result, _ := types.NewGoldenContent(h.Content)
	return result
}

func (h *Golden) GetPoster(baseDir string) *types.GoldenPoster {
	result, _ := types.NewGoldenPoster(baseDir, h.Poster)
	return result
}
