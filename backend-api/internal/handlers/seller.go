package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"shop/backend-api/internal/models"
)

func (c *Client) NewProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
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
	id, _ := strconv.Atoi(str)
	product.UserId = id
	product.Rating = 0.00
	err = c.productService.CreateNewProduct(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Client) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	user, err := c.authService.GetUserById(userId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	product, err := c.productService.GetProductById(productId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.Role == "admin" {
		err := c.productService.DeleteProduct(productId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if userId == product.UserId {
		err := c.productService.DeleteProduct(productId)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("user forbidden")
		w.WriteHeader(http.StatusForbidden)
		return
	}

}
