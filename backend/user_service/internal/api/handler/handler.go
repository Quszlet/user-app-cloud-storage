package handler

import (
	"github.com/gorilla/mux"
	"github.com/Quszlet/user_service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	users := r.PathPrefix("/users").Subrouter()
	users.HandleFunc("/create", h.CreateUser).Methods("POST")
	users.HandleFunc("/update/{id:[0-9]+}", h.UpdateUser).Methods("UPDATE")
	users.HandleFunc("/{id:[0-9]+}", h.GetUser).Methods("GET")
	users.HandleFunc("", h.GetAllUsers).Methods("GET")
	users.HandleFunc("/delete/{id:[0-9]+}", h.DeleteUser).Methods("DELETE")
	return r
}