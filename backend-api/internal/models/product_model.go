package models

type Product struct {
	ProductId         int              `json:"product_id"`
	ProductName       string           `json:"product_name"`
	Description       string           `json:"description"`
	Price             int              `json:"price"`
	Rating            int              `json:"rating"`
	ImageUrl          string           `json:"image_url"`
	UserId            int              `json:"user_id"`
	Comment           []CommentProduct `json:"comment"`
	AvailableQuantity int              `json:"available_quantity"`
	TotalQuantitySold int              `json:"total_quantity_sold"`
}
