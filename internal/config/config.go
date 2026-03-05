package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config holds the application configuration values.
type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
	Logger    *logrus.Logger
}

// Load reads configuration values from environment variables,
// providing default values if the environment variables are not set.
// Returns a pointer to a Config struct populated with these values.
func Load() *Config {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return &Config{
		DBURL:     getEnv("DB_URL", "postgres://user:password@localhost:5432/expenses"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		Port:      getEnv("PORT", "8080"),
		Logger:    logger,
	}
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not present, it returns the provided defaultValue.м
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
