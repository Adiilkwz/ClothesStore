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

type OrderModel struct {
  DB *sql.DB
}

func (m *OrderModel) Create(userID int, items []OrderItemInput) (int, error) {
  tx, err := m.DB.Begin()
  if err != nil {
    return 0, err
  }

  // 1. Insert Order Header
  var orderID int
  stmtOrder := `INSERT INTO orders (user_id, total_price, status, created_at) 
                VALUES ($1, 0, 'Pending', NOW()) RETURNING id`
  
  err = tx.QueryRow(stmtOrder, userID).Scan(&orderID)
  if err != nil {
    tx.Rollback()
    return 0, err
  }

  // 2. Insert Items into order_items
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