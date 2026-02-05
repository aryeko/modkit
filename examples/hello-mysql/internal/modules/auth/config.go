package auth

import "time"

type Config struct {
	Secret   string
	Issuer   string
	TTL      time.Duration
	Username string
	Password string
}
