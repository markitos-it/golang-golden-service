package services_test

import (
	"os"
	"testing"

	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/types"
	internal_test "markitos-it-svc-golden/internal/testsuite/internal"
)

type MockSpyGoldenRepository struct {
	LastCreatedGoldenName     *string
	LastDeleteGoldenId        *string
	LastOneGoldenId           *string
	LastUpdatedGoldenId       *string
	LastUpdatedGoldenName     *string
	LastOneForUpdatedGoldenId *string
	LastAllHaveBeenCalled     bool
	LastUpdateHaveBeenCalled  bool
	LastSearchHaveBeenCalled  bool
}

func NewMockSpyGoldenRepository() *MockSpyGoldenRepository {
	return &MockSpyGoldenRepository{
		LastCreatedGoldenName:     nil,
		LastDeleteGoldenId:        nil,
		LastOneGoldenId:           nil,
		LastUpdatedGoldenId:       nil,
		LastUpdatedGoldenName:     nil,
		LastOneForUpdatedGoldenId: nil,
		LastAllHaveBeenCalled:     false,
		LastUpdateHaveBeenCalled:  false,
		LastSearchHaveBeenCalled:  false,
	}
}

func (m *MockSpyGoldenRepository) Create(golden *model.Golden) error {
	m.LastCreatedGoldenName = &golden.Name

	return nil
}

func (m *MockSpyGoldenRepository) CreateHaveBeenCalledWith(goldenName *string) bool {
	var result bool = m.LastCreatedGoldenName != nil && *m.LastCreatedGoldenName == *goldenName

	m.LastCreatedGoldenName = nil

	return result
}

func (m *MockSpyGoldenRepository) Delete(id *types.GoldenId) error {
	value := id.Value()
	m.LastDeleteGoldenId = &value

	return nil
}

func (m *MockSpyGoldenRepository) DeleteHaveBeenCalledWith(goldenId *string) bool {
	var result bool = m.LastDeleteGoldenId != nil && *m.LastDeleteGoldenId == *goldenId

	m.LastDeleteGoldenId = nil

	return result
}

func (m *MockSpyGoldenRepository) Update(golden *model.Golden) error {
	m.LastUpdateHaveBeenCalled = true
	m.LastUpdatedGoldenId = &golden.Id
	m.LastUpdatedGoldenName = &golden.Name
	m.LastOneForUpdatedGoldenId = &golden.Id

	return nil
}

func (m *MockSpyGoldenRepository) UpdateHaveBeenCalledWith(id, name string) bool {
	var matchCalled bool = m.LastUpdateHaveBeenCalled
	var matchId bool = *m.LastUpdatedGoldenId == id
	var matchName bool = *m.LastUpdatedGoldenName == name

	m.LastUpdatedGoldenId = nil
	m.LastUpdatedGoldenName = nil
	m.LastUpdateHaveBeenCalled = false

	return matchCalled && matchId && matchName
}

func (m *MockSpyGoldenRepository) UpdateHaveBeenCalled() bool {
	var matchCalled bool = m.LastUpdateHaveBeenCalled

	m.LastUpdateHaveBeenCalled = false
	m.LastUpdatedGoldenId = nil
	m.LastUpdatedGoldenName = nil

	return matchCalled
}

func (m *MockSpyGoldenRepository) UpdateHaveBeenCalledOneWith(id string) bool {
	var matchId bool = *m.LastOneForUpdatedGoldenId == id

	m.LastOneForUpdatedGoldenId = nil

	return matchId
}

func (m *MockSpyGoldenRepository) One(id *types.GoldenId) (*model.Golden, error) {
	value := id.Value()
	m.LastOneGoldenId = &value

	return internal_test.NewRandomGoldenWithCustomId(id), nil
}

func (m *MockSpyGoldenRepository) OneHaveBeenCalledWith(goldenId *string) bool {
	var result bool = m.LastOneGoldenId != nil && *m.LastOneGoldenId == *goldenId

	m.LastOneGoldenId = nil

	return result
}

func (m *MockSpyGoldenRepository) All() ([]*model.Golden, error) {
	m.LastAllHaveBeenCalled = true

	return nil, nil
}

func (m *MockSpyGoldenRepository) AllHaveBeenCalled() bool {
	result := m.LastAllHaveBeenCalled
	m.LastAllHaveBeenCalled = false

	return result
}

func (m *MockSpyGoldenRepository) SearchAndPaginate(
	searchTerm string,
	pageNumber int,
	pageSize int) ([]*model.Golden, error) {
	m.LastSearchHaveBeenCalled = true

	return nil, nil
}

func (m *MockSpyGoldenRepository) SearchHaveBeenCalled() bool {
	result := m.LastSearchHaveBeenCalled

	m.LastSearchHaveBeenCalled = false

	return result
}

var repository *MockSpyGoldenRepository

func TestMain(m *testing.M) {
	repository = NewMockSpyGoldenRepository()

	os.Exit(m.Run())
}

func (r *MockSpyGoldenRepository) PublishEvent(event *model.GoldenEvent) error {
	return nil
}
