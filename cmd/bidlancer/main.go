package main

import (
	"log"
	"net/http"
	"os"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/delivery/http/routes"
	"github.com/niflheimdevs/backend/wire"
)

func main() {

	var di = bootstrap.Get()

	err := os.MkdirAll(di.Const.StorageDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	app, err := wire.InitializeApplication(di)

	if err != nil {
		panic(err)
	}

	log.Printf("Application is running on port%s", di.Const.Port)
	server := &http.Server{
		Addr:    di.Const.Port,
		Handler: routes.Routes(app),
	}

	server.ListenAndServe()
}
