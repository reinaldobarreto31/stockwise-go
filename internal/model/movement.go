package model

import "time"

// MovementType defines whether stock is entering or leaving.
type MovementType string

const (
	MovementIn  MovementType = "in"
	MovementOut MovementType = "out"
)

// Movement records a stock change event.
type Movement struct {
	ID        int64        `json:"id"`
	ProductID int64        `json:"product_id"`
	Type      MovementType `json:"type"`
	Quantity  int          `json:"quantity"`
	Note      string       `json:"note"`
	CreatedAt time.Time    `json:"created_at"`
}

// CreateMovementInput is the request body for registering a movement.
type CreateMovementInput struct {
	ProductID int64        `json:"product_id"`
	Type      MovementType `json:"type"`
	Quantity  int          `json:"quantity"`
	Note      string       `json:"note"`
}
