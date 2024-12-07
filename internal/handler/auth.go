package handler

import (
	"encoding/json"
	"net/http"

	"github.com/PeterM45/go-postgres-api/internal/errors"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, errors.APIError{Code: 405, Message: "method not allowed"})
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, errors.ErrInvalidInput)
		return
	}

	user, err := h.store.VerifyUser(req.Email, req.Password)
	if err != nil {
		respondWithError(w, err)
		return
	}

	token, err := h.auth.GenerateToken(user.ID)
	if err != nil {
		respondWithError(w, errors.ErrInternalServer)
		return
	}

	respondWithJSON(w, http.StatusOK, loginResponse{Token: token})
}
