package config

import (
	"fmt"
	"os"
	"strconv"
)

type UserConfig struct {
	RequireUsername bool   // If username is needed
	RequireEmail    bool   // If email is needed
	IDField         string // Type of ID (uuid, serial, etc.)
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int64 // in hours
}

type Config struct {
	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort string
	Port   string
	User   UserConfig
	JWT    JWTConfig
}

func Load() *Config {
	return &Config{
		DBUser: getEnvOrDefault("DB_USER", "admin"),
		DBPass: getEnvOrDefault("DB_PASS", "password"),
		DBName: getEnvOrDefault("DB_NAME", "myapp"),
		DBHost: getEnvOrDefault("DB_HOST", "localhost"),
		DBPort: getEnvOrDefault("DB_PORT", "5432"),
		Port:   getEnvOrDefault("PORT", "8080"),
		User: UserConfig{
			RequireUsername: getBoolEnv("REQUIRE_USERNAME", true),
			RequireEmail:    getBoolEnv("REQUIRE_EMAIL", true),
			IDField:         getEnvOrDefault("ID_FIELD", "serial"),
		},
		JWT: JWTConfig{
			SecretKey: getEnvOrDefault("JWT_SECRET", "66757119ed5f91f079b8eccbaf31b9b7ff6823e6f0b11ff30e30d317b94fe575"),
			ExpiresIn: getInt64Env("JWT_EXPIRES_IN", 24),
		},
	}
}

func getBoolEnv(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val == "true"
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}
