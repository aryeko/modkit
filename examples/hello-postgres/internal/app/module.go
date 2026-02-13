package app

import (
	"github.com/go-modkit/modkit/modkit/data/postgres"
	"github.com/go-modkit/modkit/modkit/data/sqlmodule"
	"github.com/go-modkit/modkit/modkit/module"
)

type Module struct {
	postgres module.Module
}

func NewModule() module.Module {
	return &Module{
		postgres: postgres.NewModule(postgres.Options{}),
	}
}

func (m *Module) Definition() module.ModuleDef {
	return module.ModuleDef{
		Name: "app",
		Imports: []module.Module{
			m.postgres,
		},
		Exports: []module.Token{
			sqlmodule.TokenDB,
			sqlmodule.TokenDialect,
		},
	}
}
