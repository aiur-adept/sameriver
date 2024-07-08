package sameriver

type TimeAccumulator struct {
	Accum_ms  float64
	Period_ms float64
}

func NewTimeAccumulator(period_ms float64) TimeAccumulator {
	t := TimeAccumulator{}
	t.Accum_ms = 0
	t.Period_ms = period_ms
	return t
}

func (t *TimeAccumulator) Tick(dt_ms float64) bool {
	t.Accum_ms += dt_ms
	had_tick := false
	for t.Accum_ms >= t.Period_ms {
		modSubtract := int(t.Accum_ms / t.Period_ms)
		t.Accum_ms -= float64(modSubtract) * t.Period_ms
		had_tick = true
	}
	return had_tick
}

func (t *TimeAccumulator) Completion() float64 {
	return t.Accum_ms / t.Period_ms
}

func (t *TimeAccumulator) CompletionAfterDT(dt_ms float64) float64 {
	return (t.Accum_ms + dt_ms) / t.Period_ms
}
