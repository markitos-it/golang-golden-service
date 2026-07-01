package database

import (
	"fmt"

	model "markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/domain/types"

	"gorm.io/gorm"
)

type GoldenPostgresRepository struct {
	db *gorm.DB
}

func NewGoldenPostgresRepository(db *gorm.DB) GoldenPostgresRepository {
	return GoldenPostgresRepository{db: db}
}

func (r *GoldenPostgresRepository) Create(golden *model.Golden) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(golden).Error; err != nil {
			return err
		}

		event, err := model.NewGoldenCreatedEvent(golden)
		if err != nil {
			return fmt.Errorf("error creating Create Event JSON: %w", err)
		}
		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GoldenPostgresRepository) Delete(id *types.GoldenId) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var golden model.Golden
		if err := tx.First(&golden, "id = ?", id.Value()).Error; err != nil {
			return shared.ErrGoldenNotFound
		}

		event, err := model.NewGoldenDeletedEvent(&golden)
		if err != nil {
			return fmt.Errorf("error creating Delete Event JSON: %w", err)
		}
		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.Golden{}, "id = ?", id.Value()).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GoldenPostgresRepository) Update(golden *model.Golden) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(golden).Error; err != nil {
			return err
		}

		event, err := model.NewGoldenUpdatedEvent(golden)
		if err != nil {
			return fmt.Errorf("error creating Update Event JSON: %w", err)
		}
		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		return nil
	})
}
func (r *GoldenPostgresRepository) One(id *types.GoldenId) (*model.Golden, error) {
	var golden model.Golden
	if err := r.db.First(&golden, "id = ?", id.Value()).Error; err != nil {
		return nil, shared.ErrGoldenNotFound
	}
	return &golden, nil
}

func (r *GoldenPostgresRepository) All() ([]*model.Golden, error) {
	var goldens []*model.Golden
	if err := r.db.Find(&goldens).Error; err != nil {
		return nil, shared.ErrGoldenBadRequest
	}

	return goldens, nil
}

func (r *GoldenPostgresRepository) SearchAndPaginate(
	searchTerm string,
	pageNumber int,
	pageSize int) ([]*model.Golden, error) {
	offset := (pageNumber - 1) * pageSize
	var goldens []*model.Golden
	if err := r.db.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", searchTerm)).
		Order("name").
		Limit(pageSize).
		Offset(offset).
		Find(&goldens).Error; err != nil {
		return nil, err
	}

	return goldens, nil
}

func (r *GoldenPostgresRepository) PublishEvent(event *model.GoldenEvent) error {
	return r.db.Create(event).Error
}
