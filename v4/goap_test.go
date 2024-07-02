package sameriver

import (
	"fmt"
	"strings"
	"time"

	"testing"

	"github.com/TwiN/go-color"
)

func printWorldState(ws *GOAPWorldState) {
	if ws == nil || len(ws.vals) == 0 {
		Logger.Println("    nil")
		return
	}
	for name, val := range ws.vals {
		Logger.Printf("    %s: %d", name, val)
	}
}

func printGoal(g *GOAPGoal) {
	if g == nil || len(g.vars) == 0 {
		Logger.Println("    nil")
		return
	}
	for varName, interval := range g.vars {
		Logger.Printf("    want %s: [%.0f, %.0f]", varName, interval.A, interval.B)
	}
}

func printGoalRemaining(g *GOAPGoalRemaining) {
	if g.nUnfulfilled == 0 {
		msg := "    satisfied    "
		Logger.Printf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		Logger.Printf(color.InBlackOverGreen(msg))
		Logger.Printf(color.InBlackOverGreen(strings.Repeat(" ", len(msg))))
		return
	}
	for varName, interval := range g.goalLeft {
		msg := fmt.Sprintf("    %s: [%.0f, %.0f]    ", varName, interval.A, interval.B)

		Logger.Printf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))
		Logger.Printf(color.InBold(color.InRedOverBlack(msg)))
		Logger.Printf(color.InBlackOverBlack(strings.Repeat(" ", len(msg))))

	}
}

func printGoalRemainingSurface(s *GOAPGoalRemainingSurface) {
	if s.NUnfulfilled() == 0 {
		Logger.Println("    nil")
	} else {
		for i, tgs := range s.surface {
			if i == len(s.surface)-1 {
				Logger.Printf(color.InBold(color.InRedOverGray("main:")))

			}
			for _, tg := range tgs {
				printGoalRemaining(tg)
			}
		}
	}
}

func printDiffs(diffs map[string]float64) {
	for name, diff := range diffs {
		Logger.Printf("    %s: %.0f", name, diff)
	}
}

func TestGOAPGoalRemaining(t *testing.T) {
	doTest := func(
		g *GOAPGoal,
		ws *GOAPWorldState,
		nRemaining int,
		expectedRemaining []string,
	) {

		remaining := g.remaining(ws)

		Logger.Printf("goal:")
		printGoal(g)
		Logger.Printf("state:")
		printWorldState(ws)
		Logger.Printf("remaining:")
		printGoal(remaining.goal)
		Logger.Printf("diffs:")
		printDiffs(remaining.diffs)
		Logger.Println("-------------------")

		if len(remaining.goalLeft) != nRemaining {
			t.Fatalf("Should have had %d goals remaining, had %d", nRemaining, len(remaining.goalLeft))
		}
		for _, name := range expectedRemaining {
			if diffVal, ok := remaining.diffs[name]; !ok || diffVal == 0 {
				t.Fatalf("Should have had %s in diffs with value != 0", name)
			}
		}
	}

	doTest(
		newGOAPGoal(map[string]int{
			"hasGlove,=": 1,
			"hasAxe,=":   1,
			"atTree,=":   1,
		}),
		NewGOAPWorldState(map[string]int{
			"hasGlove": 0,
			"hasAxe":   1,
			"atTree":   1,
		}),
		1,
		[]string{"hasGlove"},
	)

	doTest(
		newGOAPGoal(map[string]int{
			"hasGlove,=": 1,
			"hasAxe,=":   1,
			"atTree,=":   1,
		}),
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   1,
		}),
		0,
		[]string{},
	)

	doTest(
		newGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		NewGOAPWorldState(map[string]int{
			"drunk": 1,
		}),
		1,
		[]string{"drunk"},
	)
}

