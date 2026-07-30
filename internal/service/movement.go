package service

import (
	"errors"
	"fmt"

	"github.com/reinaldobarreto31/stockwise-go/internal/model"
	"github.com/reinaldobarreto31/stockwise-go/internal/repository"
)

var (
	ErrMovementProductNotFound = errors.New("product not found")
	ErrInsufficientStock       = errors.New("insufficient stock for outbound movement")
	ErrInvalidQuantity         = errors.New("quantity must be greater than zero")
	ErrInvalidMovementType     = errors.New("type must be 'in' or 'out'")
)

// MovementService handles movement business logic.
type MovementService struct {
	movementRepo *repository.MovementRepository
	productRepo  *repository.ProductRepository
}

// NewMovementService creates a new MovementService.
func NewMovementService(
	movementRepo *repository.MovementRepository,
	productRepo *repository.ProductRepository,
) *MovementService {
	return &MovementService{movementRepo: movementRepo, productRepo: productRepo}
}

// Create registers a new stock movement and updates product stock atomically.
func (s *MovementService) Create(input model.CreateMovementInput) (*model.Movement, error) {
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if input.Type != model.MovementIn && input.Type != model.MovementOut {
		return nil, ErrInvalidMovementType
	}

	product, err := s.productRepo.GetByID(input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	if product == nil {
		return nil, ErrMovementProductNotFound
	}

	if input.Type == model.MovementOut && product.Stock < input.Quantity {
		return nil, ErrInsufficientStock
	}

	m, err := s.movementRepo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("creating movement: %w", err)
	}
	return m, nil
}

// ListByProduct returns all movements for a product.
func (s *MovementService) ListByProduct(productID int64) ([]model.Movement, error) {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, fmt.Errorf("fetching product: %w", err)
	}
	if product == nil {
		return nil, ErrMovementProductNotFound
	}

	movements, err := s.movementRepo.GetByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("listing movements: %w", err)
	}
	if movements == nil {
		movements = []model.Movement{}
	}
	return movements, nil
}

// ListAll returns every movement.
func (s *MovementService) ListAll() ([]model.Movement, error) {
	movements, err := s.movementRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("listing movements: %w", err)
	}
	if movements == nil {
		movements = []model.Movement{}
	}
	return movements, nil
}
