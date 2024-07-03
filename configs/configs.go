package configs

import (
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int
}

type DBConfig struct {
	Name string
}
type AuthConfig struct {
	Password string
}

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	AuthConf AuthConfig
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 7540),
		},
		DB: DBConfig{
			Name: getEnv("TODO_DBFILE", "scheduler.db"),
		},
		AuthConf: AuthConfig{
			Password: getEnv("TODO_PASSWORD", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
