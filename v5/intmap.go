package sameriver

type IntMap struct {
	M              map[string]int
	ValidIntervals map[string][2]int
}

func NewIntMap(m map[string]int) IntMap {
	return IntMap{
		M:              m,
		ValidIntervals: make(map[string][2]int),
	}
}

func (m *IntMap) CopyOf() IntMap {
	m2 := make(map[string]int)
	for key := range m.M {
		m2[key] = m.M[key]
	}
	return IntMap{m2, m.ValidIntervals}
}

func (m *IntMap) SetValidInterval(k string, a, b int) {
	m.ValidIntervals[k] = [2]int{a, b}
}

func (m *IntMap) ValCanBeSetTo(k string, v int) bool {
	if validInterval, exists := m.ValidIntervals[k]; exists {
		return v >= validInterval[0] && v <= validInterval[1]
	} else {
		return true
	}
}

func (m *IntMap) Set(k string, v int) {
	if validInterval, exists := m.ValidIntervals[k]; exists {
		if v < validInterval[0] {
			v = validInterval[0]
		} else if v > validInterval[1] {
			v = validInterval[1]
		}
	}
	m.M[k] = v
}

func (m *IntMap) Get(k string) int {
	return m.M[k]
}

func (m *IntMap) Has(k string) bool {
	_, ok := m.M[k]
	return ok
}
