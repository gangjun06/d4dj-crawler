package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvKey string
type EnvToolPath string

var (
	KeyServerPort  EnvKey = "SERVER_PORT"
	KeyEnvToolPath EnvKey = "TOOL_PATH"
	KeyAssetPath   EnvKey = "ASSET_PATH"
	KeyCrawlerPool EnvKey = "CRAWLER_POOL"
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

func GetInt(key EnvKey) int {
	num, err := strconv.Atoi(os.Getenv(string(key)))
	if err != nil {
		log.Fatal(err)
	}
	return num
}
