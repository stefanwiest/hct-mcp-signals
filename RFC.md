# RFC: HCT Coordination Signals for MCP

**Status**: Public Preview  
**Author**: Stefan Wiest  
**Created**: 2025-12-16  
**Target**: AAIF/MCP Working Group

**Implementations**:
- 🐍 **Python**: [`hct-mcp-signals`](https://pypi.org/project/hct-mcp-signals/)
- 📦 **Node.js**: [`@hct-mcp/signals`](https://www.npmjs.com/package/@hct-mcp/signals)
- 🦀 **Rust**: [`hct-mcp-signals`](https://crates.io/crates/hct-mcp-signals)
- 🐹 **Go**: [`github.com/stefanwiest/hct-mcp-signals/go`](https://pkg.go.dev/github.com/stefanwiest/hct-mcp-signals/go)

---

## Abstract

This RFC proposes extending the Model Context Protocol (MCP) with **coordination semantics** derived from Harmonic Coordination Theory (HCT). While MCP excels at tool/resource connections, it lacks standard vocabulary for expressing urgency, resource intensity, approval gates, and synchronization—concepts essential for robust, scalable multi-agent coordination.

---

## Problem Statement

MCP provides `tasks/send` and `tasks/sendSubscribe` for task dispatch, but offers no standard way to express:

1.  **Urgency/Priority**: "Process this immediately" vs "batch for later"
2.  **Resource Intensity (Dynamics)**: "Use cheap/fast model" vs "Use expensive Chain-of-Thought"
3.  **Approval Gates**: "Hold until human approves"
4.  **Synchronization**: "Wait for all agents before proceeding"
5.  **Quality Thresholds**: "Repeat until quality score > 0.9"

Without these semantics, agents must implement ad-hoc coordination logic in prompts or side-channels, leading to fragmented, non-interoperable systems.

---

## Proposal: HCT Signal Extension

We propose 7 coordination signals as an optional MCP extension (`hct_signal`), inspired by musical ensemble coordination but strictly defined for computational utility.

### Signal Definitions

| Signal | MCP Method | Semantics | Computational Utility |
|--------|------------|-----------|-----------------------|
| `cue` | `tasks/send` | Trigger agent activation | Standard task dispatch / Handoff |
| `fermata` | `tasks/sendSubscribe` | Hold for approval | Human-in-the-loop / Governance Gates |
| `attacca` | `tasks/send` | Immediate transition | Low-latency / Real-time handoff |
| `vamp` | `tasks/sendSubscribe` | Repeat until condition | Polling / Retry loops / Quality Gates |
| `caesura` | `tasks/cancel` | Full stop | Emergency shutdown / Circuit Breaker |
| `tacet` | N/A | Agent inactive | Resource conservation / Sleep |
| `downbeat` | `notifications/message` | Global sync point | Barrier synchronization |

### Performance Parameters (Metadata)

HCT signals carry "Performance" metadata that dictates *how* the task should be executed.

#### Dynamics: Resource Intensity Control
Musical dynamics mapped to **Token/Compute Budget**.

| Marking | Level | Resource Multiplier | Use Case |
|:---:|---|:---:|---|
| **pp** | Pianissimo | < 0.5x | **Cache/Lookup**: Use cheapest model, no reflection, retrieval only. |
| **p** | Piano | 0.8x | **Efficient**: Standard model, zero-shot. |
| **mf** | Mezzo-Forte | 1.0x | **Standard**: Recommended model, standard prompt chain. |
| **f** | Forte | 1.5x | **Deep**: Models with reasoning (e.g., o1), multi-shot. |
| **ff** | Fortissimo | > 2.0x | **Maximum Depth**: Ensemble voting, extensive Chain-of-Thought. |

#### Tempo: Execution Urgency
Musical tempo mapped to **Latency/Priority**.

| Marking | Level | Priority | Use Case |
|:---:|---|:---:|---|
| **Largo** | Slow | Batch/Low | Background tasks, deep analysis (no deadline). |
| **Moderato** | Moderate | Normal | Standard interactive tasks. |
| **Presto** | Fast | Real-time | Latency-critical flows, user-facing voice/chat. |

### JSON Schema Extension

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "HCT-MCP Signal Extension",
  "type": "object",
  "properties": {
    "hct_signal": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["cue", "fermata", "attacca", "vamp", "caesura", "tacet", "downbeat"]
        },
        "performance": {
          "type": "object",
          "properties": {
            "urgency": {"type": "integer", "minimum": 1, "maximum": 10},
            "tempo": {"type": "string", "enum": ["largo", "andante", "moderato", "allegro", "presto"]},
            "dynamics": {
                "type": "string",
                "enum": ["pp", "p", "mp", "mf", "f", "ff"],
                "description": "Resource intensity: pp=minimal (<0.5x), mf=standard (1.0x), ff=high depth (>1.5x)"
            },
            "timeout_ms": {"type": "integer"}
          }
        },
        "context": {
            "type": "object",
            "description": "Transport context for maintained state across handoffs",
            "properties": {
                "movement": { "type": "string" },
                "objectives": { "type": "array", "items": { "type": "string" } },
                "reference_frame": { "type": "object" },
                "prior_outputs": { "type": "array" }
            }
        },
        "conditions": {
          "type": "object",
          "properties": {
            "hold_type": {"type": "string", "enum": ["human", "governance", "resource", "quality"]},
            "repeat_until": {"type": "string"},
            "quality_threshold": {"type": "number", "minimum": 0, "maximum": 1}
          }
        }
      },
      "required": ["type"]
    }
  }
}
```

### Example: CUE with High Resource Intensity (ff)

```json
{
  "jsonrpc": "2.0",
  "method": "tasks/send",
  "params": {
    "id": "task-123",
    "message": {
      "role": "user",
      "content": "Analyze Q4 revenue anomalies"
    },
    "hct_signal": {
      "type": "cue",
      "performance": {
        "urgency": 8,
        "tempo": "allegro",
        "dynamics": "ff"
      },
      "context": {
          "movement": "financial_analysis",
          "objectives": ["identify_root_cause"]
      }
    }
  }
}
```

### Example: FERMATA with Human Approval

```json
{
  "jsonrpc": "2.0",
  "method": "tasks/sendSubscribe",
  "params": {
    "id": "task-456",
    "message": {
      "role": "user",
      "content": "Ready to deploy to production"
    },
    "hct_signal": {
      "type": "fermata",
      "conditions": {
        "hold_type": "human",
        "timeout_ms": 3600000
      }
    }
  }
}
```

---

## Common Patterns

HCT signals enable standardized implementation of common multi-agent design patterns. See the `hct-patterns` repository for full documentation.

### Human-in-the-Loop (Fermata)
When an agent requires external input or authorization, it emits a `FERMATA`. This explicitly signals a pause in the workflow until a resolution (Downbeat) is provided.
- **A2A Map**: `tasks/sendSubscribe`
- **Signal**: `type="fermata"`, `reason="..."`

### Status Streaming (Ostinato)
Long-running agents emit periodic low-priority `CUE` signals to indicate liveness and progress.
- **A2A Map**: Repeating `tasks/send` or stream chunks
- **Signal**: `type="cue"`, `performance.dynamics="pp"`

## Compatibility

- **Backward Compatible**: The `hct_signal` field is optional; existing MCP clients/servers ignore it.
- **Graceful Degradation**: Servers that don't support HCT signals process the base message normally.
- **Discovery**: Servers can advertise HCT support via `capabilities.hct_signals: true`.

---

## Rationale

### Why Musical Metaphors? (Hybrid Coordination)

As AI agents move from pilot to production, teams become **hybrid ensembles** of humans and agents. Music provides a centuries-proven vocabulary for:
1.  **Distributed Sync**: Getting 100+ actors to align without centralized locking (`downbeat`).
2.  **Shared Semantics**: Humans intuitively understand "CUE" (start) and "FERMATA" (hold/pause), bridging the gap between technical protocols and human workflow.
3.  **Nuance**: "Dynamics" allows expressing the *quality* of effort (quick sketch vs. masterpiece) in a way that rigid "sim/med/hard" flags do not capture.

### Why Extend MCP?

MCP is establishing the standard transport layer for tools. By layering coordination semantics *on top of* MCP, we avoid protocol fragmentation and allow coordination patterns to ride on existing infrastructure.

---

## References

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [HCT Paper](https://github.com/stefanwiest/hct-paper)
