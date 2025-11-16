package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
}

var Envs Config

func init() {
	_ = godotenv.Load()

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("Warning: .env file not found (using system env only)")
	}

	Envs = Config{
		PublicHost: os.Getenv("PUBLIC_HOST"),
		Port:       os.Getenv("PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}
}
