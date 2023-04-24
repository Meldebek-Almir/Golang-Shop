package service

import (
	"errors"
	"fmt"
	"log"

	"shop/backend-api/internal/models"
	"shop/backend-api/internal/repository"
)

type ProductService interface {
	CreateNewProduct(product models.Product) error
	GetAllProducts() ([]models.Product, error)
	DeleteProduct(productId int) error
	GetProductById(productId int) (models.Product, error)
	GetProductByName(productName string) ([]models.Product, error)
	FilterProducts(minPrice, maxPrice, minRaiting, maxRaiting float64) ([]models.Product, error)
	CommentProduct(comment models.CommentProduct) error
	RateProduct(rate models.Rating) error
	GetRatingByProductId(productId int) (int, error)
	CreateOrder(order models.Order) error

	GetOrdersByUserId(userId int) ([]models.Order, error)
}
type productService struct {
	productRepo repository.ProductRepository
}

func (p *productService) CreateNewProduct(product models.Product) error {
	return p.productRepo.CreateNewProduct(product)
}

func (p *productService) GetAllProducts() ([]models.Product, error) {
	products, err := p.productRepo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	for i, v := range products {
		comments, err := p.productRepo.GetAllCommentByProductId(v.ProductId)
		if err != nil {
			return nil, err
		}
		rate, err := p.GetRatingByProductId(v.ProductId)
		if err != nil {
			return nil, err
		}
		products[i].Rating = int(rate)

		products[i].Comment = comments

	}
	return products, nil
}

func (p *productService) DeleteProduct(productId int) error {
	return p.productRepo.DeleteProduct(productId)
}

func (p *productService) GetProductById(productId int) (models.Product, error) {
	product, err := p.productRepo.GetProductById(productId)
	if err != nil {
		return models.Product{}, err
	}
	comments, err := p.productRepo.GetAllCommentByProductId(product.ProductId)
	if err != nil {
		return models.Product{}, err
	}
	rate, err := p.GetRatingByProductId(productId)
	if err != nil {
		return models.Product{}, err
	}
	product.Rating = int(rate)
	product.Comment = comments
	return product, nil
}

func (p *productService) GetProductByName(productName string) ([]models.Product, error) {
	products, err := p.productRepo.GetProductByName(productName)
	if err != nil {
		return nil, err
	}
	for i, v := range products {
		comments, err := p.productRepo.GetAllCommentByProductId(v.ProductId)
		if err != nil {
			return nil, err
		}
		products[i].Comment = comments

	}
	return products, nil
}

func (p *productService) FilterProducts(minPrice, maxPrice, minRaiting, maxRaiting float64) ([]models.Product, error) {
	products, err := p.productRepo.FilterProducts(minPrice, maxPrice, minRaiting, maxRaiting)
	if err != nil {
		return nil, err
	}
	for i, v := range products {
		comments, err := p.productRepo.GetAllCommentByProductId(v.ProductId)
		if err != nil {
			return nil, err
		}
		products[i].Comment = comments

	}
	return products, nil
}

func (p *productService) CommentProduct(comment models.CommentProduct) error {
	return p.productRepo.CommentProduct(comment)
}

func (p *productService) RateProduct(rate models.Rating) error {
	exitst, err := p.productRepo.ProductExists(rate.ProductId)
	if err != nil {
		return err
	}
	if !exitst {
		return errors.New("product id doesn't exists")
	}
	return p.productRepo.RateProduct(rate)
}

func (p *productService) CreateOrder(order models.Order) error {
	return p.productRepo.CreateOrder(order)
}

func (p *productService) GetRatingByProductId(productId int) (int, error) {
	ratings, err := p.productRepo.GetAllRatingsByProductId(productId)
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, v := range ratings {
		sum += v.Rating
	}
	if len(ratings) != 0 {
		return sum / len(ratings), nil
	} else {
		return 0, nil
	}
}

func (p *productService) GetOrdersByUserId(userId int) ([]models.Order, error) {
	orders, err := p.productRepo.GetOrdersByUserId(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for i, v := range orders {
		product, err := p.GetProductById(v.ProductId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		orders[i].ProductName = product.ProductName

	}
	return orders, nil
}
