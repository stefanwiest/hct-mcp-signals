package hctmcpsignals

import (
	"encoding/json"
)

// MCPTaskSend represents an MCP tasks/send message with HCT signal.
type MCPTaskSend struct {
	JSONRPC string    `json:"jsonrpc"`
	Method  string    `json:"method"`
	Params  MCPParams `json:"params"`
}

// MCPParams contains the parameters for an MCP task.
type MCPParams struct {
	ID        string                 `json:"id"`
	Message   MCPMessage             `json:"message"`
	HCTSignal *HCTSignal             `json:"hct_signal,omitempty"`
	Extra     map[string]interface{} `json:"-"` // For additional params
}

// MCPMessage represents the message content.
type MCPMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// NewMCPTaskSend creates a new MCP tasks/send message with HCT signal.
func NewMCPTaskSend(taskID, content string, signal *HCTSignal) *MCPTaskSend {
	return &MCPTaskSend{
		JSONRPC: "2.0",
		Method:  "tasks/send",
		Params: MCPParams{
			ID: taskID,
			Message: MCPMessage{
				Role:    "user",
				Content: content,
			},
			HCTSignal: signal,
		},
	}
}

// ToJSON serializes the MCP message to JSON.
func (m *MCPTaskSend) ToJSON() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// EmbedSignal adds an HCT signal to existing JSON params.
func EmbedSignal(paramsJSON []byte, signal *HCTSignal) ([]byte, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(paramsJSON, &params); err != nil {
		return nil, err
	}
	params["hct_signal"] = signal
	return json.MarshalIndent(params, "", "  ")
}

// ExtractSignal extracts an HCT signal from JSON params.
func ExtractSignal(paramsJSON []byte) (*HCTSignal, error) {
	var wrapper struct {
		HCTSignal *HCTSignal `json:"hct_signal"`
	}
	if err := json.Unmarshal(paramsJSON, &wrapper); err != nil {
		return nil, err
	}
	return wrapper.HCTSignal, nil
}

// HasSignal checks if JSON params contain an HCT signal.
func HasSignal(paramsJSON []byte) bool {
	var wrapper struct {
		HCTSignal *HCTSignal `json:"hct_signal"`
	}
	if err := json.Unmarshal(paramsJSON, &wrapper); err != nil {
		return false
	}
	return wrapper.HCTSignal != nil
}
