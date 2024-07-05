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
		_POSITION, VEC2D, "POSITION",
		_VELOCITY, VEC2D, "VELOCITY",
		_ACCELERATION, VEC2D, "ACCELERATION",
		_MAXVELOCITY, FLOAT64, "MAXVELOCITY",
		_MOVEMENTTARGET, VEC2D, "MOVEMENTTARGET",
		_STEER, VEC2D, "STEER",
		_MASS, FLOAT64, "MASS",
	}
}

func (s *SteeringSystem) LinkWorld(w *World) {
	s.w = w
	s.movementEntities = w.GetUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			"steering",
			w.em.components.BitArrayFromIDs(
				[]ComponentID{
					_POSITION, _VELOCITY, _ACCELERATION,
					_MAXVELOCITY, _MOVEMENTTARGET, _STEER, _MASS,
				})))
}

func (s *SteeringSystem) Update(dt_ms float64) {
	for _, e := range s.movementEntities.entities {
		s.Seek(e)
		s.Apply(e)
	}
}

func (s *SteeringSystem) Seek(e *Entity) {
	p0 := e.GetVec2D(_POSITION)
	p1 := e.GetVec2D(_MOVEMENTTARGET)
	v := e.GetVec2D(_VELOCITY)
	maxV := e.GetFloat64(_MAXVELOCITY)
	st := e.GetVec2D(_STEER)
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
	v := e.GetVec2D(_VELOCITY)
	maxV := e.GetFloat64(_MAXVELOCITY)
	st := e.GetVec2D(_STEER)
	mass := e.GetFloat64(_MASS)
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / *mass)
	*v = v.Add(*st).Truncate(*maxV)
}

func (s *SteeringSystem) Expand(n int) {
	// nil?
}
