package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DSN          string `json:"DSN"`
	PortListener string `json:"PortListener"`
	SMTPHost     string `json:"SMTPHost"`
}

func NewConfig() Config {
	return Config{
		DSN:          "default_dsn",
		PortListener: "default_port",
	}
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		DSN:          os.Getenv("DSN"),
		PortListener: os.Getenv("PORT_LISTENER"),
	}
}
