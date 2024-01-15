package sameriver

type StringMap struct {
	m map[string]string
}

func NewStringMap(m map[string]string) StringMap {
	return StringMap{m}
}

func (m *StringMap) CopyOf() StringMap {
	m2 := make(map[string]string)
	for key := range m.m {
		m2[key] = m.m[key]
	}
	return NewStringMap(m2)
}
