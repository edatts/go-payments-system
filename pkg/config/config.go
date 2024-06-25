package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (db DBConfig) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}

// Checks for and parses environment variables to populate
// and return a DBConfig struct.
func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     parseEnv("DB_HOST", "http://localhost"),
		Port:     parseEnv("DB_PORT", "5432"),
		User:     parseEnv("DB_USER", "postgres"),
		Password: parseEnv("DB_PASSWORD", "postgres"),
		Name:     parseEnv("DB_NAME", "postgres"),
	}
}

func parseEnv(env, defaultValue string) (value string) {
	if value = os.Getenv(env); len(value) == 0 {
		return defaultValue
	}
	return value
}
