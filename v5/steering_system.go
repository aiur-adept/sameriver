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
		w.EntityFilterFromComponentBitArray(
			"steering",
			w.Em.ComponentsTable.BitArrayFromIDs(
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
	p0 := s.w.GetVec2D(e, POSITION_)
	p1 := s.w.GetVec2D(e, MOVEMENTTARGET_)
	v := s.w.GetVec2D(e, VELOCITY_)
	maxV := s.w.GetFloat64(e, MAXVELOCITY_)
	st := s.w.GetVec2D(e, STEER_)
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
	v := s.w.GetVec2D(e, VELOCITY_)
	maxV := s.w.GetFloat64(e, MAXVELOCITY_)
	st := s.w.GetVec2D(e, STEER_)
	mass := s.w.GetFloat64(e, MASS_)
	// TODO: define this properly
	maxSteerForce := 3.0
	*st = st.Truncate(maxSteerForce)
	*st = st.Scale(1 / *mass)
	*v = v.Add(*st).Truncate(*maxV)
}

func (s *SteeringSystem) Expand(n int) {
	// nil?
}
