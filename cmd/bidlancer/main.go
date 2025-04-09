package main

import (
	"log"
	"net/http"
	"os"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/delivery/http/routes"
	"github.com/niflheimdevs/backend/internal/infrastructure/db/driver"
	"github.com/niflheimdevs/backend/internal/wire"
)

const port = ":8080"

func main() {

	var di = bootstrap.Get()

	err := os.MkdirAll(di.Const.StorageDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	pdb := driver.ConnectSQL(di)
	rdb := driver.ConncetRedis(di)

	defer pdb.Close()
	defer rdb.Close()

	app, err := wire.InitializeApplication(di)

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
