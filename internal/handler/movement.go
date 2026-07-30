package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/reinaldobarreto31/stockwise-go/internal/model"
	"github.com/reinaldobarreto31/stockwise-go/internal/service"
)

// MovementHandler handles movement endpoints.
type MovementHandler struct {
	movementSvc *service.MovementService
}

// NewMovementHandler creates a new MovementHandler.
func NewMovementHandler(movementSvc *service.MovementService) *MovementHandler {
	return &MovementHandler{movementSvc: movementSvc}
}

// Router returns a chi.Router with all movement routes mounted.
func (h *MovementHandler) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	return r
}

// List godoc
// GET /api/movements?product_id=X
func (h *MovementHandler) List(w http.ResponseWriter, r *http.Request) {
	productIDStr := r.URL.Query().Get("product_id")

	if productIDStr != "" {
		productID, err := strconv.ParseInt(productIDStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid product_id")
			return
		}

		movements, err := h.movementSvc.ListByProduct(productID)
		if err != nil {
			if errors.Is(err, service.ErrMovementProductNotFound) {
				writeError(w, http.StatusNotFound, "product not found")
				return
			}
			writeError(w, http.StatusInternalServerError, "failed to list movements")
			return
		}
		writeJSON(w, http.StatusOK, movements)
		return
	}

	movements, err := h.movementSvc.ListAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list movements")
		return
	}
	writeJSON(w, http.StatusOK, movements)
}

// Create godoc
// POST /api/movements
func (h *MovementHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateMovementInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	m, err := h.movementSvc.Create(input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMovementProductNotFound):
			writeError(w, http.StatusNotFound, "product not found")
		case errors.Is(err, service.ErrInsufficientStock):
			writeError(w, http.StatusUnprocessableEntity, "insufficient stock for outbound movement")
		case errors.Is(err, service.ErrInvalidQuantity):
			writeError(w, http.StatusBadRequest, "quantity must be greater than zero")
		case errors.Is(err, service.ErrInvalidMovementType):
			writeError(w, http.StatusBadRequest, "type must be 'in' or 'out'")
		default:
			writeError(w, http.StatusInternalServerError, "failed to create movement")
		}
		return
	}
	writeJSON(w, http.StatusCreated, m)
}
