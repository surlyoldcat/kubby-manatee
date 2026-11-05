package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	ServerPort string
}

// Load reads configuration from a local .env file if present, then from the
// process environment. godotenv never overrides variables that are already set.
func Load() (*Config, error) {
	return fromEnv()
}

// LoadFile is like Load but reads the given .env file first (used by tests).
func LoadFile(path string) (*Config, error) {
	_ = godotenv.Load(path)
	return fromEnv()
}

func fromEnv() (*Config, error) {
	_ = godotenv.Load()
	return fromKubernetes()
}

// fromKubernetes is the placeholder for fetching DB config from the Kubernetes
// environment (Secret/ConfigMap mounted as env vars or files). For now it reads
// the same DB_* environment variables that the .env file populates.
// In the cluster these come from k8s/api.yaml (ConfigMap + the app-db-creds Secret).
func fromKubernetes() (*Config, error) {
	c := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBName:     os.Getenv("DB_DATABASE"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  getenv("DB_SSLMODE", "require"),
		ServerPort: getenv("PORT", "8083"),
	}
	for name, v := range map[string]string{
		"DB_HOST": c.DBHost, "DB_DATABASE": c.DBName, "DB_USER": c.DBUser, "DB_PASSWORD": c.DBPassword,
	} {
		if v == "" {
			return nil, fmt.Errorf("missing required setting %s", name)
		}
	}
	return c, nil
}

// DSN returns a URL-style Postgres DSN with a 10s connect timeout. URL form
// keeps special characters in the password safe.
func (c *Config) DSN() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.DBUser, c.DBPassword),
		Host:     c.DBHost + ":" + c.DBPort,
		Path:     "/" + c.DBName,
		RawQuery: url.Values{"sslmode": {c.DBSSLMode}, "connect_timeout": {"10"}}.Encode(),
	}
	return u.String()
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
