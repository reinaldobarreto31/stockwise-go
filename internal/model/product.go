package model

import "time"

// Product represents a stock item.
type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	MinStock    int       `json:"min_stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductInput is the request body for creating a product.
type CreateProductInput struct {
	Name        string  `json:"name"`
	SKU         string  `json:"sku"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	MinStock    int     `json:"min_stock"`
}

// UpdateProductInput is the request body for updating a product.
type UpdateProductInput struct {
	Name        *string  `json:"name,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	MinStock    *int     `json:"min_stock,omitempty"`
}
