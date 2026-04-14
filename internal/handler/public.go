package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/service"
)

type PublicService interface {
	GetWishlistByToken(ctx context.Context, token string) (*service.PublicWishlistResponse, error)
	ReserveItem(ctx context.Context, token string, itemID int64) (*model.Item, error)
}

type PublicHandler struct {
	publicService PublicService
}

func NewPublicHandler(publicService PublicService) *PublicHandler {
	return &PublicHandler{
		publicService: publicService,
	}
}

func (h *PublicHandler) GetWishlistByToken(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "invalid token")
		return
	}

	response, err := h.publicService.GetWishlistByToken(r.Context(), token)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPublicWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *PublicHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		writeError(w, http.StatusBadRequest, "invalid token")
		return
	}

	itemID, err := parseInt64Param(r, "itemId")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	item, err := h.publicService.ReserveItem(r.Context(), token, itemID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPublicWishlistNotFound):
			writeError(w, http.StatusNotFound, "wishlist not found")
		case errors.Is(err, service.ErrPublicItemNotFound):
			writeError(w, http.StatusNotFound, "item not found")
		case errors.Is(err, service.ErrItemAlreadyReserved):
			writeError(w, http.StatusConflict, "item already reserved")
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":          item.ID,
		"is_reserved": item.IsReserved,
		"reserved_at": item.ReservedAt,
	})
}