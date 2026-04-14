package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shar1mo/wishlist-api/internal/model"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *model.Item, userID int64) error {
	query := `
		INSERT INTO wishlist_items (wishlist_id, title, description, product_url, priority)
		SELECT w.id, $2, $3, $4, $5
		FROM wishlists w
		WHERE w.id = $1 AND w.user_id = $6
		RETURNING id, is_reserved, reserved_at, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		item.WishlistID,
		item.Title,
		item.Description,
		item.ProductURL,
		item.Priority,
		userID,
	).Scan(
		&item.ID,
		&item.IsReserved,
		&item.ReservedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *ItemRepository) ListByWishlistIDAndUserID(ctx context.Context, wishlistID, userID int64) ([]model.Item, error) {
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
		WHERE wi.wishlist_id = $1
		  AND w.user_id = $2
		ORDER BY wi.created_at DESC
	`

	rows, err := r.db.Query(ctx, query, wishlistID, userID)
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

func (r *ItemRepository) GetByIDAndWishlistIDAndUserID(ctx context.Context, itemID, wishlistID, userID int64) (*model.Item, error) {
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
		WHERE wi.id = $1
		  AND wi.wishlist_id = $2
		  AND w.user_id = $3
	`

	var item model.Item

	err := r.db.QueryRow(ctx, query, itemID, wishlistID, userID).Scan(
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

func (r *ItemRepository) Update(ctx context.Context, item *model.Item, userID int64) error {
	query := `
		UPDATE wishlist_items wi
		SET title = $1,
		    description = $2,
		    product_url = $3,
		    priority = $4,
		    updated_at = NOW()
		FROM wishlists w
		WHERE wi.id = $5
		  AND wi.wishlist_id = $6
		  AND w.id = wi.wishlist_id
		  AND w.user_id = $7
		RETURNING wi.is_reserved, wi.reserved_at, wi.created_at, wi.updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		item.Title,
		item.Description,
		item.ProductURL,
		item.Priority,
		item.ID,
		item.WishlistID,
		userID,
	).Scan(
		&item.IsReserved,
		&item.ReservedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, itemID, wishlistID, userID int64) error {
	query := `
		DELETE FROM wishlist_items wi
		USING wishlists w
		WHERE wi.id = $1
		  AND wi.wishlist_id = $2
		  AND w.id = wi.wishlist_id
		  AND w.user_id = $3
	`

	tag, err := r.db.Exec(ctx, query, itemID, wishlistID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}