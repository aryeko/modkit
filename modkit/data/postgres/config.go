package postgres

import (
	"time"

	"github.com/go-modkit/modkit/modkit/config"
	"github.com/go-modkit/modkit/modkit/module"
)

// DefaultConfigModule provides Postgres configuration from environment variables.
//
// Required:
// - POSTGRES_DSN
//
// Optional:
// - POSTGRES_MAX_OPEN_CONNS
// - POSTGRES_MAX_IDLE_CONNS
// - POSTGRES_CONN_MAX_LIFETIME
// - POSTGRES_CONNECT_TIMEOUT (default 0; disables provider ping)
func DefaultConfigModule() module.Module {
	return configModule("")
}

func configModule(name string) module.Module {
	return config.NewModule(
		config.WithModuleName(moduleName(name)+".config"),
		config.WithTyped(TokenDSN, config.ValueSpec[string]{
			Key:         "POSTGRES_DSN",
			Required:    true,
			Sensitive:   true,
			Description: "Postgres DSN.",
			Parse:       config.ParseString,
		}, true),
		config.WithTyped(TokenMaxOpenConns, config.ValueSpec[int]{
			Key:         "POSTGRES_MAX_OPEN_CONNS",
			Description: "Maximum open connections for the DB pool.",
			Parse:       config.ParseInt,
		}, true),
		config.WithTyped(TokenMaxIdleConns, config.ValueSpec[int]{
			Key:         "POSTGRES_MAX_IDLE_CONNS",
			Description: "Maximum idle connections for the DB pool.",
			Parse:       config.ParseInt,
		}, true),
		config.WithTyped(tokenMaxIdleConnsSet, config.ValueSpec[bool]{
			Key:         "POSTGRES_MAX_IDLE_CONNS",
			Description: "Whether POSTGRES_MAX_IDLE_CONNS is explicitly set.",
			Parse: func(string) (bool, error) {
				return true, nil
			},
		}, true),
		config.WithTyped(TokenConnMaxLifetime, config.ValueSpec[time.Duration]{
			Key:         "POSTGRES_CONN_MAX_LIFETIME",
			Description: "Maximum amount of time a connection may be reused.",
			Parse:       config.ParseDuration,
		}, true),
		config.WithTyped(TokenConnectTimeout, config.ValueSpec[time.Duration]{
			Key:         "POSTGRES_CONNECT_TIMEOUT",
			Description: "Optional ping timeout on provider build. 0 disables ping.",
			Parse:       config.ParseDuration,
		}, true),
	)
}
