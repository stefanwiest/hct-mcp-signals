"""
HCT-MCP Signals: Coordination Signals Extension for MCP

Adds urgency, timing, and approval semantics to Model Context Protocol.
"""

from .schema import (
    SignalType,
    Tempo,
    HoldType,
    Performance,
    Conditions,
    HCTSignal,
)
from .factory import cue, fermata, attacca, vamp, caesura, tacet, downbeat

__version__ = "0.1.0"
__all__ = [
    "SignalType",
    "Tempo", 
    "HoldType",
    "Performance",
    "Conditions",
    "HCTSignal",
    "cue",
    "fermata",
    "attacca",
    "vamp",
    "caesura",
    "tacet",
    "downbeat",
]