func TestGOAPGoalRemainingsOfPath(t *testing.T) {
	w := testingWorld()
	const (
		BOOZEAMOUNT = GENERICTAGS + 1 + iota
	)
	w.RegisterComponents([]any{
		BOOZEAMOUNT, INT, "BOOZEAMOUNT",
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:    Vec2D{0, 0},
			BOX:         Vec2D{0.5, 0.5},
			BOOZEAMOUNT: 0,
		},
	})
	Logger.Println(e)

	p := NewGOAPPlanner(e)

	// we look at hasbooze *for this test* as a component,
	// but it could just as easily be written as an inventory check
	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, BOOZEAMOUNT).(*int)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, BOOZEAMOUNT).(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, BOOZEAMOUNT, &newVal)
			}
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"node": "self", // action occurs at self
		"cost": 1,
		"pres": map[string]int{
			"EACH:hasBooze,>=": 1,
		},
		"effs": map[string]int{
			"drunk,+":    1,
			"hasBooze,-": 1,
		},
	})

	p.AddModalVals(hasBoozeModal)
	p.AddActions(drink)

	start := NewGOAPWorldState(nil)
	start.w = w // this would be done automatically in Plan()
	p.checkModalInto("hasBooze", start)

	goal := map[string]int{
		"drunk,>=": 3,
	}

	// we have already drank 2 times, and want to know what remains
	path := NewGOAPPath([]*GOAPAction{drink.Parametrized(2)})

	Logger.Printf("-------------------------------------------- 1")

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	mainGoalRemaining := path.remainings.surface[len(path.remainings.surface)-1][0]
	if path.remainings.NUnfulfilled() != 2 || len(mainGoalRemaining.goalLeft) != 1 {
		t.Fatal("Remaining was not calculated properly")
	}

	Logger.Printf("-------------------------------------------- 2")

	path = NewGOAPPath([]*GOAPAction{drink.Parametrized(3)})

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	mainGoalRemaining = path.remainings.surface[len(path.remainings.surface)-1][0]
	if path.remainings.NUnfulfilled() != 1 || len(mainGoalRemaining.goalLeft) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}

	Logger.Printf("-------------------------------------------- 3")

	booze := e.GetInt(BOOZEAMOUNT)
	*booze = 3

	p.checkModalInto("hasBooze", start)

	Logger.Printf("start: %v", start.vals)

	p.computeCostAndRemainingsOfPath(path, start, NewGOAPTemporalGoal(goal))

	Logger.Printf("%d unfulfilled", path.remainings.NUnfulfilled())
	printGoalRemainingSurface(path.remainings)
	if path.remainings.NUnfulfilled() != 0 || len(mainGoalRemaining.goalLeft) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}
}

// TODO: fix this test, since chopTree should be like getAxe,
// it automatically *provides* atTree, and it would be
// up to the implementer to make sure while in the choptree action
// we manage the simple state transition from goto->chop
func TestGOAPActionPresFulfilled(t *testing.T) {

	w := testingWorld()
	e := w.Spawn(nil)
	p := NewGOAPPlanner(e)

	doTest := func(ws *GOAPWorldState, a *GOAPAction, expected bool) {
		if p.presFulfilled(a, ws) != expected {
			Logger.Println("world state:")
			printWorldState(ws)
			Logger.Println("action.pres:")
			for _, tg := range a.pres.temporalGoals {
				printGoal(tg)
			}
			t.Fatal("Did not get expected value for action presfulfilled")
		}
	}

	// NOTE: both of these in reality should be modal
	getAxe := NewGOAPAction(map[string]any{
		"name": "getAxe",
		"node": "axe",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"hasAxe,=": 1,
		},
	})
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"node": "self",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"hasBooze,-": 1,
		},
	})

	doDrinkTest := func(has int, expected bool) {
		doTest(
			NewGOAPWorldState(map[string]int{
				"hasBooze": has,
			}),
			drink,
			expected,
		)
	}
	chopTree := NewGOAPAction(map[string]any{
		"name": "chopTree",
		"node": "tree",
		"cost": 1,
		"pres": map[string]int{
			"hasGlove,>": 0,
			"hasAxe,>":   0,
		},
		"effs": map[string]int{
			"treeFelled,=": 1,
		},
	})

	p.AddActions(getAxe, drink, chopTree)

	doDrinkTest(0, false)
	doDrinkTest(1, true)
	doDrinkTest(2, true)

	if !p.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
		})) {
		t.Fatal("chopTree pres should have been fulfilled")
	}

	if p.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   0,
		})) {
		t.Fatal("chopTree pres shouldn't have been fulfilled")
	}
}

func TestGOAPPlanSimple(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	treePos := &Vec2D{19, 19}

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, POSITION).(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, POSITION, &nearTree)
		},
	}
	goToTree := NewGOAPAction(map[string]any{
		"name": "goToTree",
		"node": "tree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})

	goal := map[string]int{
		"atTree,=": 1,
	}

	Logger.Println(*e.GetVec2D(POSITION))

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(atTreeModal)
	p.AddActions(goToTree)

	Logger.Println(p.Plan(ws, goal, 50))

}

func TestGOAPPlanSimpleIota(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE: map[string]int{
				"drunk": 0,
			},
			POSITION: Vec2D{0, 0},
		},
	})

	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, STATE).(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, STATE).(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, STATE, &state)
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"node": "self",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"drunk,+": 1,
		},
	})

	goal := newGOAPGoal(map[string]int{
		"drunk,=": 1,
	})

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(drunkModal)
	p.AddActions(drink)

	Logger.Println(p.Plan(ws, goal, 50))

	goal = newGOAPGoal(map[string]int{
		"drunk,=": 3,
	})
	Logger.Println(p.Plan(ws, goal, 50))

}

