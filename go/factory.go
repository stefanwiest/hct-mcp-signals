package hctmcpsignals

import "time"

// Option is a functional option for creating signals.
type Option func(*HCTSignal)

// WithUrgency sets the urgency level (1-10).
func WithUrgency(u int) Option {
	return func(s *HCTSignal) {
		if s.Performance == nil {
			s.Performance = &Performance{}
		}
		if u < 1 {
			u = 1
		}
		if u > 10 {
			u = 10
		}
		s.Performance.Urgency = u
	}
}

// WithTempo sets the tempo.
func WithTempo(t Tempo) Option {
	return func(s *HCTSignal) {
		if s.Performance == nil {
			s.Performance = &Performance{}
		}
		s.Performance.Tempo = t
	}
}

// WithTimeoutMs sets the timeout in milliseconds.
func WithTimeoutMs(ms int64) Option {
	return func(s *HCTSignal) {
		if s.Performance == nil {
			s.Performance = &Performance{}
		}
		s.Performance.TimeoutMs = &ms
	}
}

// WithPayload sets the payload.
func WithPayload(p map[string]interface{}) Option {
	return func(s *HCTSignal) {
		s.Payload = p
	}
}

// WithPayloadEntry adds a single payload entry.
func WithPayloadEntry(key string, value interface{}) Option {
	return func(s *HCTSignal) {
		if s.Payload == nil {
			s.Payload = make(map[string]interface{})
		}
		s.Payload[key] = value
	}
}

// WithHoldType sets the hold type for FERMATA signals.
func WithHoldType(ht HoldType) Option {
	return func(s *HCTSignal) {
		if s.Conditions == nil {
			s.Conditions = &Conditions{}
		}
		s.Conditions.HoldType = &ht
	}
}

// WithRepeatUntil sets the repeat condition for VAMP signals.
func WithRepeatUntil(condition string) Option {
	return func(s *HCTSignal) {
		if s.Conditions == nil {
			s.Conditions = &Conditions{}
		}
		s.Conditions.RepeatUntil = &condition
	}
}

// WithQualityThreshold sets the quality threshold (0-1).
func WithQualityThreshold(threshold float64) Option {
	return func(s *HCTSignal) {
		if s.Conditions == nil {
			s.Conditions = &Conditions{}
		}
		s.Conditions.QualityThreshold = &threshold
	}
}

// newSignal creates a base signal.
func newSignal(signalType SignalType, source string, targets []string, opts ...Option) *HCTSignal {
	s := &HCTSignal{
		Type:    signalType,
		Source:  source,
		Targets: targets,
		Payload: make(map[string]interface{}),
		Performance: &Performance{
			Urgency: 5,
			Tempo:   Moderato,
		},
		Timestamp: time.Now().UTC(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// NewCue creates a CUE signal to trigger agent activation.
func NewCue(source string, targets []string, opts ...Option) *HCTSignal {
	return newSignal(Cue, source, targets, opts...)
}

// NewFermata creates a FERMATA signal to hold for approval.
func NewFermata(source, reason string, opts ...Option) *HCTSignal {
	allOpts := append([]Option{
		WithPayloadEntry("reason", reason),
		WithHoldType(HoldHuman),
	}, opts...)
	return newSignal(Fermata, source, []string{"governance"}, allOpts...)
}

// NewAttacca creates an ATTACCA signal for immediate transition.
func NewAttacca(source string, targets []string, opts ...Option) *HCTSignal {
	allOpts := append([]Option{
		WithUrgency(10),
		WithTempo(Presto),
	}, opts...)
	return newSignal(Attacca, source, targets, allOpts...)
}

// NewVamp creates a VAMP signal to repeat until condition met.
func NewVamp(source, repeatUntil string, opts ...Option) *HCTSignal {
	threshold := 0.9
	timeout := int64(60000)
	allOpts := append([]Option{
		WithRepeatUntil(repeatUntil),
		WithQualityThreshold(threshold),
		WithTimeoutMs(timeout),
	}, opts...)
	return newSignal(Vamp, source, []string{source}, allOpts...)
}

// NewCaesura creates a CAESURA signal for full stop.
func NewCaesura(source, reason string, opts ...Option) *HCTSignal {
	allOpts := append([]Option{
		WithPayloadEntry("reason", reason),
		WithUrgency(10),
		WithTempo(Presto),
	}, opts...)
	return newSignal(Caesura, source, []string{"*"}, allOpts...)
}

// NewTacet creates a TACET signal to mark agent as inactive.
func NewTacet(source string, opts ...Option) *HCTSignal {
	return newSignal(Tacet, source, nil, opts...)
}

// NewDownbeat creates a DOWNBEAT signal for global synchronization.
func NewDownbeat(source, syncPoint string, opts ...Option) *HCTSignal {
	allOpts := append([]Option{
		WithPayloadEntry("sync_point", syncPoint),
	}, opts...)
	return newSignal(Downbeat, source, []string{"*"}, allOpts...)
}
