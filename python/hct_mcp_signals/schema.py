"""
HCT-MCP Signal Schema Definitions

Pydantic models for HCT coordination signals.

NOTE: These definitions are derived from the canonical specification at:
      https://github.com/stefanwiest/genesis/tree/main/hct-spec/spec.yaml
"""

import json
from datetime import datetime
from enum import Enum
from typing import Any, Optional

from pydantic import BaseModel, Field


class SignalType(str, Enum):
    """HCT coordination signal types."""

    CUE = "cue"  # Trigger agent activation
    FERMATA = "fermata"  # Hold for approval
    ATTACCA = "attacca"  # Immediate transition
    VAMP = "vamp"  # Repeat until condition
    CAESURA = "caesura"  # Full stop
    TACET = "tacet"  # Agent inactive
    DOWNBEAT = "downbeat"  # Global sync point


class Tempo(str, Enum):
    """Musical tempo indications for urgency mapping."""

    LARGO = "largo"  # Very slow (~1 min response OK)
    ANDANTE = "andante"  # Walking pace (~30s response)
    MODERATO = "moderato"  # Moderate (~15s response)
    ALLEGRO = "allegro"  # Fast (~5s response)
    PRESTO = "presto"  # Very fast (~1s response)


class DynamicsLevel(str, Enum):
    """Resource intensity dynamic levels."""

    PP = "pp"  # Pianissimo - Low cost/cache ops (<0.5x budget)
    P = "p"    # Piano - Efficient/Zero-shot (0.8x budget)
    MP = "mp"  # Mezzo-piano - Light (0.9x budget)
    MF = "mf"  # Mezzo-forte - Standard (1.0x budget)
    F = "f"    # Forte - Deep/Multi-shot (1.5x budget)
    FF = "ff"  # Fortissimo - Maximum depth/CoT (>2.0x budget)


class HoldType(str, Enum):
    """Types of holds for FERMATA signals."""

    HUMAN = "human"  # Requires human approval
    GOVERNANCE = "governance"  # Requires governance check
    RESOURCE = "resource"  # Waiting for resource
    QUALITY = "quality"  # Quality threshold not met


class HCTContext(BaseModel):
    """
    Transport context for State Maintenance across handoffs.

    Ensures the "Movement" and "Objectives" persist across agent boundaries.
    """
    movement: Optional[str] = Field(default=None, description="Current movement name")
    objectives: list[str] = Field(default_factory=list, description="Current objectives")
    reference_frame: dict[str, Any] = Field(default_factory=dict, description="Shared reference frame")
    prior_outputs: list[dict[str, Any]] = Field(default_factory=list, description="History of outputs")

    model_config = {"use_enum_values": True}

class Performance(BaseModel):
    """Performance parameters (Layer 3 in HCT)."""

    urgency: int = Field(default=5, ge=1, le=10, description="1-10 urgency scale")
    tempo: Tempo = Field(default=Tempo.MODERATO, description="Expected response timing")
    dynamics: DynamicsLevel = Field(
        default=DynamicsLevel.MF, 
        description="Resource intensity (pp=low cost, ff=high depth)"
    )
    timeout_ms: Optional[int] = Field(
        default=None, description="Timeout in milliseconds"
    )

    model_config = {"use_enum_values": True}


class Conditions(BaseModel):
    """Conditions for conditional signals (FERMATA, VAMP)."""

    hold_type: Optional[HoldType] = Field(default=None, description="Type of hold")
    repeat_until: Optional[str] = Field(default=None, description="Condition for VAMP")
    quality_threshold: Optional[float] = Field(
        default=None, ge=0, le=1, description="Quality score threshold"
    )

    model_config = {"use_enum_values": True}


class HCTSignal(BaseModel):
    """
    Complete HCT Signal for MCP extension.

    Example:
        signal = HCTSignal(
            type=SignalType.CUE,
            source="orchestrator",
            targets=["analyst"],
            payload={"task": "Analyze Q4"},
            performance=Performance(urgency=8, tempo=Tempo.ALLEGRO, dynamics=DynamicsLevel.FF)
        )
        mcp_message = signal.to_mcp()
    """

    type: SignalType = Field(..., description="Signal type")
    source: str = Field(..., description="Source agent ID")
    targets: list[str] = Field(default_factory=list, description="Target agent IDs")
    payload: dict[str, Any] = Field(default_factory=dict, description="Signal payload")
    performance: Performance = Field(
        default_factory=Performance, description="Performance params"
    )
    context: Optional[HCTContext] = Field(
        default=None, description="Maintained state context"
    )
    conditions: Optional[Conditions] = Field(
        default=None, description="Conditional params"
    )
    timestamp: datetime = Field(
        default_factory=datetime.utcnow, description="Signal timestamp"
    )

    model_config = {"use_enum_values": True}

    def to_mcp(self) -> dict[str, Any]:
        """Convert to MCP-compatible JSON structure."""
        data = {
            "type": self.type,
            "source": self.source,
            "targets": self.targets,
            "payload": self.payload,
            "performance": (
                self.performance.model_dump(exclude_none=True) if self.performance else None
            ),
            "conditions": (
                self.conditions.model_dump(exclude_none=True) if self.conditions else None
            ),
            "timestamp": self.timestamp.isoformat(),
        }
        
        if self.context:
            data["context"] = self.context.model_dump(exclude_none=True)

        return {"hct_signal": data}

    def to_json(self) -> str:
        """Serialize to JSON string."""
        return json.dumps(self.to_mcp(), indent=2)

    @classmethod
    def from_mcp(cls, mcp_data: dict[str, Any]) -> "HCTSignal":
        """Parse from MCP message with hct_signal extension."""
        sig = mcp_data.get("hct_signal", mcp_data)
        
        # Performance
        perf = None
        if sig.get("performance"):
            perf = Performance(**sig["performance"])
            
        # Conditions
        cond = None
        if sig.get("conditions"):
            cond = Conditions(**sig["conditions"])
            
        # Context
        ctx = None
        if sig.get("context"):
            ctx = HCTContext(**sig["context"])

        return cls(
            type=SignalType(sig["type"]),
            source=sig["source"],
            targets=sig.get("targets", []),
            payload=sig.get("payload", {}),
            performance=perf or Performance(),
            conditions=cond,
            context=ctx,
        )

    def is_broadcast(self) -> bool:
        """Check if this is a broadcast signal."""
        return len(self.targets) == 0

    def is_for(self, agent_id: str) -> bool:
        """Check if this signal is for a specific agent."""
        return self.is_broadcast() or agent_id in self.targets
