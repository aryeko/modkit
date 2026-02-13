package postgres

import "github.com/go-modkit/modkit/modkit/module"

const (
	// TokenDSN resolves the Postgres DSN.
	TokenDSN module.Token = "postgres.dsn" //nolint:gosec // token name, not credential
	// TokenMaxOpenConns resolves the max open connections pool setting.
	TokenMaxOpenConns module.Token = "postgres.max_open_conns" //nolint:gosec // token name, not credential
	// TokenMaxIdleConns resolves the max idle connections pool setting.
	TokenMaxIdleConns    module.Token = "postgres.max_idle_conns"     //nolint:gosec // token name, not credential
	tokenMaxIdleConnsSet module.Token = "postgres.max_idle_conns_set" //nolint:gosec // token name, not credential
	// TokenConnMaxLifetime resolves the connection max lifetime pool setting.
	TokenConnMaxLifetime module.Token = "postgres.conn_max_lifetime" //nolint:gosec // token name, not credential
	// TokenConnectTimeout resolves the optional provider ping timeout.
	TokenConnectTimeout module.Token = "postgres.connect_timeout" //nolint:gosec // token name, not credential
)
