# HCT-MCP Signals

**Coordination Signals Extension for Model Context Protocol**

Adds urgency, timing, and approval semantics to MCP-based multi-agent systems.

[![PyPI](https://img.shields.io/pypi/v/hct-mcp-signals)](https://pypi.org/project/hct-mcp-signals/)
[![npm](https://img.shields.io/npm/v/@hct-mcp/signals)](https://www.npmjs.com/package/@hct-mcp/signals)
[![crates.io](https://img.shields.io/crates/v/hct-mcp-signals)](https://crates.io/crates/hct-mcp-signals)
[![Go Reference](https://pkg.go.dev/badge/github.com/stefanwiest/hct-mcp-signals/go.svg)](https://pkg.go.dev/github.com/stefanwiest/hct-mcp-signals/go)

## Overview

MCP provides great tool/resource connections but lacks vocabulary for:
- **Urgency/Priority**: Process immediately vs batch
- **Timing Requirements**: Response time expectations
- **Approval Gates**: Hold until human approves
- **Synchronization**: Wait for all agents before proceeding

This extension adds 7 coordination signals inspired by [Harmonic Coordination Theory](https://github.com/stefanwiest/hct-paper).

## Installation

```bash
# Python
pip install hct-mcp-signals

# npm/TypeScript
npm install @hct-mcp/signals

# Rust
cargo add hct-mcp-signals

# Go
go get github.com/stefanwiest/hct-mcp-signals/go
```

## Quick Start

### Python

```python
from hct_mcp_signals import cue, fermata, attacca

# Trigger an agent with urgency
signal = cue("orchestrator", ["analyst"], {"task": "Analyze Q4"}, urgency=8)
mcp_json = signal.to_mcp()

# Hold for human approval
approval = fermata("reporter", "Ready for review", hold_type="human")

# Urgent immediate transition
urgent = attacca("coordinator", ["executor"])
```

### TypeScript

```typescript
import { cue, fermata, Tempo } from '@hct/mcp-signals';

const signal = cue({
  source: 'orchestrator',
  targets: ['analyst'],
  payload: { task: 'Analyze Q4' },
  urgency: 8,
  tempo: Tempo.Allegro
});
```

### Rust

```rust
use hct_mcp_signals::{cue, fermata, Tempo};

let signal = cue("orchestrator", ["analyst"])
    .with_urgency(8)
    .with_tempo(Tempo::Allegro)
    .with_payload_entry("task", json!("Analyze Q4"))
    .build();

let json = signal.to_mcp_json()?;
```

### Go

```go
import hct "github.com/stefanwiest/hct-mcp-signals/go"

signal := hct.NewCue("orchestrator", []string{"analyst"},
    hct.WithUrgency(8),
    hct.WithTempo(hct.Allegro),
    hct.WithPayloadEntry("task", "Analyze Q4"),
)
json, _ := signal.ToMCPJSON()
```

## The 7 Signals

| Signal | Semantics | Use Case |
|--------|-----------|----------|
| `cue` | Trigger agent activation | Task dispatch, handoffs |
| `fermata` | Hold for approval | Human-in-loop gates |
| `attacca` | Immediate transition | Urgent handoffs |
| `vamp` | Repeat until condition | Retry loops, quality gates |
| `caesura` | Full stop | Emergency shutdown |
| `tacet` | Agent inactive | Resource conservation |
| `downbeat` | Global sync point | Checkpoints, barriers |

## Framework Integrations

### LangGraph

```python
from langgraph.graph import StateGraph
from hct_mcp_signals import cue, fermata, caesura

def analyst_node(state):
    # Signal completion with approval request
    signal = fermata("analyst", "Analysis complete, ready for review")
    return {"signal": signal.to_mcp(), "result": state["result"]}

def router(state):
    signal = state.get("signal", {}).get("hct_signal", {})
    if signal.get("type") == "fermata":
        return "human_review"
    elif signal.get("type") == "caesura":
        return "end"
    return "continue"

graph = StateGraph()
graph.add_node("analyst", analyst_node)
graph.add_conditional_edges("analyst", router)
```

### CrewAI

```python
from crewai import Agent, Task, Crew
from hct_mcp_signals import cue, vamp

# Embed HCT signals in agent communication
class CoordinatedAgent(Agent):
    def execute_task(self, task):
        result = super().execute_task(task)
        
        # Signal quality gate
        if result.confidence < 0.9:
            signal = vamp("verifier", "confidence >= 0.9", quality_threshold=0.9)
            return {"result": result, "signal": signal.to_mcp()}
        
        return {"result": result, "signal": cue("analyst", ["synthesizer"]).to_mcp()}
```

### AutoGen

```python
from autogen import AssistantAgent, UserProxyAgent
from hct_mcp_signals import cue, fermata, downbeat

# Add HCT signals to AutoGen messages
def signal_aware_message(content, signal=None):
    msg = {"content": content}
    if signal:
        msg["hct_signal"] = signal.model_dump()
    return msg

# Sync point before parallel work
sync = downbeat("coordinator", "phase_2_start")
assistant.send(signal_aware_message("Starting phase 2", sync))
```

### Google ADK (Agent Development Kit)

```python
from google.adk import Agent
from hct_mcp_signals import cue, caesura

class CoordinatedAgent(Agent):
    async def on_message(self, message):
        signal = message.get("hct_signal")
        
        if signal and signal["type"] == "caesura":
            await self.emergency_shutdown(signal["payload"]["reason"])
            return
        
        result = await self.process(message)
        
        # Signal next agent
        next_signal = cue(self.name, self.downstream_agents)
        await self.emit({"result": result, "hct_signal": next_signal.model_dump()})
```

### AWS Strands / Bedrock Agents

```python
from strands import Agent
from hct_mcp_signals import cue, fermata

# Strands uses MCP natively - HCT extends it
class StrandsCoordinatedAgent(Agent):
    def handle_mcp_message(self, params):
        # Extract HCT signal from MCP params
        signal = params.get("hct_signal")
        
        if signal:
            urgency = signal.get("performance", {}).get("urgency", 5)
            if urgency >= 8:
                return self.priority_process(params)
        
        return self.normal_process(params)
```

### DSPy

```python
import dspy
from hct_mcp_signals import vamp, fermata

class QualityControlledModule(dspy.Module):
    def __init__(self):
        self.generate = dspy.ChainOfThought("question -> answer")
    
    def forward(self, question):
        # Use VAMP signal concept for retry logic
        max_attempts = 3
        for attempt in range(max_attempts):
            result = self.generate(question=question)
            
            if self.quality_check(result):
                return result
            
            # Signal: vamp until quality threshold
            signal = vamp("dspy_module", "quality >= 0.9", quality_threshold=0.9)
            # Log or emit signal for observability
        
        # Give up, request human review
        return fermata("dspy_module", "Quality threshold not met after retries")
```

### TensorZero

```python
from tensorzero import Gateway
from hct_mcp_signals import cue, fermata, downbeat

# Use HCT signals for routing decisions
gateway = Gateway()

@gateway.route
def route_with_urgency(request):
    signal = request.get("hct_signal", {})
    urgency = signal.get("performance", {}).get("urgency", 5)
    
    if urgency >= 9:
        return "fast_model"  # Route to faster model
    elif signal.get("type") == "fermata":
        return "careful_model"  # Route to more capable model
    return "default_model"
```

### Letta (MemGPT)

```python
from letta import Agent
from hct_mcp_signals import cue, tacet, downbeat

# Coordinate memory operations with HCT signals
class MemoryCoordinatedAgent(Agent):
    def coordinate_memory_sync(self, agents):
        # Signal all agents to sync memory
        sync_signal = downbeat(self.name, "memory_checkpoint")
        
        for agent in agents:
            agent.receive(sync_signal.to_mcp())
        
        # Wait for all to acknowledge
        self.wait_for_barrier("memory_checkpoint")
    
    def hibernate(self, duration_ms):
        # Signal agent is going inactive
        signal = tacet(self.name, duration_ms=duration_ms)
        self.broadcast(signal.to_mcp())
        self.sleep(duration_ms)
```

## MCP Integration

HCT signals embed as an extension field in MCP messages:

```json
{
  "jsonrpc": "2.0",
  "method": "tasks/send",
  "params": {
    "id": "task-123",
    "message": {"role": "user", "content": "Analyze Q4 earnings"},
    "hct_signal": {
      "type": "cue",
      "source": "orchestrator",
      "targets": ["analyst"],
      "payload": {"priority": "high"},
      "performance": {"urgency": 8, "tempo": "allegro"},
      "timestamp": "2025-12-16T12:00:00Z"
    }
  }
}
```

## RFC

See [RFC.md](RFC.md) for the full protocol specification.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

## License

MIT

