package model

import "markitos-it-svc-golden/internal/domain/types"

type GoldenRepository interface {
	Create(golden *Golden) error
	Delete(id *types.GoldenId) error
	One(id *types.GoldenId) (*Golden, error)
	Update(golden *Golden) error
	All() ([]*Golden, error)
	SearchAndPaginate(searchTerm string, pageNumber int, pageSize int) ([]*Golden, error)
	PublishEvent(event *GoldenEvent) error
}
