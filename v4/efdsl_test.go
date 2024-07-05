package sameriver

import (
	"fmt"
	"strings"
	"testing"
)

func TestEFDSLLexer(t *testing.T) {
	lex := func(s string) {
		fmt.Println(s)
		var l EFDSLLexer
		l.Init(strings.NewReader(s))
		for tok := l.Lex(); tok != EOF; tok = l.Lex() {
			fmt.Printf("%s: %s\n", tok, l.TokenText())
		}
		fmt.Println()
	}

	lex(`HasTag(ox) && CanBe(yoked, 1); Closest(mind.field)`)
	lex(`First(HasTag(ox) && CanBe(yoked, 1); Closest(mind.field))`)
	lex(`VillagerOf(mind.village)`)
	lex(`!VillagerOf(mind.village)`)
	lex(`Is(bb.village1.strongest)`)
	lex(`HasTag(deer); Closest(self)`)
	lex(`HasTags(ox,legendary)`)
	lex(`Closest(self)`)
	lex(`WithinDistance(self, self[hearingDistance])`)
	lex(`HasInv(mind.topWantedItem)`)
	lex(`HasInvAny(mind.wantedItems)`)
	lex(`HasInvTagsAny(tool, food, book)`)
	// a coinosseur:
	lex(`HasInvTags(tool,legendary)`)
}

func TestEFDSLParser(t *testing.T) {
	parser := &EFDSLParser{}
	// TODO: use ast from this call
	_, err := parser.Parse(`HasTag(ox) && CanBe(yoked, 1); Closest(mind.field)`)
	if err != nil {
		t.Fatalf("Why did the expression return an error? it's valid!")
	}
	ast, err := parser.Parse(`HasTag(animal) && WithinDistance(self, 100) && !TimeSinceEventGt(mind.lastAnimalMetta, 10000)`)
	if err != nil {
		t.Fatalf("Why did the expression return an error? it's valid!")
	}
	Logger.Printf("%s", ast)
}

func TestEFDSLEvaluator(t *testing.T) {
	w := testingWorld()
	ox := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION: Vec2D{0, 0},
			_BOX:      Vec2D{3, 2},
			_STATE: map[string]int{
				"yoked": 0,
			},
		},
		"tags": []string{"ox"},
	})

	parser := &EFDSLParser{}

	expression := "HasTags(ox)"

	// Parse and evaluate the expression
	ast, err := parser.Parse(expression)
	if err != nil {
		t.Fatalf("Failed to parse expression: %s", err)
	}

	resolver := &EntityResolver{e: ox}
	filter, _ := EFDSL.Evaluate(ast, resolver)

	// Filter entities using the generated filter function
	result := w.FilterAllEntities(filter)

	// Check if the filtered list contains the expected ox entity
	if len(result) != 1 || result[0] != ox {
		t.Fatalf("Failed to select ox entity: got %v", result)
	}
	Logger.Printf("result of HasTags(ox): %v", result)
}
