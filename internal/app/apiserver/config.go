package apiserver

import (
	"os"

	"github.com/http-rest-api/internal/app/store"
)

// Config ...
type Config struct {
	BindAddr string
	LogLevel string
	Store    *store.Config
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: getEnv("bind_addr", ":8080"),
		LogLevel: getEnv("log_level", "debug"),
		Store:    store.NewConfig(),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
