/**
 * HCT Signal Types - Auto-generated from hct-spec/spec.yaml
 * Generated: 2025-12-18T10:37:34.791829Z
 *
 * DO NOT EDIT MANUALLY - Edit spec.yaml and regenerate.
 */

export enum SignalType {
    /** Trigger activation */
    CUE = 'cue',
    /** Hold for approval */
    FERMATA = 'fermata',
    /** Immediate transition */
    ATTACCA = 'attacca',
    /** Repeat until condition */
    VAMP = 'vamp',
    /** Full stop */
    CAESURA = 'caesura',
    /** Agent inactive */
    TACET = 'tacet',
    /** Global sync point */
    DOWNBEAT = 'downbeat',
}

export enum Tempo {
    /** Very slow (~60s) */
    LARGO = 'largo',
    /** Walking pace (~30s) */
    ANDANTE = 'andante',
    /** Moderate (~15s) */
    MODERATO = 'moderato',
    /** Fast (~5s) */
    ALLEGRO = 'allegro',
    /** Very fast (~1s) */
    PRESTO = 'presto',
}

export enum DynamicsLevel {
    /** pianissimo */
    PP = 'pp',
    /** piano */
    P = 'p',
    /** mezzo-forte */
    MF = 'mf',
    /** forte */
    F = 'f',
    /** fortissimo */
    FF = 'ff',
}
