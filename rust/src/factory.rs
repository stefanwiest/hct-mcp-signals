//! Factory functions for creating HCT signals with builder pattern.

use crate::signals::{Conditions, HCTSignal, HoldType, Performance, SignalType, Tempo};
use chrono::Utc;
use std::collections::HashMap;

/// Builder for creating HCT signals.
#[derive(Debug, Clone)]
pub struct SignalBuilder {
    signal_type: SignalType,
    source: String,
    targets: Vec<String>,
    payload: HashMap<String, serde_json::Value>,
    performance: Performance,
    conditions: Option<Conditions>,
}

impl SignalBuilder {
    /// Create a new builder.
    fn new(signal_type: SignalType, source: impl Into<String>) -> Self {
        Self {
            signal_type,
            source: source.into(),
            targets: Vec::new(),
            payload: HashMap::new(),
            performance: Performance::default(),
            conditions: None,
        }
    }

    /// Add targets.
    #[must_use]
    pub fn with_targets(mut self, targets: impl IntoIterator<Item = impl Into<String>>) -> Self {
        self.targets = targets.into_iter().map(Into::into).collect();
        self
    }

    /// Add payload.
    #[must_use]
    pub fn with_payload(mut self, payload: HashMap<String, serde_json::Value>) -> Self {
        self.payload = payload;
        self
    }

    /// Add a single payload entry.
    #[must_use]
    pub fn with_payload_entry(mut self, key: impl Into<String>, value: serde_json::Value) -> Self {
        self.payload.insert(key.into(), value);
        self
    }

    /// Set urgency (1-10).
    #[must_use]
    pub fn with_urgency(mut self, urgency: u8) -> Self {
        self.performance.urgency = urgency.clamp(1, 10);
        self
    }

    /// Set tempo.
    #[must_use]
    pub fn with_tempo(mut self, tempo: Tempo) -> Self {
        self.performance.tempo = tempo;
        self
    }

    /// Set timeout.
    #[must_use]
    pub fn with_timeout_ms(mut self, timeout_ms: u64) -> Self {
        self.performance.timeout_ms = Some(timeout_ms);
        self
    }

    /// Set conditions.
    #[must_use]
    pub fn with_conditions(mut self, conditions: Conditions) -> Self {
        self.conditions = Some(conditions);
        self
    }

    /// Build the signal.
    #[must_use]
    pub fn build(self) -> HCTSignal {
        HCTSignal {
            signal_type: self.signal_type,
            source: self.source,
            targets: self.targets,
            payload: self.payload,
            performance: Some(self.performance),
            conditions: self.conditions,
            timestamp: Utc::now(),
        }
    }
}

/// Create a CUE signal to trigger agent activation.
#[must_use]
pub fn cue(source: impl Into<String>, targets: impl IntoIterator<Item = impl Into<String>>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Cue, source)
        .with_targets(targets)
}

/// Create a FERMATA signal to hold for approval.
#[must_use]
pub fn fermata(source: impl Into<String>, reason: impl Into<String>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Fermata, source)
        .with_targets(["governance"])
        .with_payload_entry("reason", serde_json::Value::String(reason.into()))
        .with_conditions(Conditions {
            hold_type: Some(HoldType::Human),
            ..Default::default()
        })
}

/// Create an ATTACCA signal for immediate transition.
#[must_use]
pub fn attacca(source: impl Into<String>, targets: impl IntoIterator<Item = impl Into<String>>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Attacca, source)
        .with_targets(targets)
        .with_urgency(10)
        .with_tempo(Tempo::Presto)
}

/// Create a VAMP signal to repeat until condition met.
#[must_use]
pub fn vamp(source: impl Into<String>, repeat_until: impl Into<String>) -> SignalBuilder {
    let source_str = source.into();
    SignalBuilder::new(SignalType::Vamp, source_str.clone())
        .with_targets([source_str])
        .with_conditions(Conditions {
            repeat_until: Some(repeat_until.into()),
            quality_threshold: Some(0.9),
            ..Default::default()
        })
        .with_timeout_ms(60_000)
}

/// Create a CAESURA signal for full stop.
#[must_use]
pub fn caesura(source: impl Into<String>, reason: impl Into<String>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Caesura, source)
        .with_targets(["*"])
        .with_payload_entry("reason", serde_json::Value::String(reason.into()))
        .with_urgency(10)
        .with_tempo(Tempo::Presto)
}

/// Create a TACET signal to mark agent as inactive.
#[must_use]
pub fn tacet(source: impl Into<String>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Tacet, source)
}

/// Create a DOWNBEAT signal for global synchronization.
#[must_use]
pub fn downbeat(source: impl Into<String>, sync_point: impl Into<String>) -> SignalBuilder {
    SignalBuilder::new(SignalType::Downbeat, source)
        .with_targets(["*"])
        .with_payload_entry("sync_point", serde_json::Value::String(sync_point.into()))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cue_builder() {
        let signal = cue("orchestrator", ["analyst", "synthesizer"])
            .with_urgency(8)
            .with_tempo(Tempo::Allegro)
            .with_payload_entry("task", serde_json::json!("Analyze Q4"))
            .build();

        assert_eq!(signal.signal_type, SignalType::Cue);
        assert_eq!(signal.source, "orchestrator");
        assert_eq!(signal.targets, vec!["analyst", "synthesizer"]);
        assert_eq!(signal.performance.as_ref().map(|p| p.urgency), Some(8));
    }

    #[test]
    fn test_fermata_builder() {
        let signal = fermata("reporter", "Ready for review").build();

        assert_eq!(signal.signal_type, SignalType::Fermata);
        assert!(signal.payload.contains_key("reason"));
        assert!(signal.conditions.is_some());
    }

    #[test]
    fn test_attacca_is_urgent() {
        let signal = attacca("agent", ["next"]).build();

        assert_eq!(signal.performance.as_ref().map(|p| p.urgency), Some(10));
        assert_eq!(signal.performance.as_ref().map(|p| p.tempo), Some(Tempo::Presto));
    }

    #[test]
    fn test_caesura_broadcasts() {
        let signal = caesura("governance", "Budget exceeded").build();

        assert_eq!(signal.targets, vec!["*"]);
    }
}
