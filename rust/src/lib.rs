//! HCT-MCP Signals: Coordination Signals Extension for MCP
//!
//! This crate provides Rust types and utilities for HCT coordination signals,
//! enabling urgency, timing, and approval semantics in MCP-based multi-agent systems.
//!
//! # Quick Start
//!
//! ```rust
//! use hct_mcp_signals::{cue, SignalType, Tempo};
//!
//! let signal = cue("orchestrator", vec!["analyst"])
//!     .with_urgency(8)
//!     .with_tempo(Tempo::Allegro)
//!     .build();
//!
//! let json = signal.to_mcp_json().unwrap();
//! ```
//!
//! # Signal Types
//!
//! - `Cue` - Trigger agent activation
//! - `Fermata` - Hold for approval
//! - `Attacca` - Immediate transition
//! - `Vamp` - Repeat until condition
//! - `Caesura` - Full stop
//! - `Tacet` - Agent inactive
//! - `Downbeat` - Global sync point

#![forbid(unsafe_code)]
#![warn(clippy::all, clippy::pedantic, clippy::nursery)]
#![allow(clippy::module_name_repetitions)]

mod factory;
mod mcp;
mod signals;

pub use factory::{attacca, caesura, cue, downbeat, fermata, tacet, vamp, SignalBuilder};
pub use mcp::{embed_signal, extract_signal, McpTaskSend};
pub use signals::{Conditions, HCTSignal, HoldType, Performance, SignalType, Tempo};
