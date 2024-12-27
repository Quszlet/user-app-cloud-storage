package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Quszlet/file_system_service/internal/models"
	"github.com/Quszlet/file_system_service/pkg"
)


func (h *Handler) CreateFile(w http.ResponseWriter, r *http.Request) {
	metadataFile := models.MetadataFile{}
	err := json.Parse(r, &metadataFile)
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = metadataFile.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	id, err := h.services.MetadataFile.Create(metadataFile)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed created file")
		return
	}

	message := fmt.Sprintf("File created with id %d", id)

	json.Response(w, http.StatusCreated, message)
}


func (h *Handler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	metadataFile := models.MetadataFile{Id: id}
	err := json.Parse(r, &metadataFile)

	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Failed parse JSON")
		return
	}

	err = metadataFile.Validate()
	if err != nil {
		json.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid JSON")
		return
	}

	err = h.services.MetadataFile.Update(metadataFile)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed update file")
		return
	}

	message := fmt.Sprintf("File with id %d updated", id)
	
	json.Response(w, http.StatusOK, message)
}

func (h *Handler) GetFileByName(w http.ResponseWriter, r *http.Request) {
	name, _ := mux.Vars(r)["name"]
	user, err := h.services.MetadataFile.GetByName(name)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get file")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetFileById(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	user, err := h.services.MetadataFile.GetById(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get file")
		return
	}

	json.Response(w, http.StatusOK, user)
}

func (h *Handler) GetAllFiles(w http.ResponseWriter, r *http.Request) {
	users, err := h.services.MetadataFile.GetAll()
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed get files")
		return
	}

	json.Response(w, http.StatusOK, users)
}

func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := h.services.MetadataFile.Delete(id)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed delete file")
		return
	}

	message := fmt.Sprintf("User delete with id %d", id)
	
	json.Response(w, http.StatusOK, message)
}

func (h *Handler) MoveFile(w http.ResponseWriter, r *http.Request) {
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

	err = h.services.MetadataFile.Move(info.MovedObjId, info.DirectoryId)
	if err != nil {
		json.ErrorResponse(w, http.StatusInternalServerError, err.Error(), "Failed move file")
		return
	}

	message := fmt.Sprintf("File move in directory with id %d",info.DirectoryId)
	
	json.Response(w, http.StatusOK, message)
}