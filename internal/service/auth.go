package service

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/shar1mo/wishlist-api/internal/auth"
	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type ValidationErrors map[string]string

func (v ValidationErrors) Error() string {
	return "validation failed"
}

type AuthService struct {
	userRepo repository.UserRepository
	jwt      *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*model.User, error) {
	if errs := validateRegisterInput(email, password); len(errs) > 0 {
		return nil, errs
	}

	email = strings.TrimSpace(strings.ToLower(email))

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := auth.CheckPasswordHash(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.jwt.Generate(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func validateRegisterInput(email, password string) ValidationErrors {
	errs := ValidationErrors{}

	email = strings.TrimSpace(email)

	if email == "" {
		errs["email"] = "email is required"
	} else if _, err := mail.ParseAddress(email); err != nil {
		errs["email"] = "invalid email format"
	}

	if password == "" {
		errs["password"] = "password is required"
	} else if len(password) < 8 {
		errs["password"] = "password must be at least 8 characters"
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}