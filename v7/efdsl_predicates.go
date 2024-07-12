package sameriver

func EFDSLPredicatesBase(e *EFDSLEvaluator) EFDSLPredicateMap {

	return EFDSLPredicateMap{

		//TODO
		// generic Eq
		// generic Lt, Le, Gt, Ge
		/*
			signature IdentResolve<int>,IdentResolve<int>
			"Gt(self<martialarts.skill>, mind.martialarts.prospectiveOpponent<martialarts.skill>)"

			signature IdentResolve<int>,IdentResolve<int>
			Lt(mind.trading.lowestBargain, mind.trading.other.offer)
		*/

		// we want to be able to do something like:
		// Eq(self<martialarts.skill>, <martialarts.skill>)
		// and overload for whether these are bool,bool, int,int, etc.
		/*
			"Eq": e.Predicate(
				"what EFDSL signature string goes here?"
				func(    TODO: what func signature goes here?     ) func(*Entity) bool {
					return func(x *Entity) bool {
						// what code happens in here?
					}
				},
			),
		*/

		"CanBe": e.Predicate(
			"string, int",
			func(k string, v int) func(*Entity) bool {
				return func(x *Entity) bool {
					return e.w.EntityHasComponent(x, STATE_) && e.w.GetIntMap(x, STATE_).ValCanBeSetTo(k, v)
				}
			},
		),

		"State": e.Predicate(
			"string, int",
			func(k string, v int) func(*Entity) bool {
				return func(x *Entity) bool {
					return e.w.EntityHasComponent(x, STATE_) && e.w.GetIntMap(x, STATE_).Get(k) == v
				}
			},
		),

		"HasComponent": e.Predicate(
			"string",
			func(componentStr string) func(*Entity) bool {
				return func(x *Entity) bool {
					// do a little odd access pattern since we only have
					// HasComponent for ComponentID (int) not strings.
					ct := e.w.Em.ComponentsTable
					componentID := ct.StringsRev[componentStr]
					return e.w.EntityHasComponent(x, componentID)
				}
			},
		),

		"HasTag": e.Predicate(
			"string",
			func(tag string) func(*Entity) bool {
				return func(x *Entity) bool {
					return e.w.EntityHasTag(x, tag)
				}
			},
		),

		"HasTags": e.Predicate(
			"[]string",
			func(tags []string) func(*Entity) bool {
				return func(x *Entity) bool {
					return e.w.EntityHasTags(x, tags...)
				}
			},
		),

		"Is": e.Predicate(
			"IdentResolve<int>",
			func(yID int) func(*Entity) bool {
				y := e.w.GetEntity(yID)
				return func(x *Entity) bool {
					return x == y
				}
			},
		),

		"WithinDistance": e.Predicate(
			"IdentResolve<int>, float64",
			func(yID int, d float64) func(*Entity) bool {
				y := e.w.GetEntity(yID)
				return func(x *Entity) bool {
					pos := e.w.GetVec2D(x, POSITION_)
					box := e.w.GetVec2D(x, BOX_)
					return e.w.EntityDistanceFromRect(y, *pos, *box) < d
				}
			},
			"IdentResolve<*Vec2D>, IdentResolve<*Vec2D>, float64",
			func(pos *Vec2D, box *Vec2D, d float64) func(*Entity) bool {
				return func(x *Entity) bool {
					return e.w.EntityDistanceFromRect(x, *pos, *box) < d
				}
			},
		),

		"RectOverlap": func(args []string, resolver IdentifierResolver) func(*Entity) bool {
			argsTyped, err := DSLAssertArgTypes("Vec2D, Vec2D, Vec2D, Vec2D", args, resolver)
			if err != nil {
				logDSLError("%s", err)
			}
			pos := argsTyped[0].(*Vec2D)
			box := argsTyped[1].(*Vec2D)

			return func(x *Entity) bool {
				ePos := e.w.GetVec2D(x, POSITION_)
				eBox := e.w.GetVec2D(x, BOX_)

				return RectIntersectsRect(*pos, *box, *ePos, *eBox)
			}
		},
		// TODO: withinpolygon, overlapspolygon

	}
}
