package driver

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/redis/go-redis/v9"
)

func ConnectSQL(di *bootstrap.Di) *pgxpool.Pool {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		di.Env.DB.DB_User, di.Env.DB.DB_Pass, di.Env.DB.DB_Host, di.Env.DB.DB_Port, di.Env.DB.DB_Name)

	poolconfig, err := pgxpool.ParseConfig(dsn)
	poolconfig.MaxConnIdleTime = di.Const.Database.MaxIdleDbConn
	poolconfig.MaxConnLifetime = di.Const.Database.MaxDbLifeTime
	poolconfig.MaxConns = di.Const.Database.MaxOpenDbConn
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.NewWithConfig(ctx, poolconfig)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return pool
}

func ConncetRedis(di *bootstrap.Di) *redis.Client {
	ctx := context.Background()
	dsn := fmt.Sprintf("%s:%s", di.Env.DB.RDB_Addr, di.Env.DB.RDB_Port)
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: di.Env.DB.RDB_Password,
		Username: di.Env.DB.RDB_User,
		DB:       0,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client
}
