package httpserver

import (
	"net/http"

	"github.com/aryeko/modkit/examples/hello-mysql/internal/modules/app"
	modkithttp "github.com/aryeko/modkit/modkit/http"
	"github.com/aryeko/modkit/modkit/kernel"
)

func BuildHandler(opts app.Options) (http.Handler, error) {
	mod := app.NewModule(opts)
	boot, err := kernel.Bootstrap(mod)
	if err != nil {
		return nil, err
	}

	router := modkithttp.NewRouter()
	if err := modkithttp.RegisterRoutes(modkithttp.AsRouter(router), boot.Controllers); err != nil {
		return nil, err
	}

	return router, nil
}
