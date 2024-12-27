package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Quszlet/file_system_service/internal/models"
	json "github.com/Quszlet/file_system_service/pkg"
)

func (h *Handler) CreateDirectory(w http.ResponseWriter, r *http.Request) {
	directory := models.Directory{}
	err := json.Parse(r, &directory)
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = directory.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	id, err := h.services.Directory.Create(directory)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed created directory")
		return
	}

	message := fmt.Sprintf("Directory created with id %d", id)

	json.Response(w, http.StatusCreated, message)
}

func (h *Handler) UpdateDirectory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	directory := models.Directory{Id: id}
	err := json.Parse(r, &directory)

	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = directory.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	err = h.services.Directory.Update(directory)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed update directory")
		return
	}

	message := fmt.Sprintf("Directory with id %d updated", id)

	json.Response(w, http.StatusOK, message)
}

func (h *Handler) GetDirectoryByName(w http.ResponseWriter, r *http.Request) {
	name, _ := mux.Vars(r)["name"]
	user, err := h.services.Directory.GetByName(name)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get directory")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetDirectoryById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.services.Directory.GetById(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get directory")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetHomeDirectoryByUserId(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.services.Directory.GetHomeDirByUserId(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get directory")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetAllDirectories(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.Directory.GetAll()
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get users")
		return
	}

	json.Response(w, http.StatusOK, users)
}

func (h *Handler) DeleteDirectory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := h.services.Directory.Delete(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed delete user")
		return
	}

	message := fmt.Sprintf("User delete with id %d", id)

	json.Response(w, http.StatusOK, message)
}

func (h *Handler) MoveDirectory(w http.ResponseWriter, r *http.Request) {
	info := models.MoveInfo{}
	err := json.Parse(r, &info)
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = info.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	err = h.services.Directory.Move(info.MovedObjId, info.DirectoryId)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed move file")
		return
	}

	message := fmt.Sprintf("Directory move in directory with id %d",info.DirectoryId)
	
	json.Response(w, http.StatusOK, message)
}