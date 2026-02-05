package auth

import "github.com/go-modkit/modkit/modkit/module"

const (
	TokenMiddleware module.Token = "auth.middleware"
	TokenHandler    module.Token = "auth.handler"
)

type Options struct{}

type Module struct {
	opts Options
}

type AuthModule = Module

func NewModule(opts Options) module.Module {
	return &Module{opts: opts}
}

func (m Module) Definition() module.ModuleDef {
	return module.ModuleDef{
		Name: "auth",
		Providers: []module.ProviderDef{
			{
				Token: TokenMiddleware,
				Build: func(r module.Resolver) (any, error) {
					return nil, nil
				},
			},
			{
				Token: TokenHandler,
				Build: func(r module.Resolver) (any, error) {
					return nil, nil
				},
			},
		},
	}
}
