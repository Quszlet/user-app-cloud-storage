package handler

import (
	"github.com/Quszlet/authentication_service/internal/service"
	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/authentication", h.SingIn).Methods("POST")
	r.HandleFunc("/health", h.CheckHealth).Methods("GET")

	return r
}
