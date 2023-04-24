package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (c *Client) Welcome(w http.ResponseWriter, r *http.Request) {
	message := make(map[string]string)

	message["status"] = "ok"
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&message)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (c *Client) AllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := c.productService.GetAllProducts()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Println(products)
	err = json.NewEncoder(w).Encode(&products)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Client) GetProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "product_id"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := c.productService.GetProductById(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Client) SearchProduct(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "query")
	product, err := c.productService.GetProductByName(name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Client) Filter(w http.ResponseWriter, r *http.Request) {
	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)
	minRating, _ := strconv.ParseFloat(r.URL.Query().Get("min_rate"), 64)
	maxRating, _ := strconv.ParseFloat(r.URL.Query().Get("max_rate"), 64)

	if maxPrice == 0 {
		maxPrice = math.MaxInt32
	}
	if maxRating == 0 {
		maxRating = 10
	}
	products, err := c.productService.FilterProducts(
		minPrice,
		maxPrice,
		minRating,
		maxRating,
	)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&products)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
