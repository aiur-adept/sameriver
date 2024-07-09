package sameriver

func testingWorld() *World {
	w := NewWorld(map[string]any{
		"width":  1024,
		"height": 1024,
	})
	return w
}

func testingWorldWithAllLogicTypes() (*World, *testSystem, *int) {
	w := testingWorld()
	// add system
	ts := newTestSystem()
	w.RegisterSystems(ts)
	// add world logic
	worldUpdates := 0
	name := "logic"
	w.AddLogic(name, func(dt_ms float64) { worldUpdates += 1 })
	w.ActivateLogic(name)
	return w, ts, &worldUpdates
}
