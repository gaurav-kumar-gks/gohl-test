package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gks.com/gohl-test/internal/models"
	"gks.com/gohl-test/internal/repo"
	"go.uber.org/zap"
)

type UserHandler struct {
	repo *repo.UserRepository
	logger *zap.Logger
}

func NewUserHandler(repo *repo.UserRepository, logger *zap.Logger) *UserHandler {
	return &UserHandler{repo: repo, logger: logger}
}

func (u *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.repo.ListUsers(r.Context())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch users")
		u.logger.Error("List users failed", zap.Error(err))
		return
	}
	u.logger.Info("found users", zap.Any("users", users))
	jsonResponse(w, http.StatusOK, users)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := models.ValidateUser(&user); err != nil {
		message := fmt.Sprintf("invalid user: %v", err)
		errorResponse(w, http.StatusBadRequest, message)
		return
	}

	userModel, err := u.repo.CreateUser(r.Context(), &user)
	if err != nil {
		message := fmt.Sprintf("failed to create user: %v", err)
		errorResponse(w, http.StatusInternalServerError, message)
		return
	}
	jsonResponse(w, http.StatusCreated, userModel)	
}

func (u *UserHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := u.repo.GetUser(r.Context(), id)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to get user balance")
		u.logger.Error("Get user balance failed", zap.Error(err))
		return
	}
	jsonResponse(w, http.StatusOK, user)

}

func errorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}