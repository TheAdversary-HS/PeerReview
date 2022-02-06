package config

import (
	"github.com/joho/godotenv"
	"os"
)

var (
	ServerPort string

	Address string
	SubPath string

	DatabaseFile string
	FrontendDir  string
)

func init() {
	godotenv.Load()

	ServerPort = os.Getenv("SERVER_PORT")

	Address = os.Getenv("ADDRESS")
	SubPath = os.Getenv("SUBPATH")

	DatabaseFile = os.Getenv("DATABASE_FILE")
	FrontendDir = os.Getenv("FRONTEND_DIR")
}
