package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/PeterM45/go-postgres-api/internal/auth"
	"github.com/PeterM45/go-postgres-api/internal/database"
	"github.com/PeterM45/go-postgres-api/internal/errors"
)

type UserHandler struct {
	store database.UserStore
	auth  *auth.JWT
}

type updateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

func NewUserHandler(store database.UserStore, auth *auth.JWT) *UserHandler {
	return &UserHandler{store: store, auth: auth}
}

func respondWithError(w http.ResponseWriter, err error) {
	var apiErr errors.APIError
	if e, ok := err.(errors.APIError); ok {
		apiErr = e
	} else if strings.Contains(err.Error(), "duplicate key") {
		apiErr = errors.ErrUserExists
	} else {
		apiErr = errors.ErrInternalServer
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Code)
	json.NewEncoder(w).Encode(apiErr)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// GetUsers handles GET /api/users
func (h *UserHandler) GetUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.store.GetUsers()
	if err != nil {
		respondWithError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// GetUser handles GET /api/users/[id]
func (h *UserHandler) GetUser(w http.ResponseWriter, _ *http.Request, id string) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, errors.APIError{Code: 400, Message: "invalid user ID"})
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		respondWithError(w, errors.ErrUserNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUser handles POST /api/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, errors.ErrInvalidInput)
		return
	}

	if req.Password == "" {
		respondWithError(w, errors.APIError{Code: 400, Message: "password is required"})
		return
	}

	user, err := h.store.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser handles PUT /api/users/[id]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, errors.APIError{Code: 400, Message: "invalid user ID"})
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, errors.ErrInvalidInput)
		return
	}

	user, err := h.store.UpdateUser(userID, req.Username, req.Email)
	if err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /api/users/[id]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(w, errors.APIError{Code: 400, Message: "invalid user ID"})
		return
	}

	if err := h.store.DeleteUser(userID); err != nil {
		respondWithError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted"})
}
