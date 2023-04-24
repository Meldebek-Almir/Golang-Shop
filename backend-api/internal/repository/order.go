package repository

import (
	"errors"
	"fmt"
	"log"

	"shop/backend-api/internal/models"
)

func (p *productRepository) CreateOrder(order models.Order) error {
	// Начало транзакций
	tx, err := p.db.Begin()
	if err != nil {
		log.Println(err)
		return errors.New("failed to start transaction")
	}

	// Проверяем наличие продукта и достаточность его количества для покупки
	var availableQuantity int
	err = tx.QueryRow("SELECT available_quantity FROM products WHERE product_id = $1 FOR UPDATE", order.ProductId).Scan(&availableQuantity)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to retrieve product")
	}

	if availableQuantity < order.Quantity {
		tx.Rollback()
		return models.ErrNotEnoughQuantity
	}

	// Обновляем количество продуктов и количество проданных продуктов
	_, err = tx.Exec("UPDATE products SET available_quantity = available_quantity - $1, total_quantity_sold = total_quantity_sold + $1 WHERE product_id = $2", order.Quantity, order.ProductId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to update product quantity")
	}

	// Создаем заказ
	_, err = tx.Exec("INSERT INTO orders (product_id, user_id, quantity) VALUES ($1, $2, $3)", order.ProductId, order.UserId, order.Quantity)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return errors.New("failed to create order")
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return errors.New("failed to commit transaction")
	}

	return nil
}

func (p *productRepository) GetOrdersByUserId(userId int) ([]models.Order, error) {
	var orders []models.Order

	rows, err := p.db.Query("SELECT * FROM orders WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders for user %d: %w", userId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.OrderId, &order.ProductId, &order.UserId, &order.Quantity, &order.TotalPrice, &order.OrderDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}
		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to read order rows: %w", err)
	}

	return orders, nil
}
