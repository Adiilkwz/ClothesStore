package models

import (
	"database/sql"
	"fmt"
	"strings"
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

func (m *ProductModel) Insert(p Product) (int, error) {
	if p.Price < 0 {
		return 0, fmt.Errorf("price cannot be negative")
	}
	if p.StockQuantity < 0 {
		return 0, fmt.Errorf("stock quantity cannot be negative")
	}

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
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *ProductModel) GetAll(category, size string, minPrice, maxPrice, limit int) ([]Product, error) {
	baseQuery := `SELECT id, name, description, price_kzt, size, category, image_url, stock_quantity FROM products`

	var args []interface{}
	var whereClauses []string
	argCounter := 1

	if category != "" && category != "All" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, category)
		argCounter++
	}

	if size != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("size LIKE $%d", argCounter))
		args = append(args, "%"+size+"%")
		argCounter++
	}

	if minPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price_kzt >= $%d", argCounter))
		args = append(args, minPrice)
		argCounter++
	}

	if maxPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price_kzt <= $%d", argCounter))
		args = append(args, maxPrice)
		argCounter++
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += " ORDER BY id DESC"

	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, limit)
	}

	rows, err := m.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Size, &p.Category, &p.ImageURL, &p.StockQuantity)
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

func (m *ProductModel) Update(id int, p *Product) error {
	stmt := `UPDATE products 
             SET name = $1, price_kzt = $2, size = $3, category = $4, image_url = $5, stock_quantity = $6
             WHERE id = $7`

	_, err := m.DB.Exec(stmt,
		p.Name,
		p.Price,
		p.Size,
		p.Category,
		p.ImageURL,
		p.StockQuantity,
		id,
	)
	return err
}

func (m *ProductModel) Delete(id int) error {
	stmt := `DELETE FROM products WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)
	return err
}
