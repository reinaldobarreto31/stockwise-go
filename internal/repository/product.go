package repository

import (
	"database/sql"

	"github.com/reinaldobarreto31/stockwise-go/internal/model"
)

// ProductRepository handles product persistence.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetAll returns all products.
func (r *ProductRepository) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, sku, category, description, price, stock, min_stock, created_at, updated_at
		FROM products ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.SKU, &p.Category, &p.Description,
			&p.Price, &p.Stock, &p.MinStock, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetByID returns a single product by ID.
func (r *ProductRepository) GetByID(id int64) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(`
		SELECT id, name, sku, category, description, price, stock, min_stock, created_at, updated_at
		FROM products WHERE id = $1
	`, id).Scan(
		&p.ID, &p.Name, &p.SKU, &p.Category, &p.Description,
		&p.Price, &p.Stock, &p.MinStock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

// Create inserts a new product and returns the created record.
func (r *ProductRepository) Create(input model.CreateProductInput) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(`
		INSERT INTO products (name, sku, category, description, price, stock, min_stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, sku, category, description, price, stock, min_stock, created_at, updated_at
	`, input.Name, input.SKU, input.Category, input.Description,
		input.Price, input.Stock, input.MinStock,
	).Scan(
		&p.ID, &p.Name, &p.SKU, &p.Category, &p.Description,
		&p.Price, &p.Stock, &p.MinStock, &p.CreatedAt, &p.UpdatedAt,
	)
	return &p, err
}

// Delete removes a product by ID.
func (r *ProductRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	return err
}
