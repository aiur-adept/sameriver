package sameriver

type EntityPredicate func(*Entity) bool

// just for fun, let's define a useless function. let's be luxuriant about
// modern storage space. we can waste the bytes.
func NullEntityPredicate(e *Entity) bool {
	return false
}
func AllEntityPredicate(e *Entity) bool {
	return true
}

type EntityFilter struct {
	Name      string
	Predicate func(e *Entity) bool
}

func NewEntityFilter(
	name string, f func(e *Entity) bool) EntityFilter {
	return EntityFilter{Name: name, Predicate: f}
}

func (q EntityFilter) Test(e *Entity) bool {
	return q.Predicate(e)
}
