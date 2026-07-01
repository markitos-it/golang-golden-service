package database_test

import (
	"log"
	"testing"
	"time"

	"markitos-it-svc-golden/internal/domain/model"
	domain "markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/types"
	"markitos-it-svc-golden/internal/infrastructure/database"
	"markitos-it-svc-golden/internal/testsuite/infrastructure/testdb"
	internal_test "markitos-it-svc-golden/internal/testsuite/internal"

	"github.com/stretchr/testify/require"
)

func TestGoldenCreate(t *testing.T) {
	repository := database.NewGoldenPostgresRepository(testdb.GetDB())
	var golden *domain.Golden = internal_test.NewRandomGolden()

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", golden.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ?", golden.Id, golden.GetEntityName())
	})

	err := repository.Create(golden)
	require.NoError(t, err)

	var result domain.Golden
	err = db.First(&result, "id = ?", golden.Id).Error
	require.NoError(t, err, "La entidad Golden debería existir en la DB")
	require.Equal(t, golden.Id, result.Id)
	require.Equal(t, golden.Name, result.Name)
	require.WithinDuration(t, golden.CreatedAt, result.CreatedAt, time.Second)
	require.WithinDuration(t, golden.UpdatedAt, result.UpdatedAt, time.Second)

	var createdEvent model.GoldenEvent
	err = db.Where("entity_id = ? AND status = ? AND name = ?", golden.Id, model.GOLDEN_EVENT_STATUS_CREATED, model.GOLDEN_EVENT_CREATED).First(&createdEvent).Error
	require.NoError(t, err, "El evento outbox con name 'golden-created' debería existir con status 'created'")
	require.Equal(t, "golden", createdEvent.EntityName)
	require.NotEmpty(t, createdEvent.Payload)
}

func TestGoldenDelete(t *testing.T) {
	var golden *domain.Golden = internal_test.NewRandomGolden()
	err := testdb.GetRepository().Create(golden)
	require.NoError(t, err)

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", golden.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ?", golden.Id, golden.GetEntityName())
	})

	repository := database.NewGoldenPostgresRepository(db)
	id, err := types.NewGoldenId(golden.Id)
	require.NoError(t, err)

	err = repository.Delete(id)
	require.NoError(t, err)

	var createdEvent model.GoldenEvent
	err = db.Where("entity_id = ? AND status = ? AND name = ?", golden.Id, model.GOLDEN_EVENT_STATUS_CREATED, model.GOLDEN_EVENT_CREATED).First(&createdEvent).Error
	require.NoError(t, err, "El evento con status 'created' debería existir en la DB")
	require.Equal(t, "golden", createdEvent.EntityName)
	require.NotEmpty(t, createdEvent.Payload)

	var deletedEvent model.GoldenEvent
	err = db.Where("entity_id = ? AND status = ? AND name = ?", golden.Id, model.GOLDEN_EVENT_STATUS_CREATED, model.GOLDEN_EVENT_DELETED).First(&deletedEvent).Error
	require.NoError(t, err, "El evento con status 'created' debería existir en la DB")
	require.Equal(t, "golden", deletedEvent.EntityName)
	require.NotEmpty(t, deletedEvent.Payload)
}

func TestGoldenOne(t *testing.T) {
	var golden *domain.Golden = internal_test.NewRandomGolden()
	err := testdb.GetRepository().Create(golden)
	require.NoError(t, err)

	db := testdb.GetDB()
	repository := database.NewGoldenPostgresRepository(db)

	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", golden.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ?", golden.Id, golden.GetEntityName())
	})

	result, err := repository.One(golden.GetId())
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, golden.Id, result.Id)
	require.Equal(t, golden.Name, result.Name)
}

func TestGoldenUpdate(t *testing.T) {
	repository := database.NewGoldenPostgresRepository(testdb.GetDB())
	var golden *domain.Golden = internal_test.NewRandomGolden()
	log.Println("golden")
	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", golden.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ?", golden.Id, golden.GetEntityName())
	})

	err := repository.Create(golden)
	require.NoError(t, err)

	nuevoNombre := golden.Name + "Updated"
	golden.Name = nuevoNombre

	err = repository.Update(golden)
	require.NoError(t, err)

	var result domain.Golden
	err = db.First(&result, "id = ?", golden.Id).Error
	require.NoError(t, err, "La entidad Golden debería existir tras el Update")
	require.Equal(t, nuevoNombre, result.Name)

	var updatedEvent model.GoldenEvent
	err = db.Where("entity_id = ? AND status = ? AND name = ?", golden.Id, model.GOLDEN_EVENT_STATUS_CREATED, model.GOLDEN_EVENT_UPDATED).First(&updatedEvent).Error
	require.NoError(t, err, "El evento outbox con name 'golden-updated' debería existir con status 'created'")
	require.Equal(t, "golden", updatedEvent.EntityName)
	require.NotEmpty(t, updatedEvent.Payload)
	require.Contains(t, updatedEvent.Payload, nuevoNombre)
}
