package hctmcpsignals

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPTaskSend(t *testing.T) {
	signal := NewCue("orch", []string{"analyst"})
	msg := NewMCPTaskSend("task-123", "Analyze Q4", signal)

	assert.Equal(t, "2.0", msg.JSONRPC)
	assert.Equal(t, "tasks/send", msg.Method)
	assert.Equal(t, "task-123", msg.Params.ID)
	assert.Equal(t, "user", msg.Params.Message.Role)
	assert.Equal(t, "Analyze Q4", msg.Params.Message.Content)
	assert.NotNil(t, msg.Params.HCTSignal)
}

func TestMCPTaskSendToJSON(t *testing.T) {
	signal := NewCue("orch", nil)
	msg := NewMCPTaskSend("task-123", "test", signal)

	data, err := msg.ToJSON()
	require.NoError(t, err)

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"jsonrpc": "2.0"`)
	assert.Contains(t, jsonStr, `"method": "tasks/send"`)
	assert.Contains(t, jsonStr, "hct_signal")
}

func TestMCPTaskSendRoundTrip(t *testing.T) {
	signal := NewCue("orch", []string{"analyst"}, WithUrgency(8))
	msg := NewMCPTaskSend("task-123", "Analyze Q4", signal)

	data, err := msg.ToJSON()
	require.NoError(t, err)

	var parsed MCPTaskSend
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, msg.JSONRPC, parsed.JSONRPC)
	assert.Equal(t, msg.Method, parsed.Method)
	assert.Equal(t, msg.Params.ID, parsed.Params.ID)
}

func TestEmbedSignal(t *testing.T) {
	signal := NewCue("orch", nil)
	params := []byte(`{"id": "task-123", "other": "data"}`)

	result, err := EmbedSignal(params, signal)
	require.NoError(t, err)

	resultStr := string(result)
	assert.Contains(t, resultStr, "hct_signal")
	assert.Contains(t, resultStr, "task-123")
	assert.Contains(t, resultStr, "other")
}

func TestEmbedSignalInvalidJSON(t *testing.T) {
	signal := NewCue("orch", nil)
	params := []byte(`not valid json`)

	_, err := EmbedSignal(params, signal)
	assert.Error(t, err)
}

func TestExtractSignal(t *testing.T) {
	signal := NewFermata("reporter", "test reason")
	data, err := signal.ToMCPJSON()
	require.NoError(t, err)

	extracted, err := ExtractSignal(data)
	require.NoError(t, err)
	require.NotNil(t, extracted)

	assert.Equal(t, Fermata, extracted.Type)
	assert.Equal(t, "reporter", extracted.Source)
}

func TestExtractSignalNoSignal(t *testing.T) {
	data := []byte(`{"id": "task-123", "other": "data"}`)

	extracted, err := ExtractSignal(data)
	require.NoError(t, err)
	assert.Nil(t, extracted)
}

func TestExtractSignalInvalidJSON(t *testing.T) {
	data := []byte(`not valid json`)

	_, err := ExtractSignal(data)
	assert.Error(t, err)
}

func TestHasSignal(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "with signal",
			input:    []byte(`{"hct_signal": {"type": "cue", "source": "test"}}`),
			expected: true,
		},
		{
			name:     "without signal",
			input:    []byte(`{"id": "task-123"}`),
			expected: false,
		},
		{
			name:     "null signal",
			input:    []byte(`{"hct_signal": null}`),
			expected: false,
		},
		{
			name:     "invalid json",
			input:    []byte(`not valid`),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSignal(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMCPMessage(t *testing.T) {
	msg := MCPMessage{
		Role:    "assistant",
		Content: "Hello, world!",
	}

	data, err := json.Marshal(msg)
	require.NoError(t, err)

	var parsed MCPMessage
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, msg.Role, parsed.Role)
	assert.Equal(t, msg.Content, parsed.Content)
}
