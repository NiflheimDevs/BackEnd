package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/driver"
	"github.com/niflheimdevs/backend/internal/routes"
)

const port = ":8080"

func main() {
	var di = bootstrap.Get()

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		di.Env.DB.DB_Host, di.Env.DB.DB_Port, di.Env.DB.DB_User,
		di.Env.DB.DB_Pass, di.Env.DB.DB_Name)

	log.Printf("Connecting to database...")

	db := driver.ConnectSQL(dsn, di)

	defer db.Close()

	log.Printf("Application is running on port%s", port)

	server := &http.Server{
		Addr:    port,
		Handler: routes.Routes(),
	}

	server.ListenAndServe()
}
