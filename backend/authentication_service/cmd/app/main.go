package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	// handler := handler.NewHandler(service)

	// routes := handler.InitRoutes()

	srv := &http.Server{
		Addr:           ":8081",
		// Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("User Service is running!")

	log.Fatal(srv.ListenAndServe())
}
