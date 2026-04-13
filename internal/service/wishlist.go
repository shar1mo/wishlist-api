package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/repository"
)

var ErrWishlistNotFound = errors.New("wishlist not found")

type WishlistService struct {
	wishlistRepo repository.WishlistRepository
}

func NewWishlistService(wishlistRepo repository.WishlistRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
	}
}

func (s *WishlistService) Create(ctx context.Context, userID int64, title, description, eventDate string) (*model.Wishlist, error) {
	if errs := validateWishlistInput(title, description, eventDate); len(errs) > 0 {
		return nil, errs
	}

	token, err := generatePublicToken()
	if err != nil {
		return nil, err
	}

	wishlist := &model.Wishlist{
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		EventDate:   eventDate,
		PublicToken: token,
	}

	if err := s.wishlistRepo.Create(ctx, wishlist); err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistService) ListByUserID(ctx context.Context, userID int64) ([]model.Wishlist, error) {
	return s.wishlistRepo.ListByUserID(ctx, userID)
}

func (s *WishlistService) GetByID(ctx context.Context, wishlistID, userID int64) (*model.Wishlist, error) {
	wishlist, err := s.wishlistRepo.GetByIDAndUserID(ctx, wishlistID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistService) Update(ctx context.Context, wishlistID, userID int64, title, description, eventDate string) (*model.Wishlist, error) {
	if errs := validateWishlistInput(title, description, eventDate); len(errs) > 0 {
		return nil, errs
	}

	wishlist, err := s.wishlistRepo.GetByIDAndUserID(ctx, wishlistID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	wishlist.Title = strings.TrimSpace(title)
	wishlist.Description = strings.TrimSpace(description)
	wishlist.EventDate = eventDate

	if err := s.wishlistRepo.Update(ctx, wishlist); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	return wishlist, nil
}

func (s *WishlistService) Delete(ctx context.Context, wishlistID, userID int64) error {
	err := s.wishlistRepo.Delete(ctx, wishlistID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrWishlistNotFound
		}
		return err
	}

	return nil
}

func validateWishlistInput(title, description, eventDate string) ValidationErrors {
	errs := ValidationErrors{}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)

	if title == "" {
		errs["title"] = "title is required"
	} else if len(title) > 255 {
		errs["title"] = "title must be at most 255 characters"
	}

	if len(description) > 1000 {
		errs["description"] = "description must be at most 1000 characters"
	}

	if eventDate == "" {
		errs["event_date"] = "event_date is required"
	} else if _, err := time.Parse("2006-01-02", eventDate); err != nil {
		errs["event_date"] = "event_date must be in YYYY-MM-DD format"
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func generatePublicToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}