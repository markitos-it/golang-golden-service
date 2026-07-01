package services_test

import (
	"markitos-it-svc-golden/internal/domain/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanGetAllResources(t *testing.T) {
	var service services.GoldenAllService = services.NewGoldenAllService(repository)
	golden, err := service.Do()

	assert.Nil(t, err)
	assert.IsType(t, services.GoldenAllResponse{}, *golden)
	assert.True(t, repository.AllHaveBeenCalled())
}
