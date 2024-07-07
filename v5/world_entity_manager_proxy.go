package sameriver

func (w *World) MaxEntities() int {
	return w.Em.MaxEntities()
}

func (w *World) Components() *ComponentTable {
	return w.Em.components
}

func (w *World) Spawn(spec map[string]any) *Entity {
	return w.Em.Spawn(spec)
}

func (w *World) QueueSpawn(spec map[string]any) {
	w.Em.QueueSpawn(spec)
}

func (w *World) Despawn(e *Entity) {
	w.Em.Despawn(e)
	w.RemoveAllEntityLogics(e)
}

func (w *World) DespawnAll() {
	for _, e := range w.Em.GetCurrentEntitiesSet() {
		w.RemoveAllEntityLogics(e)
	}
	w.Em.DespawnAll()
}

func (w *World) Activate(e *Entity) {
	w.Em.Activate(e)
}

func (w *World) Deactivate(e *Entity) {
	w.Em.Deactivate(e)
}

func (w *World) GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return w.Em.GetUpdatedEntityList(q)
}

func (w *World) GetSortedUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return w.Em.GetSortedUpdatedEntityList(q)
}

func (w *World) GetUpdatedEntityListByName(name string) *UpdatedEntityList {
	return w.Em.GetUpdatedEntityListByName(name)
}

func (w *World) GetUpdatedEntityListByComponents(names []ComponentID) *UpdatedEntityList {
	return w.Em.GetUpdatedEntityListByComponents(names)
}

func (w *World) UniqueTaggedEntity(tag string) (*Entity, error) {
	return w.Em.UniqueTaggedEntity(tag)
}

func (w *World) UpdatedEntitiesWithTag(tag string) *UpdatedEntityList {
	return w.Em.UpdatedEntitiesWithTag(tag)
}

func (w *World) TagEntity(e *Entity, tags ...string) {
	w.Em.TagEntity(e, tags...)
}

func (w *World) TagEntities(entities []*Entity, tag string) {
	w.Em.TagEntities(entities, tag)
}

func (w *World) UntagEntity(e *Entity, tag string) {
	w.Em.UntagEntity(e, tag)
}

func (w *World) UntagEntities(entities []*Entity, tag string) {
	w.Em.UntagEntities(entities, tag)
}

func (w *World) NumEntities() (total int, active int) {
	return w.Em.NumEntities()
}

func (w *World) GetActiveEntitiesSet() map[*Entity]bool {
	return w.Em.GetActiveEntitiesSet()
}

func (w *World) GetCurrentEntitiesSet() map[int]*Entity {
	return w.Em.GetCurrentEntitiesSet()
}

func (w *World) DumpEntities() string {
	return w.Em.DumpEntities()
}
