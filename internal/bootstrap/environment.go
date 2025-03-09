package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB Database
}

type Database struct {
	DB_Host      string
	DB_Name      string
	DB_Port      string
	DB_Root_Pass string
	DB_User      string
	DB_Pass      string

	RDB_Port     string
	RDB_Addr     string
	RDB_User     string
	RDB_Password string
}

func NewEnvironment() *Env {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	return &Env{
		DB: Database{
			DB_Host:      os.Getenv("DB_HOST"),
			DB_Name:      os.Getenv("DB_NAME"),
			DB_Port:      os.Getenv("DB_PORT"),
			DB_Root_Pass: os.Getenv("DB_ROOT_PASS"),
			DB_User:      os.Getenv("DB_USER"),
			DB_Pass:      os.Getenv("DB_PASS"),

			RDB_Port:     os.Getenv("REDIS_PORT"),
			RDB_Addr:     os.Getenv("REDIS_ADDR"),
			RDB_User:     os.Getenv("REDIS_USER"),
			RDB_Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}
