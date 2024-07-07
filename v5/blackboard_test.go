package sameriver

import (
	"fmt"
	"testing"

	"math/rand"
	"time"
)

func TestBlackboardWorldEntities(t *testing.T) {
	w := testingWorld()

	bname := "village-12"
	bb := w.Blackboard(bname)

	setupVillageBB := func() {
		bb.Set("roles", []string{"farmer", "baker", "fisher", "crafts"})
	}

	spawnVillager := func(rolePreference string) {
		e := testingSpawnSimple(w)

		var villageEvents *EventChannel

		reactThreat := func(data map[string]any) {
			// ( interrupt current plan if we were in GOAP )
			// add ourselves to responders
			if bb.Has("threatResponders") {
				bb.Set("threatResponders", make(map[*Entity]bool))
			}
			responders := bb.Get("threatResponders").(map[*Entity]bool)
			responders[e] = true
			// more GOAP:
			// check bb threat location and add it to our mind to be used
			// as part of the "attack" or "flee" action executors
		}

		reactMorning := func() {
			Logger.Printf("villager %d reacting to morning :)", e.ID)
			todayRoles := bb.Get("todayRoles").(map[*Entity]string)
			unfilledRoles := bb.Get("unfilledRoles").(map[string]bool)
			selectRole := func(role string) {
				todayRoles[e] = role
				unfilledRoles[role] = false
			}
			// try to select our preferred role
			if unfilled := unfilledRoles[rolePreference]; unfilled {
				selectRole(rolePreference)
				return
			}
			// otherwise, select the first unfilled role
			for role, unfilled := range unfilledRoles {
				if unfilled {
					selectRole(role)
					return
				}
			}
			// otherwise, if all roles are already filled, select one at random
			allRoles := bb.Get("roles").([]string)
			randomRole := allRoles[rand.Intn(len(allRoles))]
			selectRole(randomRole)
		}

		w.AddWorldLogic(fmt.Sprintf("village-blackboard-%d", e.ID), func(dt_ms float64) {
			// subscribe to blackboard events
			if villageEvents == nil {
				villageEvents = bb.Events.Subscribe(SimpleEventFilter("village-events"))
			}
			// handle blackboard events
			select {
			case ev := <-villageEvents.C:
				data := ev.Data.(map[string]any)
				switch data["kind"].(string) {
				case "threat":
					reactThreat(data)
				case "morning":
					reactMorning()
				}
			default:
			}
		})
	}

	setupVillageWorldLogic := func() {

		morningTimer := NewTimeAccumulator(500)

		villageBBMorning := func() {
			bb.Set("unfilledRoles", map[string]bool{
				"farmer": true,
				"baker":  true,
				"fisher": true,
				"crafts": true,
			})
			bb.Set("todayRoles", make(map[*Entity]string))
			bb.Events.Publish("village-events", map[string]any{
				"kind": "morning",
			})
		}

		w.AddWorldLogic("village-morning", func(dt_ms float64) {
			if morningTimer.Tick(dt_ms) {
				Logger.Println("world logic announcing morning!")
				villageBBMorning()
			}
		})
	}

	setupVillageBB()
	spawnVillager("farmer")
	spawnVillager("farmer")
	spawnVillager("fisher")
	spawnVillager("fisher")
	spawnVillager("baker")
	spawnVillager("baker")
	spawnVillager("baker")
	setupVillageWorldLogic()

	Logger.Println("Update loop beginning...")
	w.Update(FRAME_MS / 2)
	time.Sleep(550 * time.Millisecond)
	w.Update(FRAME_MS / 2)
	w.Update(FRAME_MS / 2)
	for e, role := range bb.Get("todayRoles").(map[*Entity]string) {
		Logger.Printf("%d will be doing '%s'", e.ID, role)
	}
}
