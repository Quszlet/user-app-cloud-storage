package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Quszlet/authentication_service/internal/api/handler"
	"github.com/Quszlet/authentication_service/internal/service"
)

func main() {


	service := service.NewService(&service.AuthService{TimeOut: int64(time.Hour)})
	handler := handler.NewHandler(service)

	routes := handler.InitRoutes()

	srv := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Authentication Service is running!")

	log.Fatal(srv.ListenAndServe())
}
