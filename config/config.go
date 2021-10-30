package config

import (
	goquickenv "github.com/panospet/go-quick-env"
	"os"
)

type Config struct {
	Port  int
	DbDsn string
}

func NewConfig() (Config, error) {
	if err := goquickenv.LoadFile(".env"); err != nil {
		return Config{}, err
	}
	dbDsn := goquickenv.GetEnvAsString(
		"DB_PATH",
		"postgres://127.0.0.1/weeklyprogrammedb?sslmode=disable&user=admin&password=password",
	)
	port := goquickenv.GetEnvAsInt("PORT", 6000)
	return Config{
		Port:  port,
		DbDsn: dbDsn,
	}, nil
}

func envOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
