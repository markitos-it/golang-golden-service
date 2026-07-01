package services_test

import (
	"testing"

	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/services"
	"markitos-it-svc-golden/internal/domain/shared"

	"github.com/stretchr/testify/assert"
)

func TestCanCreateAUser(t *testing.T) {
	var golden model.Golden = model.Golden{
		Name: shared.RandomPersonalName(),
		/* ___CUSTOM_TEST_FIELDS___*/
	}
	var request services.GoldenCreateRequest = services.GoldenCreateRequest{
		Name: golden.Name,
		/* ___CUSTOM_TEST_FIELDS___*/
	}

	var service services.GoldenCreateService = services.NewGoldenCreateService(repository, "/tmp/test")
	response, err := service.Do(request)

	assert.Nil(t, err)
	assert.True(t, repository.CreateHaveBeenCalledWith(&request.Name))
	assert.Equal(t, response.Name, request.Name)
	assert.NotEmpty(t, response.Id)
}

func TestCantCreateWithoutName(t *testing.T) {
	var request services.GoldenCreateRequest = services.GoldenCreateRequest{
		/* ___CUSTOM_TEST_FIELDS___*/
	}

	var service services.GoldenCreateService = services.NewGoldenCreateService(repository, "/tmp/test")
	_, err := service.Do(request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrGoldenBadRequest)
	assert.False(t, repository.CreateHaveBeenCalledWith(&request.Name))
}

func TestCantCreateWithoutValidName(t *testing.T) {
	var request services.GoldenCreateRequest = services.GoldenCreateRequest{
		Name: "",
		/* ___CUSTOM_TEST_FIELDS___*/
	}

	var service services.GoldenCreateService = services.NewGoldenCreateService(repository, "/tmp/test")
	_, err := service.Do(request)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, shared.ErrGoldenBadRequest)
	assert.False(t, repository.CreateHaveBeenCalledWith(&request.Name))
}
