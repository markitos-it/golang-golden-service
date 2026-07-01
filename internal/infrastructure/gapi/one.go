package gapi

import (
	context "context"
	"encoding/base64"
	"log"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/types"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *Server) GetGolden(ctx context.Context, in *GetGoldenRequest) (*GetGoldenResponse, error) {
	if _, err := types.NewGoldenId(in.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	request := services.GoldenOneRequest{Id: in.Id}

	var service services.GoldenOneService = services.NewGoldenOneService(s.repository)
	response, err := service.Do(request)
	if err != nil {
		return nil, status.Error(s.GetGRPCCode(err), err.Error())

	}

	posterData := response.Data.Poster
	if posterData != "" {
		baseDir := os.Getenv("GOLDEN_UPLOADS_BASEDIR")
		filePath := filepath.Join(baseDir, posterData)
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
