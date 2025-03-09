package main

import (
	"log"
	"net/http"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/driver"
	"github.com/niflheimdevs/backend/internal/routes"
	"github.com/niflheimdevs/backend/internal/wire"
)

const port = ":8080"

func main() {
	var di = bootstrap.Get()

	log.Printf("Connecting to database...")

	db := driver.ConnectSQL(di)

	defer db.Close()

	app, err := wire.InitializeApplication(di, db)
	if err != nil {
		panic(err)
	}

	log.Printf("Application is running on port%s", port)

	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(app),
	}

	server.ListenAndServe()
}
