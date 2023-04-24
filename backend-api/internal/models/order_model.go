package models

import "time"

type Order struct {
	OrderId     int       `json:"order_id"`
	ProductId   int       `json:"product_id"`
	UserId      int       `json:"user_id"`
	Quantity    int       `json:"quantity"`
	TotalPrice  int       `json:"total_price"`
	OrderDate   time.Time `json:"order_date"`
	ProductName string    `json:"product_name"`
}
