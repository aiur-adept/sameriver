package sameriver

type SteeringSystem struct {
	w                *World
	movementEntities *UpdatedEntityList
}

func NewSteeringSystem() *SteeringSystem {
	return &SteeringSystem{}
}

func (s *SteeringSystem) GetComponentDeps() []any {
	return []any{
		POSITION_, VEC2D, "POSITION",
		VELOCITY_, VEC2D, "VELOCITY",
		ACCELERATION_, VEC2D, "ACCELERATION",
		MAXVELOCITY_, FLOAT64, "MAXVELOCITY",
		MOVEMENTTARGET_, VEC2D, "MOVEMENTTARGET",
		STEER_, VEC2D, "STEER",
		MASS_, FLOAT64, "MASS",
	}
}

func (s *SteeringSystem) LinkWorld(w *World) {
	s.w = w
	s.movementEntities = w.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"steering",
			w.em.components.BitArrayFromIDs(
				[]ComponentID{
					POSITION_, VELOCITY_, ACCELERATION_,
					MAXVELOCITY_, MOVEMENTTARGET_, STEER_, MASS_,
				})))
}

func (s *SteeringSystem) Update(dt_ms float64) {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *Entity) {
	p0 := e.GetVec2D(POSITION_)
	p1 := e.GetVec2D(MOVEMENTTARGET_)
	v := e.GetVec2D(VELOCITY_)
	maxV := e.GetFloat64(MAXVELOCITY_)
	st := e.GetVec2D(STEER_)
	desired := p1.Sub(*p0)
	distance := desired.Magnitude()
	desired = desired.Unit()
	// do slowing for arrival behavior
	// TODO: define this properly
	slowingRadius := 30.0
	if distance <= slowingRadius {
		desired = desired.Scale(*maxV * distance / slowingRadius)
	} else {
		desired = desired.Scale(*maxV)
	}
	force := desired.Sub(*v)
	st.Inc(force)
}

func (s *SteeringSystem) Apply(e *Entity) {
	v := e.GetVec2D(VELOCITY_)
	maxV := e.GetFloat64(MAXVELOCITY_)
	st := e.GetVec2D(STEER_)
	mass := e.GetFloat64(MASS_)
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / *mass)
	*v = v.Add(*st).Truncate(*maxV)
}

func (s *SteeringSystem) Expand(n int) {
	// nil?
}
