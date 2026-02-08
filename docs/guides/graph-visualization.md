# Graph Visualization

This guide shows how to export a bootstrapped modkit dependency graph into Mermaid or DOT for architecture inspection.

Graph export is read-only serialization of the existing kernel graph. It does not instantiate providers or mutate graph state.

## API

Use the kernel exporters:

```go
func ExportGraph(g *Graph, format GraphFormat) (string, error)
func ExportAppGraph(app *App, format GraphFormat) (string, error)
```

Formats:

- `kernel.GraphFormatMermaid`
- `kernel.GraphFormatDOT`

Error behavior:

- `kernel.ErrNilApp` from `ExportAppGraph(nil, ...)`
- `kernel.ErrNilGraph` for nil graph input
- `*kernel.UnsupportedGraphFormatError` for unknown formats

## Export After Bootstrap

```go
package main

import (
    "fmt"
    "log"

    "github.com/go-modkit/modkit/modkit/kernel"
)

func main() {
    app, err := kernel.Bootstrap(newAppModule())
    if err != nil {
        log.Fatal(err)
    }

    mermaid, err := kernel.ExportAppGraph(app, kernel.GraphFormatMermaid)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(mermaid)
}
```

## Mermaid Output

Example:

```text
graph TD
    m0["app"]
    m1["auth"]
    m2["db"]
    m3["users"]
    m0 --> m1
    m0 --> m3
    m3 --> m2
    classDef root stroke-width:3px;
    class m0 root;
```

Notes:

- Node IDs are deterministic (`m0`, `m1`, ...), allocated from sorted module names.
- Labels are escaped, so module names with quotes/backslashes stay valid.
- Root module is annotated with class `root`.

## DOT Output

Example:

```text
digraph modkit {
    rankdir=LR;
    "app";
    "app" [shape=doublecircle];
    "users";
    "app" -> "users";
}
```

Notes:

- Node IDs and edges are always quoted.
- Root node uses `shape=doublecircle`.

## Determinism and Edge Semantics

- Nodes are emitted by sorted module name.
- Imports for each node are emitted in sorted order.
- An edge `A -> B` means module `A` directly imports module `B`.
- Re-export visibility does not add new graph edges.

## Related Docs

- [Modules](modules.md)
- [Providers](providers.md)
- [Architecture](../architecture.md)
- [P2 Design Spec](../specs/design-p2-graph-visualization-and-devtools.md)
