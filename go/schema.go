package hctmcpsignals

import "encoding/json"

// JSONSchema returns the JSON Schema for HCTSignal as a map.
func JSONSchema() map[string]interface{} {
	return map[string]interface{}{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://github.com/stefanwiest/hct-mcp-signals/schema/hct-signal.json",
		"title":       "HCTSignal",
		"description": "HCT Coordination Signal for MCP",
		"type":        "object",
		"required":    []string{"type", "source", "targets", "payload", "timestamp"},
		"properties": map[string]interface{}{
			"type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"cue", "fermata", "attacca", "vamp", "caesura", "tacet", "downbeat"},
				"description": "Signal type",
			},
			"source": map[string]interface{}{
				"type":        "string",
				"description": "Source agent ID",
			},
			"targets": map[string]interface{}{
				"type":        "array",
				"items":       map[string]interface{}{"type": "string"},
				"description": "Target agent IDs",
			},
			"payload": map[string]interface{}{
				"type":                 "object",
				"additionalProperties": true,
				"description":          "Signal payload",
			},
			"performance": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"urgency":    map[string]interface{}{"type": "integer", "minimum": 1, "maximum": 10},
					"tempo":      map[string]interface{}{"type": "string", "enum": []string{"largo", "andante", "moderato", "allegro", "presto"}},
					"timeout_ms": map[string]interface{}{"type": "integer", "minimum": 0},
				},
			},
			"conditions": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"hold_type":         map[string]interface{}{"type": "string", "enum": []string{"human", "governance", "resource", "quality"}},
					"repeat_until":      map[string]interface{}{"type": "string"},
					"quality_threshold": map[string]interface{}{"type": "number", "minimum": 0, "maximum": 1},
				},
			},
			"timestamp": map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
		},
	}
}

// JSONSchemaString returns the JSON Schema as a formatted string.
func JSONSchemaString() (string, error) {
	data, err := json.MarshalIndent(JSONSchema(), "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// MCPExtensionSchema returns the schema for the MCP extension wrapper.
func MCPExtensionSchema() map[string]interface{} {
	return map[string]interface{}{
		"$schema":     "https://json-schema.org/draft/2020-12/schema",
		"$id":         "https://github.com/stefanwiest/hct-mcp-signals/schema/mcp-extension.json",
		"title":       "HCT-MCP Extension",
		"description": "HCT Signal extension for MCP params",
		"type":        "object",
		"properties": map[string]interface{}{
			"hct_signal": JSONSchema(),
		},
	}
}
