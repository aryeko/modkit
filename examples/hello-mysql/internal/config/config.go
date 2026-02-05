package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr           string
	MySQLDSN           string
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	RateLimitPerSecond float64
	RateLimitBurst     int
}

func Load() Config {
	return Config{
		HTTPAddr:           envOrDefault("HTTP_ADDR", ":8080"),
		MySQLDSN:           envOrDefault("MYSQL_DSN", "root:password@tcp(localhost:3306)/app?parseTime=true&multiStatements=true"),
		CORSAllowedOrigins: splitEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		CORSAllowedMethods: splitEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE"),
		RateLimitPerSecond: envFloat("RATE_LIMIT_PER_SECOND", 5),
		RateLimitBurst:     envInt("RATE_LIMIT_BURST", 10),
	}
}

func envOrDefault(key, def string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	return val
}

func splitEnv(key, def string) []string {
	raw := envOrDefault(key, def)
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func envFloat(key string, def float64) float64 {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return def
	}
	return parsed
}

func envInt(key string, def int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	parsed, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return parsed
}
