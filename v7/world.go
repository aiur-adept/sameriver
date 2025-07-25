package sameriver

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"

	"github.com/TwiN/go-color"
	"github.com/golang-collections/go-datastructures/bitarray"
)

type World struct {

	// rand.Seed for this world's run
	Seed int

	Width  float64
	Height float64

	Events *EventBus `json:"-"`
	Em     *EntityManager

	EFDSL *EFDSLEvaluator `json:"-"`

	IDGen IDGenerator

	// systems registered
	systems map[string]System
	// this is needed to associate ID's with Systems, since System is an
	// interface, not a struct type like LogicUnit that can have a name field
	systemsIDs map[System]int

	// logics for each system
	systemLogics map[string]*LogicUnit

	// logics invoked regularly by RuntimeSharer
	worldLogics map[string]*LogicUnit

	// entity logics
	entityLogics map[int][]*LogicUnit

	// funcs that can be called by name with data and get a result,
	// or to produce an effect
	funcs *FuncSet

	// Blackboards that entity's can join to share events and state
	Blackboards map[string]Blackboard

	// for sharing runtime among the various runtimelimiter kinds
	// and contains the RuntimeLimiters to which we Add() LogicUnits
	runtimeSharer *RuntimeLimitSharer
	// special runtime limiters for oneshots and interval logics
	oneshots  *RuntimeLimiter
	intervals *RuntimeLimiter

	// for statistics tracking - the avg ms used to run World.Update()
	totalRuntimeAvg_ms *float64

	// used for entity distance queries
	SpatialHasher       *SpatialHasher `json:"-"`
	DistanceHasherGridX int
	DistanceHasherGridY int
}

type WorldSpec struct {
	Seed                int
	Width               int
	Height              int
	DistanceHasherGridX int
	DistanceHasherGridY int
}

func destructureWorldSpec(spec map[string]any) WorldSpec {
	var seed, width, height int
	var distanceHasherGridX, distanceHasherGridY int
	if _, ok := spec["seed"].(int); ok {
		seed = spec["seed"].(int)
	} else {
		seed = 108
	}
	if _, ok := spec["width"].(int); ok {
		width = spec["width"].(int)
	} else {
		width = 100
	}
	if _, ok := spec["height"].(int); ok {
		height = spec["height"].(int)
	} else {
		height = 100
	}
	if _, ok := spec["distanceHasherGridX"].(int); ok {
		distanceHasherGridX = spec["distanceHasherGridX"].(int)
	} else {
		distanceHasherGridX = 32
	}
	if _, ok := spec["distanceHasherGridY"].(int); ok {
		distanceHasherGridY = spec["distanceHasherGridY"].(int)
	} else {
		distanceHasherGridY = 32
	}

	return WorldSpec{
		Seed:                seed,
		Width:               width,
		Height:              height,
		DistanceHasherGridX: distanceHasherGridX,
		DistanceHasherGridY: distanceHasherGridY,
	}
}

func NewWorld(spec map[string]any) *World {
	// seed a random number from [1,108]
	destructured := destructureWorldSpec(spec)
	Logger.Println(color.InBold(color.InWhiteOverCyan(fmt.Sprintf("[world seed: %d]", int(destructured.Seed)))))

	w := &World{
		Seed:          int(destructured.Seed),
		IDGen:         NewIDGenerator(),
		Width:         float64(destructured.Width),
		Height:        float64(destructured.Height),
		Events:        NewEventBus("world"),
		systems:       make(map[string]System),
		systemLogics:  make(map[string]*LogicUnit),
		systemsIDs:    make(map[System]int),
		worldLogics:   make(map[string]*LogicUnit),
		entityLogics:  make(map[int][]*LogicUnit),
		funcs:         NewFuncSet(nil),
		Blackboards:   make(map[string]Blackboard),
		runtimeSharer: NewRuntimeLimitSharer(),
	}

	// set up runtimesharer
	w.runtimeSharer.RegisterRunners(map[string]float64{
		"systems":        1,
		"world":          1,
		"entities":       1,
		"world-oneshot":  0.5,
		"world-interval": 0.5,
	})
	w.oneshots = w.runtimeSharer.RunnerMap["world-oneshot"]
	w.intervals = w.runtimeSharer.RunnerMap["world-interval"]

	// init entitymanager
	w.Em = NewEntityManager(w)

	// init EFDSL
	w.EFDSL = NewEFDSLEvaluator(w)

	// register basic components
	w.RegisterComponents([]any{
		GENERICTAGS_, TAGLIST, "GENERICTAGS",
		STATE_, INTMAP, "STATE",
		POSITION_, VEC2D, "POSITION",
		VELOCITY_, VEC2D, "VELOCITY",
		ACCELERATION_, VEC2D, "ACCELERATION",
		RIGIDBODY_, BOOL, "RIGIDBODY",
		BOX_, VEC2D, "BOX",
	})

	// set up distance spatial hasher
	w.SpatialHasher = NewSpatialHasher(
		destructured.DistanceHasherGridX,
		destructured.DistanceHasherGridY,
		w,
	)
	w.DistanceHasherGridX = destructured.DistanceHasherGridX
	w.DistanceHasherGridY = destructured.DistanceHasherGridY

	return w
}

