package engine

type GameEventQuery struct {
	Type      GameEventType
	Predicate func(e GameEvent) bool
}

func (q GameEventQuery) Test(e GameEvent) bool {
	return q.Type == e.Type && q.Predicate(e)
}

// Construct a new GameEventQuery which only asks about
// the Type of the event
func NewSimpleGameEventQuery(Type GameEventType) GameEventQuery {

	return GameEventQuery{Type, nil}
}

// Construct a new GameEventQuery which asks about Type and
// a user-given predicate
func NewPredicateGameEventQuery(
	Type GameEventType,
	predicate func(e GameEvent) bool) GameEventQuery {

	return GameEventQuery{Type, predicate}
}
