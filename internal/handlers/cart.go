package handlers

import (
	"clothes-store/internal/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderRequest struct {
	UserID     int                     `json:"user_id"`
	ProductIDs []models.OrderItemInput `json:"items"`
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
		"status":   "Created and Email sent to Queue",
	})
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userID := 1

	orders, err := h.OrderModel.GetAllByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusBadRequest)
		return
	}

	items, err := h.OrderModel.GetItems(orderID)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Order ID", http.StatusBadRequest)
		return
	}

	err = h.OrderModel.UpdateStatus(orderID, "Cancelled")
	if err != nil {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order Cancelled Successfully"))
}
