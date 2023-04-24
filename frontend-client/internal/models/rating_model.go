package models

type Rating struct {
	RatingId  int `json:"rating_id"`
	ProductId int `json:"product_id"`
	UserId    int `json:"user_id"`
	Rating    int `json:"rating"`
}
