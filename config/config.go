package config

import "os"

type EnvFileConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Dbname       string
	Sslmode      string
	HTTPHost     string
	HTTPPort     string
	AppLogLevel  string
	DBLogLevel   string
	EmailService string
	EmailLogin   string
	EmailPass    string
	EmailReply   string
}

type Config struct {
	Config EnvFileConfig
}

func New() *Config {
	return &Config{Config: EnvFileConfig{
		Host:         getEnv("DB_HOST"),
		Port:         getEnv("DB_PORT"),
		User:         getEnv("DB_USER"),
		Password:     getEnv("DB_PASSWORD"),
		Dbname:       getEnv("DB_NAME"),
		Sslmode:      getEnv("DB_SSLMODE"),
		HTTPHost:     getEnv("HTTP_HOST"),
		HTTPPort:     getEnv("HTTP_PORT"),
		AppLogLevel:  getEnv("APP_LOG_LEVEL"),
		DBLogLevel:   getEnv("DB_LOG_LEVEL"),
		EmailService: getEnv("EMAIL_SERVICE"),
		EmailLogin:   getEnv("EMAIL_LOGIN"),
		EmailPass:    getEnv("EMAIL_PASS"),
		EmailReply:   getEnv("EMAIL_REPLY"),
	}}
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return ""
}
