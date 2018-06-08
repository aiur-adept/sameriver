/*
 *
 * Allocates each of the big ass blocks of memory that each component has
 * its data living inside. This is a scientific terminology of game engine
 * design.
 *
 */

package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ComponentsTable struct {
	accessLocks [N_COMPONENT_TYPES]*ComponentAccessLock
	valueLocks  [N_COMPONENT_TYPES][MAX_ENTITIES]*ComponentValueLock

	Box      [MAX_ENTITIES]sdl.Rect
	Sprite   [MAX_ENTITIES]Sprite
	TagList  [MAX_ENTITIES]TagList
	Velocity [MAX_ENTITIES][2]float32
}

func (t *ComponentsTable) lock(component ComponentType) {
	// lock the accessLock with Lock() (write-lock), causing calls to Access()
	// to enter the wait queue until we unlock
	t.accessLocks[component].Lock()
}

func (t *ComponentsTable) unlock(component ComponentType) {
	// allows all waiting calls to Access to proceed
	t.accessLocks[component].Unlock()
}

func (t *ComponentsTable) access(component ComponentType) {
	// enter a queue if the accessLock is currently Locked, otherwise
	// we get access because all copies of this method instantly RUnlock
	// after acquiring
	t.accessLocks[component].RLock()
	t.accessLocks[component].RUnlock()
}

func (ct *ComponentsTable) Init(em *EntityManager) {
	ct.Color = &ColorComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.Health = &HealthComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.HitBox = &HitBoxComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.Mind = &MindComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.Position = &PositionComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.Sprite = &SpriteComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.TagList = &TagListComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
	ct.Velocity = &VelocityComponent{
		em:         em,
		accessLock: NewComponentAccessLock(),
	}
}
