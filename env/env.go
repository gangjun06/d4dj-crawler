package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvKey string
type EnvToolPath string

var (
	KeyServerPort  EnvKey = "SERVER_PORT"
	KeyEnvToolPath EnvKey = "TOOL_PATH"
	KeyAssetPath   EnvKey = "ASSET_PATH"
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
