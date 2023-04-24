package models

type CommentProduct struct {
	CommentId int    `json:"comment_id"`
	UserId    int    `json:"user_id"`
	ProductId int    `json:"product_id"`
	Message   string `json:"message"`
	Nickname  string `json:"nickname"`
}