func (w *World) Update(allowance_ms float64) (overunder_ms float64) {
	t0 := time.Now()
	// process entity manager and spatial hash before anything
	w.Em.Update(allowance_ms / 8)
	w.SpatialHasher.Update()
	remaining_ms := allowance_ms - float64(time.Since(t0).Nanoseconds())/1e6
	w.runtimeSharer.Share(remaining_ms)

	// maintain total runtime moving average
	total := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if w.totalRuntimeAvg_ms == nil {
		w.totalRuntimeAvg_ms = &total
	} else {
		*w.totalRuntimeAvg_ms = (*w.totalRuntimeAvg_ms + total) / 2.0
	}
	return overunder_ms
}

func (w *World) RegisterComponents(components []any) {
	if len(components)%3 != 0 {
		panic("malformed components specification given to RegisterComponents()")
	}
	// register given specs
	for i := 0; i < len(components); i += 3 {
		name := components[i].(ComponentID)
		kind := components[i+1].(ComponentKind)
		str := components[i+2].(string)
		if w.Em.ComponentsTable.ComponentExists(name) {
			Logger.Printf("[component %s already exists. Skipping...]", str)
			continue
		} else {
			Logger.Printf("%s%s%s", color.InGreen("[registering component: "), fmt.Sprintf("%s,%s", str, componentKindStrings[kind]), color.InGreen("]"))
			w.Em.ComponentsTable.addComponent(kind, name, str)
		}
	}
}

func (w *World) RegisterSystems(systems ...System) {
	// add all systems
	for _, s := range systems {
		systemName := reflect.TypeOf(s).Elem().Name()
		if !strings.HasSuffix(systemName, "System") {
			panic(fmt.Sprintf("System names must end with System; got %s", systemName))
		}
		componentDeps := s.GetComponentDeps()
		if len(componentDeps)%3 != 0 {
			panic("malformed GetComponentDeps()")
		}
		for i := 0; i < len(componentDeps); i += 3 {
			name := componentDeps[i].(ComponentID)
			kind := componentDeps[i+1].(ComponentKind)
			str := componentDeps[i+2].(string)
			if w.Em.ComponentsTable.ComponentExists(name) {
				Logger.Printf("System %s depends on component %s, which is already registered.", systemName, str)
				continue
			}
			Logger.Printf("Creating component %d of kind %s wanted by system %s", name, componentKindStrings[kind], systemName)
			w.RegisterComponents([]any{
				name, kind, str,
			})
		}
		w.addSystem(s)
	}
	// link up all systems' dependencies
	for _, s := range systems {
		w.linkSystemDependencies(s)
	}
}

func (w *World) SetSystemSchedule(systemName string, period_ms float64) {
	Logger.Printf("Setting %s period_ms %f", systemName, period_ms)
	s := w.systems[systemName]
	name := fmt.Sprintf("%s.Update()", reflect.TypeOf(s).Elem().Name())
	w.runtimeSharer.RunnerMap["systems"].SetSchedule(name, period_ms)
}

