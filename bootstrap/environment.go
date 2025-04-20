package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PGDB    PGDatabase
	RDB     RDatabase
	Server  Server
	SMS     SMS
	Storage S3
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

type S3 struct {
	Buckets   BucketName
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string
}

type BucketName struct {
	ProfilePic string
	Resume     string
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
		Storage: S3{
			Buckets: BucketName{
				ProfilePic: os.Getenv("PROFILE_PIC_BUCKET_NAME"),
				Resume:     os.Getenv("RESUME_BUCKET_NAME"),
			},
			Region:    os.Getenv("BUCKET_REGION"),
			AccessKey: os.Getenv("BUCKET_ACCESS_key"),
			SecretKey: os.Getenv("BUCKET_SECRET_key"),
			Endpoint:  os.Getenv("BUCKET_ENDPOINT"),
		},
	}
}
