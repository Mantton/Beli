package env

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Environment Loaded")
}

func Port() string {
	return os.Getenv("PORT")
}

func Redis() string {
	return os.Getenv(("REDIS_URL"))
}
