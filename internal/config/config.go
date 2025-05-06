package config

import "os"

type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
}

func Load() *Config {
	return &Config{
		DBURL:     getEnv("DB_URL", "postgres://user:password@localhost:5432/expenses"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
