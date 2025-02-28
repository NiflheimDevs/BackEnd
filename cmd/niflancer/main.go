package main

import (
	"log"
	"net/http"

	"github.com/niflheimdevs/backend/internal/routes"
)

const port = ":8080"

func main() {
	log.Printf("Application is running on port%s", port)

	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(),
	}

	server.ListenAndServe()
}
