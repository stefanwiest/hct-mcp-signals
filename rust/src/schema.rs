//! JSON Schema export utilities.
//!
//! Provides JSON Schema (Draft 2020-12) for HCT signals,
//! enabling validation in non-Rust systems.

#[allow(unused_imports)]
use crate::signals::{HCTSignal, HoldType, Performance};
#[allow(unused_imports)]
use crate::spec::{SignalType, Tempo};
use serde_json::json;

/// Get the JSON Schema for `HCTSignal` as a `serde_json::Value`.
#[must_use]
pub fn get_json_schema() -> serde_json::Value {
    json!({
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://github.com/stefanwiest/hct-mcp-signals/schema/hct-signal.json",
        "title": "HCTSignal",
        "description": "HCT Coordination Signal for MCP",
        "type": "object",
        "required": ["type", "source", "targets", "payload", "timestamp"],
        "properties": {
            "type": {
                "type": "string",
                "enum": ["cue", "fermata", "attacca", "vamp", "caesura", "tacet", "downbeat"],
                "description": "Signal type"
            },
            "source": {
                "type": "string",
                "description": "Source agent ID"
            },
            "targets": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Target agent IDs"
            },
            "payload": {
                "type": "object",
                "additionalProperties": true,
                "description": "Signal payload"
            },
            "performance": {
                "type": "object",
                "properties": {
                    "urgency": {"type": "integer", "minimum": 1, "maximum": 10},
                    "tempo": {"type": "string", "enum": ["largo", "andante", "moderato", "allegro", "presto"]},
                    "timeout_ms": {"type": "integer", "minimum": 0}
                }
            },
            "conditions": {
                "type": "object",
                "properties": {
                    "hold_type": {"type": "string", "enum": ["human", "governance", "resource", "quality"]},
                    "repeat_until": {"type": "string"},
                    "quality_threshold": {"type": "number", "minimum": 0, "maximum": 1}
                }
            },
            "timestamp": {
                "type": "string",
                "format": "date-time"
            }
        }
    })
}

/// Get the JSON Schema as a formatted string.
///
/// # Errors
///
/// Returns an error if serialization fails.
pub fn get_json_schema_string() -> Result<String, serde_json::Error> {
    serde_json::to_string_pretty(&get_json_schema())
}

/// Get the MCP extension wrapper schema.
#[must_use]
pub fn get_mcp_extension_schema() -> serde_json::Value {
    json!({
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://github.com/stefanwiest/hct-mcp-signals/schema/mcp-extension.json",
        "title": "HCT-MCP Extension",
        "description": "HCT Signal extension for MCP params",
        "type": "object",
        "properties": {
            "hct_signal": get_json_schema()
        }
    })
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_schema_has_required_fields() {
        let schema = get_json_schema();
        let required = schema["required"].as_array().unwrap();
        assert!(required.iter().any(|v| v == "type"));
        assert!(required.iter().any(|v| v == "source"));
    }

    #[test]
    fn test_schema_string() {
        let schema_str = get_json_schema_string().unwrap();
        assert!(schema_str.contains("HCTSignal"));
        assert!(schema_str.contains("cue"));
    }

    #[test]
    fn test_mcp_extension_schema() {
        let schema = get_mcp_extension_schema();
        assert!(schema["properties"]["hct_signal"].is_object());
    }
}
