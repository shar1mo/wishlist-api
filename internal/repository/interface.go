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

type ItemRepository interface {
	Create(ctx context.Context, item *model.Item, userID int64) error
	ListByWishlistIDAndUserID(ctx context.Context, wishlistID, userID int64) ([]model.Item, error)
	GetByIDAndWishlistIDAndUserID(ctx context.Context, itemID, wishlistID, userID int64) (*model.Item, error)
	Update(ctx context.Context, item *model.Item, userID int64) error
	Delete(ctx context.Context, itemID, wishlistID, userID int64) error
}

type PublicRepository interface {
	GetWishlistByToken(ctx context.Context, token string) (*model.Wishlist, error)
	ListItemsByWishlistToken(ctx context.Context, token string) ([]model.Item, error)
	ReserveItemByToken(ctx context.Context, token string, itemID int64) (*model.Item, error)
	GetItemByTokenAndID(ctx context.Context, token string, itemID int64) (*model.Item, error)
}