package hctmcpsignals

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONSchema(t *testing.T) {
	schema := JSONSchema()

	assert.Equal(t, "HCTSignal", schema["title"])
	assert.Equal(t, "object", schema["type"])

	required, ok := schema["required"].([]string)
	assert.True(t, ok)
	assert.Contains(t, required, "type")
	assert.Contains(t, required, "source")
	assert.Contains(t, required, "targets")
	assert.Contains(t, required, "payload")
	assert.Contains(t, required, "timestamp")
}

func TestJSONSchemaProperties(t *testing.T) {
	schema := JSONSchema()
	props, ok := schema["properties"].(map[string]interface{})
	assert.True(t, ok)

	// Check type property
	typeProp, ok := props["type"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "string", typeProp["type"])

	typeEnum, ok := typeProp["enum"].([]string)
	assert.True(t, ok)
	assert.Contains(t, typeEnum, "cue")
	assert.Contains(t, typeEnum, "fermata")
	assert.Contains(t, typeEnum, "attacca")
	assert.Contains(t, typeEnum, "vamp")
	assert.Contains(t, typeEnum, "caesura")
	assert.Contains(t, typeEnum, "tacet")
	assert.Contains(t, typeEnum, "downbeat")
}

func TestJSONSchemaPerformance(t *testing.T) {
	schema := JSONSchema()
	props := schema["properties"].(map[string]interface{})
	perf := props["performance"].(map[string]interface{})
	perfProps := perf["properties"].(map[string]interface{})

	urgency := perfProps["urgency"].(map[string]interface{})
	assert.Equal(t, "integer", urgency["type"])
	assert.Equal(t, 1, urgency["minimum"])
	assert.Equal(t, 10, urgency["maximum"])

	tempo := perfProps["tempo"].(map[string]interface{})
	tempoEnum := tempo["enum"].([]string)
	assert.Contains(t, tempoEnum, "largo")
	assert.Contains(t, tempoEnum, "presto")
}

func TestJSONSchemaConditions(t *testing.T) {
	schema := JSONSchema()
	props := schema["properties"].(map[string]interface{})
	cond := props["conditions"].(map[string]interface{})
	condProps := cond["properties"].(map[string]interface{})

	holdType := condProps["hold_type"].(map[string]interface{})
	holdEnum := holdType["enum"].([]string)
	assert.Contains(t, holdEnum, "human")
	assert.Contains(t, holdEnum, "governance")

	qualThresh := condProps["quality_threshold"].(map[string]interface{})
	assert.Equal(t, "number", qualThresh["type"])
	assert.Equal(t, 0, qualThresh["minimum"])
	assert.Equal(t, 1, qualThresh["maximum"])
}

func TestJSONSchemaString(t *testing.T) {
	schemaStr, err := JSONSchemaString()
	assert.NoError(t, err)
	assert.Contains(t, schemaStr, "HCTSignal")
	assert.Contains(t, schemaStr, "json-schema.org")
	assert.Contains(t, schemaStr, "cue")
}

func TestMCPExtensionSchema(t *testing.T) {
	schema := MCPExtensionSchema()

	assert.Equal(t, "HCT-MCP Extension", schema["title"])
	assert.Equal(t, "object", schema["type"])

	props := schema["properties"].(map[string]interface{})
	hctSignal := props["hct_signal"].(map[string]interface{})
	assert.Equal(t, "HCTSignal", hctSignal["title"])
}
