/**
 * Tests for HCT-MCP Signals npm package
 */

import {
    SignalType,
    Tempo,
    HoldType,
    HCTSignal,
    cue,
    fermata,
    attacca,
    vamp,
    caesura,
    tacet,
    downbeat,
    embedSignal,
    extractSignal,
    hasSignal,
    createMCPTaskSend,
} from '../src';

describe('SignalType', () => {
    it('should have all 7 signal types', () => {
        expect(Object.keys(SignalType)).toHaveLength(7);
        expect(SignalType.CUE).toBe('cue');
        expect(SignalType.FERMATA).toBe('fermata');
        expect(SignalType.CAESURA).toBe('caesura');
    });
});

describe('Factory Functions', () => {
    describe('cue', () => {
        it('should create a CUE signal', () => {
            const signal = cue({
                source: 'orchestrator',
                targets: ['analyst'],
                payload: { task: 'analyze' },
                urgency: 8,
            });

            expect(signal.type).toBe(SignalType.CUE);
            expect(signal.source).toBe('orchestrator');
            expect(signal.targets).toContain('analyst');
            expect(signal.performance?.urgency).toBe(8);
        });

        it('should use defaults for optional params', () => {
            const signal = cue({ source: 'test', targets: [] });
            expect(signal.performance?.urgency).toBe(5);
            expect(signal.performance?.tempo).toBe(Tempo.MODERATO);
        });
    });

    describe('fermata', () => {
        it('should create a FERMATA signal', () => {
            const signal = fermata({
                source: 'reporter',
                reason: 'Ready for review',
                holdType: HoldType.HUMAN,
            });

            expect(signal.type).toBe(SignalType.FERMATA);
            expect(signal.conditions?.hold_type).toBe(HoldType.HUMAN);
            expect(signal.payload.reason).toBe('Ready for review');
        });
    });

    describe('attacca', () => {
        it('should create urgent immediate signal', () => {
            const signal = attacca({ source: 'agent', targets: ['next'] });

            expect(signal.type).toBe(SignalType.ATTACCA);
            expect(signal.performance?.urgency).toBe(10);
            expect(signal.performance?.tempo).toBe(Tempo.PRESTO);
        });
    });

    describe('vamp', () => {
        it('should create repeat-until signal', () => {
            const signal = vamp({
                source: 'verifier',
                repeatUntil: 'score > 0.9',
                qualityThreshold: 0.95,
            });

            expect(signal.type).toBe(SignalType.VAMP);
            expect(signal.conditions?.repeat_until).toBe('score > 0.9');
            expect(signal.conditions?.quality_threshold).toBe(0.95);
        });
    });

    describe('caesura', () => {
        it('should create full-stop signal', () => {
            const signal = caesura({ source: 'governance', reason: 'Budget exceeded' });

            expect(signal.type).toBe(SignalType.CAESURA);
            expect(signal.targets).toContain('*');
        });
    });

    describe('tacet', () => {
        it('should create inactive signal', () => {
            const signal = tacet({ source: 'agent', duration_ms: 60000 });

            expect(signal.type).toBe(SignalType.TACET);
            expect(signal.payload.duration_ms).toBe(60000);
        });
    });

    describe('downbeat', () => {
        it('should create sync signal', () => {
            const signal = downbeat({ source: 'conductor', syncPoint: 'daily_sync' });

            expect(signal.type).toBe(SignalType.DOWNBEAT);
            expect(signal.payload.sync_point).toBe('daily_sync');
        });
    });
});

describe('MCP Integration', () => {
    describe('embedSignal', () => {
        it('should embed signal into params', () => {
            const signal = cue({ source: 'test', targets: [] });
            const params = { id: 'task-123' };

            const result = embedSignal(params, signal);

            expect(result.id).toBe('task-123');
            expect(result.hct_signal).toBeDefined();
        });
    });

    describe('extractSignal', () => {
        it('should extract signal from message', () => {
            const signal = cue({ source: 'test', targets: [] });
            const msg = { id: 'task-123', hct_signal: signal };

            const extracted = extractSignal(msg);

            expect(extracted?.type).toBe(SignalType.CUE);
        });

        it('should return undefined if no signal', () => {
            const extracted = extractSignal({ id: 'task-123' });
            expect(extracted).toBeUndefined();
        });
    });

    describe('hasSignal', () => {
        it('should detect signal presence', () => {
            expect(hasSignal({ hct_signal: {} })).toBe(true);
            expect(hasSignal({ other: 'data' })).toBe(false);
        });
    });

    describe('createMCPTaskSend', () => {
        it('should create complete MCP message', () => {
            const signal = cue({ source: 'orch', targets: ['analyst'] });
            const msg = createMCPTaskSend('task-123', 'Analyze Q4', signal);

            expect(msg.jsonrpc).toBe('2.0');
            expect(msg.method).toBe('tasks/send');
            expect(msg.params.id).toBe('task-123');
            expect(msg.params.hct_signal).toBeDefined();
        });
    });
});