func TestGOAPPlanSimpleEnough(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)

	const (
		STATE = GENERICTAGS + 1 + iota
	)

	w.RegisterComponents([]any{
		STATE, INTMAP, "STATE",
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE: map[string]int{
				"drunk": 0,
			},
			POSITION: Vec2D{0, 0},
		},
	})

	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, STATE).(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, STATE).(*IntMap).CopyOf()
			if op == "+" {
				state.m["drunk"] += x
			}
			ws.SetModal(e, STATE, &state)
		},
	}
	drink := NewGOAPAction(map[string]any{
		"name": "drink",
		"node": "self",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"drunk,+": 1,
		},
	})
	purifyOneself := NewGOAPAction(map[string]any{
		"name": "purifyOneself",
		"node": "self",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"rituallyPure,=": 1,
		},
	})

	ws := NewGOAPWorldState(nil)

	p := NewGOAPPlanner(e)
	p.AddModalVals(drunkModal)
	p.AddActions(drink, purifyOneself)

	goal := newGOAPGoal(map[string]int{
		"drunk,=":        10,
		"rituallyPure,=": 1,
	})
	Logger.Println(p.Plan(ws, goal, 50))
}

func TestGOAPPlanClassic(t *testing.T) {
	w := testingWorld()

	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "axe",
		"displayName": "axe",
		"flavourText": "a nice axe for chopping trees",
		"properties": map[string]int{
			"value":     10,
			"sharpness": 2,
		},
		"tags": []string{"tool"},
	})
	items.CreateArchetype(map[string]any{
		"name":        "glove",
		"displayName": "glove",
		"flavourText": "good hand protection",
		"properties": map[string]int{
			"value": 2,
		},
		"tags": []string{"wearable"},
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			BOX:       Vec2D{1, 1},
			INVENTORY: inventories.Create(nil),
		},
	})

	// spawn tree
	w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{6, 6},
			BOX:      Vec2D{1, 1},
		},
		"tags": []string{"tree"},
	})

	// spawn glove and axe
	items.SpawnItemEntity(Vec2D{3, 3}, items.CreateItemSimple("glove"))
	items.SpawnItemEntity(Vec2D{5, 5}, items.CreateItemSimple("axe"))

	// verify there is an entity tagged with glove
	Logger.Println(w.ClosestEntityFilter(Vec2D{0, 0}, Vec2D{1, 1}, func(e *Entity) bool {
		return e.HasTag("glove")
	}))

	hasModal := func(name string, archetype string, tags ...string) GOAPModalVal {
		return GOAPModalVal{
			name: fmt.Sprintf("has%s", name),
			check: func(ws *GOAPWorldState) int {
				inv := ws.GetModal(e, INVENTORY).(*Inventory)
				return inv.CountName(archetype)
			},
			effModalSet: func(ws *GOAPWorldState, op string, x int) {
				inv := ws.GetModal(e, INVENTORY).(*Inventory).CopyOf()
				if op == "-" {
					inv.DebitNTags(x, archetype)
				}
				if op == "=" {
					count := inv.CountTags(tags...)
					if count == 0 {
						inv.Credit(items.CreateStackSimple(x, archetype))
					} else {
						inv.SetCountName(x, archetype)
					}
				}
				if op == "+" {
					count := inv.CountName(archetype)
					if count == 0 {
						inv.Credit(items.CreateStackSimple(x, archetype))
					} else {
						inv.SetCountName(count+x, archetype)
					}
				}
				ws.SetModal(e, INVENTORY, inv)
			},
		}
	}

	hasAxeModal := hasModal("Axe", "axe")
	hasGloveModal := hasModal("Glove", "glove")

	get := func(name string) *GOAPAction {
		return NewGOAPAction(map[string]any{
			"name": fmt.Sprintf("get%s", name),
			"node": strings.ToLower(name),
			"cost": 1,
			"pres": nil,
			"effs": map[string]int{
				fmt.Sprintf("has%s,+", name): 1,
			},
		})
	}

	getAxe := get("Axe")
	getGlove := get("Glove")

	chopTree := NewGOAPAction(map[string]any{
		"name": "chopTree",
		"node": "tree",
		"cost": 1,
		"pres": []any{
			map[string]int{
				"hasGlove,=": 1,
				"hasAxe,=":   1,
			},
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddModalVals(hasGloveModal, hasAxeModal)
	p.AddActions(getAxe, getGlove, chopTree)

	ws := NewGOAPWorldState(nil)

	goal := map[string]int{
		"woodChopped,=": 3,
	}
	t0 := time.Now()
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)
}

func TestGOAPPlanResponsibleFridgeUsage(t *testing.T) {
	w := testingWorld()

	e := w.Spawn(nil)

	openFridge := NewGOAPAction(map[string]any{
		"name": "openFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 1,
		},
	})
	closeFridge := NewGOAPAction(map[string]any{
		"name": "closeFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 0,
		},
	})
	getFoodFromFridge := NewGOAPAction(map[string]any{
		"name": "getFoodFromFridge",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 1,
		},
		"effs": map[string]int{
			"food,+": 1,
		},
	})

	p := NewGOAPPlanner(e)

	p.AddActions(openFridge, getFoodFromFridge, closeFridge)

	ws := NewGOAPWorldState(map[string]int{
		"fridgeOpen": 0,
	})

	goal := map[string]int{
		"fridgeOpen,=": 0,
		"food,=":       1,
	}
	t0 := time.Now()
	plan, ok := p.Plan(ws, goal, 500)
	if !ok {
		t.Fatal("Should've found a solution")
	}
	Logger.Println(color.InGreenOverWhite(GOAPPathToString(plan)))
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	Logger.Printf("Took %f ms to find solution", dt_ms)

}

