package templates

import "testing"

func TestFSIncludesTemplates(t *testing.T) {
	fs := FS()
	paths := []string{
		"go.mod.tpl",
		"main.go.tpl",
		"module.go.tpl",
		"provider.go.tpl",
		"controller.go.tpl",
		"app_module.go.tpl",
	}

	for _, p := range paths {
		if _, err := fs.ReadFile(p); err != nil {
			t.Fatalf("missing embedded template %s: %v", p, err)
		}
	}
}
