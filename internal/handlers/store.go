package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"clothes-store/internal/models"
)

type StoreHandler struct {
	ProductModel *models.ProductModel
	Logger       *log.Logger
}

// GetAll returns all products as JSON
func (sh *StoreHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := sh.ProductModel.GetAll()
	if err != nil {
		sh.Logger.Printf("Error fetching products: %v", err)
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Create inserts a new product and returns 201 Created
func (sh *StoreHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Parse JSON input
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		sh.Logger.Printf("Error decoding product: %v", err)
		http.Error(w, "Invalid product data", http.StatusBadRequest)
		return
	}

	// Insert product into database
	id, err := sh.ProductModel.Insert(product)
	if err != nil {
		sh.Logger.Printf("Error inserting product: %v", err)
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	// Set the ID on the product and respond with 201 Created
	product.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
