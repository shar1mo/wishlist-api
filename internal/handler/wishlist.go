package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/service"
)

type WishlistService interface {
	Create(ctx context.Context, userID int64, title, description, eventDate string) (*model.Wishlist, error)
	ListByUserID(ctx context.Context, userID int64) ([]model.Wishlist, error)
	GetByID(ctx context.Context, wishlistID, userID int64) (*model.Wishlist, error)
	Update(ctx context.Context, wishlistID, userID int64, title, description, eventDate string) (*model.Wishlist, error)
	Delete(ctx context.Context, wishlistID, userID int64) error
}

type WishlistHandler struct {
	wishlistService WishlistService
	getUserID       func(ctx context.Context) (int64, bool)
}

type wishlistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	EventDate   string `json:"event_date"`
}

func NewWishlistHandler(wishlistService WishlistService, getUserID func(ctx context.Context) (int64, bool)) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
		getUserID:       getUserID,
	}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req wishlistRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wishlist, err := h.wishlistService.Create(r.Context(), userID, req.Title, req.Description, req.EventDate)
	if err != nil {
		var validationErrs service.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			writeValidationError(w, validationErrs)
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, wishlist)
}

func (h *WishlistHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlists, err := h.wishlistService.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"wishlists": wishlists,
	})
}

func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseInt64Param(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	wishlist, err := h.wishlistService.GetByID(r.Context(), wishlistID, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, wishlist)
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseInt64Param(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var req wishlistRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wishlist, err := h.wishlistService.Update(r.Context(), wishlistID, userID, req.Title, req.Description, req.EventDate)
	if err != nil {
		var validationErrs service.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			writeValidationError(w, validationErrs)
			return
		case errors.Is(err, service.ErrWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, wishlist)
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.getUserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := parseInt64Param(r, "wishlistId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.Delete(r.Context(), wishlistID, userID); err != nil {
		switch {
		case errors.Is(err, service.ErrWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}