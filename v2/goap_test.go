package sameriver

import (
	"fmt"
	//	"time"

	"testing"
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
	if g == nil || len(g.goals) == 0 {
		Logger.Println("    nil")
		return
	}
	for varName, interval := range g.goals {
		Logger.Printf("    want %s: [%.0f, %.0f]", varName, interval.A, interval.B)
	}
}

func printDiffs(diffs map[string]float64) {
	for name, diff := range diffs {
		Logger.Printf("    %s: %.0f", name, diff)
	}
}

func printPath(p *GOAPPath) {
	Logger.Printf(GOAPPathToString(p))
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

		if len(remaining.goal.goals) != nRemaining {
			t.Fatal(fmt.Sprintf("Should have had %d goals remaining, had %d", nRemaining, len(remaining.goal.goals)))
		}
		for _, name := range expectedRemaining {
			if diffVal, ok := remaining.diffs[name]; !ok || diffVal == 0 {
				t.Fatal(fmt.Sprintf("Should have had %s in diffs with value != 0", name))
			}
		}
	}

	doTest(
		NewGOAPGoal(map[string]int{
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
		NewGOAPGoal(map[string]int{
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
		NewGOAPGoal(map[string]int{
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
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	w.RegisterComponents([]string{"Int,BoozeAmount"})

	e := testingSpawnPhysics(w)

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    1,
			"hasBooze,-": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(drink)

	start := NewGOAPWorldState(nil)
	p.eval.PopulateModalStartState(start)

	goal := NewGOAPGoal(map[string]int{
		"drunk,>=": 3,
	})

	path := NewGOAPPath([]*GOAPAction{drink, drink}, GOAP_PATH_PREPEND)

	remaining := p.eval.remainingsOfPath(path, start, goal)

	Logger.Printf("%d unfulfilled", remaining.nUnfulfilled)
	printGoal(remaining.main.goal)
	for _, pre := range remaining.pres {
		printGoal(pre.goal)
	}
	if remaining.nUnfulfilled != 3 || len(remaining.main.goal.goals) != 1 {
		t.Fatal("Remaining was not calculated properly")
	}

	path = NewGOAPPath([]*GOAPAction{drink, drink, drink}, GOAP_PATH_PREPEND)

	remaining = p.eval.remainingsOfPath(path, start, goal)

	Logger.Printf("%d unfulfilled", remaining.nUnfulfilled)
	printGoal(remaining.main.goal)
	for _, pre := range remaining.pres {
		printGoal(pre.goal)
	}
	if remaining.nUnfulfilled != 3 || len(remaining.main.goal.goals) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}

	booze := e.GetInt("BoozeAmount")
	*booze = 3
	p.eval.PopulateModalStartState(start)

	remaining = p.eval.remainingsOfPath(path, start, goal)

	Logger.Printf("%d unfulfilled", remaining.nUnfulfilled)
	printGoal(remaining.main.goal)
	for _, pre := range remaining.pres {
		printGoal(pre.goal)
	}

	if remaining.nUnfulfilled != 0 || len(remaining.main.goal.goals) != 0 {
		t.Fatal("Remaining was not calculated properly")
	}
}

func TestGOAPRemainingIsLess(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	w.RegisterComponents([]string{"Int,BoozeAmount"})

	e := testingSpawnPhysics(w)

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			debugGOAPPrintf("                checked hasBooze: %d", *amount)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
			if op == "+" {
				debugGOAPPrintf("                adding to hasBooze: +%d", x)
				newVal := *amount + x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
		},
	}
	getBooze := NewGOAPAction(map[string]interface{}{
		"name": "getBooze",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"hasBooze,+": 1,
		},
	})
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    1,
			"hasBooze,-": 1,
		},
	})
	purifyOneself := NewGOAPAction(map[string]interface{}{
		"name": "purifyOneself",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,=": 0,
		},
		"effs": map[string]int{
			"rituallyPure,=": 1,
		},
	})
	chopTree := NewGOAPAction(map[string]interface{}{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"hasGlove,>": 0,
			"hasAxe,>":   0,
			"atTree,=":   1,
		},
		"effs": map[string]int{
			"treeFelled,=": 1,
		},
	})
	openFridge := NewGOAPAction(map[string]interface{}{
		"name": "openFridge",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 0,
		},
		"effs": map[string]int{
			"fridgeOpen,=": 1,
		},
	})
	closeFridge := NewGOAPAction(map[string]interface{}{
		"name": "closeFridge",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 1,
		},
		"effs": map[string]int{
			"fridgeOpen,=": 0,
		},
	})
	getFoodFromFridge := NewGOAPAction(map[string]interface{}{
		"name": "getFoodFromFridge",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 1,
		},
		"effs": map[string]int{
			"hasFood,+": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(getBooze, drink, purifyOneself, chopTree, openFridge, getFoodFromFridge, closeFridge)

	doTest := func(g *GOAPGoal, start *GOAPWorldState, before, after *GOAPPath, expect bool) {
		Logger.Println("=================================================================")
		Logger.Printf("Before: %s", GOAPPathToString(before))
		Logger.Printf("After: %s", GOAPPathToString(after))
		beforeRemaining := p.eval.remainingsOfPath(before, start, g)
		afterRemaining := p.eval.remainingsOfPath(after, start, g)
		Logger.Println("computing isCloser()...")
		less := afterRemaining.isCloser(beforeRemaining)
		if less != expect {
			Logger.Println("!!!")
			Logger.Println("!!!")
			Logger.Println("!!!")
			Logger.Println("Didn't get expected result for path after remainingIsLess than path before")
			Logger.Println("!!!")
			Logger.Println("!!!")
			Logger.Println("!!!")
			t.Fatal("Didn't get expected result for remainingIsLess()")
		}
	}

	start := NewGOAPWorldState(map[string]int{
		"drunk": 1,
	})
	p.eval.PopulateModalStartState(start)

	before := NewGOAPPath([]*GOAPAction{drink}, GOAP_PATH_PREPEND)
	after := before.prepended(drink)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		start,
		before,
		after,
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		start,
		before,
		before,
		false,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		start,
		before,
		before.appended(purifyOneself),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=":       3,
			"rituallyPure,=": 1,
		}),
		start,
		before,
		before.prepended(purifyOneself),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		start,
		before,
		before.prepended(chopTree),
		false,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		start,
		before,
		before.prepended(getBooze),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"drunk,>=": 3,
		}),
		start,
		NewGOAPPath([]*GOAPAction{drink, drink, drink}, GOAP_PATH_PREPEND),
		NewGOAPPath([]*GOAPAction{drink, drink, drink, drink}, GOAP_PATH_PREPEND),
		false,
	)

	start.vals["fridgeOpen"] = 0

	doTest(
		NewGOAPGoal(map[string]int{
			"hasFood,>=":   1,
			"fridgeOpen,=": 0,
		}),
		start,
		NewGOAPPath([]*GOAPAction{openFridge, getFoodFromFridge}, GOAP_PATH_PREPEND),
		NewGOAPPath([]*GOAPAction{openFridge, getFoodFromFridge, closeFridge}, GOAP_PATH_APPEND),
		true,
	)

	doTest(
		NewGOAPGoal(map[string]int{
			"hasFood,>=":   1,
			"fridgeOpen,=": 0,
		}),
		start,
		NewGOAPPath([]*GOAPAction{openFridge, getFoodFromFridge}, GOAP_PATH_PREPEND),
		NewGOAPPath([]*GOAPAction{closeFridge, openFridge, getFoodFromFridge}, GOAP_PATH_PREPEND),
		false,
	)

}

