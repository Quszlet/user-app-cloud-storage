package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Quszlet/user_service/internal/api/handler"
	"github.com/Quszlet/user_service/internal/repository"
	"github.com/Quszlet/user_service/internal/service"
)

func main() {
	db, err := repository.NewMySqlDB(repository.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "127.0.0.1",
		Port:   "3306",
		DBName: "users",
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := handler.NewHandler(service)

	routes := handler.InitRoutes()

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("User Service is running!")

	log.Fatal(srv.ListenAndServe())
}
