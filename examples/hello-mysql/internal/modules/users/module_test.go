// Package users provides the users domain module.
package users

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/auth"
	"github.com/go-modkit/modkit/examples/hello-mysql/internal/modules/database"
	"github.com/go-modkit/modkit/modkit/module"
)

// TestUsersModule_Definition_WiresAuth tests that the users module correctly declares its dependencies.
func TestUsersModule_Definition_WiresAuth(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(*Module).Definition()

	if def.Name != "users" {
		t.Fatalf("name = %q", def.Name)
	}
	if len(def.Imports) != 2 {
		t.Fatalf("imports = %d", len(def.Imports))
	}
}

type stubResolver struct {
	values map[module.Token]any
	errors map[module.Token]error
}

func (r stubResolver) Get(token module.Token) (any, error) {
	if err := r.errors[token]; err != nil {
		return nil, err
	}
	if value, ok := r.values[token]; ok {
		return value, nil
	}
	return nil, nil
}

type serviceStub struct{}

func (serviceStub) GetUser(ctx context.Context, id int64) (User, error) {
	return User{}, nil
}

func (serviceStub) CreateUser(ctx context.Context, input CreateUserInput) (User, error) {
	return User{}, nil
}

func (serviceStub) ListUsers(ctx context.Context) ([]User, error) {
	return nil, nil
}

func (serviceStub) UpdateUser(ctx context.Context, id int64, input UpdateUserInput) (User, error) {
	return User{}, nil
}

func (serviceStub) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

func (serviceStub) LongOperation(ctx context.Context) error {
	return nil
}

func TestUsersModule_ControllerBuildErrors(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(*Module).Definition()
	controller := def.Controllers[0]

	_, err := controller.Build(stubResolver{
		errors: map[module.Token]error{
			TokenService: errors.New("missing service"),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing service")
	}

	_, err = controller.Build(stubResolver{
		values: map[module.Token]any{
			TokenService: serviceStub{},
		},
		errors: map[module.Token]error{
			auth.TokenMiddleware: errors.New("missing middleware"),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing middleware")
	}

	_, err = controller.Build(stubResolver{
		values: map[module.Token]any{
			TokenService:         serviceStub{},
			auth.TokenMiddleware: func(next http.Handler) http.Handler { return next },
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("type mismatch service", func(t *testing.T) {
		_, err := controller.Build(stubResolver{
			values: map[module.Token]any{
				TokenService: "wrong type",
			},
		})
		if err == nil {
			t.Fatal("expected error for type mismatch")
		}
		if !strings.Contains(err.Error(), "expected users.Service") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})

	t.Run("type mismatch middleware", func(t *testing.T) {
		_, err := controller.Build(stubResolver{
			values: map[module.Token]any{
				TokenService:         serviceStub{},
				auth.TokenMiddleware: "wrong type",
			},
		})
		if err == nil {
			t.Fatal("expected error for type mismatch")
		}
		if !strings.Contains(err.Error(), "expected func(http.Handler) http.Handler") {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

func TestUsersModule_RepositoryBuildError(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(*Module).Definition()
	provider := def.Providers[0] // TokenRepository

	_, err := provider.Build(stubResolver{
		errors: map[module.Token]error{
			database.TokenDB: errors.New("database connection failed"),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing database")
	}
	if !strings.Contains(err.Error(), "database connection failed") {
		t.Fatalf("expected 'database connection failed' error, got %q", err.Error())
	}
}

func TestUsersModule_ServiceBuildError(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(*Module).Definition()
	provider := def.Providers[1] // TokenService

	_, err := provider.Build(stubResolver{
		errors: map[module.Token]error{
			TokenRepository: errors.New("repository not found"),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing repository")
	}
	if !strings.Contains(err.Error(), "repository not found") {
		t.Fatalf("expected 'repository not found' error, got %q", err.Error())
	}
}

func TestUsersModule_ProviderBuildSuccess(t *testing.T) {
	mod := NewModule(Options{Database: &database.Module{}, Auth: auth.NewModule(auth.Options{})})
	def := mod.(*Module).Definition()

	t.Run("repository success", func(t *testing.T) {
		provider := def.Providers[0]
		_, err := provider.Build(stubResolver{
			values: map[module.Token]any{
				database.TokenDB: &sql.DB{},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("service success", func(t *testing.T) {
		provider := def.Providers[1]
		_, err := provider.Build(stubResolver{
			values: map[module.Token]any{
				TokenRepository: &mysqlRepo{},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
