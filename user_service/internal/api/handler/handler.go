package handler

import (
	"net/http"

	"github.com/Quszlet/user_service/internal/service"
	json "github.com/Quszlet/user_service/pkg"
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

	users := r.PathPrefix("/users").Subrouter()
	users.HandleFunc("/create", h.CreateUser).Methods("POST")
	users.HandleFunc("/update/{id:[0-9]+}", h.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id:[0-9]+}", h.GetUserById).Methods("GET")
	users.HandleFunc("/{login:[a-zA-Z0-9_]+}", h.GetUserByLogin).Methods("GET")
	users.HandleFunc("", h.GetAllUsers).Methods("GET")
	users.HandleFunc("/delete/{id:[0-9]+}", h.DeleteUser).Methods("DELETE")
	return r
}