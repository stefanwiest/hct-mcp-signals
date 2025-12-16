package hctmcpsignals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test all factory function options for 100% coverage

func TestWithUrgencyOption(t *testing.T) {
	signal := NewCue("test", nil, WithUrgency(7))
	assert.Equal(t, 7, signal.Performance.Urgency)
}

func TestWithTempoOption(t *testing.T) {
	signal := NewCue("test", nil, WithTempo(Largo))
	assert.Equal(t, Largo, signal.Performance.Tempo)
}

func TestWithTimeoutMsOption(t *testing.T) {
	signal := NewCue("test", nil, WithTimeoutMs(30000))
	assert.NotNil(t, signal.Performance.TimeoutMs)
	assert.Equal(t, int64(30000), *signal.Performance.TimeoutMs)
}

func TestWithPayloadOption(t *testing.T) {
	payload := map[string]interface{}{"key": "value"}
	signal := NewCue("test", nil, WithPayload(payload))
	assert.Equal(t, "value", signal.Payload["key"])
}

func TestWithPayloadEntryOption(t *testing.T) {
	signal := NewCue("test", nil, WithPayloadEntry("task", "analyze"))
	assert.Equal(t, "analyze", signal.Payload["task"])
}

func TestWithPayloadEntryNilPayload(t *testing.T) {
	// Start with signal that has nil payload
	signal := &HCTSignal{}
	WithPayloadEntry("key", "value")(signal)
	assert.NotNil(t, signal.Payload)
	assert.Equal(t, "value", signal.Payload["key"])
}

func TestWithHoldTypeOption(t *testing.T) {
	signal := NewFermata("test", "reason", WithHoldType(HoldGovernance))
	assert.NotNil(t, signal.Conditions)
	assert.NotNil(t, signal.Conditions.HoldType)
	assert.Equal(t, HoldGovernance, *signal.Conditions.HoldType)
}

func TestWithHoldTypeNilConditions(t *testing.T) {
	signal := &HCTSignal{}
	WithHoldType(HoldResource)(signal)
	assert.NotNil(t, signal.Conditions)
	assert.Equal(t, HoldResource, *signal.Conditions.HoldType)
}

func TestWithRepeatUntilOption(t *testing.T) {
	signal := NewVamp("test", "score > 0.9", WithRepeatUntil("custom condition"))
	assert.NotNil(t, signal.Conditions)
	assert.NotNil(t, signal.Conditions.RepeatUntil)
	assert.Equal(t, "custom condition", *signal.Conditions.RepeatUntil)
}

func TestWithRepeatUntilNilConditions(t *testing.T) {
	signal := &HCTSignal{}
	WithRepeatUntil("condition")(signal)
	assert.NotNil(t, signal.Conditions)
	assert.Equal(t, "condition", *signal.Conditions.RepeatUntil)
}

func TestWithQualityThresholdOption(t *testing.T) {
	signal := NewVamp("test", "cond", WithQualityThreshold(0.95))
	assert.NotNil(t, signal.Conditions)
	assert.NotNil(t, signal.Conditions.QualityThreshold)
	assert.Equal(t, 0.95, *signal.Conditions.QualityThreshold)
}

func TestWithQualityThresholdNilConditions(t *testing.T) {
	signal := &HCTSignal{}
	WithQualityThreshold(0.8)(signal)
	assert.NotNil(t, signal.Conditions)
	assert.Equal(t, 0.8, *signal.Conditions.QualityThreshold)
}

func TestWithUrgencyNilPerformance(t *testing.T) {
	signal := &HCTSignal{}
	WithUrgency(5)(signal)
	assert.NotNil(t, signal.Performance)
	assert.Equal(t, 5, signal.Performance.Urgency)
}

func TestWithTempoNilPerformance(t *testing.T) {
	signal := &HCTSignal{}
	WithTempo(Allegro)(signal)
	assert.NotNil(t, signal.Performance)
	assert.Equal(t, Allegro, signal.Performance.Tempo)
}

func TestWithTimeoutMsNilPerformance(t *testing.T) {
	signal := &HCTSignal{}
	WithTimeoutMs(5000)(signal)
	assert.NotNil(t, signal.Performance)
	assert.Equal(t, int64(5000), *signal.Performance.TimeoutMs)
}

func TestAllHoldTypes(t *testing.T) {
	types := []HoldType{HoldHuman, HoldGovernance, HoldResource, HoldQuality}
	expected := []string{"human", "governance", "resource", "quality"}

	for i, ht := range types {
		assert.Equal(t, HoldType(expected[i]), ht)
	}
}

func TestAllTempos(t *testing.T) {
	tempos := []Tempo{Largo, Andante, Moderato, Allegro, Presto}
	expected := []string{"largo", "andante", "moderato", "allegro", "presto"}

	for i, tempo := range tempos {
		assert.Equal(t, Tempo(expected[i]), tempo)
	}
}

func TestAllSignalTypes(t *testing.T) {
	types := []SignalType{Cue, Fermata, Attacca, Vamp, Caesura, Tacet, Downbeat}
	expected := []string{"cue", "fermata", "attacca", "vamp", "caesura", "tacet", "downbeat"}

	for i, st := range types {
		assert.Equal(t, SignalType(expected[i]), st)
	}
}

func TestMultipleOptions(t *testing.T) {
	signal := NewCue("orch", []string{"a", "b"},
		WithUrgency(9),
		WithTempo(Presto),
		WithTimeoutMs(1000),
		WithPayloadEntry("task", "test"),
		WithPayloadEntry("priority", "high"),
	)

	assert.Equal(t, 9, signal.Performance.Urgency)
	assert.Equal(t, Presto, signal.Performance.Tempo)
	assert.Equal(t, int64(1000), *signal.Performance.TimeoutMs)
	assert.Equal(t, "test", signal.Payload["task"])
	assert.Equal(t, "high", signal.Payload["priority"])
}