func TestGOAPPlanFarmer2000(t *testing.T) {

	//
	// world init
	//
	w := testingWorld()
	ps := NewPhysicsSystem()
	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, items, inventories)

	//
	// item system init
	//
	items.CreateArchetype(map[string]any{
		"name":        "yoke",
		"displayName": "a yoke for cattle",
		"flavourText": "one of mankind's greatest inventions... an ancestral gift!",
		"properties": map[string]int{
			"value": 25,
		},
		"tags": []string{"item.agricultural"},
		"entity": map[string]any{
			"sprite": "yoke",
			"box":    [2]float64{0.2, 0.2},
		},
	})

	//
	// spawn entities
	//

	// NOTE: all spawns are on x = 0
	// villager
	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION:  Vec2D{0, 0},
			BOX:       Vec2D{1, 1},
			INVENTORY: inventories.Create(nil),
		},
	})
	// yoke
	yoke := items.CreateItemSimple("yoke")
	items.SpawnItemEntity(Vec2D{0, 5}, yoke)
	// oxen
	spawnOxen := func(positions []Vec2D) (oxen []*Entity) {
		oxen = make([]*Entity, len(positions))
		for i := 0; i < len(positions); i++ {
			oxen[i] = w.Spawn(map[string]any{
				"components": map[ComponentID]any{
					POSITION: positions[i],
					BOX:      Vec2D{3, 2},
					STATE: map[string]int{
						"yoked": 0,
					},
				},
				"tags": []string{"ox"},
			})
		}
		return oxen
	}
	// field
	field := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION: Vec2D{0, 100},
			BOX:      Vec2D{100, 100},
			STATE: map[string]int{
				"tilled": 0,
			},
		},
		"tags": []string{"field"},
	})

	//
	// GOAP actions
	//
	leadOxToField := NewGOAPAction(map[string]any{
		"name":           "leadOxToField",
		"node":           "ox",
		"travelWithNode": true,
		"cost":           1,
		"pres":           nil,
		"effs": map[string]int{
			"ox.in(field),=": 1,
		},
	})
	getYoke := NewGOAPAction(map[string]any{
		"name": "getYoke",
		"node": "yoke",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"self.inventoryHas(yoke),+": 1,
		},
	})
	yokeOxplow := NewGOAPAction(map[string]any{
		"name":       "yokeOxplow",
		"node":       "ox",
		"otherNodes": []string{"field"},
		"cost":       1,
		"pres": []any{
			map[string]int{
				"self.inventoryHas(yoke),=": 1,
			},
		},
		"effs": map[string]int{
			"ox.yoked,=": 1,
		},
	})
	oxplow := NewGOAPAction(map[string]any{
		"name": "oxplow",
		"node": "ox",
		"cost": 1,
		"pres": []any{
			map[string]int{
				"ox.in(field),=": 1,
			},
			map[string]int{
				"ox.yoked,=": 1,
			},
		},
		"effs": map[string]int{
			"field.tilled,=": 1,
		},
	})

	//
	// PLANNER INIT
	//

	// in this test, only the yoke selector gets used as a generic fallback since
	// we don't bind any more specific selector for "yoke" before any Plan() call.
	// *all* of these would be used for more general planning than this specific constrained
	// example where we have a field in mind allowing us to choose a closer ox.
	// so really, this would happen not before the planning as the BindEntitySelectors() call below,
	// but this RegisterGenericEntitySelectors() call would happen on setup of the planner itself
	p := NewGOAPPlanner(e)

	p.AddActions(leadOxToField, getYoke, yokeOxplow, oxplow)

	//
	// bb workplan
	//

	// NOTE: we'd *get* the currently active bb work plan for the field rather than
	// generate it if someone was already doing plant
	tillPlanBB := func() {
		e.SetMind("plan.field", field)
		planField := e.GetMind("plan.field").(*Entity)
		// this would really be a filtering not of all entities but of perception
		closestOxToField := e.World.ClosestEntityFilter(
			*planField.GetVec2D(POSITION),
			*planField.GetVec2D(BOX),
			func(e *Entity) bool {
				return e.HasTag("ox") && e.GetIntMap(STATE).ValCanBeSetTo("yoked", 1)
			})
		if closestOxToField != nil {
			Logger.Printf("closest ox to field: (position: %v)%v", *closestOxToField.GetVec2D(POSITION), closestOxToField)
		}
		e.SetMind("plan.ox", closestOxToField)
	}
	tillPlanBindEntities := func() {
		p.BindEntitySelectors(map[string]any{
			// ox from blackboard plan - the closest to the field
			"ox": "mind.plan.ox",
			// the field from the blackboard plan
			"field": "mind.plan.field",
		})
	}
	mockMakeTillPlan := func() {
		tillPlanBB()
		tillPlanBindEntities()
	}

	//
	// initial world state
	//
	// TODO: this would be a perception system thing - we would get the current
	// state of the world from perception/memory
	ws := NewGOAPWorldState(nil)

	// TODO: this would derive from a utility, not be hardcoded
	goal := map[string]int{
		"field.tilled,=": 1,
	}

	runAPlan := func(expect bool) (dt_ms float64) {
		mockMakeTillPlan()
		t0 := time.Now()
		plan, ok := p.Plan(ws, goal, 500)
		if ok != expect {
			t.Fatalf("should have had ok: %t", expect)
		}
		if ok {
			Logger.Println(color.InGreen(plan.String()))
		}
		return float64(time.Since(t0).Nanoseconds()) / 1.0e6
	}

	// first run with no oxen
	Logger.Println("No oxen")
	dt_ms := runAPlan(false)
	Logger.Printf("Took %f ms to fail", dt_ms)

	// second run with oxen
	// spawn them (note: one is in the field already)
	oxen := spawnOxen([]Vec2D{Vec2D{0, 100}, Vec2D{0, 20}, Vec2D{0, -100}})
	dt_ms = runAPlan(true)
	Logger.Printf("Took %f ms to find solution", dt_ms)

	// third run with oxen all out of the field by despawning the one we found in
	w.Despawn(e.GetMind("plan.ox").(*Entity))
	Logger.Println("All oxen are outside field")
	dt_ms = runAPlan(true)
	Logger.Printf("Took %f ms to find solution", dt_ms)

	// we will want to use {0, 20}, so let's make it unyokable
	const BECOME_UNGOVERNABLE = true
	if BECOME_UNGOVERNABLE {
		oxen[1].GetIntMap(STATE).SetValidInterval("yoked", 0, 0)
		// inside runAPlan, when we plan the bb, the bound selector should check
		// for yokable on state intmap
		dt_ms = runAPlan(true)
		Logger.Printf("Took %f ms to find solution", dt_ms)

		// all oxen either despawned or unyokable
		oxen[2].GetIntMap(STATE).SetValidInterval("yoked", 0, 0)
		Logger.Println("No *yokable* oxen")
		dt_ms = runAPlan(false)
		Logger.Printf("Took %f ms to fail", dt_ms)

		// restore the humility of these brave beasts, make them fit to work!
		oxen[1].GetIntMap(STATE).SetValidInterval("yoked", 0, 1)
		oxen[2].GetIntMap(STATE).SetValidInterval("yoked", 0, 1)
		Logger.Println("Pick the good ox!")
		dt_ms = runAPlan(true)
		if !e.GetMind("plan.ox").(*Entity).GetVec2D(POSITION).Equals(Vec2D{0, 20}) {
			t.Fatalf("Didn't grandpappy learn ya right? Always pick the best ox!!! Ya done picked %v", e.GetMind("plan.ox").(*Entity))
		}
		Logger.Printf("Took %f ms to find solution", dt_ms)

	}
}
