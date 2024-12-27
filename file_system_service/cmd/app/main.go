package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Quszlet/file_system_service/internal/api/handler"
	"github.com/Quszlet/file_system_service/internal/repository"
	"github.com/Quszlet/file_system_service/internal/service"
)

func main() {
	db, err := repository.NewMySqlDB(repository.Config{
		User:   os.Getenv("USER_DB"),
		Passwd: os.Getenv("PASSWORD_DB"),
		Net:    os.Getenv("NET_DB"),
		Addr:   os.Getenv("ADDRESS_DB"),
		Port:   os.Getenv("PORT_DB"),
		DBName: os.Getenv("NAME_DB"),
	})

	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)

	repos.Directory.Migration()
	repos.MetadataFile.Migration()

	service := service.NewService(repos)
	handler := handler.NewHandler(service)

	routes := handler.InitRoutes()

	srv := &http.Server{
		Addr:           os.Getenv("PORT"),
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("File system service is running!")

	log.Fatal(srv.ListenAndServe())
}
