# HCT-MCP Signals (Go)

<div align="center">

### The Coordination Layer for Model Context Protocol

**Express Urgency, Timing, and Synchronization in Multi-Agent Systems.**

[![Go Reference](https://img.shields.io/badge/go-reference-blue?style=for-the-badge&logo=go&logoColor=white)](https://pkg.go.dev/github.com/stefanwiest/hct-mcp-signals/go)

</div>

---

## 💡 What is HCT?

**HCT-MCP Signals** extends the protocol with 7 musical primitives proven to coordinate complex ensembles without a central conductor.

| Signal | Meanings |
|---|---|
| **CUE** | Act Now / Task dispatch |
| **FERMATA** | Hold / Human-in-the-loop |
| **ATTACCA** | Immediate / Real-time |
| **VAMP** | Loop / Retry |
| **CAESURA** | Stop / Emergency |
| **TACET** | Silent / Sleep |
| **DOWNBEAT**| Sync / Barrier |

## 📦 Installation

```bash
go get github.com/stefanwiest/hct-mcp-signals/go
```

## 🚀 Quick Start

```go
package main

import (
	"fmt"
	hct "github.com/stefanwiest/hct-mcp-signals/go"
)

func main() {
	// 1. Create a CUE signal
	signal := hct.NewCue("orch", []string{"analyst"},
		hct.WithUrgency(9),
		hct.WithTempo(hct.Presto),
	)

	// 2. Convert to MCP JSON
	json, _ := signal.ToMCPJSON()
	fmt.Println(string(json))
}
```

## 🔗 Links

- [Main Repository / Documentation](https://github.com/stefanwiest/hct-mcp-signals)
- [Protocol RFC](https://github.com/stefanwiest/hct-mcp-signals/blob/main/RFC.md)

## License

Apache-2.0
