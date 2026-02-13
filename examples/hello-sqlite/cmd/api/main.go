package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/go-modkit/modkit/examples/hello-sqlite/internal/app"
	mkhttp "github.com/go-modkit/modkit/modkit/http"
	"github.com/go-modkit/modkit/modkit/kernel"
	"github.com/go-modkit/modkit/modkit/module"
)

type HealthController struct{}

func (c *HealthController) RegisterRoutes(r mkhttp.Router) {
	r.Handle(http.MethodGet, "/health", http.HandlerFunc(c.health))
}

func (c *HealthController) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

type RootModule struct{}

func (m *RootModule) Definition() module.ModuleDef {
	return module.ModuleDef{
		Name: "api",
		Imports: []module.Module{
			app.NewModule(),
		},
		Controllers: []module.ControllerDef{{
			Name: "HealthController",
			Build: func(_ module.Resolver) (any, error) {
				return &HealthController{}, nil
			},
		}},
	}
}

func main() {
	app, err := kernel.Bootstrap(&RootModule{})
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}

	router := mkhttp.NewRouter()
	if err := mkhttp.RegisterRoutes(mkhttp.AsRouter(router), app.Controllers); err != nil {
		log.Fatalf("routes: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":8080"
	srv := &http.Server{Addr: addr, Handler: router}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	log.Printf("Server starting on http://localhost%s", addr)
	log.Println("Try: curl http://localhost:8080/health")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
}
