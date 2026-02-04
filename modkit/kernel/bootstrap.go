package kernel

import "github.com/aryeko/modkit/modkit/module"

type App struct {
	Graph       *Graph
	container   *Container
	Controllers map[string]any
}

func Bootstrap(root module.Module) (*App, error) {
	graph, err := BuildGraph(root)
	if err != nil {
		return nil, err
	}

	visibility, err := buildVisibility(graph)
	if err != nil {
		return nil, err
	}

	container, err := newContainer(graph, visibility)
	if err != nil {
		return nil, err
	}

	controllers := make(map[string]any)
	for _, node := range graph.Modules {
		resolver := container.resolverFor(node.Name)
		for _, controller := range node.Def.Controllers {
			if _, exists := controllers[controller.Name]; exists {
				return nil, &DuplicateControllerNameError{Name: controller.Name}
			}
			instance, err := controller.Build(resolver)
			if err != nil {
				return nil, &ControllerBuildError{Module: node.Name, Controller: controller.Name, Err: err}
			}
			controllers[controller.Name] = instance
		}
	}

	return &App{
		Graph:       graph,
		container:   container,
		Controllers: controllers,
	}, nil
}

// Resolver returns a root-scoped resolver that enforces module visibility.
func (a *App) Resolver() module.Resolver {
	return a.container.resolverFor(a.Graph.Root)
}

// Get resolves a token from the root module scope.
func (a *App) Get(token module.Token) (any, error) {
	return a.Resolver().Get(token)
}
