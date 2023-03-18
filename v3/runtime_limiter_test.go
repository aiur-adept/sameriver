package sameriver

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestRuntimeLimiterAdd(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			name:        name,
			worldID:     i,
			f:           func(dt_ms float64) {},
			active:      true,
			runSchedule: nil}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.worldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}

func TestRuntimeLimiterAddDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should have panic'd")
		}
	}()
	r := NewRuntimeLimiter()
	logic := &LogicUnit{
		name:        "logic",
		worldID:     0,
		f:           func(dt_ms float64) {},
		active:      true,
		runSchedule: nil}
	r.Add(logic)
	r.Add(logic)
	t.Fatal("should have panic'd")
}

func TestRuntimeLimiterRun(t *testing.T) {
	r := NewRuntimeLimiter()
	x := 0
	name := "l1"
	r.Add(&LogicUnit{
		name:        name,
		worldID:     0,
		f:           func(dt_ms float64) { x += 1 },
		active:      true,
		runSchedule: nil})
	for i := 0; i < 32; i++ {
		r.Run(FRAME_MS, false)
		time.Sleep(FRAME_DURATION)
	}
	Logger.Println(x)
	if x != 32 {
		t.Fatal("didn't run logic")
	}
	if !r.Finished() {
		t.Fatal("should have returned finished = true when running sole " +
			"logic within time limit")
	}
}

func TestRuntimeLimiterOverrun(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:        "logic",
		worldID:     0,
		f:           func(dt_ms float64) { time.Sleep(150 * time.Millisecond) },
		active:      true,
		runSchedule: nil})
	remaining_ms := r.Run(100, false)
	if remaining_ms > 0 {
		t.Fatal("overrun time not calculated properly")
	}
	if !r.overrun {
		t.Fatal("didn't set overrun flag")
	}
}

func TestRuntimeLimiterUnderrun(t *testing.T) {
	r := NewRuntimeLimiter()
	r.Add(&LogicUnit{
		name:        "logic",
		worldID:     0,
		f:           func(dt_ms float64) { time.Sleep(100 * time.Millisecond) },
		active:      true,
		runSchedule: nil})
	remaining_ms := r.Run(300, false)
	if !(remaining_ms > 0 && remaining_ms <= 200) {
		t.Fatal("underrun time not calculated properly")
	}
}

func TestRuntimeLimiterLimiting(t *testing.T) {
	r := NewRuntimeLimiter()
	fastRan := false
	r.Add(&LogicUnit{
		name:        "logic-slow",
		worldID:     0,
		f:           func(dt_ms float64) { time.Sleep(10 * time.Millisecond) },
		active:      true,
		runSchedule: nil})
	r.Add(&LogicUnit{
		name:        "logic-slow",
		worldID:     1,
		f:           func(dt_ms float64) { fastRan = true },
		active:      true,
		runSchedule: nil})
	r.Run(2, false)
	if fastRan {
		t.Fatal("continued running logic despite using up allowed milliseconds")
	}
}

