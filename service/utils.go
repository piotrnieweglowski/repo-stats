package service

import "math"

// Round is responsible for rounding a number with given precision
func Round(x float64, unit float64) float64 {
	return math.Round(x/unit) * unit
}
