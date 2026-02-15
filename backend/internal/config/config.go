package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Host        string
	DBPath      string
	CORSOrigins string
	Env         string
}

// Apparantly it is conventional for loaders to return pointers
// but I don't see a point in that, configs should be immutable
func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		Port:        getEnv("PORT", "8080"),
		Host:        getEnv("HOST", "localhost"),
		DBPath:      getEnv("TICKER_DB_PATH", ""),
		CORSOrigins: getEnv("CORS_ORIGINS", "http://localhost:5173"),
		Env:         getEnv("ENV", "development"),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c Config) IsDevelopment() bool {
	return c.Env == "development"
}
