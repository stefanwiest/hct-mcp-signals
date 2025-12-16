"""
Comprehensive tests for HCT-MCP Signals package.
"""

import pytest

from hct_mcp_signals import (
    Conditions,
    HCTSignal,
    HoldType,
    MCPSignalExtension,
    Performance,
    SignalType,
    Tempo,
    attacca,
    caesura,
    create_mcp_task_send,
    cue,
    downbeat,
    fermata,
    get_json_schema,
    tacet,
    vamp,
)


class TestSignalTypes:
    """Test signal type enums."""

    def test_all_signal_types_exist(self):
        """All 7 HCT signals should be defined."""
        assert SignalType.CUE == "cue"
        assert SignalType.FERMATA == "fermata"
        assert SignalType.ATTACCA == "attacca"
        assert SignalType.VAMP == "vamp"
        assert SignalType.CAESURA == "caesura"
        assert SignalType.TACET == "tacet"
        assert SignalType.DOWNBEAT == "downbeat"

    def test_tempo_values(self):
        """Tempo enum should have 5 levels."""
        assert len(Tempo) == 5
        assert Tempo.LARGO.value == "largo"
        assert Tempo.PRESTO.value == "presto"

    def test_hold_types(self):
        """HoldType should have 4 categories."""
        assert HoldType.HUMAN == "human"
        assert HoldType.GOVERNANCE == "governance"
        assert HoldType.RESOURCE == "resource"
        assert HoldType.QUALITY == "quality"


class TestPerformance:
    """Test Performance parameters."""

    def test_default_performance(self):
        """Default performance should have sensible defaults."""
        perf = Performance()
        assert perf.urgency == 5
        assert perf.tempo == Tempo.MODERATO
        assert perf.timeout_ms is None

    def test_urgency_bounds(self):
        """Urgency should be bounded 1-10."""
        with pytest.raises(ValueError):
            Performance(urgency=0)
        with pytest.raises(ValueError):
            Performance(urgency=11)

    def test_custom_performance(self):
        """Custom performance values should work."""
        perf = Performance(urgency=8, tempo=Tempo.ALLEGRO, timeout_ms=30000)
        assert perf.urgency == 8
        assert perf.tempo == Tempo.ALLEGRO
        assert perf.timeout_ms == 30000


class TestHCTSignal:
    """Test HCTSignal model."""

    def test_minimal_signal(self):
        """Signal with only required fields."""
        sig = HCTSignal(type=SignalType.CUE, source="test")
        assert sig.type == SignalType.CUE
        assert sig.source == "test"
        assert sig.targets == []
        assert sig.payload == {}

    def test_full_signal(self):
        """Signal with all fields."""
        sig = HCTSignal(
            type=SignalType.FERMATA,
            source="verifier",
            targets=["human"],
            payload={"reason": "Need approval"},
            performance=Performance(urgency=9),
            conditions=Conditions(hold_type=HoldType.HUMAN),
        )
        assert sig.type == SignalType.FERMATA
        assert "human" in sig.targets
        assert sig.conditions.hold_type == HoldType.HUMAN

    def test_to_mcp(self):
        """Signal should serialize to MCP format."""
        sig = HCTSignal(type=SignalType.CUE, source="orchestrator", targets=["analyst"])
        mcp = sig.to_mcp()
        assert "hct_signal" in mcp
        assert mcp["hct_signal"]["type"] == "cue"
        assert mcp["hct_signal"]["source"] == "orchestrator"

    def test_from_mcp(self):
        """Signal should deserialize from MCP format."""
        mcp_data = {
            "hct_signal": {
                "type": "fermata",
                "source": "reporter",
                "targets": ["governance"],
                "payload": {"report_id": "123"},
                "performance": {"urgency": 7, "tempo": "allegro"},
            }
        }
        sig = HCTSignal.from_mcp(mcp_data)
        assert sig.type == SignalType.FERMATA
        assert sig.source == "reporter"
        assert sig.performance.urgency == 7


class TestFactory:
    """Test signal factory functions."""

    def test_cue(self):
        """cue() should create CUE signal."""
        sig = cue("orch", ["analyst"], {"task": "analyze"}, urgency=8)
        assert sig.type == SignalType.CUE
        assert sig.source == "orch"
        assert "analyst" in sig.targets
        assert sig.performance.urgency == 8

    def test_fermata(self):
        """fermata() should create FERMATA signal."""
        sig = fermata("reporter", "Ready for review", hold_type="human")
        assert sig.type == SignalType.FERMATA
        assert sig.conditions.hold_type == HoldType.HUMAN

    def test_attacca(self):
        """attacca() should create urgent immediate signal."""
        sig = attacca("source", ["target"])
        assert sig.type == SignalType.ATTACCA
        assert sig.performance.urgency == 10
        assert sig.performance.tempo == Tempo.PRESTO

    def test_vamp(self):
        """vamp() should create repeat-until signal."""
        sig = vamp("agent", "score > 0.9", quality_threshold=0.9)
        assert sig.type == SignalType.VAMP
        assert sig.conditions.quality_threshold == 0.9

    def test_caesura(self):
        """caesura() should create full-stop signal."""
        sig = caesura("governance", "Budget exceeded")
        assert sig.type == SignalType.CAESURA
        assert sig.targets == ["*"]

    def test_tacet(self):
        """tacet() should create inactive signal."""
        sig = tacet("agent", duration_ms=60000)
        assert sig.type == SignalType.TACET
        assert sig.payload["duration_ms"] == 60000

    def test_downbeat(self):
        """downbeat() should create sync signal."""
        sig = downbeat("conductor", "daily_sync")
        assert sig.type == SignalType.DOWNBEAT
        assert sig.payload["sync_point"] == "daily_sync"


class TestMCPIntegration:
    """Test MCP integration utilities."""

    def test_embed_signal(self):
        """MCPSignalExtension should embed signals."""
        ext = MCPSignalExtension()
        sig = cue("orch", ["analyst"])
        params = {"id": "task-123", "message": {"content": "Hello"}}

        result = ext.embed_signal(params, sig)

        assert "hct_signal" in result
        assert result["id"] == "task-123"  # Original preserved

    def test_extract_signal(self):
        """MCPSignalExtension should extract signals."""
        ext = MCPSignalExtension()
        msg = {
            "id": "task-123",
            "hct_signal": {
                "type": "cue",
                "source": "orch",
                "targets": [],
                "payload": {},
                "performance": {"urgency": 5, "tempo": "moderato"},
            },
        }

        sig = ext.extract_signal(msg)

        assert sig is not None
        assert sig.type == SignalType.CUE

    def test_has_signal(self):
        """has_signal should detect presence."""
        ext = MCPSignalExtension()
        assert ext.has_signal({"hct_signal": {}}) is True
        assert ext.has_signal({"other": "data"}) is False

    def test_create_mcp_task_send(self):
        """create_mcp_task_send should build complete message."""
        sig = cue("orch", ["analyst"], urgency=8)
        msg = create_mcp_task_send("task-123", "Analyze Q4", sig)

        assert msg["jsonrpc"] == "2.0"
        assert msg["method"] == "tasks/send"
        assert msg["params"]["id"] == "task-123"
        assert "hct_signal" in msg["params"]


class TestJSONSchema:
    """Test JSON Schema export."""

    def test_get_json_schema(self):
        """get_json_schema should return valid schema."""
        schema = get_json_schema()
        assert "$defs" in schema or "properties" in schema
        assert schema.get("title") or schema.get("$defs")
