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

// ProductHandler handles product endpoints.
type ProductHandler struct {
	productSvc *service.ProductService
}

// NewProductHandler creates a new ProductHandler.
func NewProductHandler(productSvc *service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc: productSvc}
}

// Router returns a chi.Router with all product routes mounted.
func (h *ProductHandler) Router() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
	return r
}

// List godoc
// GET /api/products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.productSvc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list products")
		return
	}
	writeJSON(w, http.StatusOK, products)
}

// Get godoc
// GET /api/products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	p, err := h.productSvc.Get(id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch product")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// Create godoc
// POST /api/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" || input.SKU == "" {
		writeError(w, http.StatusBadRequest, "name and sku are required")
		return
	}

	p, err := h.productSvc.Create(input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create product")
		return
	}
	writeJSON(w, http.StatusCreated, p)
}

// Update godoc
// PUT /api/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var input model.UpdateProductInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.productSvc.Update(id, input)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update product")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

// Delete godoc
// DELETE /api/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.productSvc.Delete(id); err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			writeError(w, http.StatusNotFound, "product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseIDParam reads a URL param and converts it to int64.
func parseIDParam(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}
