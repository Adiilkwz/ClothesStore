package models

import (
	"database/sql"
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
	ProductID int
	Quantity  int
}

type OrderItem struct {
	ProductName string `json:"product_name"`
	Price       int    `json:"price_kzt"`
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

	var totalPrice int
	for _, item := range items {
		var price int
		err := tx.QueryRow("SELECT price_kzt FROM products WHERE id = $1", item.ProductID).Scan(&price)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		totalPrice += price * item.Quantity
	}

	var orderID int
	stmtOrder := `INSERT INTO orders (user_id, total_price, status, created_at) 
                VALUES ($1, $2, 'Pending', NOW()) RETURNING id`

	err = tx.QueryRow(stmtOrder, userID, totalPrice).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmtItems := `INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`

	for _, item := range items {
		_, err := tx.Exec(stmtItems, orderID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return orderID, tx.Commit()
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
    SELECT p.name, p.price_kzt, oi.quantity
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
