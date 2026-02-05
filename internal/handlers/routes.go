package handlers

import (
	"clothes-store/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, auth *AuthHandler, store *StoreHandler, order *OrderHandler) {
	r.HandleFunc("/signup", auth.SignUp).Methods("POST")
	r.HandleFunc("/login", auth.Login).Methods("POST")

	r.HandleFunc("/products", store.GetAll).Methods("GET")
	r.HandleFunc("/products/{id}", store.GetOne).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.RequireAuth)

	api.HandleFunc("/users/me", auth.GetProfile).Methods("GET")
	api.HandleFunc("/users/me", auth.UpdateProfile).Methods("PUT")

	api.HandleFunc("/orders", order.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", order.GetMyOrders).Methods("GET")
	api.HandleFunc("/orders/{id}", order.GetOrderDetails).Methods("GET")
	api.HandleFunc("/orders/{id}/cancel", order.CancelOrder).Methods("PUT")

	api.HandleFunc("/products", store.Create).Methods("POST")
	api.HandleFunc("/products/{id}", store.Update).Methods("PUT")
	api.HandleFunc("/products/{id}", store.Delete).Methods("DELETE")
}
