package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/service"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthHandler struct {
	authService AuthService
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		var validationErrs service.ValidationErrors

		switch {
		case errors.As(err, &validationErrs):
			writeValidationError(w, validationErrs)
			return
		case errors.Is(err, service.ErrEmailAlreadyExists):
			writeError(w, http.StatusConflict, "email already exists")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}
	}

	writeJSON(w, http.StatusOK, loginResponse{
		Token: token,
	})
}