func (w *World) addSystem(s System) {
	name := reflect.TypeOf(s).Elem().Name()
	if _, ok := w.systems[name]; ok {
		panic(fmt.Sprintf("double-add of system %s", name))
	}
	w.systems[name] = s
	ID := w.IDGen.Next()
	w.systemsIDs[s] = ID
	s.LinkWorld(w)
	// add logic immediately rather than wait for RuntimeSharer.Share() to
	// process the add/remove logic channel so that if we call SetSystemSchedule()
	// immediately after RegisterSystems(), the LogicUnit will be in the runner
	// to set the runSchedule on
	l := &LogicUnit{
		name:        fmt.Sprintf("%s.Update()", name),
		f:           s.Update,
		active:      true,
		runSchedule: nil,
	}
	w.systemLogics[name] = l
	w.runtimeSharer.RunnerMap["systems"].addLogicImmediately(l)
}

func (w *World) assertSystemTypeValid(t reflect.Type) {
	if t.Kind() != reflect.Ptr {
		panic("Implementers of engine.System must be pointer-receivers")
	}
	typeName := t.Elem().Name()
	validName, _ := regexp.MatchString(".+System$", typeName)
	if !validName {
		panic(fmt.Sprintf("Implementers of System must have a name "+
			"matching regexp .+System$. %s did not", typeName))
	}
}

func (w *World) linkSystemDependencies(s System) {
	// foreach field of the underlying struct,
	// check if it has the tag `sameriver-system-dependency`
	// if it does, search for the system with the same type as that
	// field and assign it as a pointer, cast to the expected type,
	// to that field
	//
	// sType is going to be something like *CollisionSystem
	sType := reflect.TypeOf(s).Elem()
	// get a type to represent the System interface (to ensure dependencies
	// are to implementers of System)
	systemInterface := reflect.TypeOf((*System)(nil)).Elem()
	for i := 0; i < sType.NumField(); i++ {
		// for each field of the struct
		// f would be something like sh *SpatialHashSystem, possibly with a tag
		f := sType.Field(i)
		tagVal := f.Tag.Get("sameriver-system-dependency")
		if tagVal != "" {
			// check that tagged field implements System and is a valid System
			// implemented
			isSystem := f.Type.Implements(systemInterface)
			if !isSystem {
				panic(fmt.Sprintf("fields tagged sameriver-system-dependency "+
					"must implement engine.System "+
					"(field %s %v of %s did not pass this requirement",
					f.Name, f.Type, sType.Name()))
			}
			w.assertSystemTypeValid(f.Type)
			// iterate through the other systems and find one whose type matches
			// the field's type
			var foundSystem System
			for _, otherSystem := range w.systems {
				if otherSystem == s {
					continue
				}
				if reflect.TypeOf(otherSystem) == f.Type {
					foundSystem = otherSystem
					break
				}
			}
			if foundSystem == nil {
				if tagVal == "optional" {
					continue
				} else {
					panic(fmt.Sprintf("%s %v of %s dependency could not be "+
						"resolved. No system found of type %v.",
						f.Name, f.Type, sType.Name(), f.Type))
				}
			}
			// now that we have found the system which corresponds to the
			// dependency, we will assign it to the place it should be
			//
			// thank you to feilengcui008 from golang-nuts for this method of
			// assigning to an unexported pointer field whose value is nil
			//
			// since vf is nil value, vf.Elem() will be the zero value, and
			// since the zero value is not addressable or settable, we
			// need to allocate a new settable value at the same address
			v := reflect.Indirect(reflect.ValueOf(s))
			vf := v.Field(i)
			vf = reflect.NewAt(vf.Type(), unsafe.Pointer(vf.UnsafeAddr())).Elem()
			vf.Set(reflect.ValueOf(foundSystem))
		}
	}
}

func (w *World) SetTimeout(F func(), ms float64) {
	var l *LogicUnit
	schedule := NewTimeAccumulator(ms)
	l = &LogicUnit{
		name: fmt.Sprintf("oneshot-%d", w.IDGen.Next()),
		f: func(dt_ms float64) {
			F()
			w.oneshots.Remove(l)
		},
		active:      true,
		runSchedule: &schedule,
	}
	w.oneshots.Add(l)
}

func (w *World) SetInterval(F func(), ms float64) (interval string) {
	schedule := NewTimeAccumulator(ms)
	name := fmt.Sprintf("interval-%d", w.IDGen.Next())
	l := &LogicUnit{
		name: name,
		f: func(dt_ms float64) {
			F()
		},
		active:      true,
		runSchedule: &schedule,
	}
	w.intervals.Add(l)
	return name
}

