"""
HCT Signal Types - Auto-generated from hct-spec/spec.yaml
Generated: 2025-12-18T10:37:34.791657Z

DO NOT EDIT MANUALLY - Edit spec.yaml and regenerate.
"""

from enum import Enum


class SignalType(str, Enum):
    """HCT coordination signal types."""
    CUE = "cue"  # Trigger activation
    FERMATA = "fermata"  # Hold for approval
    ATTACCA = "attacca"  # Immediate transition
    VAMP = "vamp"  # Repeat until condition
    CAESURA = "caesura"  # Full stop
    TACET = "tacet"  # Agent inactive
    DOWNBEAT = "downbeat"  # Global sync point


class Tempo(str, Enum):
    """Musical tempo indications."""
    LARGO = "largo"  # Very slow (~60s)
    ANDANTE = "andante"  # Walking pace (~30s)
    MODERATO = "moderato"  # Moderate (~15s)
    ALLEGRO = "allegro"  # Fast (~5s)
    PRESTO = "presto"  # Very fast (~1s)


class DynamicsLevel(str, Enum):
    """Resource intensity levels."""
    PP = "pp"  # pianissimo
    P = "p"  # piano
    MF = "mf"  # mezzo-forte
    F = "f"  # forte
    FF = "ff"  # fortissimo
