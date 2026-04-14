package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shar1mo/wishlist-api/internal/model"
)

type PublicRepository struct {
	db *pgxpool.Pool
}

func NewPublicRepository(db *pgxpool.Pool) *PublicRepository {
	return &PublicRepository{db: db}
}

func (r *PublicRepository) GetWishlistByToken(ctx context.Context, token string) (*model.Wishlist, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			description,
			TO_CHAR(event_date, 'YYYY-MM-DD') AS event_date,
			public_token,
			created_at,
			updated_at
		FROM wishlists
		WHERE public_token = $1
	`

	var wishlist model.Wishlist

	err := r.db.QueryRow(ctx, query, token).Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.Title,
		&wishlist.Description,
		&wishlist.EventDate,
		&wishlist.PublicToken,
		&wishlist.CreatedAt,
		&wishlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &wishlist, nil
}

func (r *PublicRepository) ListItemsByWishlistToken(ctx context.Context, token string) ([]model.Item, error) {
	query := `
		SELECT
			wi.id,
			wi.wishlist_id,
			wi.title,
			wi.description,
			wi.product_url,
			wi.priority,
			wi.is_reserved,
			wi.reserved_at,
			wi.created_at,
			wi.updated_at
		FROM wishlist_items wi
		JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE w.public_token = $1
		ORDER BY wi.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]model.Item, 0)

	for rows.Next() {
		var item model.Item

		if err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Title,
			&item.Description,
			&item.ProductURL,
			&item.Priority,
			&item.IsReserved,
			&item.ReservedAt,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PublicRepository) ReserveItemByToken(ctx context.Context, token string, itemID int64) (*model.Item, error) {
	query := `
		UPDATE wishlist_items wi
		SET is_reserved = TRUE,
		    reserved_at = NOW(),
		    updated_at = NOW()
		FROM wishlists w
		WHERE wi.id = $1
		  AND w.id = wi.wishlist_id
		  AND w.public_token = $2
		  AND wi.is_reserved = FALSE
		RETURNING
			wi.id,
			wi.wishlist_id,
			wi.title,
			wi.description,
			wi.product_url,
			wi.priority,
			wi.is_reserved,
			wi.reserved_at,
			wi.created_at,
			wi.updated_at
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, itemID, token).Scan(
		&item.ID,
		&item.WishlistID,
		&item.Title,
		&item.Description,
		&item.ProductURL,
		&item.Priority,
		&item.IsReserved,
		&item.ReservedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &item, nil
}

func (r *PublicRepository) GetItemByTokenAndID(ctx context.Context, token string, itemID int64) (*model.Item, error) {
	query := `
		SELECT
			wi.id,
			wi.wishlist_id,
			wi.title,
			wi.description,
			wi.product_url,
			wi.priority,
			wi.is_reserved,
			wi.reserved_at,
			wi.created_at,
			wi.updated_at
		FROM wishlist_items wi
		JOIN wishlists w ON w.id = wi.wishlist_id
		WHERE w.public_token = $1
		  AND wi.id = $2
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, token, itemID).Scan(
		&item.ID,
		&item.WishlistID,
		&item.Title,
		&item.Description,
		&item.ProductURL,
		&item.Priority,
		&item.IsReserved,
		&item.ReservedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &item, nil
}