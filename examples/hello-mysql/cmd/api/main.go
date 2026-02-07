package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-modkit/modkit/examples/hello-mysql/docs"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/httpserver"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/lifecycle"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/app"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/auth"
	configmodule "github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/config"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/platform/logging"
	modkithttp "github.com/go-modkit/modkit/modkit/http"
	"github.com/go-modkit/modkit/modkit/kernel"
	"github.com/go-modkit/modkit/modkit/module"
)

// @title hello-mysql API
// @version 0.1
// @description Example modkit service with MySQL.
// @BasePath /api/v1
func main() {
	opts, err := loadAppOptions()
	if err != nil {
		log.Fatalf("config load failed: %v", err)
	}

	boot, handler, err := httpserver.BuildAppHandler(opts)
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	logger := logging.New()
	logStartup(logger, opts.HTTPAddr)

	server := &http.Server{
		Addr:    opts.HTTPAddr,
		Handler: handler,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigCh)
		close(sigCh)
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	hooks := buildShutdownHooks(boot)
	if err := runServer(modkithttp.ShutdownTimeout, server, sigCh, errCh, hooks); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func loadAppOptions() (app.Options, error) {
	cfgModule := configmodule.NewModule(configmodule.Options{})
	boot, err := kernel.Bootstrap(cfgModule)
	if err != nil {
		return app.Options{}, err
	}

	httpAddr, err := module.Get[string](boot, configmodule.TokenHTTPAddr)
	if err != nil {
		return app.Options{}, err
	}
	mySQLDSN, err := module.Get[string](boot, configmodule.TokenMySQLDSN)
	if err != nil {
		return app.Options{}, err
	}
	jwtSecret, err := module.Get[string](boot, configmodule.TokenJWTSecret)
	if err != nil {
		return app.Options{}, err
	}
	jwtIssuer, err := module.Get[string](boot, configmodule.TokenJWTIssuer)
	if err != nil {
		return app.Options{}, err
	}
	jwtTTL, err := module.Get[time.Duration](boot, configmodule.TokenJWTTTL)
	if err != nil {
		return app.Options{}, err
	}
	authUsername, err := module.Get[string](boot, configmodule.TokenAuthUsername)
	if err != nil {
		return app.Options{}, err
	}
	authPassword, err := module.Get[string](boot, configmodule.TokenAuthPassword)
	if err != nil {
		return app.Options{}, err
	}
	corsAllowedOrigins, err := module.Get[[]string](boot, configmodule.TokenCORSAllowedOrigins)
	if err != nil {
		return app.Options{}, err
	}
	corsAllowedMethods, err := module.Get[[]string](boot, configmodule.TokenCORSAllowedMethods)
	if err != nil {
		return app.Options{}, err
	}
	corsAllowedHeaders, err := module.Get[[]string](boot, configmodule.TokenCORSAllowedHeaders)
	if err != nil {
		return app.Options{}, err
	}
	rateLimitPerSecond, err := module.Get[float64](boot, configmodule.TokenRateLimitPerSecond)
	if err != nil {
		return app.Options{}, err
	}
	rateLimitBurst, err := module.Get[int](boot, configmodule.TokenRateLimitBurst)
	if err != nil {
		return app.Options{}, err
	}

	return app.Options{
		HTTPAddr: httpAddr,
		MySQLDSN: mySQLDSN,
		Auth: auth.Config{
			Secret:   jwtSecret,
			Issuer:   jwtIssuer,
			TTL:      jwtTTL,
			Username: authUsername,
			Password: authPassword,
		},
		CORSAllowedOrigins: corsAllowedOrigins,
		CORSAllowedMethods: corsAllowedMethods,
		CORSAllowedHeaders: corsAllowedHeaders,
		RateLimitPerSecond: rateLimitPerSecond,
		RateLimitBurst:     rateLimitBurst,
	}, nil
}

type shutdownServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type appLifecycle interface {
	CleanupHooks() []func(context.Context) error
	CloseContext(context.Context) error
}

func buildShutdownHooks(app appLifecycle) []lifecycle.CleanupHook {
	hooks := lifecycle.FromFuncs(app.CleanupHooks())
	return append([]lifecycle.CleanupHook{app.CloseContext}, hooks...)
}

func runServer(shutdownTimeout time.Duration, server shutdownServer, sigCh <-chan os.Signal, errCh <-chan error, hooks []lifecycle.CleanupHook) error {
	select {
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	case <-sigCh:
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		shutdownErr := lifecycle.ShutdownServer(ctx, server, hooks)

		err := <-errCh
		if err == http.ErrServerClosed {
			err = nil
		}
		if shutdownErr != nil {
			return shutdownErr
		}
		return err
	}
}
