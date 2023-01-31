package engine

import (
	"fmt"
)

type RuntimeLimitSharer struct {
	runIX     int
	runners   []*RuntimeLimiter
	runnerMap map[string]*RuntimeLimiter
}

func NewRuntimeLimitSharer() *RuntimeLimitSharer {
	r := &RuntimeLimitSharer{
		runners:   make([]*RuntimeLimiter, 0),
		runnerMap: make(map[string]*RuntimeLimiter),
	}
	return r
}

func (r *RuntimeLimitSharer) RegisterRunner(name string) {
	if _, ok := r.runnerMap[name]; ok {
		panic(fmt.Sprintf("Trying to double-add RuntimeLimiter %s", name))
	}
	runner := NewRuntimeLimiter()
	r.runners = append(r.runners, runner)
	r.runnerMap[name] = runner
}

func (r *RuntimeLimitSharer) AddLogic(runnerName string, l *LogicUnit) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to add to runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].Add(l)
}

func (r *RuntimeLimitSharer) RemoveLogic(runnerName string, l *LogicUnit) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to remove from runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].Remove(l)
}

func (r *RuntimeLimitSharer) ActivateAll(runnerName string) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to activate all in runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].ActivateAll()
}

func (r *RuntimeLimitSharer) DeactivateAll(runnerName string) {
	if _, ok := r.runnerMap[runnerName]; !ok {
		panic(fmt.Sprintf("Trying to deactivate all in runtimeLimiter with name %s - doesn't exist", runnerName))
	}
	r.runnerMap[runnerName].DeactivateAll()
}

func (r *RuntimeLimitSharer) Share(allowance float64) (remaining float64, starved int) {
	remaining = allowance
	ran := 0
	perRunner := allowance / float64(len(r.runners))
	// while we have allowance, try to cycle all the way around
	for allowance >= 0 && ran < len(r.runners) {
		runner := r.runners[r.runIX]
		overunder := runner.Run(perRunner)
		used := perRunner - overunder
		allowance -= used
		// increment to run next runner even if runner.Finished() isn't true
		// this means it will get another chance to finish itself when its
		// turn comes back around
		r.runIX = (r.runIX + 1) % len(r.runners)
		ran++
	}
	starved = len(r.runners) - ran
	return allowance, starved
}

func (r *RuntimeLimitSharer) DumpStats() map[string](map[string]float64) {
	stats := make(map[string](map[string]float64))
	stats["totals"] = make(map[string]float64)
	for name, r := range r.runnerMap {
		runnerStats, totals := r.DumpStats()
		stats[name] = runnerStats
		stats["totals"][name] = totals
	}
	return stats
}