// setinterval but it is guaranteed to run n times
func (w *World) SetNInterval(F func(), ms float64, n int) (interval string) {
	schedule := NewTimeAccumulator(ms)
	name := fmt.Sprintf("interval-%d", w.IDGen.Next())
	ran := 0
	var l *LogicUnit
	l = &LogicUnit{
		name: name,
		f: func(dt_ms float64) {
			F()
			ran++
			if ran == n {
				w.intervals.Remove(l)
			}
		},
		active:      true,
		runSchedule: &schedule,
	}
	w.intervals.Add(l)
	return name
}

func (w *World) ClearInterval(interval string) {
	w.intervals.Remove(w.intervals.logicUnitsMap[interval])
}

func (w *World) AddEntityLogic(e *Entity, Name string, F func(dt_ms float64)) *LogicUnit {
	l := &LogicUnit{
		name:    e.LogicUnitName(Name),
		f:       F,
		active:  true,
		worldID: e.ID,
	}
	if _, ok := w.entityLogics[e.ID]; !ok {
		w.entityLogics[e.ID] = make([]*LogicUnit, 0)
	}
	w.entityLogics[e.ID] = append(w.entityLogics[e.ID], l)
	w.runtimeSharer.RunnerMap["entities"].Add(l)
	return l
}

func (w *World) RemoveEntityLogic(e *Entity, Name string) {
	if logicUnits, ok := w.entityLogics[e.ID]; ok {
		for i, logic := range logicUnits {
			if logic.name == e.LogicUnitName(Name) {
				w.runtimeSharer.RunnerMap["entities"].Remove(logic)
				// remove the logicunit from the slice
				w.entityLogics[e.ID] = append(w.entityLogics[e.ID][:i], w.entityLogics[e.ID][i+1:]...)
				w.IDGen.Free(logic.worldID)
				return
			}
		}
	}
}

func (w *World) AddLogic(Name string, F func(dt_ms float64)) *LogicUnit {
	if _, ok := w.worldLogics[Name]; ok {
		panic(fmt.Sprintf("double-add of world logic %s", Name))
	}
	l := &LogicUnit{
		name:    Name,
		f:       F,
		active:  true,
		worldID: w.IDGen.Next(),
	}
	w.worldLogics[Name] = l
	w.runtimeSharer.RunnerMap["world"].Add(l)
	return l
}

func (w *World) AddLogicWithSchedule(Name string, F func(dt_ms float64), period_ms float64) *LogicUnit {
	l := w.AddLogic(Name, F)
	runSchedule := NewTimeAccumulator(period_ms)
	l.runSchedule = &runSchedule
	return l
}

func (w *World) RemoveLogic(Name string) {
	if logic, ok := w.worldLogics[Name]; ok {
		w.runtimeSharer.RunnerMap["world"].Remove(logic)
		delete(w.worldLogics, Name)
		w.IDGen.Free(logic.worldID)
	}
}

func (w *World) ActivateAllLogics() {
	w.runtimeSharer.RunnerMap["world"].ActivateAll()
}

func (w *World) DeactivateAllLogics() {
	w.runtimeSharer.RunnerMap["world"].DeactivateAll()
}

func (w *World) ActivateLogic(name string) {
	if logic, ok := w.worldLogics[name]; ok {
		logic.Activate()
	}
}

func (w *World) DeactivateLogic(name string) {
	if logic, ok := w.worldLogics[name]; ok {
		logic.Deactivate()
	}
}

func (w *World) AddFuncs(funcs map[string](func(any) any)) {
	for name, f := range funcs {
		w.funcs.Add(name, f)
	}
}

func (w *World) AddFunc(name string, f func(any) any) {
	w.funcs.Add(name, f)
}

func (w *World) RemoveFunc(name string) {
	w.funcs.Remove(name)
}

func (w *World) GetFunc(name string) func(any) any {
	return w.funcs.funcs[name]
}

func (w *World) HasFunc(name string) bool {
	return w.funcs.Has(name)
}

