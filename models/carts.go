package models

import "time"

type CartItem struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type Cart struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Items      []*CartItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
}
