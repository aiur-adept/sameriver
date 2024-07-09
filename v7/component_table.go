package sameriver

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

var componentKindStrings = map[ComponentKind]string{
	VEC2D:           "VEC2D",
	BOOL:            "BOOL",
	INT:             "INT",
	FLOAT64:         "FLOAT64",
	TIME:            "TIME",
	TIMEACCUMULATOR: "TIMEACCUMULATOR",
	STRING:          "STRING",
	SPRITE:          "SPRITE",
	TAGLIST:         "TAGLIST",
	INTMAP:          "INTMAP",
	FLOATMAP:        "FLOATMAP",
	STRINGMAP:       "STRINGMAP",
	ITEM:            "ITEM",
	INVENTORY:       "INVENTORY",
}

type ComponentTable struct {
	// the size of the tables
	Capacity int

	NextIx     int                           `json:"nextIx"`
	Ixs        map[ComponentID]int           `json:"ixs"`
	IxsRev     map[int]ComponentID           `json:"ixsRev"`
	Strings    map[ComponentID]string        `json:"strings"`
	StringsRev map[string]ComponentID        `json:"stringsRev"`
	Kinds      map[ComponentID]ComponentKind `json:"kinds"`

	ComponentStrings   []map[string]bool
	ComponentBitArrays []bitarray.BitArray `json:"-"`

	// data storage
	Vec2DMap           map[ComponentID][]Vec2D           `json:"vec2DMap"`
	BoolMap            map[ComponentID][]bool            `json:"boolMap"`
	IntMap             map[ComponentID][]int             `json:"intMap"`
	Float64Map         map[ComponentID][]float64         `json:"float64Map"`
	TimeMap            map[ComponentID][]time.Time       `json:"timeMap"`
	TimeAccumulatorMap map[ComponentID][]TimeAccumulator `json:"timeAccumulatorMap"`
	StringMap          map[ComponentID][]string          `json:"stringMap"`
	SpriteMap          map[ComponentID][]Sprite          `json:"spriteMap"`
	TagListMap         map[ComponentID][]TagList         `json:"tagListMap"`
	IntMapMap          map[ComponentID][]IntMap          `json:"intMapMap"`
	FloatMapMap        map[ComponentID][]FloatMap        `json:"floatMapMap"`
	StringMapMap       map[ComponentID][]StringMap       `json:"stringMapMap"`
	ItemMap            map[ComponentID][]Item            `json:"itemMap"`
	InventoryMap       map[ComponentID][]Inventory       `json:"inventoryMap"`
}

func NewComponentTable(capacity int) ComponentTable {
	return ComponentTable{
		Capacity: capacity,

		Ixs:        make(map[ComponentID]int),
		IxsRev:     make(map[int]ComponentID),
		Strings:    make(map[ComponentID]string),
		StringsRev: make(map[string]ComponentID),
		Kinds:      make(map[ComponentID]ComponentKind),

		ComponentStrings:   make([]map[string]bool, capacity),
		ComponentBitArrays: make([]bitarray.BitArray, capacity),

		Vec2DMap:           make(map[ComponentID][]Vec2D),
		BoolMap:            make(map[ComponentID][]bool),
		IntMap:             make(map[ComponentID][]int),
		Float64Map:         make(map[ComponentID][]float64),
		TimeMap:            make(map[ComponentID][]time.Time),
		TimeAccumulatorMap: make(map[ComponentID][]TimeAccumulator),
		StringMap:          make(map[ComponentID][]string),
		SpriteMap:          make(map[ComponentID][]Sprite),
		TagListMap:         make(map[ComponentID][]TagList),
		IntMapMap:          make(map[ComponentID][]IntMap),
		FloatMapMap:        make(map[ComponentID][]FloatMap),
		StringMapMap:       make(map[ComponentID][]StringMap),
		ItemMap:            make(map[ComponentID][]Item),
		InventoryMap:       make(map[ComponentID][]Inventory),
	}
}

