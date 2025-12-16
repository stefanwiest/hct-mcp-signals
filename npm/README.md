# HCT-MCP Signals (TypeScript/Node.js)

<div align="center">

### The Coordination Layer for Model Context Protocol

**Express Urgency, Timing, and Synchronization in Multi-Agent Systems.**

[![npm](https://img.shields.io/npm/v/@hct-mcp/signals?style=for-the-badge&logo=npm&logoColor=white)](https://www.npmjs.com/package/@hct-mcp/signals)

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
npm install @hct-mcp/signals
```

## 🚀 Quick Start

```typescript
import { cue, Tempo, embedSignal } from '@hct-mcp/signals';

// 1. Create a CUE signal
const signal = cue({
    source: 'orch',
    targets: ['analyst'],
    urgency: 9,
    tempo: Tempo.PRESTO
});

// 2. Embed into MCP Tool Call Params
let params = { id: 'task-123' };
const paramsWithSignal = embedSignal(params, signal);
```

## 🔗 Links

- [Main Repository / Documentation](https://github.com/stefanwiest/hct-mcp-signals)
- [Protocol RFC](https://github.com/stefanwiest/hct-mcp-signals/blob/main/RFC.md)

## License

Apache-2.0
