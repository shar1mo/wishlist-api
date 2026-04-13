package repository

import (
	"context"

	"github.com/shar1mo/wishlist-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *model.Wishlist) error
	ListByUserID(ctx context.Context, userID int64) ([]model.Wishlist, error)
	GetByIDAndUserID(ctx context.Context, wishlistID, userID int64) (*model.Wishlist, error)
	Update(ctx context.Context, wishlist *model.Wishlist) error
	Delete(ctx context.Context, wishlistID, userID int64) error
}