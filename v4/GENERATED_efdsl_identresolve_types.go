package sameriver

/*
Heed this warning. Do not edit this file by hand; instead use sameriver-efdsl-gen. And yes, it is horrifying. Blame Rob Pike! My revenge is allowing overloading of predicate/sort func signatures.
*/

var IdentResolveTypeAssertMap = map[string]DSLArgTypeAssertionFunc{
	"bool": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[bool](resolver.Resolve(arg), "bool")
	},
	"int": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[int](resolver.Resolve(arg), "int")
	},
	"float64": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[float64](resolver.Resolve(arg), "float64")
	},
	"string": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[string](resolver.Resolve(arg), "string")
	},
	"[]string": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[[]string](resolver.Resolve(arg), "[]string")
	},
	"*Entity": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Entity](resolver.Resolve(arg), "*Entity")
	},
	"[]*Entity": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[[]*Entity](resolver.Resolve(arg), "[]*Entity")
	},
	"*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[*Vec2D](resolver.Resolve(arg), "*Vec2D")
	},
	"[]*Vec2D": func(arg string, resolver IdentifierResolver) (any, error) {
		return AssertT[[]*Vec2D](resolver.Resolve(arg), "[]*Vec2D")
	},
}