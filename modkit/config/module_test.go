package config_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-modkit/modkit/modkit/config"
	"github.com/go-modkit/modkit/modkit/kernel"
	"github.com/go-modkit/modkit/modkit/module"
)

type mapSource map[string]string

func (m mapSource) Lookup(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

type testModule struct {
	def module.ModuleDef
}

func (m *testModule) Definition() module.ModuleDef {
	return m.def
}

func mod(
	name string,
	imports []module.Module,
	controllers []module.ControllerDef,
) module.Module {
	return &testModule{def: module.ModuleDef{
		Name:        name,
		Imports:     imports,
		Providers:   nil,
		Controllers: controllers,
		Exports:     nil,
	}}
}

func TestWithTyped_DefaultAndParse(t *testing.T) {
	const token module.Token = "config.jwt_ttl"
	def := 1 * time.Hour

	cfgModule := config.NewModule(
		config.WithSource(mapSource{"JWT_TTL": " 90s "}),
		config.WithTyped(token, config.ValueSpec[time.Duration]{
			Key:      "JWT_TTL",
			Default:  &def,
			Parse:    config.ParseDuration,
			Required: false,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	gotAny, err := app.Get(token)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	got, ok := gotAny.(time.Duration)
	if !ok {
		t.Fatalf("got type %T", gotAny)
	}
	if got != 90*time.Second {
		t.Fatalf("got %v", got)
	}
}

func TestWithTyped_UsesDefaultWhenUnset(t *testing.T) {
	const token module.Token = "config.http_addr"
	def := ":8080"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{}),
		config.WithTyped(token, config.ValueSpec[string]{
			Key:      "HTTP_ADDR",
			Default:  &def,
			Parse:    config.ParseString,
			Required: false,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	gotAny, err := app.Get(token)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if gotAny != ":8080" {
		t.Fatalf("got %v", gotAny)
	}
}

func TestWithTyped_OptionalUnsetReturnsZeroWithoutParsing(t *testing.T) {
	const token module.Token = "config.optional_int"
	called := false

	cfgModule := config.NewModule(
		config.WithSource(mapSource{}),
		config.WithTyped(token, config.ValueSpec[int]{
			Key:      "OPTIONAL_INT",
			Required: false,
			Parse: func(_ string) (int, error) {
				called = true
				return 0, fmt.Errorf("parser should not run for unset optional value")
			},
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	gotAny, err := app.Get(token)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if called {
		t.Fatalf("parser was called")
	}

	got, ok := gotAny.(int)
	if !ok {
		t.Fatalf("got type %T", gotAny)
	}
	if got != 0 {
		t.Fatalf("got %d", got)
	}
}

func TestWithTyped_MissingRequired(t *testing.T) {
	const token module.Token = "config.jwt_secret"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{}),
		config.WithTyped(token, config.ValueSpec[string]{
			Key:       "JWT_SECRET",
			Required:  true,
			Sensitive: true,
			Parse:     config.ParseString,
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

	var buildErr *kernel.ProviderBuildError
	if !errors.As(err, &buildErr) {
		t.Fatalf("expected ProviderBuildError, got %T", err)
	}

	var missingErr *config.MissingRequiredError
	if !errors.As(err, &missingErr) {
		t.Fatalf("expected MissingRequiredError, got %T", err)
	}

	if missingErr.Key != "JWT_SECRET" {
		t.Fatalf("unexpected key: %q", missingErr.Key)
	}
	if !missingErr.Sensitive {
		t.Fatalf("expected sensitive key")
	}
}

func TestWithTyped_ParseError(t *testing.T) {
	const token module.Token = "config.rate_limit_burst"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{"RATE_LIMIT_BURST": "NaN"}),
		config.WithTyped(token, config.ValueSpec[int]{
			Key:      "RATE_LIMIT_BURST",
			Required: true,
			Parse:    config.ParseInt,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	_, err = app.Get(token)
	if err == nil {
		t.Fatalf("expected parse error")
	}

	var parseErr *config.ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected ParseError, got %T", err)
	}
	if parseErr.Key != "RATE_LIMIT_BURST" {
		t.Fatalf("unexpected key: %q", parseErr.Key)
	}
	if parseErr.Type != "int" {
		t.Fatalf("unexpected type: %q", parseErr.Type)
	}
}

func TestWithTyped_InvalidSpec(t *testing.T) {
	t.Run("empty key", func(t *testing.T) {
		const token module.Token = "config.foo"
		cfgModule := config.NewModule(
			config.WithSource(mapSource{"X": "1"}),
			config.WithTyped(token, config.ValueSpec[int]{
				Parse: config.ParseInt,
			}, true),
		)

		root := mod("root", []module.Module{cfgModule}, nil)
		app, err := kernel.Bootstrap(root)
		if err != nil {
			t.Fatalf("bootstrap failed: %v", err)
		}

		_, err = app.Get(token)
		if err == nil {
			t.Fatalf("expected invalid spec error")
		}

		var specErr *config.InvalidSpecError
		if !errors.As(err, &specErr) {
			t.Fatalf("expected InvalidSpecError, got %T", err)
		}
	})

	t.Run("nil parse", func(t *testing.T) {
		const token module.Token = "config.bar"
		cfgModule := config.NewModule(
			config.WithSource(mapSource{"BAR": "1"}),
			config.WithTyped(token, config.ValueSpec[int]{
				Key: "BAR",
			}, true),
		)

		root := mod("root", []module.Module{cfgModule}, nil)
		app, err := kernel.Bootstrap(root)
		if err != nil {
			t.Fatalf("bootstrap failed: %v", err)
		}

		_, err = app.Get(token)
		if err == nil {
			t.Fatalf("expected invalid spec error")
		}

		var specErr *config.InvalidSpecError
		if !errors.As(err, &specErr) {
			t.Fatalf("expected InvalidSpecError, got %T", err)
		}
	})
}

func TestWithTyped_SensitiveErrorDoesNotLeakValue(t *testing.T) {
	const token module.Token = "config.jwt_secret"

	cfgModule := config.NewModule(
		config.WithSource(mapSource{"JWT_SECRET": "super-secret-value"}),
		config.WithTyped(token, config.ValueSpec[int]{
			Key:       "JWT_SECRET",
			Required:  true,
			Sensitive: true,
			Parse:     config.ParseInt,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	_, err = app.Get(token)
	if err == nil {
		t.Fatalf("expected error")
	}

	if strings.Contains(err.Error(), "super-secret-value") {
		t.Fatalf("error leaked secret value: %v", err)
	}
}

func TestWithSourceNil(t *testing.T) {
	const token module.Token = "config.foo"

	cfgModule := config.NewModule(
		config.WithSource(nil),
		config.WithTyped(token, config.ValueSpec[int]{
			Key:      "FOO",
			Parse:    config.ParseInt,
			Required: true,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	_, err = app.Get(token)
	if err == nil {
		t.Fatalf("expected invalid spec error")
	}

	var specErr *config.InvalidSpecError
	if !errors.As(err, &specErr) {
		t.Fatalf("expected InvalidSpecError, got %T", err)
	}
}

func TestNoReflectionMagic_CustomParser(t *testing.T) {
	const token module.Token = "config.custom"

	parseCustom := func(raw string) (string, error) {
		if raw != "expected" {
			return "", fmt.Errorf("unsupported value")
		}
		return "ok", nil
	}

	cfgModule := config.NewModule(
		config.WithSource(mapSource{"CUSTOM": "expected"}),
		config.WithTyped(token, config.ValueSpec[string]{
			Key:      "CUSTOM",
			Required: true,
			Parse:    parseCustom,
		}, true),
	)

	root := mod("root", []module.Module{cfgModule}, nil)
	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	got, err := app.Get(token)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got != "ok" {
		t.Fatalf("got %v", got)
	}
}

func TestNewModule_DefaultNamesAvoidCollisions(t *testing.T) {
	const tokenA module.Token = "config.a"
	const tokenB module.Token = "config.b"

	cfgA := config.NewModule(
		config.WithSource(mapSource{"A": "value-a"}),
		config.WithTyped(tokenA, config.ValueSpec[string]{
			Key:      "A",
			Required: true,
			Parse:    config.ParseString,
		}, true),
	)

	cfgB := config.NewModule(
		config.WithSource(mapSource{"B": "value-b"}),
		config.WithTyped(tokenB, config.ValueSpec[string]{
			Key:      "B",
			Required: true,
			Parse:    config.ParseString,
		}, true),
	)

	root := mod("root", []module.Module{cfgA, cfgB}, []module.ControllerDef{{
		Name: "UsesTwoConfigModules",
		Build: func(r module.Resolver) (any, error) {
			a, err := module.Get[string](r, tokenA)
			if err != nil {
				return nil, err
			}
			b, err := module.Get[string](r, tokenB)
			if err != nil {
				return nil, err
			}
			return a + ":" + b, nil
		},
	}})

	app, err := kernel.Bootstrap(root)
	if err != nil {
		t.Fatalf("bootstrap failed: %v", err)
	}

	if got := app.Controllers["root:UsesTwoConfigModules"]; got != "value-a:value-b" {
		t.Fatalf("unexpected controller value: %v", got)
	}
}
