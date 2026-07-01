package services_test

import (
	"testing"

	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestCanUpdateAGolden(t *testing.T) {
	var request services.GoldenUpdateRequest = services.GoldenUpdateRequest{
		Id:   shared.UUIDv4(),
		Name: shared.RandomPersonalName(),
	}

	var service services.GoldenUpdateService = services.NewGoldenUpdateService(repository, "/tmp/test")
	err := service.Do(request)

	assert.Nil(t, err)
	assert.True(t, repository.UpdateHaveBeenCalledWith(request.Id, request.Name))
	assert.True(t, repository.UpdateHaveBeenCalledOneWith(request.Id))
}

func TestCantUpdateWithoutName(t *testing.T) {
	var request services.GoldenUpdateRequest = services.GoldenUpdateRequest{
		Id: shared.UUIDv4(),
	}

	var service services.GoldenUpdateService = services.NewGoldenUpdateService(repository, "/tmp/test")
	err := service.Do(request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrGoldenBadRequest)
	assert.False(t, repository.UpdateHaveBeenCalled())
}

func TestCantUpdateWithoutId(t *testing.T) {
	var request services.GoldenUpdateRequest = services.GoldenUpdateRequest{
		Name: shared.RandomPersonalName(),
	}

	var service services.GoldenUpdateService = services.NewGoldenUpdateService(repository, "/tmp/test")
	err := service.Do(request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrGoldenBadRequest)
	assert.False(t, repository.UpdateHaveBeenCalled())
}

func TestCantUpdateWithInvalidId(t *testing.T) {
	var request services.GoldenUpdateRequest = services.GoldenUpdateRequest{
		Id:   "invalid-id",
		Name: shared.RandomPersonalName(),
	}

	var service services.GoldenUpdateService = services.NewGoldenUpdateService(repository, "/tmp/test")
	err := service.Do(request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrGoldenBadRequest)
	assert.False(t, repository.UpdateHaveBeenCalled())
}
