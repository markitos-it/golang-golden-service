package model

import (
	"fmt"
	"markitos-it-svc-golden/internal/domain/shared"
	"time"
)

const (
	GOLDEN_EVENT_CREATED = "golden-created"
	GOLDEN_EVENT_DELETED = "golden-deleted"
	GOLDEN_EVENT_UPDATED = "golden-updated"
)

const (
	GOLDEN_EVENT_STATUS_CREATED = "created"
)

type GoldenEvent struct {
	Id         string    `json:"id" binding:"required,uuid"`
	Name       string    `json:"name" binding:"required"`
	EntityId   string    `json:"entity_id" binding:"required,uuid"`
	EntityName string    `json:"entity_name" binding:"required"`
	Payload    string    `json:"payload" binding:"required"`
	Status     string    `json:"status" binding:"required" default:"created"`
	CreatedAt  time.Time `json:"created_at" binding:"required,datetime" default:"now"`
	UpdatedAt  time.Time `json:"updated_at" binding:"required,datetime" default:"now"`
}

func NewGoldenCreatedEvent(entity *Golden) (*GoldenEvent, error) {
	return createdGoldenEvent(GOLDEN_EVENT_CREATED, entity)
}

func NewGoldenDeletedEvent(entity *Golden) (*GoldenEvent, error) {
	return createdGoldenEvent(GOLDEN_EVENT_DELETED, entity)
}

func NewGoldenUpdatedEvent(entity *Golden) (*GoldenEvent, error) {
	return createdGoldenEvent(GOLDEN_EVENT_UPDATED, entity)
}

func createdGoldenEvent(eventName string, entity *Golden) (*GoldenEvent, error) {
	jsonPayload, err := entity.ToJSONString()
	if err != nil {
		return nil, fmt.Errorf("error creating Update Payload JSON: %w", err)
	}

	return &GoldenEvent{
		Id:         shared.UUIDv4(),
		Name:       eventName,
		EntityId:   entity.Id,
		EntityName: entity.GetEntityName(),
		Payload:    string(jsonPayload),
		Status:     GOLDEN_EVENT_STATUS_CREATED,
	}, nil
}
