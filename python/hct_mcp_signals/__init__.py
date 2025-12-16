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
from .json_schema import (
    get_json_schema,
    get_json_schema_string,
    export_json_schema,
    HCT_SIGNAL_EXTENSION_SCHEMA,
)
from .mcp import (
    MCPSignalExtension,
    create_mcp_task_send,
    create_mcp_task_subscribe,
)

__version__ = "0.1.0"

__all__ = [
    # Schema
    "SignalType",
    "Tempo", 
    "HoldType",
    "Performance",
    "Conditions",
    "HCTSignal",
    # Factory
    "cue",
    "fermata",
    "attacca",
    "vamp",
    "caesura",
    "tacet",
    "downbeat",
    # JSON Schema
    "get_json_schema",
    "get_json_schema_string",
    "export_json_schema",
    "HCT_SIGNAL_EXTENSION_SCHEMA",
    # MCP Integration
    "MCPSignalExtension",
    "create_mcp_task_send",
    "create_mcp_task_subscribe",
]
