package sameriver

import (
	"fmt"
	"sort"
)

// notice the sortf returned by Evaluate() is a closure that wants the result string so it can actually use i, j int
// params: EFDSLEval takes the expression and a resolver (for identifiers)
// returns: an entity predicate and an entity sort function and possibly an error
//
//	aka (p, q, err)
func EFDSLEval(expr string, resolver IdentifierResolver, world *World) (func(*Entity) bool, func(xs []*Entity) func(i, j int) bool, error) {
	parser := &EFDSLParser{}

	ast, err := parser.Parse(expr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse expr: %s", err)
	}

	filter, sort := world.EFDSL.Evaluate(ast, resolver)

	return filter, sort, nil
}

func EFDSLFilter(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filter, _, err := EFDSLEval(expr, resolver, world)
	if err != nil {
		return nil, err
	}
	result := world.FilterAllEntities(filter)
	return result, nil
}

func EFDSLFilterSort(expr string, resolver IdentifierResolver, world *World) ([]*Entity, error) {
	filterf, sortf, err := EFDSLEval(expr, resolver, world)
	if err != nil {
		return nil, err
	}
	result := world.FilterAllEntities(filterf)
	if sortf != nil {
		sort.Slice(result, sortf(result))
	}
	return result, nil
}

func (w *World) EFDSLFilterEntity(e *Entity, expr string) ([]*Entity, error) {
	resolver := &EntityResolver{e: e, w: w}
	return EFDSLFilter(expr, resolver, w)
}

func (w *World) EFDSLFilter(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return EFDSLFilter(expr, resolver, w)
}

func (w *World) EFDSLFilterSortEntity(e *Entity, expr string) ([]*Entity, error) {
	resolver := &EntityResolver{e: e, w: w}
	return EFDSLFilterSort(expr, resolver, w)
}

func (w *World) EFDSLFilterSort(expr string) ([]*Entity, error) {
	resolver := &WorldResolver{w: w}
	return EFDSLFilterSort(expr, resolver, w)
}
