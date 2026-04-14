package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/service"
)

type ItemService interface {
	Create(ctx context.Context, userID, wishlistID int64, title, description, productURL string, priority int) (*model.Item, error)
	ListByWishlistID(ctx context.Context, userID, wishlistID int64) ([]model.Item, error)
	GetByID(ctx context.Context, userID, wishlistID, itemID int64) (*model.Item, error)
	Update(ctx context.Context, userID, wishlistID, itemID int64, title, description, productURL string, priority int) (*model.Item, error)
	Delete(ctx context.Context, userID, wishlistID, itemID int64) error
}

type ItemHandler struct {
	itemService ItemService
	getUserID   func(ctx context.Context) (int64, bool)
}

type itemRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
}

func NewItemHandler(itemService ItemService, getUserID func(ctx context.Context) (int64, bool)) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
		getUserID:   getUserID,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseIDParam(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var req itemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.itemService.Create(
		r.Context(),
		userID,
		wishlistID,
		req.Title,
		req.Description,
		req.ProductURL,
		req.Priority,
	)
	if err != nil {
		var validationErrs service.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			writeValidationError(w, validationErrs)
		case errors.Is(err, service.ErrWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseIDParam(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	items, err := h.itemService.ListByWishlistID(r.Context(), userID, wishlistID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
	})
}

func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseIDParam(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	itemID, err := parseIDParam(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.itemService.GetByID(r.Context(), userID, wishlistID, itemID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			writeError(w, http.StatusNotFound, "item not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseIDParam(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	itemID, err := parseIDParam(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	var req itemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.itemService.Update(
		r.Context(),
		userID,
		wishlistID,
		itemID,
		req.Title,
		req.Description,
		req.ProductURL,
		req.Priority,
	)
	if err != nil {
		var validationErrs service.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			writeValidationError(w, validationErrs)
		case errors.Is(err, service.ErrItemNotFound):
			writeError(w, http.StatusNotFound, "item not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseIDParam(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	itemID, err := parseIDParam(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.itemService.Delete(r.Context(), userID, wishlistID, itemID); err != nil {
		switch {
		case errors.Is(err, service.ErrItemNotFound):
			writeError(w, http.StatusNotFound, "item not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(r *http.Request, key string) (int64, error) {
	value := chi.URLParam(r, key)
	return strconv.ParseInt(value, 10, 64)
}