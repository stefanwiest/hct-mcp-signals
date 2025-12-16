package hctmcpsignals

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignalTypes(t *testing.T) {
	assert.Equal(t, SignalType("cue"), Cue)
	assert.Equal(t, SignalType("fermata"), Fermata)
	assert.Equal(t, SignalType("caesura"), Caesura)
}

func TestTempoValues(t *testing.T) {
	assert.Equal(t, Tempo("largo"), Largo)
	assert.Equal(t, Tempo("presto"), Presto)
}

func TestNewCue(t *testing.T) {
	signal := NewCue("orchestrator", []string{"analyst", "synthesizer"},
		WithUrgency(8),
		WithTempo(Allegro),
		WithPayloadEntry("task", "Analyze Q4"),
	)

	assert.Equal(t, Cue, signal.Type)
	assert.Equal(t, "orchestrator", signal.Source)
	assert.Contains(t, signal.Targets, "analyst")
	assert.Equal(t, 8, signal.Performance.Urgency)
	assert.Equal(t, Allegro, signal.Performance.Tempo)
	assert.Equal(t, "Analyze Q4", signal.Payload["task"])
}

func TestNewFermata(t *testing.T) {
	signal := NewFermata("reporter", "Ready for review")

	assert.Equal(t, Fermata, signal.Type)
	assert.Equal(t, "Ready for review", signal.Payload["reason"])
	require.NotNil(t, signal.Conditions)
	require.NotNil(t, signal.Conditions.HoldType)
	assert.Equal(t, HoldHuman, *signal.Conditions.HoldType)
}

func TestNewAttacca(t *testing.T) {
	signal := NewAttacca("agent", []string{"next"})

	assert.Equal(t, Attacca, signal.Type)
	assert.Equal(t, 10, signal.Performance.Urgency)
	assert.Equal(t, Presto, signal.Performance.Tempo)
}

func TestNewVamp(t *testing.T) {
	signal := NewVamp("verifier", "score > 0.9")

	assert.Equal(t, Vamp, signal.Type)
	require.NotNil(t, signal.Conditions)
	require.NotNil(t, signal.Conditions.RepeatUntil)
	assert.Equal(t, "score > 0.9", *signal.Conditions.RepeatUntil)
}

func TestNewCaesura(t *testing.T) {
	signal := NewCaesura("governance", "Budget exceeded")

	assert.Equal(t, Caesura, signal.Type)
	assert.Contains(t, signal.Targets, "*")
	assert.Equal(t, "Budget exceeded", signal.Payload["reason"])
}

func TestNewTacet(t *testing.T) {
	signal := NewTacet("agent")

	assert.Equal(t, Tacet, signal.Type)
	assert.Empty(t, signal.Targets)
}

func TestNewDownbeat(t *testing.T) {
	signal := NewDownbeat("conductor", "daily_sync")

	assert.Equal(t, Downbeat, signal.Type)
	assert.Equal(t, "daily_sync", signal.Payload["sync_point"])
}

func TestUrgencyClamping(t *testing.T) {
	signal := NewCue("test", nil, WithUrgency(15))
	assert.Equal(t, 10, signal.Performance.Urgency)

	signal = NewCue("test", nil, WithUrgency(0))
	assert.Equal(t, 1, signal.Performance.Urgency)
}

func TestToJSON(t *testing.T) {
	signal := NewCue("orch", []string{"analyst"})
	data, err := signal.ToJSON()

	require.NoError(t, err)
	assert.Contains(t, string(data), `"type": "cue"`)
	assert.Contains(t, string(data), `"source": "orch"`)
}

func TestToMCPJSON(t *testing.T) {
	signal := NewFermata("reporter", "test")
	data, err := signal.ToMCPJSON()

	require.NoError(t, err)
	assert.Contains(t, string(data), `"hct_signal"`)
}

func TestFromJSON(t *testing.T) {
	original := NewCue("orch", []string{"analyst"}, WithUrgency(7))
	data, err := original.ToJSON()
	require.NoError(t, err)

	parsed, err := FromJSON(data)
	require.NoError(t, err)
	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.Source, parsed.Source)
}

func TestMCPTaskSend(t *testing.T) {
	signal := NewCue("orch", []string{"analyst"})
	msg := NewMCPTaskSend("task-123", "Analyze Q4", signal)

	assert.Equal(t, "2.0", msg.JSONRPC)
	assert.Equal(t, "tasks/send", msg.Method)
	assert.Equal(t, "task-123", msg.Params.ID)
	assert.NotNil(t, msg.Params.HCTSignal)

	data, err := msg.ToJSON()
	require.NoError(t, err)
	assert.Contains(t, string(data), "hct_signal")
}

func TestEmbedSignal(t *testing.T) {
	signal := NewCue("orch", nil)
	params := []byte(`{"id": "task-123"}`)

	result, err := EmbedSignal(params, signal)
	require.NoError(t, err)
	assert.Contains(t, string(result), "hct_signal")
}

func TestExtractSignal(t *testing.T) {
	signal := NewCue("orch", nil)
	data, _ := signal.ToMCPJSON()

	extracted, err := ExtractSignal(data)
	require.NoError(t, err)
	require.NotNil(t, extracted)
	assert.Equal(t, Cue, extracted.Type)
}

func TestHasSignal(t *testing.T) {
	withSignal := []byte(`{"hct_signal": {"type": "cue", "source": "test"}}`)
	withoutSignal := []byte(`{"id": "task-123"}`)

	assert.True(t, HasSignal(withSignal))
	assert.False(t, HasSignal(withoutSignal))
}

func TestJSONRoundTrip(t *testing.T) {
	original := NewCue("orchestrator", []string{"analyst", "synthesizer"},
		WithUrgency(8),
		WithTempo(Allegro),
		WithPayloadEntry("task", "Analyze Q4"),
		WithPayloadEntry("priority", "high"),
	)

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var parsed HCTSignal
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, original.Type, parsed.Type)
	assert.Equal(t, original.Source, parsed.Source)
	assert.Equal(t, original.Targets, parsed.Targets)
	assert.Equal(t, original.Performance.Urgency, parsed.Performance.Urgency)
}
