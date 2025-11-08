// ============================================
// config/config.go
// ============================================
package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	WorkerCount int
	BatchSize   int
}

func LoadConfig() *Config {
	viper.SetConfigType("env")
	viper.SetConfigName(".env") // name of Config file (without extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/secrets")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to load config file", err)
	}

	return &Config{
		DatabaseURL: getRequiredString("DATABASE_URL"),
		ServerPort:  getRequiredString("SERVER_PORT"),
		WorkerCount: getRequiredInt("WORKER_COUNT"),
		BatchSize:   getRequiredInt("BATCH_SIZE"),
	}
}

func getRequiredString(key string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}

	log.Fatalln(fmt.Errorf("KEY %s IS MISSING", key))
	return ""
}

func getRequiredInt(key string) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}

	panic(fmt.Errorf("KEY %s IS MISSING", key))
}
