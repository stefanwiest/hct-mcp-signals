/**
 * Signal Factory Functions
 *
 * Convenience functions for creating HCT signals.
 */

import {
    HCTSignal,
    SignalType,
    Tempo,
    HoldType,
} from './types';

/** Options for creating a CUE signal */
export interface CueOptions {
    source: string;
    targets: string[];
    payload?: Record<string, unknown>;
    urgency?: number;
    tempo?: Tempo;
}

/** Create a CUE signal to trigger agent activation */
export function cue(options: CueOptions): HCTSignal {
    return {
        type: SignalType.CUE,
        source: options.source,
        targets: options.targets,
        payload: options.payload ?? {},
        performance: {
            urgency: options.urgency ?? 5,
            tempo: options.tempo ?? Tempo.MODERATO,
        },
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating a FERMATA signal */
export interface FermataOptions {
    source: string;
    reason: string;
    holdType?: HoldType;
    timeout_ms?: number;
}

/** Create a FERMATA signal to hold for approval */
export function fermata(options: FermataOptions): HCTSignal {
    return {
        type: SignalType.FERMATA,
        source: options.source,
        targets: ['governance'],
        payload: { reason: options.reason },
        conditions: { hold_type: options.holdType ?? HoldType.HUMAN },
        performance: { timeout_ms: options.timeout_ms },
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating an ATTACCA signal */
export interface AttaccaOptions {
    source: string;
    targets: string[];
    payload?: Record<string, unknown>;
}

/** Create an ATTACCA signal for immediate transition */
export function attacca(options: AttaccaOptions): HCTSignal {
    return {
        type: SignalType.ATTACCA,
        source: options.source,
        targets: options.targets,
        payload: options.payload ?? {},
        performance: { urgency: 10, tempo: Tempo.PRESTO },
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating a VAMP signal */
export interface VampOptions {
    source: string;
    repeatUntil: string;
    qualityThreshold?: number;
    timeout_ms?: number;
}

/** Create a VAMP signal to repeat until condition met */
export function vamp(options: VampOptions): HCTSignal {
    return {
        type: SignalType.VAMP,
        source: options.source,
        targets: [options.source], // Self-loop
        payload: {},
        conditions: {
            repeat_until: options.repeatUntil,
            quality_threshold: options.qualityThreshold ?? 0.9,
        },
        performance: { timeout_ms: options.timeout_ms ?? 60000 },
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating a CAESURA signal */
export interface CaesuraOptions {
    source: string;
    reason: string;
}

/** Create a CAESURA signal for full stop */
export function caesura(options: CaesuraOptions): HCTSignal {
    return {
        type: SignalType.CAESURA,
        source: options.source,
        targets: ['*'], // Broadcast
        payload: { reason: options.reason },
        performance: { urgency: 10, tempo: Tempo.PRESTO },
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating a TACET signal */
export interface TacetOptions {
    source: string;
    duration_ms?: number;
}

/** Create a TACET signal to mark agent as inactive */
export function tacet(options: TacetOptions): HCTSignal {
    return {
        type: SignalType.TACET,
        source: options.source,
        targets: [],
        payload: options.duration_ms ? { duration_ms: options.duration_ms } : {},
        timestamp: new Date().toISOString(),
    };
}

/** Options for creating a DOWNBEAT signal */
export interface DownbeatOptions {
    source: string;
    syncPoint: string;
}

/** Create a DOWNBEAT signal for global synchronization */
export function downbeat(options: DownbeatOptions): HCTSignal {
    return {
        type: SignalType.DOWNBEAT,
        source: options.source,
        targets: ['*'], // Broadcast
        payload: { sync_point: options.syncPoint },
        timestamp: new Date().toISOString(),
    };
}
