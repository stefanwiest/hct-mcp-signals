# Rust Package: hct-mcp-signals

**Status**: Planned
**Estimated Effort**: 2-3 days

## Overview

Rust implementation of HCT-MCP signals for high-performance agent systems.

## Package Structure

```
rust/
├── Cargo.toml
├── src/
│   ├── lib.rs          # Main entry point
│   ├── signals.rs      # Signal types (enum)
│   ├── schema.rs       # Serde types
│   ├── factory.rs      # Builder functions
│   └── mcp.rs          # MCP integration
└── tests/
    └── signals_test.rs
```

## Key Dependencies

```toml
[dependencies]
serde = { version = "1.0", features = ["derive"] }
serde_json = "1.0"
chrono = { version = "0.4", features = ["serde"] }
thiserror = "1.0"

[dev-dependencies]
tokio = { version = "1", features = ["full"] }
```

## API Design

```rust
use hct_mcp_signals::{SignalType, HCTSignal, Performance, Tempo};

// Create a CUE signal
let signal = HCTSignal::cue("orchestrator", vec!["analyst"])
    .with_payload(json!({"task": "Analyze Q4"}))
    .with_urgency(8)
    .with_tempo(Tempo::Allegro)
    .build();

// Serialize to MCP format
let mcp_json = signal.to_mcp_json()?;

// Parse from MCP message
let parsed: HCTSignal = HCTSignal::from_mcp_json(&json_str)?;
```

## Type Definitions

```rust
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum SignalType {
    Cue,
    Fermata,
    Attacca,
    Vamp,
    Caesura,
    Tacet,
    Downbeat,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum Tempo {
    Largo,
    Andante,
    Moderato,
    Allegro,
    Presto,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Performance {
    pub urgency: u8,  // 1-10
    pub tempo: Tempo,
    pub timeout_ms: Option<u64>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct HCTSignal {
    #[serde(rename = "type")]
    pub signal_type: SignalType,
    pub source: String,
    pub targets: Vec<String>,
    pub payload: serde_json::Value,
    pub performance: Option<Performance>,
    pub conditions: Option<Conditions>,
    pub timestamp: chrono::DateTime<chrono::Utc>,
}
```

## Publishing

```bash
# Test
cargo test

# Build
cargo build --release

# Publish to crates.io
cargo login
cargo publish
```

## CI/CD Addition

Add to `.github/workflows/tests.yml`:

```yaml
rust:
  runs-on: ubuntu-latest
  defaults:
    run:
      working-directory: rust
  steps:
    - uses: actions/checkout@v4
    - uses: dtolnay/rust-toolchain@stable
    - run: cargo fmt --check
    - run: cargo clippy -- -D warnings
    - run: cargo test
```
