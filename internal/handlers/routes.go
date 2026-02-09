package handlers

import (
	"log"
	"net/http"
	"text/template"

	"clothes-store/internal/middleware"

	"github.com/gorilla/mux"
)

func render(w http.ResponseWriter, page string) {
	files := []string{
		"./ui/html/base.layout.html",
		"./ui/html/" + page,
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", 500)
	}
}

func RegisterRoutes(r *mux.Router, auth *AuthHandler, store *StoreHandler, order *OrderHandler) {
	fileServer := http.FileServer(http.Dir("./ui/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "home.page.html")
	}).Methods("GET")

	r.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		render(w, "product.page.html")
	}).Methods("GET")

	r.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {
		render(w, "cart.page.html")
	}).Methods("GET")

	r.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		render(w, "admin.page.html")
	}).Methods("GET")

	r.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		render(w, "profile.page.html")
	}).Methods("GET")

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

	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireAdmin)

	admin.HandleFunc("/products", store.Create).Methods("POST")
	admin.HandleFunc("/products/{id}", store.Update).Methods("PUT")
	admin.HandleFunc("/products/{id}", store.Delete).Methods("DELETE")
}