func (w *World) CreateBlackboard(name string) Blackboard {
	if _, ok := w.Blackboards[name]; !ok {
		w.Blackboards[name] = NewBlackboard(name)
	}
	return w.Blackboards[name]
}

func (w *World) ApplyComponentSet(e *Entity, spec map[ComponentID]any) {
	w.Em.ComponentsTable.ApplyComponentSet(e, spec)
}

func (w *World) EntityHasComponentString(e *Entity, name string) bool {
	b, _ := w.Em.ComponentsTable.ComponentBitArrays[e.ID].GetBit(uint64(w.Em.ComponentsTable.Ixs[w.Em.ComponentsTable.StringsRev[name]]))
	return b
}

func (w *World) EntityHasComponent(e *Entity, name ComponentID) bool {
	b, _ := w.Em.ComponentsTable.ComponentBitArrays[e.ID].GetBit(uint64(w.Em.ComponentsTable.Ixs[name]))
	return b
}

func (w *World) EntityHasComponents(e *Entity, names ...ComponentID) bool {
	has := true
	for _, name := range names {
		b, _ := w.Em.ComponentsTable.ComponentBitArrays[e.ID].GetBit(uint64(w.Em.ComponentsTable.Ixs[name]))
		has = has && b
	}
	return has
}

func (w *World) EntityFilterFromTag(tag string) EntityFilter {
	return EntityFilter{
		Name: tag,
		Predicate: func(e *Entity) bool {
			return w.GetTagList(e, GENERICTAGS_).Has(tag)
		}}
}

func (w *World) EntityFilterFromComponentBitArray(
	name string, q bitarray.BitArray) EntityFilter {
	return EntityFilter{
		Name: name,
		Predicate: func(e *Entity) bool {
			// determine if q = q&b
			// that is, if every set bit of q is set in b
			b := w.Em.ComponentsTable.ComponentBitArrays[e.ID]
			return q.Equals(q.And(b))
		}}
}

func (w *World) EntityFilterFromCanBe(canBe map[string]int) EntityFilter {
	return EntityFilter{
		Name: "canbe",
		Predicate: func(e *Entity) bool {
			for k, v := range canBe {
				if !w.GetIntMap(e, STATE_).ValCanBeSetTo(k, v) {
					return false
				}
			}
			return true
		},
	}
}

// bit of a meta filter:
// matches the closest entity to to that fulfills the given filter
func (w *World) EntityFilterFromClosest(to *Entity, filter EntityFilter) EntityFilter {
	return EntityFilter{
		Name: "closest",
		Predicate: func(e *Entity) bool {
			return e == w.ClosestEntityFilter(
				*w.GetVec2D(to, POSITION_),
				*w.GetVec2D(to, BOX_),
				filter.Predicate)
		},
	}
}

func (w *World) EntityHasTag(e *Entity, tag string) bool {
	return w.GetTagList(e, GENERICTAGS_).Has(tag)
}

func (w *World) EntityHasTags(e *Entity, tags ...string) bool {
	has := true
	for _, t := range tags {
		has = has && w.GetTagList(e, GENERICTAGS_).Has(t)
	}
	return has
}

func (w *World) EntityDistanceFrom(e *Entity, x *Entity) float64 {
	ePos, eBox := w.GetVec2D(e, POSITION_), w.GetVec2D(e, BOX_)
	xPos, xBox := w.GetVec2D(x, POSITION_), w.GetVec2D(x, BOX_)
	return RectDistance(*ePos, *eBox, *xPos, *xBox)
}

func (w *World) EntityDistanceFromRect(e *Entity, pos Vec2D, box Vec2D) float64 {
	ePos, eBox := w.GetVec2D(e, POSITION_), w.GetVec2D(e, BOX_)
	return RectDistance(*ePos, *eBox, pos, box)
}

func (w *World) String() string {
	// TODO: implement
	return "TODO"
}

func (w *World) DumpStats() map[string](map[string]float64) {
	stats := w.runtimeSharer.DumpStats()
	// add total Update() runtime avg
	if w.totalRuntimeAvg_ms != nil {
		stats["__totals"]["World.Update()"] = *w.totalRuntimeAvg_ms
	} else {
		stats["__totals"]["World.Update()"] = 0.0
	}
	return stats
}

func (w *World) DumpStatsString() string {
	stats := w.DumpStats()
	b, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
