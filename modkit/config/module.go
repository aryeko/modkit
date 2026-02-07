package config

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reflect"
	"sort"
	"strings"

	"github.com/go-modkit/modkit/modkit/module"
)

const moduleName = "config"

// Option configures the config module builder.
type Option func(*builder)

type builder struct {
	source  Source
	name    string
	entries []entry
}

type entry struct {
	token  module.Token
	export bool
	build  func(src Source) module.ProviderDef
}

// ValueSpec defines how to resolve and parse a typed config value.
type ValueSpec[T any] struct {
	Key         string
	Required    bool
	Default     *T
	Sensitive   bool
	Description string
	Parse       func(raw string) (T, error)
}

type mod struct {
	def module.ModuleDef
}

func (m *mod) Definition() module.ModuleDef {
	return m.def
}

// NewModule builds a regular modkit module that provides config values.
func NewModule(opts ...Option) module.Module {
	b := &builder{source: envSource{}}
	for _, opt := range opts {
		if opt != nil {
			opt(b)
		}
	}

	providers := make([]module.ProviderDef, 0, len(b.entries))
	exports := make([]module.Token, 0, len(b.entries))
	for _, e := range b.entries {
		providers = append(providers, e.build(b.source))
		if e.export {
			exports = append(exports, e.token)
		}
	}

	return &mod{def: module.ModuleDef{
		Name:      moduleNameForBuilder(b),
		Providers: providers,
		Exports:   exports,
	}}
}

// WithModuleName sets an explicit module name.
func WithModuleName(name string) Option {
	return func(b *builder) {
		b.name = strings.TrimSpace(name)
	}
}

// WithSource sets a custom key lookup source.
func WithSource(src Source) Option {
	return func(b *builder) {
		b.source = src
	}
}

// WithTyped registers a typed config provider.
func WithTyped[T any](token module.Token, spec ValueSpec[T], export bool) Option {
	return func(b *builder) {
		b.entries = append(b.entries, entry{
			token:  token,
			export: export,
			build: func(src Source) module.ProviderDef {
				return module.ProviderDef{
					Token: token,
					Build: func(_ module.Resolver) (any, error) {
						return resolve(token, spec, src)
					},
				}
			},
		})
	}
}

func resolve[T any](token module.Token, spec ValueSpec[T], src Source) (T, error) {
	var zero T

	if token == "" {
		return zero, &InvalidSpecError{Token: token, Reason: "token must not be empty"}
	}
	if spec.Key == "" {
		return zero, &InvalidSpecError{Token: token, Reason: "key must not be empty"}
	}
	if spec.Parse == nil {
		return zero, &InvalidSpecError{Token: token, Reason: "parse function must not be nil"}
	}
	if src == nil {
		return zero, &InvalidSpecError{Token: token, Reason: "source must not be nil"}
	}

	raw, ok := src.Lookup(spec.Key)
	if ok {
		raw = strings.TrimSpace(raw)
	}

	if !ok || raw == "" {
		if spec.Default != nil {
			return *spec.Default, nil
		}
		if spec.Required {
			return zero, &MissingRequiredError{Key: spec.Key, Token: token, Sensitive: spec.Sensitive}
		}
		return zero, nil
	}

	parsed, err := spec.Parse(raw)
	if err != nil {
		wrapped := fmt.Errorf("parse %q: %w", spec.Key, err)
		if spec.Sensitive {
			wrapped = errors.New("parse failed for sensitive key")
		}

		return zero, &ParseError{
			Key:       spec.Key,
			Token:     token,
			Type:      typeName[T](),
			Sensitive: spec.Sensitive,
			Err:       wrapped,
		}
	}

	return parsed, nil
}

func typeName[T any]() string {
	t := reflect.TypeFor[T]()
	if t == nil {
		return "<unknown>"
	}
	return t.String()
}

func moduleNameForBuilder(b *builder) string {
	if b.name != "" {
		return b.name
	}
	if len(b.entries) == 0 {
		return moduleName
	}

	tokens := make([]string, 0, len(b.entries))
	for _, e := range b.entries {
		tokens = append(tokens, string(e.token))
	}
	sort.Strings(tokens)

	h := fnv.New64a()
	for _, t := range tokens {
		_, _ = h.Write([]byte(t))
		_, _ = h.Write([]byte{0})
	}

	return fmt.Sprintf("%s.%x", moduleName, h.Sum64())
}
