package config

import (
	"strconv"
	"strings"
	"time"
)

// ParseString parses a string value.
func ParseString(raw string) (string, error) {
	return raw, nil
}

// ParseInt parses an int value.
func ParseInt(raw string) (int, error) {
	return strconv.Atoi(raw)
}

// ParseFloat64 parses a float64 value.
func ParseFloat64(raw string) (float64, error) {
	return strconv.ParseFloat(raw, 64)
}

// ParseBool parses a bool value.
func ParseBool(raw string) (bool, error) {
	return strconv.ParseBool(raw)
}

// ParseDuration parses a time.Duration value.
func ParseDuration(raw string) (time.Duration, error) {
	return time.ParseDuration(raw)
}

// ParseCSV parses a comma-separated list. Empty entries are dropped.
func ParseCSV(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out, nil
}
