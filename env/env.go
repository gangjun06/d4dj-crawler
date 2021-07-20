package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string

var (
	KeyServerPort EnvKey = "SERVER_PORT"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Get(key EnvKey) string {
	return os.Getenv(string(key))
}
