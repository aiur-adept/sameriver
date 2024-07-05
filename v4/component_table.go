package sameriver

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type ComponentKind int
type ComponentID int

const (
	VEC2D ComponentKind = iota
	BOOL
	INT
	FLOAT64
	TIME
	TIMEACCUMULATOR
	STRING
	SPRITE
	TAGLIST
	INTMAP
	FLOATMAP
	STRINGMAP
	ITEM
	INVENTORY
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
	capacity int

	NextIx     int                           `json:"nextIx"`
	Ixs        map[ComponentID]int           `json:"ixs"`
	IxsRev     map[int]ComponentID           `json:"ixsRev"`
	Strings    map[ComponentID]string        `json:"strings"`
	StringsRev map[string]ComponentID        `json:"stringsRev"`
	Kinds      map[ComponentID]ComponentKind `json:"kinds"`

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

func NewComponentTable(capacity int) *ComponentTable {
	return &ComponentTable{
		capacity: capacity,

		Ixs:        make(map[ComponentID]int),
		IxsRev:     make(map[int]ComponentID),
		Strings:    make(map[ComponentID]string),
		StringsRev: make(map[string]ComponentID),
		Kinds:      make(map[ComponentID]ComponentKind),

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
	Logger.Printf("Expanding component tables from %d to %d", ct.capacity, ct.capacity+n)
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
	ct.capacity += n
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
		ct.Vec2DMap[name] = make([]Vec2D, ct.capacity, 2*ct.capacity)
	case BOOL:
		ct.BoolMap[name] = make([]bool, ct.capacity, 2*ct.capacity)
	case INT:
		ct.IntMap[name] = make([]int, ct.capacity, 2*ct.capacity)
	case FLOAT64:
		ct.Float64Map[name] = make([]float64, ct.capacity, 2*ct.capacity)
	case TIME:
		ct.TimeMap[name] = make([]time.Time, ct.capacity, 2*ct.capacity)
	case TIMEACCUMULATOR:
		ct.TimeAccumulatorMap[name] = make([]TimeAccumulator, ct.capacity, 2*ct.capacity)
	case STRING:
		ct.StringMap[name] = make([]string, ct.capacity, 2*ct.capacity)
	case SPRITE:
		ct.SpriteMap[name] = make([]Sprite, ct.capacity, 2*ct.capacity)
	case TAGLIST:
		ct.TagListMap[name] = make([]TagList, ct.capacity, 2*ct.capacity)
	case INTMAP:
		ct.IntMapMap[name] = make([]IntMap, ct.capacity, 2*ct.capacity)
	case FLOATMAP:
		ct.FloatMapMap[name] = make([]FloatMap, ct.capacity, 2*ct.capacity)
	case STRINGMAP:
		ct.StringMapMap[name] = make([]StringMap, ct.capacity, 2*ct.capacity)
	case ITEM:
		ct.ItemMap[name] = make([]Item, ct.capacity, 2*ct.capacity)
	case INVENTORY:
		ct.InventoryMap[name] = make([]Inventory, ct.capacity, 2*ct.capacity)
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
	for name, v := range cs.vec2DMap {
		ct.Vec2DMap[name][e.ID] = v
	}
	for name, b := range cs.boolMap {
		ct.BoolMap[name][e.ID] = b
	}
	for name, i := range cs.intMap {
		ct.IntMap[name][e.ID] = i
	}
	for name, f := range cs.float64Map {
		ct.Float64Map[name][e.ID] = f
	}
	for name, t := range cs.timeMap {
		ct.TimeMap[name][e.ID] = t
	}
	for name, t := range cs.timeAccumulatorMap {
		ct.TimeAccumulatorMap[name][e.ID] = t
	}
	for name, s := range cs.stringMap {
		ct.StringMap[name][e.ID] = s
	}
	for name, s := range cs.spriteMap {
		ct.SpriteMap[name][e.ID] = s
	}
	for name, t := range cs.tagListMap {
		ct.TagListMap[name][e.ID] = t
	}
	for name, m := range cs.intMapMap {
		ct.IntMapMap[name][e.ID] = m
	}
	for name, m := range cs.floatMapMap {
		ct.FloatMapMap[name][e.ID] = m
	}
	for name, m := range cs.stringMapMap {
		ct.StringMapMap[name][e.ID] = m
	}
	for name, m := range cs.itemMap {
		ct.ItemMap[name][e.ID] = m
	}
	for name, m := range cs.inventoryMap {
		ct.InventoryMap[name][e.ID] = m
	}

	ct.orBitArrayInto(e, ct.bitArrayFromComponentSet(cs))
}

func (ct *ComponentTable) orBitArrayInto(e *Entity, b bitarray.BitArray) {
	if e.ComponentBitArray == nil {
		e.ComponentBitArray = bitarray.NewBitArray(uint64(len(ct.Ixs)))
	}
	for _, i := range ct.Ixs {
		bit, _ := b.GetBit(uint64(i))
		if bit {
			e.ComponentBitArray.SetBit(uint64(i))
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

func (ct *ComponentTable) guardInvalidComponentGet(e *Entity, name ComponentID) {
	var ix int
	var ok bool
	if ix, ok = ct.Ixs[name]; !ok {
		msg := fmt.Sprintf("Tried to access %s component; but there is no component with that name", ct.Strings[name])
		panic(msg)
	}
	bit, _ := e.ComponentBitArray.GetBit(uint64(ix))
	if !bit {
		Logger.Printf("%s", ct.Strings[name])
		msg := fmt.Sprintf("Tried to get %s component of entity without: %s", ct.Strings[name], e)
		panic(msg)
	}
}

func (e *Entity) GetVec2D(name ComponentID) *Vec2D {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.Vec2DMap[name][e.ID]
}
func (e *Entity) GetBool(name ComponentID) *bool {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.BoolMap[name][e.ID]
}
func (e *Entity) GetInt(name ComponentID) *int {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.IntMap[name][e.ID]
}
func (e *Entity) GetFloat64(name ComponentID) *float64 {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.Float64Map[name][e.ID]
}
func (e *Entity) GetTime(name ComponentID) *time.Time {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.TimeMap[name][e.ID]
}
func (e *Entity) GetTimeAccumulator(name ComponentID) *TimeAccumulator {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.TimeAccumulatorMap[name][e.ID]
}
func (e *Entity) GetString(name ComponentID) *string {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.StringMap[name][e.ID]
}
func (e *Entity) GetSprite(name ComponentID) *Sprite {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.SpriteMap[name][e.ID]
}
func (e *Entity) GetTagList(name ComponentID) *TagList {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.TagListMap[name][e.ID]
}
func (e *Entity) GetIntMap(name ComponentID) *IntMap {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.IntMapMap[name][e.ID]
}
func (e *Entity) GetFloatMap(name ComponentID) *FloatMap {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.FloatMapMap[name][e.ID]
}
func (e *Entity) GetStringMap(name ComponentID) *StringMap {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.StringMapMap[name][e.ID]
}
func (e *Entity) GetItem(name ComponentID) *Item {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.ItemMap[name][e.ID]
}
func (e *Entity) GetInventory(name ComponentID) *Inventory {
	e.World.em.components.guardInvalidComponentGet(e, name)
	return &e.World.em.components.InventoryMap[name][e.ID]
}

func (e *Entity) GetVal(name ComponentID) any {
	e.World.em.components.guardInvalidComponentGet(e, name)
	kind := e.World.em.components.Kinds[name]
	switch kind {
	case VEC2D:
		return &e.World.em.components.Vec2DMap[name][e.ID]
	case BOOL:
		return &e.World.em.components.BoolMap[name][e.ID]
	case INT:
		return &e.World.em.components.IntMap[name][e.ID]
	case FLOAT64:
		return &e.World.em.components.Float64Map[name][e.ID]
	case STRING:
		return &e.World.em.components.StringMap[name][e.ID]
	case SPRITE:
		return &e.World.em.components.SpriteMap[name][e.ID]
	case TAGLIST:
		return &e.World.em.components.TagListMap[name][e.ID]
	case INTMAP:
		return &e.World.em.components.IntMapMap[name][e.ID]
	case FLOATMAP:
		return &e.World.em.components.FloatMapMap[name][e.ID]
	case STRINGMAP:
		return &e.World.em.components.StringMapMap[name][e.ID]
	default:
		panic(fmt.Sprintf("Can't get component with ID %d - it doesn't seem to exist", name))
	}
}
