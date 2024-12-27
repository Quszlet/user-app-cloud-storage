package handler

import (
	"net/http"

	"github.com/Quszlet/authentication_service/internal/models"
	json "github.com/Quszlet/authentication_service/pkg"
)

func (h *Handler) SingIn(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.Parse(r, &user)

	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	token, err := h.services.Authentication.SignIn(user.Login, user.Password)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed created jwt token")
		return
	}

	response := map[string]string{"token": token}

	json.Response(w, http.StatusCreated, response)
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	json.Response(w, http.StatusOK, "Good")
}
