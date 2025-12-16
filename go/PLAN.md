# Go Package: hct-mcp-signals

**Status**: Planned
**Estimated Effort**: 1-2 days

## Overview

Go implementation of HCT-MCP signals for cloud-native agent systems.

## Package Structure

```
go/
├── go.mod
├── go.sum
├── signals.go      # Main types and factory
├── schema.go       # JSON schema helpers
├── mcp.go          # MCP integration
└── signals_test.go
```

## Module Name

```
github.com/stefanwiest/hct-mcp-signals/go
```

Or as standalone:
```
github.com/hct/mcp-signals-go
```

## API Design

```go
package hctmcpsignals

import (
    "time"
)

// SignalType defines the 7 HCT coordination signals
type SignalType string

const (
    Cue      SignalType = "cue"
    Fermata  SignalType = "fermata"
    Attacca  SignalType = "attacca"
    Vamp     SignalType = "vamp"
    Caesura  SignalType = "caesura"
    Tacet    SignalType = "tacet"
    Downbeat SignalType = "downbeat"
)

// Tempo defines urgency timing
type Tempo string

const (
    Largo    Tempo = "largo"
    Andante  Tempo = "andante"
    Moderato Tempo = "moderato"
    Allegro  Tempo = "allegro"
    Presto   Tempo = "presto"
)

// HoldType for FERMATA signals
type HoldType string

const (
    Human      HoldType = "human"
    Governance HoldType = "governance"
    Resource   HoldType = "resource"
    Quality    HoldType = "quality"
)

// Performance parameters
type Performance struct {
    Urgency   int    `json:"urgency,omitempty"`
    Tempo     Tempo  `json:"tempo,omitempty"`
    TimeoutMs *int64 `json:"timeout_ms,omitempty"`
}

// Conditions for conditional signals
type Conditions struct {
    HoldType         *HoldType `json:"hold_type,omitempty"`
    RepeatUntil      *string   `json:"repeat_until,omitempty"`
    QualityThreshold *float64  `json:"quality_threshold,omitempty"`
}

// HCTSignal is the main signal struct
type HCTSignal struct {
    Type        SignalType             `json:"type"`
    Source      string                 `json:"source"`
    Targets     []string               `json:"targets"`
    Payload     map[string]interface{} `json:"payload"`
    Performance *Performance           `json:"performance,omitempty"`
    Conditions  *Conditions            `json:"conditions,omitempty"`
    Timestamp   time.Time              `json:"timestamp"`
}
```

## Factory Functions

```go
// NewCue creates a CUE signal
func NewCue(source string, targets []string, opts ...CueOption) *HCTSignal {
    sig := &HCTSignal{
        Type:      Cue,
        Source:    source,
        Targets:   targets,
        Payload:   make(map[string]interface{}),
        Performance: &Performance{
            Urgency: 5,
            Tempo:   Moderato,
        },
        Timestamp: time.Now().UTC(),
    }
    for _, opt := range opts {
        opt(sig)
    }
    return sig
}

// Option pattern
type CueOption func(*HCTSignal)

func WithUrgency(u int) CueOption {
    return func(s *HCTSignal) {
        s.Performance.Urgency = u
    }
}

func WithPayload(p map[string]interface{}) CueOption {
    return func(s *HCTSignal) {
        s.Payload = p
    }
}
```

## Usage Example

```go
package main

import (
    hct "github.com/stefanwiest/hct-mcp-signals/go"
    "encoding/json"
    "fmt"
)

func main() {
    // Create a CUE signal
    signal := hct.NewCue(
        "orchestrator",
        []string{"analyst"},
        hct.WithUrgency(8),
        hct.WithPayload(map[string]interface{}{
            "task": "Analyze Q4",
        }),
    )

    // Serialize to JSON
    data, _ := json.MarshalIndent(signal, "", "  ")
    fmt.Println(string(data))
}
```

## Publishing

```bash
# Initialize module
go mod init github.com/stefanwiest/hct-mcp-signals/go

# Test
go test ./...

# Tag for release
git tag go/v0.1.0
git push origin go/v0.1.0
```

Go modules are fetched directly from Git, no central registry needed.

## CI/CD Addition

Add to `.github/workflows/tests.yml`:

```yaml
go:
  runs-on: ubuntu-latest
  defaults:
    run:
      working-directory: go
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.22'
    - run: go fmt ./...
    - run: go vet ./...
    - run: go test -v -race ./...
```
