"""
HCT-MCP Signal Schema Definitions

Pydantic models for HCT coordination signals.

Protocol types (SignalType, Tempo, DynamicsLevel) are imported from spec.py,
which is auto-generated from hct-spec/spec.yaml.

Implementation types (HoldType, Performance, Conditions, HCTSignal) are
defined here as they are specific to the MCP extension implementation.
"""

import json
from datetime import datetime, timezone
from enum import Enum
from typing import Any, Optional

from pydantic import BaseModel, Field

# Protocol types - auto-generated from hct-spec/spec.yaml
from .spec import DynamicsLevel, SignalType, Tempo


class HoldType(str, Enum):
    """
    Types of holds for FERMATA signals.
    
    Implementation-specific: not part of the HCT protocol spec.
    """

    HUMAN = "human"  # Requires human approval
    GOVERNANCE = "governance"  # Requires governance check
    RESOURCE = "resource"  # Waiting for resource
    QUALITY = "quality"  # Quality threshold not met


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
    conditions: Optional[Conditions] = Field(
        default=None, description="Conditional params"
    )
    timestamp: datetime = Field(
        default_factory=lambda: datetime.now(timezone.utc), description="Signal timestamp"
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

        return cls(
            type=SignalType(sig["type"]),
            source=sig["source"],
            targets=sig.get("targets", []),
            payload=sig.get("payload", {}),
            performance=perf or Performance(),
            conditions=cond,
        )

    def is_broadcast(self) -> bool:
        """Check if this is a broadcast signal."""
        return len(self.targets) == 0

    def is_for(self, agent_id: str) -> bool:
        """Check if this signal is for a specific agent."""
        return self.is_broadcast() or agent_id in self.targets
