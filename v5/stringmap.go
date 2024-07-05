package sameriver

type StringMap struct {
	M map[string]string
}

func NewStringMap(m map[string]string) StringMap {
	return StringMap{
		M: m,
	}
}

func (m *StringMap) CopyOf() StringMap {
	m2 := make(map[string]string)
	for key := range m.M {
		m2[key] = m.M[key]
	}
	return NewStringMap(m2)
}
