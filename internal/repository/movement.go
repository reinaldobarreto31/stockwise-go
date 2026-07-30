package repository

import (
	"database/sql"

	"github.com/reinaldobarreto31/stockwise-go/internal/model"
)

// MovementRepository handles movement persistence.
type MovementRepository struct {
	db *sql.DB
}

// NewMovementRepository creates a new MovementRepository.
func NewMovementRepository(db *sql.DB) *MovementRepository {
	return &MovementRepository{db: db}
}

// Create inserts a new movement and updates the product stock in a single transaction.
func (r *MovementRepository) Create(input model.CreateMovementInput) (*model.Movement, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	var m model.Movement
	err = tx.QueryRow(`
		INSERT INTO movements (product_id, type, quantity, note)
		VALUES ($1, $2, $3, $4)
		RETURNING id, product_id, type, quantity, note, created_at
	`, input.ProductID, input.Type, input.Quantity, input.Note).Scan(
		&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.Note, &m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Adjust stock: +quantity for "in", -quantity for "out"
	var stockDelta int
	if input.Type == model.MovementIn {
		stockDelta = input.Quantity
	} else {
		stockDelta = -input.Quantity
	}

	_, err = tx.Exec(`
		UPDATE products SET stock = stock + $1, updated_at = NOW() WHERE id = $2
	`, stockDelta, input.ProductID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &m, nil
}

// GetByProductID returns all movements for a given product, most recent first.
func (r *MovementRepository) GetByProductID(productID int64) ([]model.Movement, error) {
	rows, err := r.db.Query(`
		SELECT id, product_id, type, quantity, note, created_at
		FROM movements
		WHERE product_id = $1
		ORDER BY created_at DESC
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []model.Movement
	for rows.Next() {
		var m model.Movement
		if err := rows.Scan(
			&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.Note, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, rows.Err()
}

// GetAll returns all movements, most recent first.
func (r *MovementRepository) GetAll() ([]model.Movement, error) {
	rows, err := r.db.Query(`
		SELECT id, product_id, type, quantity, note, created_at
		FROM movements
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []model.Movement
	for rows.Next() {
		var m model.Movement
		if err := rows.Scan(
			&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.Note, &m.CreatedAt,
		); err != nil {
			return nil, err
		}
		movements = append(movements, m)
	}
	return movements, rows.Err()
}