// this is likely to be an expensive operation
func (ct *ComponentTable) expand(n int) {
	Logger.Printf("Expanding component tables from %d to %d", ct.Capacity, ct.Capacity+n)
	for name, slice := range ct.Vec2DMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]Vec2D, n)
		ct.Vec2DMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.BoolMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]bool, n)
		ct.BoolMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.IntMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]int, n)
		ct.IntMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.Float64Map {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]float64, n)
		ct.Float64Map[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.TimeMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]time.Time, n)
		ct.TimeMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.TimeAccumulatorMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]TimeAccumulator, n)
		ct.TimeAccumulatorMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.StringMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]string, n)
		ct.StringMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.SpriteMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]Sprite, n)
		ct.SpriteMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.TagListMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]TagList, n)
		ct.TagListMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.IntMapMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]IntMap, n)
		ct.IntMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.FloatMapMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]FloatMap, n)
		ct.FloatMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.StringMapMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]StringMap, n)
		ct.StringMapMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.ItemMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]Item, n)
		ct.ItemMap[name] = append(slice, extraSpace...)
	}
	for name, slice := range ct.InventoryMap {
		Logger.Printf("Expanding table of component %s,%s", componentKindStrings[ct.Kinds[name]], ct.Strings[name])
		extraSpace := make([]Inventory, n)
		ct.InventoryMap[name] = append(slice, extraSpace...)
	}
	// expand ComponentBitArrays
	ct.ComponentStrings = append(ct.ComponentStrings, make([]map[string]bool, n)...)
	ct.ComponentBitArrays = append(ct.ComponentBitArrays, make([]bitarray.BitArray, n)...)
	ct.Capacity += n
}

func (ct *ComponentTable) RegisterComponentStrings(strings map[ComponentID]string) {
	for name, str := range strings {
		ct.Strings[name] = str
		ct.StringsRev[str] = name
	}
}

func (ct *ComponentTable) index(name ComponentID) {
	// increment index and store (used for bitarray generation)
	ct.Ixs[name] = ct.NextIx
	ct.IxsRev[ct.NextIx] = name
	ct.NextIx++
}

func (ct *ComponentTable) ComponentExists(name ComponentID) bool {
	if _, ok := ct.Ixs[name]; ok {
		return true
	}
	return false
}

func (ct *ComponentTable) addComponent(kind ComponentKind, name ComponentID, str string) {
	// create table in appropriate map
	// (note we allocate with capacity 2* so that if we reach max entities the
	// first time expanding the tables won't necessarily be expensive; but
	// then again, if we do reach the NEW capacity, the slices will have to
	// be reallocated to new memory locations as they'll have totally
	// eaten up the capacity)
	switch kind {
	case VEC2D:
		ct.Vec2DMap[name] = make([]Vec2D, ct.Capacity, 2*ct.Capacity)
	case BOOL:
		ct.BoolMap[name] = make([]bool, ct.Capacity, 2*ct.Capacity)
	case INT:
		ct.IntMap[name] = make([]int, ct.Capacity, 2*ct.Capacity)
	case FLOAT64:
		ct.Float64Map[name] = make([]float64, ct.Capacity, 2*ct.Capacity)
	case TIME:
		ct.TimeMap[name] = make([]time.Time, ct.Capacity, 2*ct.Capacity)
	case TIMEACCUMULATOR:
		ct.TimeAccumulatorMap[name] = make([]TimeAccumulator, ct.Capacity, 2*ct.Capacity)
	case STRING:
		ct.StringMap[name] = make([]string, ct.Capacity, 2*ct.Capacity)
	case SPRITE:
		ct.SpriteMap[name] = make([]Sprite, ct.Capacity, 2*ct.Capacity)
	case TAGLIST:
		ct.TagListMap[name] = make([]TagList, ct.Capacity, 2*ct.Capacity)
	case INTMAP:
		ct.IntMapMap[name] = make([]IntMap, ct.Capacity, 2*ct.Capacity)
	case FLOATMAP:
		ct.FloatMapMap[name] = make([]FloatMap, ct.Capacity, 2*ct.Capacity)
	case STRINGMAP:
		ct.StringMapMap[name] = make([]StringMap, ct.Capacity, 2*ct.Capacity)
	case ITEM:
		ct.ItemMap[name] = make([]Item, ct.Capacity, 2*ct.Capacity)
	case INVENTORY:
		ct.InventoryMap[name] = make([]Inventory, ct.Capacity, 2*ct.Capacity)
	default:
		panic(fmt.Sprintf("added component of kind %s has no case in component_table.go", componentKindStrings[kind]))
	}

	// note name and kind
	ct.index(name)
	ct.Kinds[name] = kind

	// note string
	ct.Strings[name] = str
	ct.StringsRev[str] = name
}

