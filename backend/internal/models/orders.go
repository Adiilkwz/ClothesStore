package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalPrice int       `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderItemInput struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderItem struct {
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
}

type OrderModel struct {
	DB *sql.DB
}

func (m *OrderModel) Create(userID int, items []OrderItemInput) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var totalPrice int
	for _, item := range items {
		var price int
		var currentStock int

		err := tx.QueryRow("SELECT price_kzt, stock_quantity FROM products WHERE id = $1", item.ProductID).Scan(&price, &currentStock)
		if err != nil {
			return 0, err
		}

		if currentStock < item.Quantity {
			return 0, fmt.Errorf("not enough stock for product ID %d (In Stock: %d)", item.ProductID, currentStock)
		}

		totalPrice += price * item.Quantity
	}

	stmtOrder := `INSERT INTO orders (user_id, total_price, status, created_at) 
                  VALUES ($1, $2, 'Paid', NOW()) RETURNING id`

	var orderID int
	err = tx.QueryRow(stmtOrder, userID, totalPrice).Scan(&orderID)
	if err != nil {
		return 0, err
	}

	stmtItem := `INSERT INTO order_items (order_id, product_id, quantity, price)
                 VALUES ($1, $2, $3, (SELECT price_kzt FROM products WHERE id = $2))`

	stmtUpdateStock := `UPDATE products 
                        SET stock_quantity = stock_quantity - $1 
                        WHERE id = $2`

	for _, item := range items {
		_, err := tx.Exec(stmtItem, orderID, item.ProductID, item.Quantity)
		if err != nil {
			return 0, err
		}

		_, err = tx.Exec(stmtUpdateStock, item.Quantity, item.ProductID)
		if err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (m *OrderModel) GetAllByUserID(userID int) ([]Order, error) {
	stmt := `SELECT id, user_id, total_price, status, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err = rows.Scan(&o.ID, &o.UserID, &o.TotalPrice, &o.Status, &o.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (m *OrderModel) GetItems(orderID int) ([]OrderItem, error) {
	stmt := `
    SELECT p.name, oi.price, oi.quantity
    FROM order_items oi
    JOIN products p ON oi.product_id = p.id
    WHERE oi.order_id = $1`

	rows, err := m.DB.Query(stmt, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.ProductName, &i.Price, &i.Quantity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (m *OrderModel) UpdateStatus(orderID int, status string) error {
	stmt := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := m.DB.Exec(stmt, status, orderID)
	return err
}
