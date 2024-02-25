package bracket

import (
	"math"
	"testing"
)

func TestCalculateNewElo(t *testing.T) {
	pointsA, pointsB := CalculateNewElo(800, 800, 0.0)
	if pointsA >= pointsB {
		t.Fatalf(`A should have less points than B after losing. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
	if math.Abs(float64(pointsA)) != math.Abs(float64(pointsB)) {
		t.Fatalf(`Points gained should be equal to points lost. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
	pointsA, pointsB = CalculateNewElo(800, 800, 1.0)
	if pointsA <= pointsB {
		t.Fatalf(`A should have more points than B after winning. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
	if math.Abs(float64(pointsA)) != math.Abs(float64(pointsB)) {
		t.Fatalf(`Points gained should be equal to points lost. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
	pointsA, pointsB = CalculateNewElo(800, 800, 0.5)
	if pointsA != pointsB {
		t.Fatalf(`A should have the same amount of points as B after a draw. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
	if math.Abs(float64(pointsA)) != math.Abs(float64(pointsB)) {
		t.Fatalf(`Points gained should be equal to points lost. Points A: %q, Points B: %q`, pointsA, pointsB)
	}
}
