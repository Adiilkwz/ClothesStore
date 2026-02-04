package models

import (
	"database/sql"
)

type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         int    `json:"price_kzt"` // Maps to price_kzt
	Size          string `json:"size"`      // Added
	Category      string `json:"category"`
	ImageURL      string `json:"image_url"`
	StockQuantity int    `json:"stock_quantity"` // Added
}

type ProductModel struct {
	DB *sql.DB
}

// Insert adds a new product
func (m *ProductModel) Insert(p Product) (int, error) {
	stmt := `INSERT INTO products (name, description, price_kzt, size, category, image_url, stock_quantity)
  VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt,
		p.Name,
		p.Description,
		p.Price,
		p.Size,
		p.Category,
		p.ImageURL,
		p.StockQuantity,
	).Scan(&id)
	return id, err
}

// GetAll fetches all products
func (m *ProductModel) GetAll() ([]Product, error) {
	stmt := `SELECT id, name, description, price_kzt, size, category, image_url, stock_quantity FROM products`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Size,
			&p.Category,
			&p.ImageURL,
			&p.StockQuantity,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (m *ProductModel) Get(id int) (*Product, error) {
	stmt := `SELECT id, name, description, price_kzt, size, category, image_url, stock_quantity FROM products WHERE id = $1`
	p := &Product{}
	err := m.DB.QueryRow(stmt, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Size, &p.Category, &p.ImageURL, &p.StockQuantity)
	return p, err
}

func (m *ProductModel) Update(id int, price, stock int) error {
	stmt := `UPDATE products SET price_kzt = $1, stock_quantity = $2 WHERE id = $3`
	_, err := m.DB.Exec(stmt, price, stock, id)
	return err
}

func (m *ProductModel) Delete(id int) error {
	stmt := `DELETE FROM products WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	return err
}
