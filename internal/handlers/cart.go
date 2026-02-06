package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"clothes-store/internal/mailer"
	"clothes-store/internal/models"

	"github.com/gorilla/mux"
)

type OrderRequest struct {
	ProductIDs []models.OrderItemInput `json:"items"`
}

type OrderHandler struct {
	OrderModel *models.OrderModel
	UserModel  *models.UserModel
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")

	userID, ok := userIDVal.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req OrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	orderID, err := h.OrderModel.Create(userID, req.ProductIDs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user, err := h.UserModel.GetByID(userID)
	if err == nil {
		mailer.EmailQueue <- user.Email
	} else {
		log.Printf("Could not fetch email for user %d", userID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": orderID,
		"status":   "Created and Email sent to Queue",
	})
}

func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("userID")
	userID, ok := userIDVal.(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
