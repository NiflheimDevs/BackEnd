package main

import (
	"log"
	"net/http"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/delivery/routes"
	"github.com/niflheimdevs/backend/internal/infrastructure/websocket"
	"github.com/niflheimdevs/backend/wire"
)

func main() {

	var di = bootstrap.Get()

	hub := websocket.NewHub()
	go hub.Run()

	app, err := wire.InitializeApplication(di, hub)

	if err != nil {
		panic(err)
	}

	app.Seeder.SeedDaddy()

	log.Printf("Application is running on port%s", di.Const.Port)
	server := &http.Server{
		Addr:    di.Const.Port,
		Handler: routes.Routes(app),
	}

	err = server.ListenAndServeTLS(di.Const.SSLKeysPath+"/certificate.pem", di.Const.SSLKeysPath+"/privatekey.key")
	if err != nil {
		log.Fatal(err)
	}
}
