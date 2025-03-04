package driver

import (
	"database/sql"
	"log"

	"github.com/niflheimdevs/backend/internal/bootstrap"
)

func ConnectSQL(dsn string, di *bootstrap.Di) *sql.DB {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db.SetMaxOpenConns(di.Const.Database.MaxOpenDbConn)
	db.SetConnMaxIdleTime(di.Const.Database.MaxIdleDbConn)
	db.SetConnMaxLifetime(di.Const.Database.MaxDbLifeTime)

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return db
}
