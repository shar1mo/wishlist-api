package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/shar1mo/wishlist-api/internal/model"
	"github.com/shar1mo/wishlist-api/internal/repository"
)

var (
	ErrPublicWishlistNotFound = errors.New("public wishlist not found")
	ErrPublicItemNotFound     = errors.New("item not found")
	ErrItemAlreadyReserved    = errors.New("item already reserved")
)

type PublicWishlistResponse struct {
	ID          int64              `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	EventDate   string             `json:"event_date"`
	Items       []PublicItemOutput `json:"items"`
}

type PublicItemOutput struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ProductURL  string `json:"product_url"`
	Priority    int    `json:"priority"`
	IsReserved  bool   `json:"is_reserved"`
}

type PublicService struct {
	publicRepo repository.PublicRepository
}

func NewPublicService(publicRepo repository.PublicRepository) *PublicService {
	return &PublicService{
		publicRepo: publicRepo,
	}
}

func (s *PublicService) GetWishlistByToken(ctx context.Context, token string) (*PublicWishlistResponse, error) {
	wishlist, err := s.publicRepo.GetWishlistByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPublicWishlistNotFound
		}
		return nil, err
	}

	items, err := s.publicRepo.ListItemsByWishlistToken(ctx, token)
	if err != nil {
		return nil, err
	}

	response := &PublicWishlistResponse{
		ID:          wishlist.ID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		EventDate:   wishlist.EventDate,
		Items:       make([]PublicItemOutput, 0, len(items)),
	}

	for _, item := range items {
		response.Items = append(response.Items, PublicItemOutput{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			ProductURL:  item.ProductURL,
			Priority:    item.Priority,
			IsReserved:  item.IsReserved,
		})
	}

	return response, nil
}

func (s *PublicService) ReserveItem(ctx context.Context, token string, itemID int64) (*model.Item, error) {
	wishlist, err := s.publicRepo.GetWishlistByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPublicWishlistNotFound
		}
		return nil, err
	}

	_ = wishlist

	item, err := s.publicRepo.ReserveItemByToken(ctx, token, itemID)
	if err == nil {
		return item, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	existingItem, checkErr := s.publicRepo.GetItemByTokenAndID(ctx, token, itemID)
	if checkErr != nil {
		if errors.Is(checkErr, pgx.ErrNoRows) {
			return nil, ErrPublicItemNotFound
		}
		return nil, checkErr
	}

	if existingItem.IsReserved {
		return nil, ErrItemAlreadyReserved
	}

	return nil, ErrPublicItemNotFound
}