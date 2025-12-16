# HCT-MCP Signals

**HCT Coordination Signals Extension for Model Context Protocol**

Adds urgency, timing, and approval semantics to MCP-based multi-agent systems.

[![PyPI](https://img.shields.io/pypi/v/hct-mcp-signals)](https://pypi.org/project/hct-mcp-signals/)
[![npm](https://img.shields.io/npm/v/@hct/mcp-signals)](https://www.npmjs.com/package/@hct/mcp-signals)

## Overview

MCP provides great tool/resource connections, but lacks vocabulary for:
- **Urgency/Priority**: Process immediately vs batch for later
- **Timing Requirements**: Response time expectations
- **Approval Gates**: Hold until human approves
- **Synchronization**: Wait for all agents before proceeding

This extension adds 7 coordination signals from [Harmonic Coordination Theory](https://github.com/stefanwiest/hct-paper).

## Installation

```bash
# Python
pip install hct-mcp-signals

# npm
npm install @hct/mcp-signals
```

## Usage (Python)

```python
from hct_mcp_signals import HCTSignal, SignalType, Performance

# Create a CUE signal with urgency
signal = HCTSignal(
    type=SignalType.CUE,
    source="orchestrator",
    targets=["analyst"],
    payload={"task": "Analyze Q4"},
    performance=Performance(urgency=8, tempo="allegro")
)

# Convert to MCP-compatible JSON
mcp_message = signal.to_mcp()
```

## Usage (TypeScript)

```typescript
import { createCue, Tempo } from '@hct/mcp-signals';

const signal = createCue({
  source: 'orchestrator',
  targets: ['analyst'],
  payload: { task: 'Analyze Q4' },
  performance: { urgency: 8, tempo: Tempo.Allegro }
});
```

## Signals

| Signal | Semantics | Use Case |
|--------|-----------|----------|
| `cue` | Trigger agent | Handoff, task dispatch |
| `fermata` | Hold for approval | Human-in-loop gates |
| `attacca` | Immediate transition | Urgent handoffs |
| `vamp` | Repeat until condition | Retry loops |
| `caesura` | Full stop | Emergency shutdown |
| `tacet` | Agent inactive | Resource conservation |
| `downbeat` | Global sync point | Checkpoints |

## RFC

See [RFC.md](RFC.md) for the full protocol specification.

## License

MIT
