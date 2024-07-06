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

func New() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 7545),
		},
		DB: DBConfig{
			Name: getEnv("TODO_DBFILE", "scheduler.db"),
		},
		AuthConf: AuthConfig{
			Password: getEnv("TODO_PASSWORD", ""),
		},
	}
}

// Функция для чтения переменных среды или возвращает дефолтное значение
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Функция для чтения целочисленных переменных среды или возвращает дефолтное значение
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
