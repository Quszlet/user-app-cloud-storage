package handler

import (
	"net/http"

	"github.com/Quszlet/file_system_service/internal/service"
	json "github.com/Quszlet/file_system_service/pkg"
	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	json.Response(w, http.StatusOK, "Good")
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/health", h.CheckHealth).Methods("GET")

	directories := r.PathPrefix("/directories").Subrouter()
	{
		directories.HandleFunc("/create", h.CreateDirectory).Methods("POST")
		directories.HandleFunc("/update/{id:[0-9]+}", h.UpdateDirectory).Methods("PUT")
		directories.HandleFunc("/{id:[0-9]+}", h.GetDirectoryById).Methods("GET")
		directories.HandleFunc("/home/user/{id:[0-9]+}", h.GetHomeDirectoryByUserId).Methods("GET")
		directories.HandleFunc("/{name:[a-zA-Z0-9_]+}", h.GetDirectoryByName).Methods("GET")
		directories.HandleFunc("", h.GetAllDirectories).Methods("GET")
		directories.HandleFunc("/delete/{id:[0-9]+}", h.DeleteDirectory).Methods("DELETE")
		directories.HandleFunc("/move", h.MoveDirectory).Methods("PUT")
	}

	metadataFile := r.PathPrefix("/files").Subrouter()
	{
		metadataFile.HandleFunc("/create", h.CreateFile).Methods("POST")
		metadataFile.HandleFunc("/update/{id:[0-9]+}", h.UpdateFile).Methods("PUT")
		metadataFile.HandleFunc("/{id:[0-9]+}", h.GetFileById).Methods("GET")
		metadataFile.HandleFunc("/{name:[a-zA-Z0-9_]+}", h.GetDirectoryByName).Methods("GET")
		metadataFile.HandleFunc("", h.GetAllFiles).Methods("GET")
		metadataFile.HandleFunc("/delete/{id:[0-9]+}", h.DeleteFile).Methods("DELETE")
		metadataFile.HandleFunc("/move", h.MoveFile).Methods("PUT")
	}

	return r
}
