package handlers

import (
	"shop/backend-api/internal/service"
	"shop/backend-api/pkg/config"
)

type Client struct {
	config         *config.Config
	authService    service.AuthService
	productService service.ProductService
}

func NewClient(config *config.Config, authService service.AuthService, productService service.ProductService) *Client {
	return &Client{
		config:         config,
		authService:    authService,
		productService: productService,
	}
}
