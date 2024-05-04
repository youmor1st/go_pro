package models

import "time"

type Order struct {
	ID          int
	UserID      int
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
}

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
}
