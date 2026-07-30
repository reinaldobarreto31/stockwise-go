package service

import (
	"errors"
	"fmt"

	"github.com/reinaldobarreto31/stockwise-go/internal/model"
	"github.com/reinaldobarreto31/stockwise-go/internal/repository"
)

var ErrProductNotFound = errors.New("product not found")

// ProductService handles product business logic.
type ProductService struct {
	productRepo *repository.ProductRepository
}

// NewProductService creates a new ProductService.
func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// List returns all products.
func (s *ProductService) List() ([]model.Product, error) {
	products, err := s.productRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("listing products: %w", err)
	}
	if products == nil {
		products = []model.Product{}
	}
	return products, nil
}

// Get returns a single product by ID.
func (s *ProductService) Get(id int64) (*model.Product, error) {
	p, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	if p == nil {
		return nil, ErrProductNotFound
	}
	return p, nil
}

// Create creates a new product.
func (s *ProductService) Create(input model.CreateProductInput) (*model.Product, error) {
	p, err := s.productRepo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}
	return p, nil
}

// Update applies partial updates to a product.
func (s *ProductService) Update(id int64, input model.UpdateProductInput) (*model.Product, error) {
	existing, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	if existing == nil {
		return nil, ErrProductNotFound
	}

	p, err := s.productRepo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("updating product: %w", err)
	}
	return p, nil
}

// Delete removes a product by ID.
func (s *ProductService) Delete(id int64) error {
	existing, err := s.productRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("fetching product: %w", err)
	}
	if existing == nil {
		return ErrProductNotFound
	}

	if err := s.productRepo.Delete(id); err != nil {
		return fmt.Errorf("deleting product: %w", err)
	}
	return nil
}
