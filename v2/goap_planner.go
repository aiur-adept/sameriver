package sameriver

import (
	"fmt"
	"strings"
	"time"

	"github.com/TwiN/go-color"
)

type GOAPPlanner struct {
	e    *Entity
	eval *GOAPEvaluator
}

func NewGOAPPlanner(e *Entity) *GOAPPlanner {
	return &GOAPPlanner{
		e:    e,
		eval: NewGOAPEvaluator(),
	}
}

func (p *GOAPPlanner) traverseFulfillers(
	pq *GOAPPriorityQueue,
	start *GOAPWorldState,
	here *GOAPPQueueItem,
	goal *GOAPGoal,
	pathsSeen map[string]bool) {

	if DEBUG_GOAP {
		logGOAPDebug("traverse--------------------------")
		logGOAPDebug(color.InRedOverGray("remaining:"))
		debugGOAPPrintGoalRemainingSurface(here.path.remainings)
		logGOAPDebug("%d possible actions", len(p.eval.actions.set))
	}

	// determine if action is good to insert anywhere
	// consider, surface: [Apre, Bpre, Main]
	// consider inserting at 0 means fulfilling Apre
	for i, g := range here.path.remainings.surface {
		if g.nUnfulfilled == 0 {
			continue
		}
		for varName, _ := range g.goalLeft {
			for action, _ := range p.eval.varActions[varName] {
				if DEBUG_GOAP {
					logGOAPDebug("[ ] Considering action %s", action.DisplayName())
				}
				if DEBUG_GOAP {
					var toSatisfyMsg string
					if i == len(here.path.remainings.surface)-1 {
						toSatisfyMsg = "main goal"
					} else {
						toSatisfyMsg = fmt.Sprintf("pre of %s", here.path.path[i].Name)
					}
					logGOAPDebug(color.InGreenOverGray(
						fmt.Sprintf("checking if %s can be inserted at %d to satisfy %s",
							action.DisplayName(), i, toSatisfyMsg)))
				}
				scale, helpful := p.eval.actionHelpsToInsert(start, here.path, i, action)
				if helpful {
					if DEBUG_GOAP {
						logGOAPDebug("[X] %s helpful!", action.DisplayName())
					}
					var toInsert *GOAPAction
					if scale > 1 {
						toInsert = action.Parametrized(scale)
					} else {
						toInsert = action
					}
					newPath := here.path.inserted(toInsert, i)
					pathStr := newPath.String()
					if _, ok := pathsSeen[pathStr]; ok {
						logGOAPDebug(color.InRedOverGray("xxxxxxxxxxxxxxxxxxx path already seen xxxxxxxxxxxxxxxxxxxxx"))
						continue
					} else {
						pathsSeen[pathStr] = true
					}
					p.eval.computeRemainingsOfPath(newPath, start, goal)
					if DEBUG_GOAP {
						msg := fmt.Sprintf("{} - {} - {}    new path: %s     (cost %d)",
							GOAPPathToString(newPath), newPath.cost)
						logGOAPDebug(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
						logGOAPDebug(color.InWhiteOverCyan(msg))
						logGOAPDebug(color.InWhiteOverCyan(strings.Repeat(" ", len(msg))))
					}
					pq.Push(&GOAPPQueueItem{path: newPath})
				} else {
					if DEBUG_GOAP {
						logGOAPDebug("[_] %s not helpful", action.DisplayName())
					}
				}
			}
		}
	}
	logGOAPDebug("--------------------------/traverse")
}

func (p *GOAPPlanner) Plan(
	start *GOAPWorldState,
	goal *GOAPGoal,
	maxIter int) (solution *GOAPPath, ok bool) {

	// populate start state with any modal vals at start
	p.eval.PopulateModalStartState(start)

	// used to return the solution with lowest cost among solutions found
	resultPq := &GOAPPriorityQueue{}

	// used to keep track of which paths we've already seen since there's multiple ways to
	// reach a path in the insertion-based logic we use
	pathsSeen := make(map[string]bool)

	// used for the search
	pq := &GOAPPriorityQueue{}

	rootPath := NewGOAPPath(nil)
	p.eval.computeRemainingsOfPath(rootPath, start, goal)
	backtrackRoot := &GOAPPQueueItem{
		path:  rootPath,
		index: -1, // going to be set by Push()
	}
	pq.Push(backtrackRoot)

	iter := 0
	// TODO: should we just pop out the *very first result*?
	// why wait for 2 or exhausting the pq?
	t0 := time.Now()
	for iter < maxIter && pq.Len() > 0 && resultPq.Len() < 2 {
		logGOAPDebug("=== iter ===")
		here := pq.Pop().(*GOAPPQueueItem)
		if DEBUG_GOAP {
			logGOAPDebug(color.InRedOverGray("here:"))
			logGOAPDebug(color.InWhiteOverBlue(color.InBold(GOAPPathToString(here.path))))
			logGOAPDebug(color.InRedOverGray(fmt.Sprintf("(%d unfulfilled)",
				here.path.remainings.NUnfulfilled())))
		}

		if here.path.remainings.NUnfulfilled() == 0 {
			ok := p.eval.validateForward(here.path, start, goal)
			if !ok {
				logGOAPDebug(">>>>>>> potential solution rejected")
				continue
			}

			if DEBUG_GOAP {
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(fmt.Sprintf("    SOLUTION: %s", GOAPPathToString(here.path)))))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(">>>>>>>>>>>>>>>>>>>>>>")))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(GOAPPathToString(here.path))))
				logGOAPDebug(color.InGreenOverWhite(color.InBold(fmt.Sprintf("%d solutions found so far", resultPq.Len()+1))))
			}
			resultPq.Push(here)
		} else {
			p.traverseFulfillers(pq, start, here, goal, pathsSeen)
			iter++
		}
	}

	dt := float64(time.Since(t0).Nanoseconds()) / 1.0e6
	if iter >= maxIter {
		logGOAPDebug("Took %f ms to reach max iter (%d)", dt, iter)
		logGOAPDebug("================================ REACHED MAX ITER !!!")
	}
	if pq.Len() == 0 && resultPq.Len() == 0 {
		logGOAPDebug("Took %f ms to exhaust pq without solution (%d iterations)", dt, iter)
	}
	if resultPq.Len() > 0 {
		logGOAPDebug("Took %f ms to find %d solutions (%d iterations)", dt, resultPq.Len(), iter)
		if pq.Len() == 0 {
			logGOAPDebug("Exhausted pq")
		}
		for _, item := range *resultPq {
			logGOAPDebug(color.InWhiteOverBlue(item.path))
		}
		return resultPq.Pop().(*GOAPPQueueItem).path, true
	} else {
		return nil, false
	}
}