func TestGOAPActionPresFulfilled(t *testing.T) {

	eval := NewGOAPEvaluator()

	doTest := func(ws *GOAPWorldState, a *GOAPAction, expected bool) {
		if eval.presFulfilled(a, ws) != expected {
			Logger.Println("world state:")
			printWorldState(ws)
			Logger.Println("action.pres:")
			printGoal(a.pres)
			t.Fatal("Did not get expected value for action presfulfilled")
		}
	}

	// NOTE: both of these in reality should be modal
	goToAxe := NewGOAPAction(map[string]interface{}{
		"name": "goToAxe",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atAxe,=": 1,
		},
	})
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
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
	chopTree := NewGOAPAction(map[string]interface{}{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"hasGlove,>": 0,
			"hasAxe,>":   0,
			"atTree,=":   1,
		},
		"effs": map[string]int{
			"treeFelled,=": 1,
		},
	})

	eval.AddActions(goToAxe, drink, chopTree)

	doDrinkTest(0, false)
	doDrinkTest(1, true)
	doDrinkTest(2, true)

	if !eval.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   1,
		})) {
		t.Fatal("chopTree pres should have been fulfilled")
	}

	if eval.presFulfilled(
		chopTree,
		NewGOAPWorldState(map[string]int{
			"hasGlove": 1,
			"hasAxe":   1,
			"atTree":   0,
		})) {
		t.Fatal("chopTree pres shouldn't have been fulfilled")
	}
}

