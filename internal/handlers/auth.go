package handlers

import (
	"encoding/json"
	"net/http"

	"clothes-store/internal/models"
)

type AuthHandler struct {
	UserModel *models.UserModel
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Address  string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err := h.UserModel.Insert(input.Name, input.Email, input.Password, input.Address)
	if err != nil {
		if err == models.ErrDuplicateEmail {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User Created Successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	id, err := h.UserModel.Authenticate(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	response := map[string]any{
		"message": "Login succesful",
		"user_id": id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := 1
	user, err := h.UserModel.GetByID(userID)
	if err != nil {
		http.Error(w, "User not found", 404)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId := 1
	var input struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	h.UserModel.Update(userId, input.Name, input.Address)
	w.Write([]byte("Profile updated"))
}
