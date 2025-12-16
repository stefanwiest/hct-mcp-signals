// Package hctmcpsignals provides HCT Coordination Signals Extension for MCP.
//
// It adds urgency, timing, and approval semantics to Model Context Protocol
// based multi-agent systems using 7 coordination signals inspired by musical
// ensemble coordination.
//
// # Quick Start
//
//	signal := hctmcpsignals.NewCue("orchestrator", []string{"analyst"},
//	    hctmcpsignals.WithUrgency(8),
//	    hctmcpsignals.WithTempo(hctmcpsignals.Allegro),
//	)
//	json, _ := signal.ToJSON()
//
// # Signal Types
//
//   - Cue: Trigger agent activation
//   - Fermata: Hold for approval
//   - Attacca: Immediate transition
//   - Vamp: Repeat until condition
//   - Caesura: Full stop
//   - Tacet: Agent inactive
//   - Downbeat: Global sync point
package hctmcpsignals

import (
	"encoding/json"
	"time"
)

// SignalType defines the 7 HCT coordination signals.
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

// Tempo defines urgency timing.
type Tempo string

const (
	Largo    Tempo = "largo"    // Very slow (~1 min response OK)
	Andante  Tempo = "andante"  // Walking pace (~30s response)
	Moderato Tempo = "moderato" // Moderate (~15s response)
	Allegro  Tempo = "allegro"  // Fast (~5s response)
	Presto   Tempo = "presto"   // Very fast (~1s response)
)

// HoldType for FERMATA signals.
type HoldType string

const (
	HoldHuman      HoldType = "human"
	HoldGovernance HoldType = "governance"
	HoldResource   HoldType = "resource"
	HoldQuality    HoldType = "quality"
)

// Performance parameters (Layer 3 in HCT).
type Performance struct {
	Urgency   int    `json:"urgency,omitempty"`
	Tempo     Tempo  `json:"tempo,omitempty"`
	TimeoutMs *int64 `json:"timeout_ms,omitempty"`
}

// Conditions for conditional signals (FERMATA, VAMP).
type Conditions struct {
	HoldType         *HoldType `json:"hold_type,omitempty"`
	RepeatUntil      *string   `json:"repeat_until,omitempty"`
	QualityThreshold *float64  `json:"quality_threshold,omitempty"`
}

// HCTSignal is the main signal struct.
type HCTSignal struct {
	Type        SignalType             `json:"type"`
	Source      string                 `json:"source"`
	Targets     []string               `json:"targets"`
	Payload     map[string]interface{} `json:"payload"`
	Performance *Performance           `json:"performance,omitempty"`
	Conditions  *Conditions            `json:"conditions,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// ToJSON serializes the signal to JSON.
func (s *HCTSignal) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

// ToMCPJSON wraps the signal in hct_signal key for MCP embedding.
func (s *HCTSignal) ToMCPJSON() ([]byte, error) {
	wrapper := map[string]interface{}{
		"hct_signal": s,
	}
	return json.MarshalIndent(wrapper, "", "  ")
}

// FromJSON parses a signal from JSON bytes.
func FromJSON(data []byte) (*HCTSignal, error) {
	var s HCTSignal
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// FromMCPJSON parses a signal from MCP-wrapped JSON.
func FromMCPJSON(data []byte) (*HCTSignal, error) {
	var wrapper struct {
		HCTSignal *HCTSignal `json:"hct_signal"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.HCTSignal, nil
}
