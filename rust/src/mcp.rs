//! MCP integration utilities.

use crate::signals::HCTSignal;
use serde::{Deserialize, Serialize};

/// MCP task/send message with HCT signal extension.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct McpTaskSend {
    pub jsonrpc: String,
    pub method: String,
    pub params: McpParams,
}

/// MCP message parameters.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct McpParams {
    pub id: String,
    pub message: McpMessage,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hct_signal: Option<HCTSignal>,
}

/// MCP message content.
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct McpMessage {
    pub role: String,
    pub content: String,
}

impl McpTaskSend {
    /// Create a new MCP tasks/send message with HCT signal.
    #[must_use]
    pub fn new(task_id: impl Into<String>, content: impl Into<String>, signal: HCTSignal) -> Self {
        Self {
            jsonrpc: "2.0".to_string(),
            method: "tasks/send".to_string(),
            params: McpParams {
                id: task_id.into(),
                message: McpMessage {
                    role: "user".to_string(),
                    content: content.into(),
                },
                hct_signal: Some(signal),
            },
        }
    }

    /// Serialize to JSON string.
    ///
    /// # Errors
    ///
    /// Returns an error if serialization fails.
    pub fn to_json(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string_pretty(self)
    }
}

/// Embed an HCT signal into existing JSON params.
///
/// # Errors
///
/// Returns an error if the input is not a valid JSON object.
pub fn embed_signal(
    params_json: &str,
    signal: &HCTSignal,
) -> Result<String, serde_json::Error> {
    let mut params: serde_json::Value = serde_json::from_str(params_json)?;
    if let Some(obj) = params.as_object_mut() {
        obj.insert("hct_signal".to_string(), serde_json::to_value(signal)?);
    }
    serde_json::to_string_pretty(&params)
}

/// Extract an HCT signal from JSON params.
///
/// # Errors
///
/// Returns an error if parsing fails.
pub fn extract_signal(params_json: &str) -> Result<Option<HCTSignal>, serde_json::Error> {
    #[derive(Deserialize)]
    struct Wrapper {
        hct_signal: Option<HCTSignal>,
    }
    let wrapper: Wrapper = serde_json::from_str(params_json)?;
    Ok(wrapper.hct_signal)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::factory::cue;

    #[test]
    fn test_mcp_task_send() {
        let signal = cue("orch", ["analyst"]).build();
        let msg = McpTaskSend::new("task-123", "Analyze Q4", signal);

        assert_eq!(msg.jsonrpc, "2.0");
        assert_eq!(msg.method, "tasks/send");
        assert_eq!(msg.params.id, "task-123");
        assert!(msg.params.hct_signal.is_some());
    }

    #[test]
    fn test_embed_signal() {
        let signal = cue("orch", ["analyst"]).build();
        let params = r#"{"id": "task-123"}"#;
        
        let result = embed_signal(params, &signal).unwrap();
        assert!(result.contains("hct_signal"));
    }

    #[test]
    fn test_extract_signal() {
        let signal = cue("orch", ["analyst"]).build();
        let json = signal.to_mcp_json().unwrap();
        
        let extracted = extract_signal(&json).unwrap();
        assert!(extracted.is_some());
    }
}
