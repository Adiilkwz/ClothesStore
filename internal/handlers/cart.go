package handlers

import (
	"clothes-store/internal/models"
	"encoding/json"
	"net/http"
)

type OrderRequest struct {
	UserID        int					   `json:"user_id"`
	ProductIDs    []models.OrderItemInput  `json:"items"`
}

type OrderHandler struct {
	OrderModel *models.OrderModel
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	orderID, err := h.OrderModel.Create(req.UserID, req.ProductIDs)
	if err != nil {
		http.Error(w, "Failed to create order: "+err.Error(), http.StatusInternalServerError)
		return
	}

	EmailQueue <- "customer@example.com"
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"status":  "Created and Email sent to Queue", 
	})
}