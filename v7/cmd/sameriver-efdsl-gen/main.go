package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const types = "bool,int,float64,string,[]string,Vec2D,[]Vec2D"

const commentWarning = `/*
Heed this warning. Do not edit this file by hand; instead use sameriver-efdsl-gen. And yes, it is horrifying. Blame Rob Pike! My revenge is allowing overloading of predicate/sort func signatures.
*/`

func writeFile(filename string, code string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, code)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
	}
}

/*
permutationsWithRepetitions([]string{"a", "b", "c"}, 2)
yields

[[a a] [b a] [c a] [a b] [b b] [c b] [a c] [b c] [c c]]

permutationsWithRepetitions([]string{"a", "b", "c"}, 3)
yields

[[a a a] [b a a] [c a a] [a b a] [b b a] [c b a] [a c a] [b c a] [c c a] [a a b] [b a b] [c a b] [a b b] [b b b] [c b b] [a c b] [b c b] [c c b] [a a c] [b a c] [c a c] [a b c] [b b c] [c b c] [a c c] [b c c] [c c c]]
*/
func permutationsWithRepetitions(elements []string, length int) [][]string {
	if length == 1 {
		result := make([][]string, len(elements))
		for i, element := range elements {
			result[i] = []string{element}
		}
		return result
	}

	previousPermutations := permutationsWithRepetitions(elements, length-1)
	result := [][]string{}
	for _, element := range elements {
		for _, previousPermutation := range previousPermutations {
			result = append(result, append(previousPermutation, element))
		}
	}

	return result
}

//
// switch statements
//

func genCodeSigAssertSwitchCore(funcReturnType, switchFuncName, userRegisterFunc string, typesStr string) string {
	types := strings.Split(typesStr, ",")
	typePermutations := [][]string{}

	// Generate permutations with repetitions of length 1, 2, and 3
	for i := 1; i <= 3; i++ {
		typePermutations = append(typePermutations, permutationsWithRepetitions(types, i)...)
	}

	var cases []string

	for _, permutation := range typePermutations {
		caseStr := fmt.Sprintf("\tcase func(%s) %s:", strings.Join(permutation, ", "), funcReturnType)
		args := make([]string, len(permutation))
		for i, argType := range permutation {
			args[i] = fmt.Sprintf("argsTyped[%d].(%s)", i, argType)
		}
		caseStr += fmt.Sprintf("\n\t\tresult = fTyped(%s)", strings.Join(args, ", "))
		cases = append(cases, caseStr)
	}

	return fmt.Sprintf(`func (e *EFDSLEvaluator) %s(f any, argsTyped []any) %s {
	var result %s
	switch fTyped := f.(type) {
%s
	default:
		panic("No case in either engine or user-registered signatures for the given func. Use EFDSL.%s()")
	}

	return result
}`, switchFuncName, funcReturnType, funcReturnType, strings.Join(cases, "\n"), userRegisterFunc)
}

func genCodePredicateSigAssertSwitch(typesStr string) string {
	return genCodeSigAssertSwitchCore("func(*Entity) bool", "predicateSignatureAssertSwitch", "RegisterUserPredicateSignatureAsserter", typesStr)
}

func genCodeSortSigAssertSwitch(typesStr string) string {
	return genCodeSigAssertSwitchCore("func(xs []*Entity) func(i, j int) int", "sortSignatureAssertSwitch", "RegisterUserSortSignatureAsserter", typesStr)
}

func genFileSigAssertSwitches(typesStr string) {
	predicateSwitches := genCodePredicateSigAssertSwitch(typesStr)
	sortSwitches := genCodeSortSigAssertSwitch(typesStr)

	code := "package sameriver\n\n" + commentWarning + "\n\n" + predicateSwitches + "\n\n" + sortSwitches

	writeFile("GENERATED_efdsl_sig_assert_switches.go", code)
}

//
// IdentResolve type assert map
//

func genCodeIdentResolveTypeAssertMap(typesStr string) string {
	types := strings.Split(typesStr, ",")

	var entries []string
	for _, t := range types {
		entry := fmt.Sprintf("\t\"%s\": func(arg string, resolver IdentifierResolver) (any, error) {\n\t\treturn AssertT[%s](resolver.Resolve(arg), \"%s\")\n\t},", t, t, t)
		entries = append(entries, entry)
	}

	return fmt.Sprintf("var IdentResolveTypeAssertMap = map[string]DSLArgTypeAssertionFunc{\n%s\n}", strings.Join(entries, "\n"))
}

func genFileIdentResolveTypeAssertMap(types string) {
	identResolveTypeAssertMap := genCodeIdentResolveTypeAssertMap(types)

	code := "package sameriver\n\n" + commentWarning + "\n\n" + identResolveTypeAssertMap

	writeFile("GENERATED_efdsl_identresolve_types.go", code)
}

func main() {
	fmt.Println("sameriver-efdsl-gen (extremely f***** delicious spaghetti lasagna!)")

	genFileSigAssertSwitches(types)
	genFileIdentResolveTypeAssertMap(types)
}
