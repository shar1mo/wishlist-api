package service

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/repository"
)

var ErrItemNotFound = errors.New("item not found")

type ItemService struct {
	itemRepo repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) *ItemService {
	return &ItemService{
		itemRepo: itemRepo,
	}
}

func (s *ItemService) Create(
	ctx context.Context,
	userID, wishlistID int64,
	title, description, productURL string,
	priority int,
) (*model.Item, error) {
	if errs := validateItemInput(title, description, productURL, priority); len(errs) > 0 {
		return nil, errs
	}

	item := &model.Item{
		WishlistID:  wishlistID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		ProductURL:  strings.TrimSpace(productURL),
		Priority:    priority,
	}

	if err := s.itemRepo.Create(ctx, item, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWishlistNotFound
		}
		return nil, err
	}

	return item, nil
}

func (s *ItemService) ListByWishlistID(ctx context.Context, userID, wishlistID int64) ([]model.Item, error) {
	items, err := s.itemRepo.ListByWishlistIDAndUserID(ctx, wishlistID, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ItemService) GetByID(ctx context.Context, userID, wishlistID, itemID int64) (*model.Item, error) {
	item, err := s.itemRepo.GetByIDAndWishlistIDAndUserID(ctx, itemID, wishlistID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (s *ItemService) Update(
	ctx context.Context,
	userID, wishlistID, itemID int64,
	title, description, productURL string,
	priority int,
) (*model.Item, error) {
	if errs := validateItemInput(title, description, productURL, priority); len(errs) > 0 {
		return nil, errs
	}

	item, err := s.itemRepo.GetByIDAndWishlistIDAndUserID(ctx, itemID, wishlistID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	item.Title = strings.TrimSpace(title)
	item.Description = strings.TrimSpace(description)
	item.ProductURL = strings.TrimSpace(productURL)
	item.Priority = priority

	if err := s.itemRepo.Update(ctx, item, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return item, nil
}

func (s *ItemService) Delete(ctx context.Context, userID, wishlistID, itemID int64) error {
	if err := s.itemRepo.Delete(ctx, itemID, wishlistID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrItemNotFound
		}
		return err
	}

	return nil
}

func validateItemInput(title, description, productURL string, priority int) ValidationErrors {
	errs := ValidationErrors{}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	productURL = strings.TrimSpace(productURL)

	if title == "" {
		errs["title"] = "title is required"
	} else if len(title) > 255 {
		errs["title"] = "title must be at most 255 characters"
	}

	if len(description) > 1000 {
		errs["description"] = "description must be at most 1000 characters"
	}

	if productURL != "" {
		parsedURL, err := url.ParseRequestURI(productURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			errs["product_url"] = "invalid product_url"
		}
	}

	if priority < 1 || priority > 5 {
		errs["priority"] = "priority must be between 1 and 5"
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}