package model

import "time"

type Item struct {
	ID          int64      `json:"id"`
	WishlistID  int64      `json:"wishlist_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ProductURL  string     `json:"product_url"`
	Priority    int        `json:"priority"`
	IsReserved  bool       `json:"is_reserved"`
	ReservedAt  *time.Time `json:"reserved_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}