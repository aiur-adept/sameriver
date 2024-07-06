package sameriver

type InventorySystem struct {
	InventoryEntities *UpdatedEntityList
	itemSystem        *ItemSystem `sameriver-system-dependency:"-"`
}

func NewInventorySystem() *InventorySystem {
	return &InventorySystem{}
}

func (i *InventorySystem) Create(listing map[string]int) Inventory {
	result := NewInventory()
	for arch, count := range listing {
		if count != 0 {
			item := i.itemSystem.CreateStackSimple(count, arch)
			result.Credit(item)
		}
	}
	return result
}

// System funcs

func (i *InventorySystem) GetComponentDeps() []any {
	return []any{
		INVENTORY_, INVENTORY, "INVENTORY",
	}
}

func (i *InventorySystem) LinkWorld(w *World) {

	i.InventoryEntities = w.Em.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"inventory",
			w.Em.components.BitArrayFromIDs([]ComponentID{INVENTORY_})))
}

func (i *InventorySystem) Update(dt_ms float64) {
	// nil
}

func (i *InventorySystem) Expand(n int) {
	// nil
}