// simulate very overloaded conditions
func TestRuntimeLimiterLoad(t *testing.T) {
	r := NewRuntimeLimiter()

	// time.Sleep doesn't like amounts < 1ms, so we scale up the time axis
	// to allow proper sleeping
	allowance_ms := 10.0
	N_EPSILON := 10
	epsilon_factor := 0.1
	N_HEAVY := 5
	heavy_factor := 0.7

	totalLoad := float64(N_EPSILON)*epsilon_factor + float64(N_HEAVY)*heavy_factor

	Logger.Printf("allowance_ms: %f", allowance_ms)
	Logger.Printf("N_EPSILON: %v", N_EPSILON)
	Logger.Printf("epsilon_factor: %v", epsilon_factor)
	Logger.Printf("N_HEAVY: %v", N_HEAVY)
	Logger.Printf("heavy_factor: %v", heavy_factor)
	Logger.Printf("total load: %f", totalLoad)

	frame := -1
	seq := make([][]string, 0)
	markRan := func(name string) {
		seq[frame] = append(seq[frame], name)
	}
	pushFrame := func() {
		frame++
		seq = append(seq, make([]string, 0))
		Logger.Printf("------------------ frame %d ----------------------", frame)
	}
	printFrame := func() {
		b, _ := json.MarshalIndent(seq[frame], "", "\t")
		Logger.Printf(string(b))
		for _, l := range r.logicUnits {
			Logger.Printf("%s: h%d", l.name, l.hotness)
		}
	}

	for i := 0; i < N_EPSILON; i++ {
		name := fmt.Sprintf("epsilon-%d", i)
		r.Add(&LogicUnit{
			name:    name,
			worldID: i,
			f: func(dt_ms float64) {
				Logger.Printf("epsilon func sleeping %f ms", float64(time.Duration(epsilon_factor*allowance_ms*1e6)*time.Nanosecond)/1e6)
				t0 := time.Now()
				time.Sleep(time.Duration(epsilon_factor*allowance_ms*1e6) * time.Nanosecond)
				Logger.Printf("Sleep took %f ms", float64(time.Since(t0).Nanoseconds())/1e6)
				markRan(name)
			},
			active:      true,
			runSchedule: nil})
	}

	x := 0
	for i := 0; i < N_HEAVY; i++ {
		name := fmt.Sprintf("heavy-%d", i)
		r.Add(&LogicUnit{
			name:    name,
			worldID: N_EPSILON + 1 + i,
			f: func(dt_ms float64) {
				x += 1
				markRan(name)
				time.Sleep(time.Duration(heavy_factor*allowance_ms) * time.Millisecond)
			},
			active:      true,
			runSchedule: nil})
	}

	runFrame := func(allowanceScale float64) {
		t0 := time.Now()
		if allowanceScale != 1 {
			Logger.Printf("<CONSTRICTED FRAME>")
		}
		pushFrame()
		r.Run(allowanceScale*allowance_ms, false)
		elapsed := float64(time.Since(t0).Nanoseconds()) / 1.0e6
		printFrame()
		Logger.Printf("            elapsed: %f ms", elapsed)
		if allowanceScale != 1 {
			Logger.Printf("</CONSTRICTED FRAME>")
		}
	}

	// since it's never run before, running the logic will set its estimate
	runFrame(1.0)

	/*
		TODO: we need better math here to account for N_EPSILON and N_HEAVY

		heavyFirstFrame := int(math.Ceil((1.0 - (float64(N_EPSILON) * 0.1)) / heavy_factor))
		Logger.Printf("Expecting %d heavies to have run in first frame", heavyFirstFrame)
		if x != heavyFirstFrame {
			t.Fatalf("Should've run %d heavies on first frame", heavyFirstFrame)
		}

		// now we try to run it again, but give it no time to run (exceeds estimate)
		runFrame(0.1)

		if x != heavyFirstFrame+1 {
			t.Fatal("should've ran one more heavy in a constricted frame")
		}
	*/

	// what happens after the second heavy has run when we give only a tenth of the time?
	runFrame(0.1)
	runFrame(0.1)
	runFrame(0.1)

	// run a bunch of frames
	for i := 0; i < 12; i++ {
		runFrame(1.0)
	}
	// run a bunch of constricted frames
	for i := 0; i < 12; i++ {
		runFrame(0.333)
	}
}

func TestRuntimeLimiterRemove(t *testing.T) {
	r := NewRuntimeLimiter()
	// test that we can remove a logic which doens't exist idempotently
	if r.Remove(nil) != false {
		t.Fatal("somehow removed a logic which doesn't exist")
	}
	x := 0
	name := "l1"
	logic := &LogicUnit{
		name:        name,
		worldID:     0,
		f:           func(dt_ms float64) { x += 1 },
		active:      true,
		runSchedule: nil}
	r.Add(logic)
	// run logic a few times so that it has runtimeEstimate data
	for i := 0; i < 32; i++ {
		r.Run(FRAME_MS, false)
	}
	// remove it
	Logger.Printf("Removing logic: %s", logic.name)
	r.Remove(logic)
	// test if removed
	if _, ok := r.runtimeEstimates[logic]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if _, ok := r.indexes[logic.worldID]; ok {
		t.Fatal("did not delete runtimeEstimates data")
	}
	if len(r.logicUnits) != 0 {
		t.Fatal("did not remove from logicUnits list")
	}
}

