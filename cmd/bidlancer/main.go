package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/driver"
	"github.com/niflheimdevs/backend/internal/routes"
	"github.com/niflheimdevs/backend/internal/wire"
	"github.com/redis/go-redis/v9"
)

const port = ":8080"

func main() {
	var di = bootstrap.Get()

	pdb, rdb := createConnections(di)

	defer pdb.Close()
	defer rdb.Close()

	app, err := wire.InitializeApplication(di, pdb, rdb)

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

func createConnections(di *bootstrap.Di) (*pgxpool.Pool, *redis.Client) {

	pdb := driver.ConnectSQL(di)
	rdb := driver.ConncetRedis(di)

	ctx := context.Background()

	err := pdb.Ping(ctx)
	if err != nil {
		panic(err)
	}
	err = rdb.Ping(ctx).Err()

	if err != nil {
		panic(err)
	}

	return pdb, rdb
}