func TestGOAPActionModalVal(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e := testingSpawnPhysics(w)

	ws := NewGOAPWorldState(nil)
	treePos := &Vec2D{11, 11}

	eval := NewGOAPEvaluator()

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearTree)
		},
	}
	oceanPos := &Vec2D{500, 0}
	atOceanModal := GOAPModalVal{
		name: "atOcean",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*oceanPos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearOcean := oceanPos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearOcean)
		},
	}
	goToTree := NewGOAPAction(map[string]interface{}{
		"name": "goToTree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})
	chopTree := NewGOAPAction(map[string]interface{}{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"atTree,=": 1,
			"hasAxe,>": 0,
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})
	hugTree := NewGOAPAction(map[string]interface{}{
		"name": "hugTree",
		"cost": 1,
		"pres": map[string]int{
			"atTree,=": 1,
		},
		"effs": map[string]int{
			"connectionToNature,+": 2,
		},
	})
	goToOcean := NewGOAPAction(map[string]interface{}{
		"name": "goToOcean",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atOcean,=": 1,
		},
	})

	eval.AddModalVals(atTreeModal, atOceanModal)
	eval.AddActions(goToTree, hugTree, chopTree, goToOcean)

	//
	// test presfulfilled
	//
	*e.GetVec2D("Position") = *treePos

	if !eval.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 1, satisfying atTree,=: 1")
	}

	*e.GetVec2D("Position") = Vec2D{-100, -100}

	if eval.presFulfilled(hugTree, ws) {
		t.Fatal("check result of atTreeModal should have returned 0, failing to satisfy atTree,=: 1")
	}

	badWS := NewGOAPWorldState(map[string]int{
		"atTree": 0,
	})

	*e.GetVec2D("Position") = *treePos

	if !eval.presFulfilled(hugTree, badWS) {
		t.Fatal("regardless of what worldstate says, modal pre should decide and should've been true based on entity position = tree position")
	}

	axeWS := NewGOAPWorldState(map[string]int{
		"hasAxe": 1,
	})
	if !eval.presFulfilled(chopTree, axeWS) {
		t.Fatal("mix of modal and basic world state vals should fulfill pre")
	}

	//
	// test applyAction
	//

	g := NewGOAPGoal(map[string]int{
		"atTree,=": 1,
	})
	appliedState := eval.applyAction(goToTree, NewGOAPWorldState(nil))
	remaining := g.remaining(appliedState)
	Logger.Println("goal:")
	printGoal(g)
	Logger.Println("state after applying goToTree:")
	printWorldState(appliedState)
	if appliedState.vals["atTree"] != 1 {
		t.Fatal("atTree should've been 1 after goToTree")
	}
	Logger.Println("goal remaining:")
	printGoal(remaining.goal)
	if len(remaining.goal.goals) != 0 {
		t.Fatal("Goal should have been satisfied")
	}
	Logger.Println("diffs:")
	printDiffs(remaining.diffs)

	g2 := NewGOAPGoal(map[string]int{
		"atTree,=": 1,
		"drunk,>=": 10,
	})
	remaining = g2.remaining(appliedState)
	if len(remaining.goal.goals) != 1 {
		t.Fatal("drunk goal should be unfulfilled by atTree state")
	}

	//
	// test modal effect of applyAction
	//

	*e.GetVec2D("Position") = Vec2D{-100, -100}
	atTreeApplied := eval.applyAction(goToTree, NewGOAPWorldState(nil))
	Logger.Println("state after applying modal action eff of atTree:")
	printWorldState(atTreeApplied)
	if val, ok := atTreeApplied.vals["atTree"]; !ok || val != 1 {
		t.Fatal("Modal action eff should've set atTree=1")
	}
	Logger.Println("modal position of entity after modal action eff of atTree:")
	posAfter := atTreeApplied.GetModal(e, "Position").(*Vec2D)
	Logger.Printf("[%f, %f]", posAfter.X, posAfter.Y)

	//
	// test modal pre after modal set
	//

	*e.GetVec2D("Position") = Vec2D{-100, -100}
	atOceanApplied := eval.applyAction(goToOcean, NewGOAPWorldState(nil))
	Logger.Println("state after applying modal action eff of atOcean:")
	printWorldState(atOceanApplied)

	if eval.presFulfilled(hugTree, atOceanApplied) {
		t.Fatal("atTree modal pre of hugTree should fail when modal position is set at ocean")
	}

	nowGoToTreeApplied := eval.applyAction(goToTree, atOceanApplied)
	Logger.Println("state after goToOcean->goToTree:")
	printWorldState(nowGoToTreeApplied)
	if nowGoToTreeApplied.vals["atOcean"] != 0 {
		t.Fatal("Should've had atOcean=0 after goToTree")
	}
	if nowGoToTreeApplied.vals["atTree"] != 1 {
		t.Fatal("Should've had atTree=1 after goToTree")
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
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearTree)
		},
	}
	goToTree := NewGOAPAction(map[string]interface{}{
		"name": "goToTree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})

	goal := NewGOAPGoal(map[string]int{
		"atTree,=": 1,
	})

	Logger.Println(*e.GetVec2D("Position"))

	ws := NewGOAPWorldState(nil)

	planner := NewGOAPPlanner(e)
	planner.eval.AddModalVals(atTreeModal)
	planner.eval.AddActions(goToTree)

	Logger.Println(planner.Plan(ws, goal, 50))

}

func TestGOAPPlanSimpleIota(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)

	w.RegisterComponents([]string{"IntMap,State"})

	e := w.Spawn(map[string]any{
		"components": map[string]any{
			"IntMap,State": map[string]int{
				"drunk": 0,
			},
		},
	})

	drunkModal := GOAPModalVal{
		name: "drunk",
		check: func(ws *GOAPWorldState) int {
			state := ws.GetModal(e, "State").(*IntMap)
			return state.m["drunk"]
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			state := ws.GetModal(e, "State").(*IntMap).CopyOf()
			if op == "+" {
				Logger.Println("                modal setting")
				state.m["drunk"] += x
			}
			ws.SetModal(e, "State", &state)
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"drunk,+": 1,
		},
	})

	goal := NewGOAPGoal(map[string]int{
		"drunk,=": 1,
	})

	ws := NewGOAPWorldState(nil)

	planner := NewGOAPPlanner(e)
	planner.eval.AddModalVals(drunkModal)
	planner.eval.AddActions(drink)

	Logger.Println(planner.Plan(ws, goal, 50))

	goal = NewGOAPGoal(map[string]int{
		"drunk,=": 3,
	})
	Logger.Println(planner.Plan(ws, goal, 50))

}

