package service

import (
	"math"
	"testing"
)

func TestRoundIsRoundedToZero(t *testing.T) {
	epsilon := 0.0001
	expected := 3.0
	rounded := Round(3.004, 0.01)

	difference := math.Abs(expected - rounded)

	if difference > epsilon {
		t.Error("Round has not sufficient precision")
	}
}

func TestRoundIsRoundedToTen(t *testing.T) {
	epsilon := 0.0001
	expected := 3.01
	rounded := Round(3.005, 0.01)

	difference := math.Abs(expected - rounded)

	if difference > epsilon {
		t.Error("Round has not sufficient precision")
	}
}
