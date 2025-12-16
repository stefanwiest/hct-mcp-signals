//! Signal type definitions and core structs.

use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// The 7 HCT coordination signal types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum SignalType {
    /// Trigger agent activation
    Cue,
    /// Hold for approval
    Fermata,
    /// Immediate transition
    Attacca,
    /// Repeat until condition met
    Vamp,
    /// Full stop
    Caesura,
    /// Agent inactive
    Tacet,
    /// Global sync point
    Downbeat,
}

/// Musical tempo indications mapped to urgency timing.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize, Default)]
#[serde(rename_all = "lowercase")]
pub enum Tempo {
    /// Very slow (~1 min response OK)
    Largo,
    /// Walking pace (~30s response)
    Andante,
    /// Moderate (~15s response)
    #[default]
    Moderato,
    /// Fast (~5s response)
    Allegro,
    /// Very fast (~1s response)
    Presto,
}

/// Types of holds for FERMATA signals.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize, Default)]
#[serde(rename_all = "lowercase")]
pub enum HoldType {
    /// Requires human approval
    #[default]
    Human,
    /// Requires governance check
    Governance,
    /// Waiting for resource
    Resource,
    /// Quality threshold not met
    Quality,
}

/// Performance parameters (Layer 3 in HCT).
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
pub struct Performance {
    /// Urgency level 1-10
    #[serde(default = "default_urgency")]
    pub urgency: u8,
    /// Expected response timing
    #[serde(default)]
    pub tempo: Tempo,
    /// Timeout in milliseconds
    #[serde(skip_serializing_if = "Option::is_none")]
    pub timeout_ms: Option<u64>,
}

fn default_urgency() -> u8 {
    5
}

impl Performance {
    /// Create new performance with defaults.
    #[must_use]
    pub fn new() -> Self {
        Self::default()
    }

    /// Set urgency (clamped to 1-10).
    #[must_use]
    pub fn with_urgency(mut self, urgency: u8) -> Self {
        self.urgency = urgency.clamp(1, 10);
        self
    }

    /// Set tempo.
    #[must_use]
    pub fn with_tempo(mut self, tempo: Tempo) -> Self {
        self.tempo = tempo;
        self
    }

    /// Set timeout.
    #[must_use]
    pub fn with_timeout_ms(mut self, timeout_ms: u64) -> Self {
        self.timeout_ms = Some(timeout_ms);
        self
    }
}

/// Conditions for conditional signals (FERMATA, VAMP).
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize, Default)]
pub struct Conditions {
    /// Type of hold for FERMATA
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hold_type: Option<HoldType>,
    /// Condition expression for VAMP
    #[serde(skip_serializing_if = "Option::is_none")]
    pub repeat_until: Option<String>,
    /// Quality score threshold 0-1
    #[serde(skip_serializing_if = "Option::is_none")]
    pub quality_threshold: Option<f64>,
}

/// Complete HCT Signal.
#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct HCTSignal {
    /// Signal type
    #[serde(rename = "type")]
    pub signal_type: SignalType,
    /// Source agent ID
    pub source: String,
    /// Target agent IDs
    pub targets: Vec<String>,
    /// Signal payload
    pub payload: HashMap<String, serde_json::Value>,
    /// Performance parameters
    #[serde(skip_serializing_if = "Option::is_none")]
    pub performance: Option<Performance>,
    /// Conditional parameters
    #[serde(skip_serializing_if = "Option::is_none")]
    pub conditions: Option<Conditions>,
    /// Timestamp
    pub timestamp: DateTime<Utc>,
}

impl HCTSignal {
    /// Create a new signal.
    #[must_use]
    pub fn new(signal_type: SignalType, source: impl Into<String>) -> Self {
        Self {
            signal_type,
            source: source.into(),
            targets: Vec::new(),
            payload: HashMap::new(),
            performance: Some(Performance::default()),
            conditions: None,
            timestamp: Utc::now(),
        }
    }

    /// Convert to MCP-compatible JSON.
    ///
    /// # Errors
    ///
    /// Returns an error if serialization fails.
    pub fn to_mcp_json(&self) -> Result<String, serde_json::Error> {
        let wrapper = serde_json::json!({
            "hct_signal": self
        });
        serde_json::to_string_pretty(&wrapper)
    }

    /// Parse from MCP JSON.
    ///
    /// # Errors
    ///
    /// Returns an error if the JSON is invalid or missing hct_signal field.
    pub fn from_mcp_json(json: &str) -> Result<Self, serde_json::Error> {
        #[derive(Deserialize)]
        struct Wrapper {
            hct_signal: HCTSignal,
        }
        let wrapper: Wrapper = serde_json::from_str(json)?;
        Ok(wrapper.hct_signal)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_signal_type_serialization() {
        let sig_type = SignalType::Cue;
        let json = serde_json::to_string(&sig_type).unwrap();
        assert_eq!(json, "\"cue\"");
    }

    #[test]
    fn test_tempo_serialization() {
        let tempo = Tempo::Allegro;
        let json = serde_json::to_string(&tempo).unwrap();
        assert_eq!(json, "\"allegro\"");
    }

    #[test]
    fn test_performance_defaults() {
        let perf = Performance::default();
        assert_eq!(perf.urgency, 5);
        assert_eq!(perf.tempo, Tempo::Moderato);
        assert!(perf.timeout_ms.is_none());
    }

    #[test]
    fn test_performance_builder() {
        let perf = Performance::new()
            .with_urgency(8)
            .with_tempo(Tempo::Presto)
            .with_timeout_ms(5000);
        
        assert_eq!(perf.urgency, 8);
        assert_eq!(perf.tempo, Tempo::Presto);
        assert_eq!(perf.timeout_ms, Some(5000));
    }

    #[test]
    fn test_urgency_clamping() {
        let perf = Performance::new().with_urgency(15);
        assert_eq!(perf.urgency, 10);
        
        let perf = Performance::new().with_urgency(0);
        assert_eq!(perf.urgency, 1);
    }

    #[test]
    fn test_signal_creation() {
        let signal = HCTSignal::new(SignalType::Cue, "orchestrator");
        assert_eq!(signal.signal_type, SignalType::Cue);
        assert_eq!(signal.source, "orchestrator");
        assert!(signal.targets.is_empty());
    }

    #[test]
    fn test_signal_to_mcp_json() {
        let signal = HCTSignal::new(SignalType::Fermata, "verifier");
        let json = signal.to_mcp_json().unwrap();
        assert!(json.contains("\"hct_signal\""));
        assert!(json.contains("\"fermata\""));
    }
}
