package sameriver

type FloatMap struct {
	M              map[string]float64
	ValidIntervals map[string][2]float64
}

func NewFloatMap(m map[string]float64) FloatMap {
	return FloatMap{
		M:              m,
		ValidIntervals: make(map[string][2]float64),
	}
}

func (m *FloatMap) CopyOf() FloatMap {
	m2 := make(map[string]float64)
	for key := range m.M {
		m2[key] = m.M[key]
	}
	return NewFloatMap(m2)
}

func (m *FloatMap) SetValidInterval(k string, a, b float64) {
	m.ValidIntervals[k] = [2]float64{a, b}
}

func (m *FloatMap) ValCanBeSetTo(k string, v float64) bool {
	if validInterval, exists := m.ValidIntervals[k]; exists {
		return v >= validInterval[0] && v <= validInterval[1]
	} else {
		return true
	}
}

func (m *FloatMap) Set(k string, v float64) {
	if validInterval, exists := m.ValidIntervals[k]; exists {
		if v < validInterval[0] {
			v = validInterval[0]
		} else if v > validInterval[1] {
			v = validInterval[1]
		}
	}
	m.M[k] = v
}

func (m *FloatMap) Get(k string) float64 {
	return m.M[k]
}

func (m *FloatMap) Has(k string) bool {
	_, ok := m.M[k]
	return ok
}
