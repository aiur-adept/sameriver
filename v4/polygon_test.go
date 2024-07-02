package sameriver

import (
	"math"
	"testing"
)

func TestPolygonDistanceToSide(t *testing.T) {
	// Test case 1: Square centered at the origin
	squareOrigin := &Polygon{
		Vertices: []Vec2D{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}},
	}
	pointInSquare1 := Vec2D{1, 0}
	expectedDistance1 := 0.0
	actualDistance1 := squareOrigin.DistanceToSide(pointInSquare1)
	if math.Abs(actualDistance1-expectedDistance1) > 1e-6 {
		t.Errorf("Test case 1 failed: expected %f, got %f", expectedDistance1, actualDistance1)
	}

	// Test case 2: Right triangle with legs on positive x and y axis
	triangle := &Polygon{
		Vertices: []Vec2D{{0, 0}, {2, 0}, {0, 2}},
	}
	pointInTriangle := Vec2D{1, 1}
	expectedDistance2 := 0.0
	actualDistance2 := triangle.DistanceToSide(pointInTriangle)
	if math.Abs(actualDistance2-expectedDistance2) > 1e-6 {
		t.Errorf("Test case 2 failed: expected %f, got %f", expectedDistance2, actualDistance2)
	}

	// Test case 3: Square centered at 10, 10 and rotated 45 degrees
	squareRotated := &Polygon{
		Vertices: []Vec2D{{10 - math.Sqrt(2), 10 - math.Sqrt(2)}, {10 + math.Sqrt(2), 10 - math.Sqrt(2)},
			{10 + math.Sqrt(2), 10 + math.Sqrt(2)}, {10 - math.Sqrt(2), 10 + math.Sqrt(2)}},
	}
	pointInSquare2 := Vec2D{10, 10}
	expectedDistance3 := math.Sqrt(2)
	actualDistance3 := squareRotated.DistanceToSide(pointInSquare2)
	if math.Abs(actualDistance3-expectedDistance3) > 1e-6 {
		t.Errorf("Test case 3 failed: expected %f, got %f", expectedDistance3, actualDistance3)
	}
}

func TestPolygonDistanceToVertex(t *testing.T) {
	// Test case 1: Square centered at the origin
	squareOrigin := &Polygon{
		Vertices: []Vec2D{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}},
	}
	pointInSquare1 := Vec2D{0, 0}
	expectedDistance1 := math.Sqrt(2)
	actualDistance1 := squareOrigin.DistanceToVertex(pointInSquare1)
	if math.Abs(actualDistance1-expectedDistance1) > 1e-6 {
		t.Errorf("Test case 1 failed: expected %f, got %f", expectedDistance1, actualDistance1)
	}

	// Test case 2: Right triangle with legs on positive x and y axis
	triangle := &Polygon{
		Vertices: []Vec2D{{0, 0}, {2, 0}, {0, 2}},
	}
	pointInTriangle := Vec2D{1, 1}
	expectedDistance2 := math.Sqrt(2)
	actualDistance2 := triangle.DistanceToVertex(pointInTriangle)
	if math.Abs(actualDistance2-expectedDistance2) > 1e-6 {
		t.Errorf("Test case 2 failed: expected %f, got %f", expectedDistance2, actualDistance2)
	}

	// Test case 3: Square centered at 10, 10 and rotated 45 degrees
	squareRotated := &Polygon{
		Vertices: []Vec2D{{10 - math.Sqrt(2), 10 - math.Sqrt(2)}, {10 + math.Sqrt(2), 10 - math.Sqrt(2)},
			{10 + math.Sqrt(2), 10 + math.Sqrt(2)}, {10 - math.Sqrt(2), 10 + math.Sqrt(2)}},
	}
	pointInSquare2 := Vec2D{10, 10}
	expectedDistance3 := 2.0
	actualDistance3 := squareRotated.DistanceToVertex(pointInSquare2)
	if math.Abs(actualDistance3-expectedDistance3) > 1e-6 {
		t.Errorf("Test case 3 failed: expected %f, got %f", expectedDistance3, actualDistance3)
	}
}
