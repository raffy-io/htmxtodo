package conf

import (
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DBURL     string // Idiomatic PascalCase
	AuthToken string
}

func Load() (*Config ,error) {
	_ = godotenv.Load() 

	cfg := &Config{
		DBURL:     os.Getenv("DB_URL"),
		AuthToken: os.Getenv("AUTH_TOKEN"),
	}

	return cfg, nil
}