func TestRuntimeLimitShare(t *testing.T) {
	w := testingWorld()
	sharer := NewRuntimeLimitSharer()
	counters := make([]int, 0)

	const N = 3
	const M = 3
	const LOOPS = 5
	const SLEEP = 16

	sharer.RegisterRunner("basic")
	for i := 0; i < N; i++ {
		func(i int) {
			counters = append(counters, 0) // jet fuel can't melt steel beams
			sharer.AddLogic("basic", &LogicUnit{
				name:    fmt.Sprintf("basic-%d", i),
				worldID: w.IdGen.Next(),
				f: func(dt_ms float64) {
					time.Sleep(SLEEP)
					counters[i] += 1
				},
				active:      true,
				runSchedule: nil})
		}(i)
	}
	sharer.RegisterRunner("extra")
	for i := 0; i < M; i++ {
		func(i int) {
			counters = append(counters, 0) // jet fuel can't melt steel beams
			sharer.AddLogic("extra", &LogicUnit{
				name:    fmt.Sprintf("extra-%d", i),
				worldID: w.IdGen.Next(),
				f: func(dt_ms float64) {
					time.Sleep(SLEEP)
					counters[i] += 1
				},
				active:      true,
				runSchedule: nil})
		}(i)
	}
	for i := 0; i < LOOPS; i++ {
		sharer.Share((N+M)*SLEEP + 100)
		time.Sleep(FRAME_DURATION)
	}
	expected := N*LOOPS + M*LOOPS
	sum := 0
	for _, counter := range counters {
		sum += counter

	}
	for _, l := range sharer.runnerMap["basic"].logicUnits {
		Logger.Printf("basic.%s: h%d", l.name, l.hotness)
	}
	for _, l := range sharer.runnerMap["extra"].logicUnits {
		Logger.Printf("extra.%s: h%d", l.name, l.hotness)
	}
	if sum < expected {
		t.Fatalf("didn't share runtime properly; expected >= %d, got %d", expected, sum)
	}
}

func TestRuntimeLimitShareInsertWhileRunning(t *testing.T) {
	w := testingWorld()
	getCount := func(doInsert bool) int {
		sharer := NewRuntimeLimitSharer()
		counter := 0
		const N = 3
		const LOOPS = 5
		const SLEEP = 16
		sharer.RegisterRunner("basic")
		insert := func(i int) {
			sharer.AddLogic("basic", &LogicUnit{
				name:    fmt.Sprintf("basic-%d", i),
				worldID: w.IdGen.Next(),
				f: func(dt_ms float64) {
					time.Sleep(SLEEP)
					counter += 1
				},
				active:      true,
				runSchedule: nil})
		}
		for i := 0; i < N; i++ {
			insert(i)
		}
		for i := 0; i < LOOPS; i++ {
			// insert with 3 loops left to go
			if doInsert && i == LOOPS-3 {
				insert(N + i)
			}
			// ensure there's always enough time to run every one
			sharer.Share(5 * N * SLEEP)
			time.Sleep(FRAME_DURATION)
		}
		return counter
	}
	plain := getCount(false)
	inserted := getCount(true)
	if inserted <= plain {
		t.Fatalf("didn't share runtime properly, expected inserted run to have higher counter")
	}
}

func TestRuntimeLimiterInsertAppending(t *testing.T) {
	r := NewRuntimeLimiter()
	for i := 0; i < 32; i++ {
		name := fmt.Sprintf("logic-%d", i)
		logic := &LogicUnit{
			name:        name,
			worldID:     i,
			f:           func(dt_ms float64) {},
			active:      true,
			runSchedule: nil}
		r.Add(logic)
		if !(len(r.logicUnits) > 0 &&
			r.indexes[logic.worldID] == len(r.logicUnits)-1) {
			t.Fatal("was not inserted properly")
		}
	}
}
