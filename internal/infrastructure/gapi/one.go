package gapi

import (
	context "context"
	"encoding/base64"
	"errors"
	"log"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/domain/types"
	"os"
	"path/filepath"
)

func (s *Server) GetGolden(ctx context.Context, in *GetGoldenRequest) (*GetGoldenResponse, error) {
	if _, err := types.NewGoldenId(in.Id); err != nil {
		return nil, s.ToStatusError(shared.ErrInvalidGoldenId)
	}

	request := services.GoldenOneRequest{Id: in.Id}

	var service services.GoldenOneService = services.NewGoldenOneService(s.repository)
	response, err := service.Do(request)
	if err != nil {
		if errors.Is(err, shared.ErrGoldenNotFound) {
			log.Printf("[WARN] GetGolden: golden not found id=%s", in.Id)
		}
		return nil, s.ToStatusError(err)
	}

	posterData := response.Data.Poster
	if posterData != "" {
		filePath := filepath.Join(s.uploadsBaseDir(), filepath.Base(posterData))
		if fileBytes, err := os.ReadFile(filePath); err == nil {
			posterData = base64.StdEncoding.EncodeToString(fileBytes)
		} else {
			log.Printf("[WARNING] No se pudo leer el poster %s: %v\n", filePath, err)
		}
	}

	return &GetGoldenResponse{
		Id:      response.Data.Id,
		Name:    response.Data.Name,
		Poster:  posterData,
		Content: response.Data.Content,
	}, nil
}
