package store

import "os"

// Config ...
type Config struct {
	DatabaseURL  string
	DatabaseName string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		DatabaseURL:  getEnv("database_url", ""),
		DatabaseName: getEnv("database_name", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
