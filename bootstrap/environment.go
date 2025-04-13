package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PGDB   PGDatabase
	RDB    RDatabase
	Server Server
	SMS    SMS
}

type PGDatabase struct {
	DB_Host string
	DB_Name string
	DB_Port string
	DB_User string
	DB_Pass string
}

type RDatabase struct {
	RDB_Port     string
	RDB_Addr     string
	RDB_User     string
	RDB_Password string
}

type SMS struct {
	IPAddr        string
	APIKey        string
	APISandboxKey string
}

type Server struct {
	IP_Addr string
}

func NewEnvironment() *Env {
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}
	return &Env{
		PGDB: PGDatabase{
			DB_Host: os.Getenv("DB_HOST"),
			DB_Name: os.Getenv("DB_NAME"),
			DB_Port: os.Getenv("DB_PORT"),
			DB_User: os.Getenv("DB_USER"),
			DB_Pass: os.Getenv("DB_PASS"),
		},
		RDB: RDatabase{
			RDB_Port:     os.Getenv("REDIS_PORT"),
			RDB_Addr:     os.Getenv("REDIS_ADDR"),
			RDB_User:     os.Getenv("REDIS_USER"),
			RDB_Password: os.Getenv("REDIS_PASSWORD"),
		},
		Server: Server{
			IP_Addr: os.Getenv("IP_ADDR"),
		},
		SMS: SMS{
			IPAddr:        os.Getenv("SMS_URL"),
			APIKey:        os.Getenv("SMS_API_KEY"),
			APISandboxKey: os.Getenv("SMS_API_KEY_SANDBOX"),
		},
	}
}
