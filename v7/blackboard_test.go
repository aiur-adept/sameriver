package sameriver

import (
	"encoding/json"
	"fmt"
	"reflect"
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

		w.AddLogic(fmt.Sprintf("village-blackboard-%d", e.ID), func(dt_ms float64) {
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

		w.AddLogic("village-morning", func(dt_ms float64) {
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

func TestBlackboardSaveLoad(t *testing.T) {
	w := testingWorld()

	bname := "village-12"
	bb := w.Blackboard(bname)
	bb.Set("roles", []string{"farmer", "baker", "fisher", "crafts"})
	bb.Set("number", 12)
	bb.Set("bool", true)
	bb.Set("numbers", []int{1, 2, 3, 4, 5})

	jsonStr, err := json.Marshal(bb)
	if err != nil {
		t.Fatalf("error marshalling blackboard: %v", err)
	}

	fmt.Println(string(jsonStr))

	bb2 := NewBlackboard(bname + "-reloaded")
	err = json.Unmarshal(jsonStr, bb2)
	if err != nil {
		t.Fatalf("error unmarshalling blackboard: %v", err)
	}

	// test if bb2 contains the same roles
	roles1, ok1 := bb.Get("roles").([]string)
	roles2, ok2 := bb2.Get("roles").([]string)

	fmt.Printf("Type of roles1: %T, Type of roles2: %T\n", bb.Get("roles"), bb2.Get("roles"))
	fmt.Printf("roles1: %v, roles2: %v\n", roles1, roles2)
	fmt.Printf("Type assertion ok1: %v, ok2: %v\n", ok1, ok2)

	if !ok1 || !ok2 || !reflect.DeepEqual(roles1, roles2) {
		t.Fatal("roles array not in bb2 or type mismatch")
	}

	number1, ok1 := bb.Get("number").(float64)
	number2, ok2 := bb2.Get("number").(float64)
	fmt.Printf("Type of number1: %T, Type of number2: %T\n", bb.Get("number"), bb2.Get("number"))
	fmt.Printf("number1: %v, number2: %v\n", number1, number2)
	fmt.Printf("Type assertion ok1: %v, ok2: %v\n", ok1, ok2)
	if !ok1 || !ok2 || number1 != number2 {
		t.Fatal("number not in bb2 or type mismatch")
	}

	bool1, ok1 := bb.Get("bool").(bool)
	bool2, ok2 := bb2.Get("bool").(bool)
	if !ok1 || !ok2 || bool1 != bool2 {
		t.Fatal("bool not in bb2 or type mismatch")
	}

	numbers1, ok1 := bb.Get("numbers").([]float64)
	numbers2, ok2 := bb2.Get("numbers").([]float64)
	fmt.Printf("Type of numbers1: %T, Type of numbers2: %T\n", bb.Get("numbers"), bb2.Get("numbers"))
	fmt.Printf("numbers1: %v, numbers2: %v\n", numbers1, numbers2)
	fmt.Printf("Type assertion ok1: %v, ok2: %v\n", ok1, ok2)
	if !ok1 || !ok2 || !reflect.DeepEqual(numbers1, numbers2) {
		t.Fatal("numbers not in bb2 or type mismatch")
	}
}
