package config_test

import (
	"errors"
	"testing"

	"github.com/go-modkit/modkit/modkit/config"
	"github.com/go-modkit/modkit/modkit/kernel"
	"github.com/go-modkit/modkit/modkit/module"
)

func TestIntegration_NonExportedTokenNotVisible(t *testing.T) {
	const token module.Token = "config.private"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{"PRIVATE": "secret"}),
		config.WithTyped(token, config.ValueSpec[string]{
			Key:      "PRIVATE",
			Required: true,
			Parse:    config.ParseString,
		}, false),
	)

	appModule := mod("app", []module.Module{cfgModule}, []module.ControllerDef{{
		Name: "NeedsPrivateConfig",
		Build: func(r module.Resolver) (any, error) {
			_, err := module.Get[string](r, token)
			if err != nil {
				return nil, err
			}
			return "ok", nil
		},
	}})

	_, err := kernel.Bootstrap(appModule)
	if err == nil {
		t.Fatalf("expected visibility error")
	}

	var visErr *kernel.TokenNotVisibleError
	if !errors.As(err, &visErr) {
		t.Fatalf("expected TokenNotVisibleError, got %T", err)
	}
	if visErr.Token != token {
		t.Fatalf("unexpected token: %q", visErr.Token)
	}
}

func TestIntegration_FirstResolutionFailsForMissingRequired(t *testing.T) {
	const token module.Token = "config.required"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{}),
		config.WithTyped(token, config.ValueSpec[string]{
			Key:      "REQUIRED_KEY",
			Required: true,
			Parse:    config.ParseString,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	_, err = app.Get(token)
	if err == nil {
		t.Fatalf("expected missing required error")
	}

	var missingErr *config.MissingRequiredError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected MissingRequiredError, got %T", err)
	}
}
