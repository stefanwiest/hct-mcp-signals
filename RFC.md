# RFC: HCT Coordination Signals for MCP

**Status**: Draft  
**Author**: Stefan Wiest  
**Created**: 2025-12-16  
**Target**: AAIF/MCP Working Group

---

## Abstract

This RFC proposes extending the Model Context Protocol (MCP) with **coordination semantics** derived from Harmonic Coordination Theory (HCT). While MCP excels at tool/resource connections, it lacks vocabulary for expressing urgency, timing, approval gates, and synchronization—concepts essential for robust multi-agent coordination.

---

## Problem Statement

MCP provides `tasks/send` and `tasks/sendSubscribe` for task dispatch, but offers no standard way to express:

1. **Urgency/Priority**: "Process this immediately" vs "batch for later"
2. **Timing Requirements**: Expected response time, deadline awareness
3. **Approval Gates**: "Hold until human approves"
4. **Synchronization**: "Wait for all agents before proceeding"
5. **Quality Thresholds**: "Repeat until quality score > 0.9"

Without these semantics, agents must implement ad-hoc coordination, leading to fragmented, non-interoperable systems.

---

## Proposal: HCT Signal Extension

We propose 7 coordination signals as an optional MCP extension, inspired by musical ensemble coordination:

### Signal Definitions

| Signal | MCP Method | Semantics | Use Case |
|--------|------------|-----------|----------|
| `cue` | `tasks/send` | Trigger agent activation | Handoff, task dispatch |
| `fermata` | `tasks/sendSubscribe` | Hold for approval | Human-in-loop gates |
| `attacca` | `tasks/send` | Immediate transition | Urgent handoffs |
| `vamp` | `tasks/sendSubscribe` | Repeat until condition | Retry loops, polling |
| `caesura` | `tasks/cancel` | Full stop | Emergency shutdown |
| `tacet` | N/A | Agent inactive | Resource conservation |
| `downbeat` | `notifications/message` | Global sync point | Coordination checkpoints |

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
            "timeout_ms": {"type": "integer"}
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

### Example: CUE with Urgency

```json
{
  "jsonrpc": "2.0",
  "method": "tasks/send",
  "params": {
    "id": "task-123",
    "message": {
      "role": "user",
      "content": "Analyze Q4 revenue"
    },
    "hct_signal": {
      "type": "cue",
      "performance": {
        "urgency": 8,
        "tempo": "allegro",
        "timeout_ms": 30000
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

## Compatibility

- **Backward Compatible**: The `hct_signal` field is optional; existing MCP clients/servers ignore it
- **Graceful Degradation**: Servers that don't support HCT signals process the base message normally
- **Discovery**: Servers can advertise HCT support via `capabilities.hct_signals: true`

---

## Reference Implementation

Available at: https://github.com/stefanwiest/hct-core

```python
from hct.coordination import cue, fermata
from hct.mcp import HCTSignalClient

async with HCTSignalClient() as client:
    # Send CUE with urgency
    await client.emit(
        signal_type="cue",
        source="orchestrator",
        targets=["analyst"],
        payload={"task": "Analyze Q4"},
        performance={"urgency": 8, "tempo": "allegro"}
    )
    
    # Create FERMATA hold
    await client.emit(
        signal_type="fermata",
        source="synthesizer",
        targets=["human"],
        payload={"reason": "Report requires approval"}
    )
```

---

## Rationale

### Why Musical Metaphors?

Music represents humanity's most sophisticated solution to distributed coordination:
- **Orchestras** coordinate 100+ musicians without centralized state
- **Jazz ensembles** adapt in real-time to unexpected inputs
- **Chamber groups** self-organize without conductors

HCT borrows vocabulary proven over centuries of ensemble practice.

### Why Extend MCP?

MCP is rapidly becoming the standard for agent-tool connections. By extending MCP rather than creating a competing protocol, HCT signals can be adopted incrementally without requiring ecosystem changes.

---

## Next Steps

1. Community feedback on signal definitions
2. Implementation in major MCP libraries (Python, TypeScript)
3. Integration with A2A protocol for end-to-end coordination

---

## References

- [MCP Specification](https://spec.modelcontextprotocol.io/)
- [HCT Paper](https://github.com/stefanwiest/hct-paper)
- [A2A Protocol](https://github.com/google/A2A)
