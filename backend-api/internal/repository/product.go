package repository

import (
	"database/sql"
	"errors"
	"log"

	"shop/backend-api/internal/models"
)

type ProductRepository interface {
	CreateNewProduct(product models.Product) error
	DeleteProduct(productId int) error
	CommentProduct(comment models.CommentProduct) error
	RateProduct(rate models.Rating) error

	GetAllProducts() ([]models.Product, error)
	GetAllCommentByProductId(productId int) ([]models.CommentProduct, error)
	GetAllRaitValuesByProductId(productId int) ([]models.Rating, error)
	GetProductById(productId int) (models.Product, error)
	GetProductByName(productName string) ([]models.Product, error)

	FilterProducts(minPrice, maxPrice, minRaiting, maxRaiting float64) ([]models.Product, error)

	ProductExists(productId int) (bool, error)
	CreateOrder(order models.Order) error
	GetAllRatingsByProductId(productId int) ([]models.Rating, error)

	GetOrdersByUserId(userId int) ([]models.Order, error)
}

type productRepository struct {
	db *sql.DB
}

func (p *productRepository) CreateNewProduct(product models.Product) error {
	query := "INSERT INTO products (product_name, description, price, rating, image_url, user_id, available_quantity) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := p.db.Exec(query, product.ProductName, product.Description, product.Price, product.Rating, product.ImageUrl, product.UserId, product.AvailableQuantity)
	if err != nil {
		log.Println(err)
		return errors.New("failed to create new product")
	}
	return nil
}

func (p *productRepository) CommentProduct(comment models.CommentProduct) error {
	query := "INSERT INTO comments (user_id, product_id, message) VALUES ($1, $2, $3)"
	res, err := p.db.Exec(query, comment.UserId, comment.ProductId, comment.Message)
	if err != nil {
		log.Println(err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.New("failed to get rows affected")
	}
	if rowsAffected == 0 {
		return errors.New("product or user not found")
	}
	return nil
}

func (p *productRepository) DeleteProduct(productId int) error {
	// Удалить связанные с продуктом рейтинги, комментарии и заказы.
	if _, err := p.db.Exec("DELETE FROM ratings WHERE product_id = $1", productId); err != nil {
		log.Println(err)
		return errors.New("failed to delete ratings")
	}
	if _, err := p.db.Exec("DELETE FROM comments WHERE product_id = $1", productId); err != nil {
		log.Println(err)
		return errors.New("failed to delete comments")
	}
	if _, err := p.db.Exec("DELETE FROM orders WHERE product_id = $1", productId); err != nil {
		log.Println(err)
		return errors.New("failed to delete orders")
	}

	// Удалить сам продукт.
	query := "DELETE FROM products WHERE product_id = $1"
	res, err := p.db.Exec(query, productId)
	if err != nil {
		log.Println(err)
		return errors.New("failed to delete product")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.New("failed to get rows affected")
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (p *productRepository) RateProduct(rate models.Rating) error {
	query := `
	INSERT INTO ratings (user_id, product_id, rating)
	VALUES ($1, $2, $3)
	ON CONFLICT (user_id, product_id) DO UPDATE
	SET rating = EXCLUDED.rating
	`
	res, err := p.db.Exec(query, rate.UserId, rate.ProductId, rate.Rating)
	if err != nil {
		log.Println(err)
		return errors.New("failed to rate product")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.New("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (p *productRepository) GetAllProducts() ([]models.Product, error) {
	products := []models.Product{}

	rows, err := p.db.Query("SELECT * FROM products")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ProductId, &product.ProductName, &product.Description, &product.Price, &product.Rating, &product.ImageUrl, &product.UserId, &product.AvailableQuantity, &product.TotalQuantitySold)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) GetAllCommentByProductId(productId int) ([]models.CommentProduct, error) {
	comments := []models.CommentProduct{}

	rows, err := p.db.Query(`
        SELECT c.comment_id, c.user_id, c.product_id, c.message, u.nickname
        FROM comments c
        JOIN users u ON c.user_id = u.user_id
        WHERE c.product_id=$1
    `, productId)
	if err != nil {
		log.Println(err)
		return comments, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.CommentProduct
		err := rows.Scan(&comment.CommentId, &comment.UserId, &comment.ProductId, &comment.Message, &comment.Nickname)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return comments, err
	}

	return comments, nil
}

func (p *productRepository) GetAllRaitValuesByProductId(productId int) ([]models.Rating, error) {
	var ratings []models.Rating

	rows, err := p.db.Query("SELECT rating_id, product_id, user_id, rating FROM ratings WHERE product_id=$1", productId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(&rating.RatingId, &rating.ProductId, &rating.UserId, &rating.Rating)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return ratings, nil
}

func (p *productRepository) GetProductById(productId int) (models.Product, error) {
	var product models.Product

	err := p.db.QueryRow("SELECT *  FROM products WHERE product_id=$1", productId).Scan(&product.ProductId, &product.ProductName, &product.Description, &product.Price, &product.Rating, &product.ImageUrl, &product.UserId, &product.AvailableQuantity, &product.TotalQuantitySold)
	if err != nil {
		log.Println(err)
		return models.Product{}, err
	}

	return product, nil
}

func (p *productRepository) GetProductByName(productName string) ([]models.Product, error) {
	products := []models.Product{}

	rows, err := p.db.Query("SELECT * FROM products WHERE product_name=$1", productName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ProductId, &product.ProductName, &product.Description, &product.Price, &product.Rating, &product.ImageUrl, &product.UserId, &product.AvailableQuantity, &product.TotalQuantitySold)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) FilterProducts(minPrice, maxPrice, minRaiting, maxRaiting float64) ([]models.Product, error) {
	var products []models.Product
	query := "SELECT * FROM products WHERE price BETWEEN $1 AND $2 AND rating BETWEEN $3 AND $4"

	rows, err := p.db.Query(query, minPrice, maxPrice, minRaiting, maxRaiting)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		err = rows.Scan(&product.ProductId, &product.ProductName, &product.Description, &product.Price, &product.Rating, &product.ImageUrl, &product.UserId, &product.AvailableQuantity, &product.TotalQuantitySold)
		if err != nil {
			return nil, err
		}
		products = append(products, product)

	}
	return products, nil
}

func (p *productRepository) ProductExists(productId int) (bool, error) {
	query := "SELECT COUNT(*) FROM products WHERE product_id = $1"
	var count int
	err := p.db.QueryRow(query, productId).Scan(&count)
	if err != nil {
		log.Println(err)
		return false, errors.New("failed to check if product exists")
	}
	return count > 0, nil
}

func (p *productRepository) GetAllRatingsByProductId(productId int) ([]models.Rating, error) {
	ratings := []models.Rating{}

	rows, err := p.db.Query("SELECT * FROM ratings WHERE product_id=$1", productId)
	if err != nil {
		log.Println(err)
		return ratings, err
	}
	defer rows.Close()

	for rows.Next() {
		var rating models.Rating
		err := rows.Scan(&rating.RatingId, &rating.ProductId, &rating.UserId, &rating.Rating)
		if err != nil {
			return ratings, err
		}
		ratings = append(ratings, rating)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return ratings, err
	}

	return ratings, nil
}