func (ct *ComponentTable) AssertValidComponentSet(cs ComponentSet) {
	for name := range cs.vec2DMap {
		if _, ok := ct.Vec2DMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in vec2DMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.boolMap {
		if _, ok := ct.BoolMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in boolMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.intMap {
		if _, ok := ct.IntMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.float64Map {
		if _, ok := ct.Float64Map[name]; !ok {
			panic(fmt.Sprintf("%s not found in float64Map - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.timeMap {
		if _, ok := ct.TimeMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in timeMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.stringMap {
		if _, ok := ct.StringMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in stringMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.spriteMap {
		if _, ok := ct.SpriteMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in spriteMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.tagListMap {
		if _, ok := ct.TagListMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in tagListMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.intMapMap {
		if _, ok := ct.IntMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in intMapMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.floatMapMap {
		if _, ok := ct.FloatMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in floatMapMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.stringMapMap {
		if _, ok := ct.StringMapMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in stringMapMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.itemMap {
		if _, ok := ct.ItemMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in itemMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
	for name := range cs.inventoryMap {
		if _, ok := ct.InventoryMap[name]; !ok {
			panic(fmt.Sprintf("%s not found in inventoryMap - maybe not registered yet?", ct.Strings[name]))
		}
	}
}

func (ct *ComponentTable) ApplyComponentSet(e *Entity, spec map[ComponentID]any) {
	ct.applyComponentSet(e, ct.makeComponentSet(spec))
}

func (ct *ComponentTable) applyComponentSet(e *Entity, cs ComponentSet) {
	ct.AssertValidComponentSet(cs)
	ct.ComponentStrings[e.ID] = make(map[string]bool)
	for name, v := range cs.vec2DMap {
		ct.Vec2DMap[name][e.ID] = v
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, b := range cs.boolMap {
		ct.BoolMap[name][e.ID] = b
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, i := range cs.intMap {
		ct.IntMap[name][e.ID] = i
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, f := range cs.float64Map {
		ct.Float64Map[name][e.ID] = f
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, t := range cs.timeMap {
		ct.TimeMap[name][e.ID] = t
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, t := range cs.timeAccumulatorMap {
		ct.TimeAccumulatorMap[name][e.ID] = t
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, s := range cs.stringMap {
		ct.StringMap[name][e.ID] = s
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, s := range cs.spriteMap {
		ct.SpriteMap[name][e.ID] = s
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, t := range cs.tagListMap {
		ct.TagListMap[name][e.ID] = t
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, m := range cs.intMapMap {
		ct.IntMapMap[name][e.ID] = m
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, m := range cs.floatMapMap {
		ct.FloatMapMap[name][e.ID] = m
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, m := range cs.stringMapMap {
		ct.StringMapMap[name][e.ID] = m
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, m := range cs.itemMap {
		ct.ItemMap[name][e.ID] = m
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}
	for name, m := range cs.inventoryMap {
		ct.InventoryMap[name][e.ID] = m
		e.Components = append(e.Components, ct.Strings[name])
		ct.ComponentStrings[e.ID][ct.Strings[name]] = true
	}

	ct.orBitArrayInto(e, ct.bitArrayFromComponentSet(cs))
}

func (ct *ComponentTable) orStringIntoBitArray(eid int, component string) {
	if ct.ComponentBitArrays[eid] == nil {
		ct.ComponentBitArrays[eid] = bitarray.NewBitArray(uint64(len(ct.Ixs)))
	}
	ix := ct.Ixs[ct.StringsRev[component]]
	ct.ComponentBitArrays[eid].SetBit(uint64(ix))
	ct.ComponentStrings[eid][component] = true
}

func (ct *ComponentTable) orBitArrayInto(e *Entity, b bitarray.BitArray) {
	if ct.ComponentBitArrays[e.ID] == nil {
		ct.ComponentBitArrays[e.ID] = bitarray.NewBitArray(uint64(len(ct.Ixs)))
	}
	for _, i := range ct.Ixs {
		bit, _ := b.GetBit(uint64(i))
		if bit {
			ct.ComponentBitArrays[e.ID].SetBit(uint64(i))
		}
	}
}

func (ct *ComponentTable) BitArrayFromIDs(IDs []ComponentID) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.Ixs)))
	for _, name := range IDs {
		b.SetBit(uint64(ct.Ixs[name]))
	}
	return b
}

func (ct *ComponentTable) BitArrayFromComponentSet(spec map[ComponentID]any) bitarray.BitArray {
	return ct.bitArrayFromComponentSet(ct.makeComponentSet(spec))
}

func (ct *ComponentTable) bitArrayFromComponentSet(cs ComponentSet) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.Ixs)))
	for name := range cs.names {
		b.SetBit(uint64(ct.Ixs[name]))
	}
	return b
}

// BitArrayToString prints a string representation of a component bitarray as a set with
// string representations of each component type whose bit is set
func (ct *ComponentTable) BitArrayToString(b bitarray.BitArray) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	names := make([]string, 0)
	for name, ix := range ct.Ixs {
		bit, _ := b.GetBit(uint64(ix))
		// the index into the array is the component type int from the
		// iota const block in component_enum.go
		if bit {
			names = append(names, ct.Strings[name])
		}
	}
	for i, name := range names {
		buf.WriteString(name)
		if i != len(names)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func (ct *ComponentTable) BitArrayToStringArray(b bitarray.BitArray) []string {
	names := make([]string, 0)
	for name, ix := range ct.Ixs {
		bit, _ := b.GetBit(uint64(ix))
		// the index into the array is the component type int from the
		// iota const block in component_enum.go
		if bit {
			names = append(names, ct.Strings[name])
		}
	}
	return names
}

func (ct *ComponentTable) BitArrayFromStrings(names []string) bitarray.BitArray {
	b := bitarray.NewBitArray(uint64(len(ct.Ixs)))
	for _, name := range names {
		b.SetBit(uint64(ct.Ixs[ct.StringsRev[name]]))
	}
	return b
}
func (ct *ComponentTable) guardInvalidComponentGet(e *Entity, name ComponentID) {
	var ix int
	var ok bool
	if ix, ok = ct.Ixs[name]; !ok {
		msg := fmt.Sprintf("Tried to access %s component; but there is no component with that name", ct.Strings[name])
		panic(msg)
	}
	bit, _ := ct.ComponentBitArrays[e.ID].GetBit(uint64(ix))
	if !bit {
		Logger.Printf("%s", ct.Strings[name])
		msg := fmt.Sprintf("Tried to get %s component of entity without: %s", ct.Strings[name], e)
		panic(msg)
	}
}

func (w *World) GetVec2D(e *Entity, name ComponentID) *Vec2D {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.Vec2DMap[name][e.ID]
}
func (w *World) GetBool(e *Entity, name ComponentID) *bool {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.BoolMap[name][e.ID]
}
func (w *World) GetInt(e *Entity, name ComponentID) *int {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.IntMap[name][e.ID]
}
func (w *World) GetFloat64(e *Entity, name ComponentID) *float64 {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.Float64Map[name][e.ID]
}
func (w *World) GetTime(e *Entity, name ComponentID) *time.Time {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.TimeMap[name][e.ID]
}
func (w *World) GetTimeAccumulator(e *Entity, name ComponentID) *TimeAccumulator {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.TimeAccumulatorMap[name][e.ID]
}
func (w *World) GetString(e *Entity, name ComponentID) *string {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.StringMap[name][e.ID]
}
func (w *World) GetSprite(e *Entity, name ComponentID) *Sprite {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.SpriteMap[name][e.ID]
}
func (w *World) GetTagList(e *Entity, name ComponentID) *TagList {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.TagListMap[name][e.ID]
}
func (w *World) GetIntMap(e *Entity, name ComponentID) *IntMap {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.IntMapMap[name][e.ID]
}
func (w *World) GetFloatMap(e *Entity, name ComponentID) *FloatMap {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.FloatMapMap[name][e.ID]
}
func (w *World) GetStringMap(e *Entity, name ComponentID) *StringMap {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.StringMapMap[name][e.ID]
}
func (w *World) GetItem(e *Entity, name ComponentID) *Item {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.ItemMap[name][e.ID]
}
func (w *World) GetInventory(e *Entity, name ComponentID) *Inventory {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	return &w.Em.ComponentsTable.InventoryMap[name][e.ID]
}

func (w *World) GetVal(e *Entity, name ComponentID) any {
	w.Em.ComponentsTable.guardInvalidComponentGet(e, name)
	kind := w.Em.ComponentsTable.Kinds[name]
	switch kind {
	case VEC2D:
		return &w.Em.ComponentsTable.Vec2DMap[name][e.ID]
	case BOOL:
		return &w.Em.ComponentsTable.BoolMap[name][e.ID]
	case INT:
		return &w.Em.ComponentsTable.IntMap[name][e.ID]
	case FLOAT64:
		return &w.Em.ComponentsTable.Float64Map[name][e.ID]
	case STRING:
		return &w.Em.ComponentsTable.StringMap[name][e.ID]
	case SPRITE:
		return &w.Em.ComponentsTable.SpriteMap[name][e.ID]
	case TAGLIST:
		return &w.Em.ComponentsTable.TagListMap[name][e.ID]
	case INTMAP:
		return &w.Em.ComponentsTable.IntMapMap[name][e.ID]
	case FLOATMAP:
		return &w.Em.ComponentsTable.FloatMapMap[name][e.ID]
	case STRINGMAP:
		return &w.Em.ComponentsTable.StringMapMap[name][e.ID]
	case ITEM:
		return &w.Em.ComponentsTable.ItemMap[name][e.ID]
	case INVENTORY:
		return &w.Em.ComponentsTable.InventoryMap[name][e.ID]
	default:
		panic(fmt.Sprintf("Can't get component with ID %d - it doesn't seem to exist", name))
	}
}
