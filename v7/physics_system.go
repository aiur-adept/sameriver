package sameriver

import (
	"runtime"
	"sync"
)

// moves entities according to their velocity
type PhysicsSystem struct {
	granularity     int
	w               *World
	physicsEntities *UpdatedEntityList
	h               *SpatialHasher
	c               *CollisionSystem `sameriver-system-dependency:"-"`
}

func NewPhysicsSystem() *PhysicsSystem {
	return NewPhysicsSystemWithGranularity(1)
}

func NewPhysicsSystemWithGranularity(granularity int) *PhysicsSystem {
	return &PhysicsSystem{
		granularity: granularity,
	}
}

func (p *PhysicsSystem) GetComponentDeps() []any {
	// TODO: do something with mass
	// TODO: impart momentum to collided objects?
	return []any{
		POSITION_, VEC2D, "POSITION",
		VELOCITY_, VEC2D, "VELOCITY",
		ACCELERATION_, VEC2D, "ACCELERATION",
		BOX_, VEC2D, "BOX",
		MASS_, FLOAT64, "MASS",
		RIGIDBODY_, BOOL, "RIGIDBODY",
	}
}

func (p *PhysicsSystem) LinkWorld(w *World) {
	p.w = w
	p.physicsEntities = w.Em.GetSortedUpdatedEntityList(
		w.EntityFilterFromComponentBitArray(
			"physical",
			w.Em.ComponentsTable.BitArrayFromIDs([]ComponentID{POSITION_, VELOCITY_, ACCELERATION_, BOX_, MASS_, RIGIDBODY_})))
	p.h = NewSpatialHasher(10, 10, w)
}

func (p *PhysicsSystem) Update(dt_ms float64) {
	p.h.Update()
	sum_dt := 0.0
	for i := 0; i < p.granularity; i++ {
		p.ParallelUpdate(dt_ms / float64(p.granularity))
		sum_dt += dt_ms / float64(p.granularity)
	}
}

func (p *PhysicsSystem) collision(e *Entity, pos, box *Vec2D, dx, dy float64) bool {
	rigidBody := p.w.GetBool(e, RIGIDBODY_)
	if !*rigidBody {
		return false
	}

	collision := false

	// check collisions using spatial hasher
	testCollision := func(i *Entity, j *Entity) bool {
		iPos := p.w.GetVec2D(i, POSITION_)
		iBox := p.w.GetVec2D(i, BOX_)
		jPos := p.w.GetVec2D(j, POSITION_)
		jBox := p.w.GetVec2D(j, BOX_)

		return RectIntersectsRect(*iPos, *iBox, *jPos, *jBox)
	}
	cellX0, cellX1, cellY0, cellY1 := p.h.CellRangeOfRect(*pos, *box)
	for y := cellY0; y <= cellY1; y++ {
		for x := cellX0; x <= cellX1; x++ {
			if x < 0 || x >= p.h.GridX || y < 0 || y >= p.h.GridY {
				continue
			}
			entities := p.h.Entities(x, y)
			for i := 0; i < len(entities); i++ {
				other := entities[i]
				if other.ID == e.ID {
					continue
				}
				if !p.w.EntityHasComponent(other, RIGIDBODY_) {
					continue
				}
				otherRigidBody := p.w.GetBool(other, RIGIDBODY_)
				if !*otherRigidBody {
					continue
				}
				if testCollision(e, other) {
					collision = true
					if e.ID < other.ID {
						p.c.DoCollide(e, other)
					} else {
						p.c.DoCollide(other, e)
					}
				}
			}
		}
	}
	return collision
}

func (p *PhysicsSystem) physics(e *Entity, dt_ms float64) {

	// the logic is simpler to read that way
	pos := p.w.GetVec2D(e, POSITION_)
	box := p.w.GetVec2D(e, BOX_)

	// calculate velocity
	acc := p.w.GetVec2D(e, ACCELERATION_)
	vel := p.w.GetVec2D(e, VELOCITY_)
	vel.X += acc.X * dt_ms
	vel.Y += acc.Y * dt_ms
	dx := vel.X * dt_ms
	dy := vel.Y * dt_ms

	halfWidth := box.X / 2
	halfHeight := box.Y / 2

	// motion in x
	if pos.X+dx-halfWidth < 0 || pos.X+dx+halfWidth > float64(p.w.Width) {
		dx = 0
	} else {
		pos.X += dx
	}
	if p.collision(e, pos, box, dx, dy) {
		pos.X -= dx
	}

	// motion in y
	if pos.Y+dy-halfHeight < 0 || pos.Y+dy+halfHeight > float64(p.w.Height) {
		dy = 0
	} else {
		pos.Y += dy
	}
	if p.collision(e, pos, box, dx, dy) {
		pos.Y -= dy
	}
}

func (p *PhysicsSystem) ParallelUpdate(dt_ms float64) {
	// divide the entities into N segments,
	// where N is the number of CPU cores
	numWorkers := runtime.NumCPU()
	entitiesPerWorker := len(p.physicsEntities.entities) / numWorkers
	remainder := len(p.physicsEntities.entities) % numWorkers

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		startIndex := i * entitiesPerWorker
		endIndex := (i + 1) * entitiesPerWorker
		if i == numWorkers-1 {
			endIndex += remainder
		}

		go func(startIndex, endIndex int) {
			for j := startIndex; j < endIndex; j++ {
				e := p.physicsEntities.entities[j]
				p.physics(e, dt_ms)
			}
			wg.Done()
		}(startIndex, endIndex)
	}

	wg.Wait()
}

func (p *PhysicsSystem) SingleThreadUpdate(dt_ms float64) {
	// note: there are no function calls in the below, so we won't
	// be preempted while computing physics (this is very good, get it over with)
	for i := range p.physicsEntities.entities {
		e := p.physicsEntities.entities[i]
		p.physics(e, dt_ms)
	}
}

func (p *PhysicsSystem) Expand(n int) {
	// nil?
}
