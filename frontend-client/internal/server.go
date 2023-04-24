package internal

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"shop/frontend-client/internal/handlers"
)

func Server() *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	fs := http.FileServer(http.Dir("./templates/static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fs))

	router.HandleFunc("/sign-in", handlers.Signin)
	router.HandleFunc("/sign-up", handlers.Signup)
	router.Get("/logout", handlers.Logout)

	router.Get("/profile", handlers.Profile)
	router.Get("/product/{product_id}", handlers.More)
	router.Get("/delete-product/{product_id}", handlers.DeleteProduct)
	router.Get("/create-product", handlers.CreateProduct)
	router.Post("/create-product", handlers.CreateProduct)

	router.Get("/my-orders", handlers.GetOrders)
	router.Post("/comment", handlers.CreateComment)
	router.Post("/rating", handlers.Rating)

	router.Post("/order", handlers.Order)

	router.Get("/search", handlers.Search)
	router.HandleFunc("/", handlers.Welcome)
	return &http.Server{
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 30,
		Addr:         ":8080",
		Handler:      router,
	}
}
