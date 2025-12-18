// HCT Signal Types - Auto-generated from hct-spec/spec.yaml
// Generated: 2025-12-18T10:37:34.791964Z
//
// DO NOT EDIT MANUALLY - Edit spec.yaml and regenerate.

package hctmcpsignals

// SignalType defines the HCT coordination signals.
type SignalType string

const (
	Cue SignalType = "cue" // Trigger activation
	Fermata SignalType = "fermata" // Hold for approval
	Attacca SignalType = "attacca" // Immediate transition
	Vamp SignalType = "vamp" // Repeat until condition
	Caesura SignalType = "caesura" // Full stop
	Tacet SignalType = "tacet" // Agent inactive
	Downbeat SignalType = "downbeat" // Global sync point
)

// Tempo defines urgency timing.
type Tempo string

const (
	Largo Tempo = "largo" // Very slow (~60s)
	Andante Tempo = "andante" // Walking pace (~30s)
	Moderato Tempo = "moderato" // Moderate (~15s)
	Allegro Tempo = "allegro" // Fast (~5s)
	Presto Tempo = "presto" // Very fast (~1s)
)

// DynamicsLevel defines resource intensity.
type DynamicsLevel string

const (
	DynamicsPP DynamicsLevel = "pp" // pianissimo
	DynamicsP DynamicsLevel = "p" // piano
	DynamicsMF DynamicsLevel = "mf" // mezzo-forte
	DynamicsF DynamicsLevel = "f" // forte
	DynamicsFF DynamicsLevel = "ff" // fortissimo
)
