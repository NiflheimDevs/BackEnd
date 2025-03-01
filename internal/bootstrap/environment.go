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
}

func NewEnvironment() *Env {
	godotenv.Load(".env")
	return &Env{
		DB: Database{
			DB_Host:      os.Getenv("DB_HOST"),
			DB_Name:      os.Getenv("DB_NAME"),
			DB_Port:      os.Getenv("DB_PORT"),
			DB_Root_Pass: os.Getenv("DB_ROOT_PASS"),
			DB_User:      os.Getenv("DB_USER"),
			DB_Pass:      os.Getenv("DB_PASS"),
		},
	}
}
