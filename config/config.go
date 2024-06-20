package config

import "os"

type EnvFileConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Sslmode  string
	HTTPHost string
	HTTPPort string
}

type Config struct {
	Config EnvFileConfig
}

func New() *Config {
	return &Config{Config: EnvFileConfig{
		Host:     getEnv("DB_HOST"),
		Port:     getEnv("DB_PORT"),
		User:     getEnv("DB_USER"),
		Password: getEnv("DB_PASSWORD"),
		Dbname:   getEnv("DB_NAME"),
		Sslmode:  getEnv("DB_SSLMODE"),
		HTTPHost: getEnv("HTTP_HOST"),
		HTTPPort: getEnv("HTTP_PORT"),
	}}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
