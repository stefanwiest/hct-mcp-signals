/**
 * HCT Signal Type Definitions
 *
 * Core types for HCT coordination signals.
 */

/** Signal types corresponding to musical coordination signals */
export enum SignalType {
    /** Trigger agent activation */
    CUE = 'cue',
    /** Hold for approval */
    FERMATA = 'fermata',
    /** Immediate transition */
    ATTACCA = 'attacca',
    /** Repeat until condition met */
    VAMP = 'vamp',
    /** Full stop */
    CAESURA = 'caesura',
    /** Agent inactive */
    TACET = 'tacet',
    /** Global sync point */
    DOWNBEAT = 'downbeat',
}

/** Musical tempo indications mapped to urgency */
export enum Tempo {
    /** Very slow (~1 min response OK) */
    LARGO = 'largo',
    /** Walking pace (~30s response) */
    ANDANTE = 'andante',
    /** Moderate (~15s response) */
    MODERATO = 'moderato',
    /** Fast (~5s response) */
    ALLEGRO = 'allegro',
    /** Very fast (~1s response) */
    PRESTO = 'presto',
}

/** Types of holds for FERMATA signals */
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
    timeout_ms: 30000,
};
