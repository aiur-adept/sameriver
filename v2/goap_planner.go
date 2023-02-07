package sameriver

import (
	"bytes"
	"os"
	"strings"
)

func debugGOAPPrintf(s string, args ...any) {
	if val, ok := os.LookupEnv("DEBUG_GOAP"); ok && val == "true" {
		Logger.Printf(s, args...)
	}
}

func GOAPPathToString(plan []*GOAPAction) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, action := range plan {
		buf.WriteString(action.name)
		if i != len(plan)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

func debugPrintGOAPGoal(g *GOAPGoal) {
	if g == nil || len(g.goals) == 0 {
		debugGOAPPrintf("    nil")
		return
	}
	for spec, interval := range g.goals {
		split := strings.Split(spec, ",")
		varName := split[0]
		debugGOAPPrintf("    %s: [%.0f, %.0f]", varName, interval.a, interval.b)
	}
}

type GOAPPlanner struct {
	e    *Entity
	eval *GOAPEvaluator
}

type GOAPPathAndRemaining struct {
	path      []*GOAPAction
	remaining *GOAPGoal
}

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:    e,
		eval: NewGOAPEvaluator(),
	}
}

func (p *GOAPPlanner) deepen(
	start *GOAPWorldState,
	path []*GOAPAction,
	goal *GOAPGoal) (newPaths []*GOAPPathAndRemaining) {

	newPaths = make([]*GOAPPathAndRemaining, 0)

	prepend := func(a *GOAPAction, path []*GOAPAction) []*GOAPAction {
		prepended := make([]*GOAPAction, len(path))
		copy(prepended, path)
		prepended = append([]*GOAPAction{a}, path...)
		return prepended
	}

	pathResult := p.eval.applyPath(path, start)
	for _, action := range p.eval.actions.set {
		newPath := prepend(action, path)
		newResult := p.eval.applyPath(newPath, start)
		closer, remaining := goal.stateCloserInSomeVar(newResult, pathResult)
		if closer {
			newPaths = append(newPaths, &GOAPPathAndRemaining{
				path:      newPath,
				remaining: remaining,
			})
		}
	}
	return newPaths
}

func (p *GOAPPlanner) traverseFulfillers(
	pq *GOAPPriorityQueue,
	start *GOAPWorldState,
	path []*GOAPAction,
	goal *GOAPGoal) {

	debugGOAPPrintf("--------------------------")
	debugGOAPPrintf("backtrack path so far: ")
	debugGOAPPrintf(GOAPPathToString(path))
	debugGOAPPrintf("traversing fulfillers of goal:")
	debugPrintGOAPGoal(goal)
	newPaths := p.deepen(start, path, goal)
	debugGOAPPrintf("newPaths:")
	for _, pathAndRemaining := range newPaths {
		newPath := pathAndRemaining.path
		debugGOAPPrintf("---")
		debugGOAPPrintf("    %s", GOAPPathToString(newPath))

		debugGOAPPrintf("    remaining for this path:")
		debugPrintGOAPGoal(pathAndRemaining.remaining)
		mergedGoal, ok := goal.prependingMerge(pathAndRemaining.remaining)
		debugGOAPPrintf("    new merged goal:")
		debugPrintGOAPGoal(mergedGoal)
		if ok {
			pq.Push(NewGOAPPQueueItem(newPath, mergedGoal))
		}
	}
}

func (p *GOAPPlanner) Plan(
	start *GOAPWorldState,
	goal *GOAPGoal,
	maxIter int) (solution []*GOAPAction, ok bool) {

	resultPq := &GOAPPriorityQueue{}

	pq := &GOAPPriorityQueue{}

	p.traverseFulfillers(pq, start, []*GOAPAction{}, goal)

	iter := 0
	for iter < maxIter && pq.Len() > 0 && resultPq.Len() < 2 {
		here := pq.Pop().(*GOAPPQueueItem)
		debugGOAPPrintf("len(here.goal.goals) == %d", len(here.goal.goals))
		if len(here.goal.goals) == 0 {
			// potential solution!
			if p.validateForward(start, here.path, goal) {
				// we push to a pqueue so we can, at the end, pop the
				// solution with the least cost
				resultPq.Push(here)
			} else {
				debugGOAPPrintf("found an invalid solution on validateForward()")
			}
		} else {
			p.traverseFulfillers(pq, start, here.path, here.goal)
			iter++
		}
	}

	if resultPq.Len() > 0 {
		return resultPq.Pop().(*GOAPPQueueItem).path, true
	} else {
		return nil, false
	}
}

func (p *GOAPPlanner) validateForward(
	start *GOAPWorldState,
	path []*GOAPAction,
	goal *GOAPGoal) bool {

	world := start
	for _, action := range path {
		if !p.eval.presFulfilled(action, world) {
			return false
		}
		world = p.eval.applyAction(action, world)
	}

	remaining, _ := goal.remaining(world)
	return len(remaining.goals) == 0
}
