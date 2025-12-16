/**
 * MCP Integration Utilities
 *
 * Helpers for integrating HCT signals with MCP messages.
 */

import { HCTSignal, SignalType } from './types';

/** Key used for HCT signal extension in MCP messages */
export const HCT_SIGNAL_KEY = 'hct_signal';

/** MCP message structure */
export interface MCPMessage {
    jsonrpc: '2.0';
    method: string;
    params: Record<string, unknown>;
}

/**
 * Embed an HCT signal into MCP params.
 *
 * @param params - Existing MCP params
 * @param signal - HCT signal to embed
 * @returns Updated params with hct_signal extension
 */
export function embedSignal(
    params: Record<string, unknown>,
    signal: HCTSignal
): Record<string, unknown> {
    return {
        ...params,
        [HCT_SIGNAL_KEY]: signal,
    };
}

/**
 * Extract HCT signal from MCP message.
 *
 * @param message - MCP message or params
 * @returns HCTSignal if present, undefined otherwise
 */
export function extractSignal(
    message: Record<string, unknown>
): HCTSignal | undefined {
    const signal = message[HCT_SIGNAL_KEY] as HCTSignal | undefined;
    return signal;
}

/**
 * Check if MCP message contains HCT signal.
 */
export function hasSignal(message: Record<string, unknown>): boolean {
    return HCT_SIGNAL_KEY in message;
}

/**
 * Get signal type from MCP message without full parsing.
 */
export function getSignalType(
    message: Record<string, unknown>
): SignalType | undefined {
    const signal = message[HCT_SIGNAL_KEY] as { type?: string } | undefined;
    if (!signal?.type) return undefined;
    return signal.type as SignalType;
}

/**
 * Create a complete MCP tasks/send message with HCT signal.
 */
export function createMCPTaskSend(
    taskId: string,
    content: string,
    signal: HCTSignal,
    role: string = 'user'
): MCPMessage {
    return {
        jsonrpc: '2.0',
        method: 'tasks/send',
        params: {
            id: taskId,
            message: { role, content },
            [HCT_SIGNAL_KEY]: signal,
        },
    };
}

/**
 * Create an MCP tasks/sendSubscribe message with HCT signal.
 * Used for FERMATA and VAMP signals that need response streaming.
 */
export function createMCPTaskSubscribe(
    taskId: string,
    content: string,
    signal: HCTSignal,
    role: string = 'user'
): MCPMessage {
    return {
        jsonrpc: '2.0',
        method: 'tasks/sendSubscribe',
        params: {
            id: taskId,
            message: { role, content },
            [HCT_SIGNAL_KEY]: signal,
        },
    };
}
