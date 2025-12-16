# HCT-MCP Signals (Rust)

<div align="center">

### The Coordination Layer for Model Context Protocol

**Express Urgency, Timing, and Synchronization in Multi-Agent Systems.**

[![crates.io](https://img.shields.io/crates/v/hct-mcp-signals?style=for-the-badge&logo=rust&logoColor=white)](https://crates.io/crates/hct-mcp-signals)

</div>

---

## 💡 What is HCT?

**HCT-MCP Signals** extends the protocol with 7 musical primitives proven to coordinate complex ensembles without a central conductor.

| Signal | Meanings |
|---|---|
| **CUE** | Act Now / Task dispatch |
| **FERMATA** | Hold / Human-in-the-loop |
| **ATTACCA** | Immediate / Real-time |
| **VAMP** | Loop / Retry |
| **CAESURA** | Stop / Emergency |
| **TACET** | Silent / Sleep |
| **DOWNBEAT**| Sync / Barrier |

## 📦 Installation

```bash
cargo add hct-mcp-signals
```

## 🚀 Quick Start

```rust
use hct_mcp_signals::{cue, Tempo};

fn main() {
    // 1. Create a CUE signal using Builder pattern
    let signal = cue("orch", ["analyst"])
        .with_urgency(9)
        .with_tempo(Tempo::Presto)
        .build();

    // 2. Convert to MCP JSON
    let json = signal.to_mcp_json().unwrap();
    println!("{}", json);
}
```

## 🔗 Links

- [Main Repository / Documentation](https://github.com/stefanwiest/hct-mcp-signals)
- [Protocol RFC](https://github.com/stefanwiest/hct-mcp-signals/blob/main/RFC.md)

## License

Apache-2.0
