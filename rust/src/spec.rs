//! HCT Signal Types - Auto-generated from hct-spec/spec.yaml
//! Generated: 2025-12-18T10:37:34.792098Z
//!
//! DO NOT EDIT MANUALLY - Edit spec.yaml and regenerate.

use serde::{Deserialize, Serialize};

/// The HCT coordination signal types.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum SignalType {
    /// Trigger activation
    Cue,
    /// Hold for approval
    Fermata,
    /// Immediate transition
    Attacca,
    /// Repeat until condition
    Vamp,
    /// Full stop
    Caesura,
    /// Agent inactive
    Tacet,
    /// Global sync point
    Downbeat,
}

/// Musical tempo indications.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize, Default)]
#[serde(rename_all = "lowercase")]
pub enum Tempo {
    /// Very slow (~60s)
    Largo,
    /// Walking pace (~30s)
    Andante,
    /// Moderate (~15s)
    #[default]
    Moderato,
    /// Fast (~5s)
    Allegro,
    /// Very fast (~1s)
    Presto,
}

/// Resource intensity levels.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash, Serialize, Deserialize, Default)]
#[serde(rename_all = "lowercase")]
pub enum DynamicsLevel {
    /// pianissimo
    PP,
    /// piano
    P,
    /// mezzo-forte
    #[default]
    MF,
    /// forte
    F,
    /// fortissimo
    FF,
}
