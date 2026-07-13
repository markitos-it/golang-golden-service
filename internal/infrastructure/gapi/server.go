package gapi

import (
	"encoding/base64"
	"errors"
	"log"
	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/infrastructure/configuration"
	"os"
	"path/filepath"
	"strings"

	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	UnimplementedGoldenserviceServer
	address    string
	repository model.GoldenRepository
	config     configuration.GoldenConfiguration
}

func NewServer(address string, repository model.GoldenRepository, config configuration.GoldenConfiguration) *Server {
	apiGRPC := &Server{
		address:    address,
		repository: repository,
		config:     config,
	}

	return apiGRPC
}

func (s *Server) Repository() model.GoldenRepository {
	return s.repository
}

func (s *Server) GetGRPCCode(err error) codes.Code {
	var code codes.Code = codes.Internal

	switch {
	case errors.Is(err, shared.ErrGoldenNotFound):
		code = codes.NotFound
	case errors.Is(err, shared.ErrInvalidGoldenId),
		errors.Is(err, shared.ErrInvalidGoldenName),
		errors.Is(err, shared.ErrInvalidPageNumber),
		errors.Is(err, shared.ErrInvalidPageSize),
		strings.Contains(err.Error(), "invalid"),
		strings.Contains(err.Error(), "Invalid"),
		strings.Contains(err.Error(), "illegal"),
		strings.Contains(err.Error(), "bad request"):
		code = codes.InvalidArgument
	}

	return code
}

func (s *Server) ToStatusError(err error) error {
	code := s.GetGRPCCode(err)

	switch code {
	case codes.InvalidArgument:
		return status.Error(code, msgInvalidRequest)
	case codes.NotFound:
		return status.Error(code, msgResourceNotFound)
	default:
		return status.Error(code, msgInternalServerErr)
	}
}

func (s *Server) uploadsBaseDir() string {
	baseDir := os.Getenv("GOLDEN_UPLOADS_BASEDIR")
	if baseDir == "" {
		baseDir = s.config.BaseDir
	}
	if baseDir == "" {
		baseDir = "/tmp/goldens"
	}

	return filepath.Clean(baseDir)
}

func (s *Server) ToProtos(domainGoldens []*model.Golden) []*Golden {
	var protoGoldens []*Golden
	for _, golden := range domainGoldens {
		protoGoldens = append(protoGoldens, s.ToProto(golden))
	}

	return protoGoldens
}

func (s *Server) ToProto(domainGolden *model.Golden) *Golden {
	posterData := domainGolden.Poster
	if posterData != "" {
		filePath := filepath.Join(s.uploadsBaseDir(), filepath.Base(posterData))
		if fileBytes, err := os.ReadFile(filePath); err == nil {
			posterData = base64.StdEncoding.EncodeToString(fileBytes)
		} else {
			log.Printf("[WARNING] ToProto: No se pudo leer el poster %s: %v\n", filePath, err)
		}
	}

	return &Golden{
		Id:        domainGolden.Id,
		Name:      domainGolden.Name,
		Content:   domainGolden.Content,
		Poster:    posterData,
		CreatedAt: timestamppb.New(domainGolden.CreatedAt),
		UpdatedAt: timestamppb.New(domainGolden.UpdatedAt),
	}
}
