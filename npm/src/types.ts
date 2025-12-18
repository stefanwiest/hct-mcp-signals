/**
 * HCT Signal Type Definitions
 *
 * Protocol types (SignalType, Tempo, DynamicsLevel) are imported from spec.ts,
 * which is auto-generated from hct-spec/spec.yaml.
 *
 * Implementation types (HoldType, Performance, Conditions, HCTSignal) are
 * defined here as they are specific to the MCP extension implementation.
 */

// Protocol types - auto-generated from hct-spec/spec.yaml
export { SignalType, Tempo, DynamicsLevel } from './spec';
import { SignalType, Tempo, DynamicsLevel } from './spec';

/** Types of holds for FERMATA signals (implementation-specific) */
export enum HoldType {
    /** Requires human approval */
    HUMAN = 'human',
    /** Requires governance check */
    GOVERNANCE = 'governance',
    /** Waiting for resource */
    RESOURCE = 'resource',
    /** Quality threshold not met */
    QUALITY = 'quality',
}

/** Performance parameters (Layer 3 in HCT) */
export interface Performance {
    /** Urgency level 1-10 */
    urgency?: number;
    /** Expected response timing */
    tempo?: Tempo;
    /** Resource intensity */
    dynamics?: DynamicsLevel;
    /** Timeout in milliseconds */
    timeout_ms?: number;
}

/** Conditions for conditional signals */
export interface Conditions {
    /** Type of hold for FERMATA */
    hold_type?: HoldType;
    /** Condition expression for VAMP */
    repeat_until?: string;
    /** Quality score threshold 0-1 */
    quality_threshold?: number;
}

/** Complete HCT Signal */
export interface HCTSignal {
    /** Signal type */
    type: SignalType;
    /** Source agent ID */
    source: string;
    /** Target agent IDs */
    targets: string[];
    /** Signal payload */
    payload: Record<string, unknown>;
    /** Performance parameters */
    performance?: Performance;
    /** Conditional parameters */
    conditions?: Conditions;
    /** ISO timestamp */
    timestamp?: string;
}

/** Default performance values */
export const DEFAULT_PERFORMANCE: Required<Performance> = {
    urgency: 5,
    tempo: Tempo.MODERATO,
    dynamics: DynamicsLevel.MF,
    timeout_ms: 30000,
};