/*
func TestGOAPPlannerDeepen(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			// simulate infinite booze supply
			return 1
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    1,
			"hasBooze,-": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(drink)

	start := NewGOAPWorldState(nil)
	p.eval.PopulateModalStartState(start)
	goal := NewGOAPGoal(map[string]int{
		"drunk,>=": 3,
	})
	backtrackRoot := &GOAPPQueueItem{
		path:          []*GOAPAction{},
		presRemaining: make(map[string]*GOAPGoal),
		remaining:     goal,
		nUnfulfilled:  len(goal.goals),
		endState:      start,
		cost:          0,
		index:         -1, // going to be set by Push()
	}

	newPaths := p.deepen(start, backtrackRoot)
	if len(newPaths) != 1 {
		t.Fatal("Should have found 1 path")
	}
	if len(newPaths[0].remaining.goals) == 0 {
		t.Fatal("Should not have fulfilled the goal")
	}

	start = NewGOAPWorldState(nil)
	goal = NewGOAPGoal(map[string]int{
		"drunk,=": 1,
	})
	newPaths = p.deepen(start, backtrackRoot)
	if len(newPaths) != 1 && len(newPaths[0].remaining.goals) == 0 {
		t.Fatal("Should have found a path (drink) and had goal fulfilled")
	}
}

func TestGOAPPlannerBasic(t *testing.T) {

	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			// simulate infinite booze supply
			return 1
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    2,
			"hasBooze,-": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(drink)

	start := NewGOAPWorldState(nil)
	goal := NewGOAPGoal(map[string]int{
		"drunk,>=": 5,
	})
	solution, ok := p.Plan(start, goal, 50)

	if ok {
		Logger.Println(GOAPPathToString(solution))
	} else {
		Logger.Println("Didn't find a solution.")
	}
}

func TestGOAPPlannerAlanWatts(t *testing.T) {

	w := testingWorld()
	w.RegisterComponents([]string{"Int,BoozeAmount"})

	e, _ := w.Spawn([]string{}, MakeComponentSet(map[string]interface{}{
		"Int,BoozeAmount": 10,
	}))

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    2,
			"hasBooze,-": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(drink)

	start := NewGOAPWorldState(nil)
	goal := NewGOAPGoal(map[string]int{
		"drunk,>=": 20,
	})
	solution, ok := p.Plan(start, goal, 50)
	Logger.Println("Alan Watt's plan:")
	Logger.Println(GOAPPathToString(solution))

	if !ok {
		t.Fatal("Should have found a solution with ten booze")
	}

	*e.GetInt("BoozeAmount") = 5
	solution, ok = p.Plan(start, goal, 50)

	if ok {
		t.Fatal("Should not have found a plan with five booze")
	}
}

func TestGOAPPlannerPurifyOneself(t *testing.T) {

	w := testingWorld()
	w.RegisterComponents([]string{"Int,BoozeAmount"})

	e, _ := w.Spawn([]string{}, MakeComponentSet(map[string]interface{}{
		"Int,BoozeAmount": 10,
	}))

	p := NewGOAPPlanner(e)

	hasBoozeModal := GOAPModalVal{
		name: "hasBooze",
		check: func(ws *GOAPWorldState) int {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			return *amount
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			amount := ws.GetModal(e, "BoozeAmount").(*int)
			if op == "-" {
				newVal := *amount - x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
			if op == "=" {
				newVal := x
				ws.SetModal(e, "BoozeAmount", &newVal)
			}
		},
	}
	drink := NewGOAPAction(map[string]interface{}{
		"name": "drink",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"drunk,+":    2,
			"hasBooze,-": 1,
		},
	})
	dropAllBooze := NewGOAPAction(map[string]interface{}{
		"name": "dropAllBooze",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,>": 0,
		},
		"effs": map[string]int{
			"hasBooze,=": 0,
		},
	})
	purifyOneself := NewGOAPAction(map[string]interface{}{
		"name": "purifyOneself",
		"cost": 1,
		"pres": map[string]int{
			"hasBooze,=": 0,
		},
		"effs": map[string]int{
			"rituallyPure,=": 1,
		},
	})
	enterTemple := NewGOAPAction(map[string]interface{}{
		"name": "enterTemple",
		"cost": 1,
		"pres": map[string]int{
			"rituallyPure,=": 1,
		},
		"effs": map[string]int{
			"templeEntered,=": 1,
		},
	})

	p.eval.AddModalVals(hasBoozeModal)
	p.eval.AddActions(drink, dropAllBooze, purifyOneself, enterTemple)
	start := NewGOAPWorldState(nil)
	p.eval.PopulateModalStartState(start)

	goal := NewGOAPGoal(map[string]int{
		"drunk,>=":        3,
		"templeEntered,=": 1,
	})
	solution, ok := p.Plan(start, goal, 50)
	Logger.Println("Alan Watt's plan:")
	Logger.Println(GOAPPathToString(solution))

	if !ok {
		t.Fatal("Should have found a solution")
	}
}

func TestGOAPPlannerResponsibleFridgeUsage(t *testing.T) {

	w := testingWorld()
	w.RegisterComponents([]string{"Int,FoodAmount"})

	e, _ := testingSpawnSimple(w)

	p := NewGOAPPlanner(e)

	openFridge := NewGOAPAction(map[string]interface{}{
		"name": "openFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 1,
		},
	})
	getFood := NewGOAPAction(map[string]interface{}{
		"name": "getFood",
		"cost": 1,
		"pres": map[string]int{
			"fridgeOpen,=": 1,
		},
		"effs": map[string]int{
			"hasFood,+": 1,
		},
	})
	closeFridge := NewGOAPAction(map[string]interface{}{
		"name": "closeFridge",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"fridgeOpen,=": 0,
		},
	})

	p.eval.AddActions(openFridge, getFood, closeFridge)

	start := NewGOAPWorldState(map[string]int{
		"hasFood":    0,
		"fridgeOpen": 0,
	})
	goal := NewGOAPGoal(map[string]int{
		"hasFood,>":    0,
		"fridgeOpen,=": 0,
	})
	solution, ok := p.Plan(start, goal, 50)
	Logger.Println("Responsible fridge use:")
	Logger.Println(GOAPPathToString(solution))

	if !ok {
		t.Fatal("Should have found a solution")
	}
}

func TestGOAPPlannerWoodsmanByTheSea(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps)
	e, _ := testingSpawnPhysics(w)

	treePos := &Vec2D{11, 11}
	axePos := &Vec2D{-20, 20}
	glovePos := &Vec2D{-20, 5}
	oceanPos := &Vec2D{0, -10}

	atTreeModal := GOAPModalVal{
		name: "atTree",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*treePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearTree := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearTree)
		},
	}
	atOceanModal := GOAPModalVal{
		name: "atOcean",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*oceanPos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearOcean := oceanPos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearOcean)
		},
	}
	atAxeModal := GOAPModalVal{
		name: "atAxe",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*axePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearAxe := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearAxe)
		},
	}
	atGloveModal := GOAPModalVal{
		name: "atGlove",
		check: func(ws *GOAPWorldState) int {
			ourPos := ws.GetModal(e, "Position").(*Vec2D)
			_, _, d := ourPos.Distance(*glovePos)
			if d < 2 {
				return 1
			} else {
				return 0
			}
		},
		effModalSet: func(ws *GOAPWorldState, op string, x int) {
			nearGlove := treePos.Add(Vec2D{1, 0})
			ws.SetModal(e, "Position", &nearGlove)
		},
	}
	goToTree := NewGOAPAction(map[string]interface{}{
		"name": "goToTree",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atTree,=": 1,
		},
	})
	goToGlove := NewGOAPAction(map[string]interface{}{
		"name": "goToGlove",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atGlove,=": 1,
		},
	})
	goToAxe := NewGOAPAction(map[string]interface{}{
		"name": "goToAxe",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atAxe,=": 1,
		},
	})
	getGlove := NewGOAPAction(map[string]interface{}{
		"name": "getGlove",
		"cost": 1,
		"pres": map[string]int{
			"atGlove,=": 1,
		},
		"effs": map[string]int{
			"hasGlove,+": 1,
		},
	})
	getAxe := NewGOAPAction(map[string]interface{}{
		"name": "getAxe",
		"cost": 1,
		"pres": map[string]int{
			"atAxe,=": 1,
		},
		"effs": map[string]int{
			"hasAxe,+": 1,
		},
	})
	chopTree := NewGOAPAction(map[string]interface{}{
		"name": "chopTree",
		"cost": 1,
		"pres": map[string]int{
			"atTree,=":   1,
			"hasAxe,>":   0,
			"hasGlove,>": 0,
		},
		"effs": map[string]int{
			"woodChopped,+": 1,
		},
	})
	goToOcean := NewGOAPAction(map[string]interface{}{
		"name": "goToOcean",
		"cost": 1,
		"pres": nil,
		"effs": map[string]int{
			"atOcean,=": 1,
		},
	})

	p := NewGOAPPlanner(e)
	p.eval.AddModalVals(atTreeModal, atOceanModal, atAxeModal, atGloveModal)
	p.eval.AddActions(goToTree, chopTree, goToOcean, goToGlove, goToAxe, getGlove, getAxe)

	*e.GetVec2D("Position") = *oceanPos
	start := NewGOAPWorldState(nil)

	p.eval.PopulateModalStartState(start)

	start = NewGOAPWorldState(map[string]int{
		"woodChopped": 3,
	})
	p.eval.PopulateModalStartState(start)
	goal := NewGOAPGoal(map[string]int{
		"woodChopped,=": 3,
	})
	solution, ok := p.Plan(start, goal, 50)
	Logger.Println(GOAPPathToString(solution))

	if !ok {
		t.Fatal("Should have found a solution")
	}

}

*/
