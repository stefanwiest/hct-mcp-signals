/**
 * HCT-MCP Signals: Coordination Signals Extension for MCP
 * 
 * TypeScript implementation of HCT coordination signals.
 */

export enum SignalType {
    CUE = 'cue',
    FERMATA = 'fermata',
    ATTACCA = 'attacca',
    VAMP = 'vamp',
    CAESURA = 'caesura',
    TACET = 'tacet',
    DOWNBEAT = 'downbeat'
}

export enum Tempo {
    LARGO = 'largo',
    ANDANTE = 'andante',
    MODERATO = 'moderato',
    ALLEGRO = 'allegro',
    PRESTO = 'presto'
}

export enum HoldType {
    HUMAN = 'human',
    GOVERNANCE = 'governance',
    RESOURCE = 'resource',
    QUALITY = 'quality'
}

export interface Performance {
    urgency?: number;      // 1-10
    tempo?: Tempo;
    timeout_ms?: number;
}

export interface Conditions {
    hold_type?: HoldType;
    repeat_until?: string;
    quality_threshold?: number;
}

export interface HCTSignal {
    type: SignalType;
    source: string;
    targets: string[];
    payload: Record<string, any>;
    performance?: Performance;
    conditions?: Conditions;
    timestamp?: string;
}

// Factory functions
export function createCue(params: {
    source: string;
    targets: string[];
    payload?: Record<string, any>;
    performance?: Performance;
}): HCTSignal {
    return {
        type: SignalType.CUE,
        source: params.source,
        targets: params.targets,
        payload: params.payload || {},
        performance: params.performance || { urgency: 5, tempo: Tempo.MODERATO },
        timestamp: new Date().toISOString()
    };
}

export function createFermata(params: {
    source: string;
    reason: string;
    holdType?: HoldType;
    timeout_ms?: number;
}): HCTSignal {
    return {
        type: SignalType.FERMATA,
        source: params.source,
        targets: ['governance'],
        payload: { reason: params.reason },
        conditions: { hold_type: params.holdType || HoldType.HUMAN },
        performance: { timeout_ms: params.timeout_ms },
        timestamp: new Date().toISOString()
    };
}

export function createCaesura(params: {
    source: string;
    reason: string;
}): HCTSignal {
    return {
        type: SignalType.CAESURA,
        source: params.source,
        targets: ['*'],
        payload: { reason: params.reason },
        performance: { urgency: 10, tempo: Tempo.PRESTO },
        timestamp: new Date().toISOString()
    };
}

// Utility to embed in MCP message
export function toMCPExtension(signal: HCTSignal): { hct_signal: HCTSignal } {
    return { hct_signal: signal };
}
