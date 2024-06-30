package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var Envs envs

type envs struct {
	DB_HOST                string
	DB_PORT                string
	DB_USER                string
	DB_PASSWORD            string
	DB_NAME                string
	JWT_EXPIRATION_SECONDS int
	JWT_SECRET             string
}

func init() {
	Envs = envs{
		// DB_HOST:                parseEnv("DB_HOST", "localhost"),
		DB_HOST:                parseEnv("DB_HOST", "db"),
		DB_PORT:                parseEnv("DB_PORT", "5432"),
		DB_USER:                parseEnv("DB_USER", "postgres"),
		DB_PASSWORD:            parseEnv("DB_PASSWORD", "postgres"),
		DB_NAME:                parseEnv("DB_NAME", "postgres"),
		JWT_EXPIRATION_SECONDS: parseEnvAsInt("JWT_EXPIRATION_SECONDS", 300), // 5 minutes
		JWT_SECRET:             parseEnv("JWT_SECRET", ""),
	}

	if Envs.JWT_SECRET == "" {
		// panic("No JWT secret provided, base64 encoded JWT secret must be provided.")
	}
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (db DBConfig) PostgresURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?connect_timeout=5", db.User, db.Password, db.Host, db.Port, db.Name)
}

// Checks for and parses environment variables to populate
// and return a DBConfig struct.
func GetDBConfig() DBConfig {
	return DBConfig{
		Host:     Envs.DB_HOST,
		Port:     Envs.DB_PORT,
		User:     Envs.DB_USER,
		Password: Envs.DB_PASSWORD,
		Name:     Envs.DB_NAME,
	}
}

// func GetDBConfig() DBConfig {
// 	return DBConfig{
// 		Host:     parseEnv("DB_HOST", "localhost"),
// 		Port:     parseEnv("DB_PORT", "5432"),
// 		User:     parseEnv("DB_USER", "postgres"),
// 		Password: parseEnv("DB_PASSWORD", "postgres"),
// 		Name:     parseEnv("DB_NAME", "postgres"),
// 	}
// }

func parseEnv(env, defaultValue string) (value string) {
	if value = os.Getenv(env); len(value) == 0 {
		return defaultValue
	}
	return value
}

func parseEnvAsInt(env string, defaultValue int) int {
	if str := os.Getenv(env); len(str) != 0 {
		value, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("failed converting env string to int: %s", err)
			return defaultValue
		}

		return value
	}

	return defaultValue
}
