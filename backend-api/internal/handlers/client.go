package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"shop/backend-api/internal/models"
)

func (c *Client) CreateComment(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	strId := r.Context().Value(keyUserType(keyUser))

	str, ok := strId.(string)
	if !ok {
		log.Println("empty strId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, _ := strconv.Atoi(str)
	var comment models.CommentProduct
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	comment.UserId = userId
	comment.ProductId = productId

	err = c.productService.CommentProduct(comment)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Client) RateProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strId := r.Context().Value(keyUserType(keyUser))

	str, ok := strId.(string)
	if !ok {
		log.Println("empty strId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, _ := strconv.Atoi(str)

	var rating models.Rating
	err = json.NewDecoder(r.Body).Decode(&rating)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	rating.UserId = userId
	rating.ProductId = productId
	if rating.Rating <= 0 || rating.Rating > 10 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = c.productService.RateProduct(rating)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Client) CreateOrder(w http.ResponseWriter, r *http.Request) {

	productId, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println("1", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	strId := r.Context().Value(keyUserType(keyUser))

	str, ok := strId.(string)
	if !ok {
		log.Println("empty strId")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userId, err := strconv.Atoi(str)
	if err != nil {
		log.Println(err)
	}
	var order models.Order

	fmt.Println(order)
	err = json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		log.Println("2", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order.UserId = userId
	order.ProductId = productId
	err = c.productService.CreateOrder(order)
	if err == models.ErrNotEnoughQuantity {
		log.Println(err)
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Println("3", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Client) GetOrders(w http.ResponseWriter, r *http.Request) {
	strId := r.Context().Value(keyUserType(keyUser))

	str, ok := strId.(string)
	if !ok {
		log.Println("empty strId")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, _ := strconv.Atoi(str)
	orders, err := c.productService.GetOrdersByUserId(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&orders)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
