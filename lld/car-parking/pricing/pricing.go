package pricing

import (
	"math"
	"time"
)

// PricingStrategy computes the parking fee from entry and exit timestamps.
// For an active stay, exit may be the zero time; strategies that need an exit
// should return 0 in that case.
type PricingStrategy interface {
	Compute(entry, exit time.Time) float64
}

type hourly struct{}

func NewHourly() PricingStrategy { return hourly{} }

func (hourly) Compute(entry, exit time.Time) float64 {
	if exit.IsZero() {
		return 0
	}
	v := float64(exit.Sub(entry).Seconds()) * 10
	return math.Round(v*100) / 100
}

type flat struct{}

func NewFlat() PricingStrategy { return flat{} }

func (flat) Compute(_, _ time.Time) float64 {
	return 100
}

type dynamic struct{}

func NewDynamic() PricingStrategy { return dynamic{} }

func (dynamic) Compute(_, _ time.Time) float64 {
	return 0
}
