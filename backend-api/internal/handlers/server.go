package handlers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (c *Client) Server() *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Post("/register", AuthRedirectMiddleware(c.Register))
	router.Post("/login", AuthRedirectMiddleware(c.Login))
	router.Post("/api/products", SellerMiddleware(c.NewProduct))
	router.Post("/api/products/{product_id}/comments", ClientMiddleware(c.CreateComment))
	router.Post("/api/products/{product_id}/rate", ClientMiddleware(c.RateProduct))
	router.Post("/api/products/{product_id}/order", ClientMiddleware(c.CreateOrder))
	router.Post("/api/products/{product_id}/delete", SellerMiddleware(c.DeleteProduct))

	router.Get("/api/orders", ClientMiddleware(c.GetOrders))
	router.Get("/api/products", c.AllProducts)
	router.Get("/api/products/{product_id}", c.GetProduct)
	router.Get("/api/search/{query}", c.SearchProduct)
	router.Get("/api/filter", c.Filter)
	router.Get("/", c.Welcome)

	router.Post("/api/user/info", AuthNeedMiddleware(c.GetInfoUser))

	return &http.Server{
		ReadTimeout:  time.Second * time.Duration(c.config.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(c.config.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(c.config.IdleTimeout),
		Addr:         c.config.Port,
		Handler:      router,
	}
}
