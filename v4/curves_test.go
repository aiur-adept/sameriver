package sameriver

import (
	"math"
	"testing"
)

func TestCurvesClimb(t *testing.T) {
	expect := []float64{
		Curves.Sigmoid(0.5, 1)(-1), 0,
		Curves.Sigmoid(0.5, 1)(0), 0,
		Curves.Sigmoid(0.5, 1)(0.5), 0.5,
		Curves.Sigmoid(0.5, 1)(1), 1,
		Curves.Sigmoid(0.5, 1)(2), 1,

		Curves.Shelf(0.5, 1)(-1), 0,
		Curves.Shelf(0.5, 1)(0), 0,
		Curves.Shelf(0.5, 1)(0.5), 0,
		Curves.Shelf(0.5, 1)(1), 1,
		Curves.Shelf(0.5, 1)(2), 1,

		Curves.Shelf(1, 1)(-1), 0,
		Curves.Shelf(1, 1)(0), 0,
		Curves.Shelf(1, 1)(0.5), 0,
		Curves.Shelf(1, 1)(1), 1,
		Curves.Shelf(1, 1)(2), 1,

		Curves.Quad(0.5, 1)(-1), 0,
		Curves.Quad(0.5, 1)(0), 0,
		Curves.Quad(0.5, 1)(0.5), 0,
		Curves.Quad(0.5, 1)(1), 1,
		Curves.Quad(0.5, 1)(2), 1,

		Curves.Cubi(0.5, 1)(-1), 0,
		Curves.Cubi(0.5, 1)(0), 0,
		Curves.Cubi(0.5, 1)(0.5), 0,
		Curves.Cubi(0.5, 1)(1), 1,
		Curves.Cubi(0.5, 1)(2), 1,

		Curves.Lin(-1), 0,
		Curves.Lin(0), 0,
		Curves.Lin(0.5), 0.5,
		Curves.Lin(1), 1,
		Curves.Lin(2), 1,

		Curves.Lint(0.5, 0.2)(-1), 0,
		Curves.Lint(0.5, 0.2)(0), 0,
		Curves.Lint(0.5, 0.2)(0.5), 0.5,
		Curves.Lint(0.5, 0.2)(1), 1,
		Curves.Lint(0.5, 0.2)(2), 1,

		Curves.Exp(5)(-1), 0,
		Curves.Exp(5)(0), 0,
		Curves.Exp(5)(1), 1,
		Curves.Exp(5)(2), 1,
	}
	Logger.Println(Curves.Lin(0.5))
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}

func TestCurvesPeaks(t *testing.T) {
	expect := []float64{
		Curves.Bell(0.5, 1)(-1), 0.058,
		Curves.Bell(0.5, 1)(0), 0.058,
		Curves.Bell(0.5, 1)(0.5), 1,
		Curves.Bell(0.5, 1)(1), 0.058,
		Curves.Bell(0.5, 1)(2), 0.058,

		Curves.BellPinned(0.5)(-1), 0,
		Curves.BellPinned(0.5)(0), 0,
		Curves.BellPinned(0.5)(0.5), 1,
		Curves.BellPinned(0.5)(1), 0,
		Curves.BellPinned(0.5)(2), 0,

		Curves.Plateau(4)(-1), 0,
		Curves.Plateau(4)(0), 0,
		Curves.Plateau(4)(0.5), 1,
		Curves.Plateau(4)(1), 0,
		Curves.Plateau(4)(2), 0,
	}
	Logger.Println(Curves.Lin(0.5))
	for i := 0; i < len(expect); i += 2 {
		// "close enough" since for example sigmoid(0.5, 1)(1) isn't exactly 1
		if math.Abs(expect[i]-expect[i+1]) > 0.001 {
			t.Fatalf("condition %d resulted in %f, not %f", i/2, expect[i], expect[i+1])
		}
	}
}
