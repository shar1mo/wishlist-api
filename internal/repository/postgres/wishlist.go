package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shar1mo/wishlist-api/internal/model"
)

type WishlistRepository struct {
	db *pgxpool.Pool
}

func NewWishlistRepository(db *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, wishlist *model.Wishlist) error {
	query := `
		INSERT INTO wishlists (user_id, title, description, event_date, public_token)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		wishlist.UserID,
		wishlist.Title,
		wishlist.Description,
		wishlist.EventDate,
		wishlist.PublicToken,
	).Scan(
		&wishlist.ID,
		&wishlist.CreatedAt,
		&wishlist.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *WishlistRepository) ListByUserID(ctx context.Context, userID int64) ([]model.Wishlist, error) {
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
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []model.Wishlist

	for rows.Next() {
		var wishlist model.Wishlist

		if err := rows.Scan(
			&wishlist.ID,
			&wishlist.UserID,
			&wishlist.Title,
			&wishlist.Description,
			&wishlist.EventDate,
			&wishlist.PublicToken,
			&wishlist.CreatedAt,
			&wishlist.UpdatedAt,
		); err != nil {
			return nil, err
		}

		wishlists = append(wishlists, wishlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (r *WishlistRepository) GetByIDAndUserID(ctx context.Context, wishlistID, userID int64) (*model.Wishlist, error) {
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
		WHERE id = $1 AND user_id = $2
	`

	var wishlist model.Wishlist

	err := r.db.QueryRow(ctx, query, wishlistID, userID).Scan(
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

func (r *WishlistRepository) Update(ctx context.Context, wishlist *model.Wishlist) error {
	query := `
		UPDATE wishlists
		SET title = $1,
		    description = $2,
		    event_date = $3,
		    updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		wishlist.Title,
		wishlist.Description,
		wishlist.EventDate,
		wishlist.ID,
		wishlist.UserID,
	).Scan(&wishlist.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *WishlistRepository) Delete(ctx context.Context, wishlistID, userID int64) error {
	query := `
		DELETE FROM wishlists
		WHERE id = $1 AND user_id = $2
	`

	tag, err := r.db.Exec(ctx, query, wishlistID, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}