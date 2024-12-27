package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"github.com/Quszlet/user_service/internal/models"
	"github.com/Quszlet/user_service/pkg"
	"github.com/gorilla/mux"
)


func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.Parse(r, &user)
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = user.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	id, err := h.services.User.Create(user)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed created user")
		return
	}


	json.Response(w, http.StatusCreated, map[string]interface{}{"id": id})
}


func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user := models.User{Id: id}
	err := json.Parse(r, &user)

	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = user.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	err = h.services.User.Update(user)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed update user")
		return
	}

	message := fmt.Sprintf("User with id %d updated", id)
	
	json.Response(w, http.StatusOK, message)
}

func (h *Handler) GetUserByLogin(w http.ResponseWriter, r *http.Request) {
	login, _ := mux.Vars(r)["login"]
	user, err := h.services.User.GetByLogin(login)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get user")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.services.User.GetById(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get user")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.User.GetAll()
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get users")
		return
	}

	json.Response(w, http.StatusOK, users)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := h.services.User.Delete(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed delete user")
		return
	}

	message := fmt.Sprintf("User delete with id %d", id)
	
	json.Response(w, http.StatusOK, message)
}