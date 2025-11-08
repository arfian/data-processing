// ============================================
// config/config.go
// ============================================
package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	WorkerCount int
	BatchSize   int
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/dbname?sslmode=disable"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		WorkerCount: getEnvInt("WORKER_COUNT", 5),
		BatchSize:   getEnvInt("BATCH_SIZE", 100),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
