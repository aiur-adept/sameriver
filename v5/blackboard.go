package sameriver

type Blackboard struct {
	Name   string
	State  map[string]any
	Events *EventBus `json:"-"`
}

func NewBlackboard(name string) *Blackboard {
	return &Blackboard{
		Name:   name,
		State:  make(map[string]any),
		Events: NewEventBus("blackboard-" + name),
	}
}

func (b *Blackboard) Has(k string) bool {
	_, ok := b.State[k]
	return ok
}

func (b *Blackboard) Get(k string) any {
	return b.State[k]
}

func (b *Blackboard) Set(k string, v any) {
	b.State[k] = v
}

func (b *Blackboard) Remove(k string) {
	delete(b.State, k)
}
