// Copyright 2021 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package substrate

import (
	"fmt"
	"math/big"
)

// Dot wraps a *big.Int and provides conversion and formatting for Dot values.
type Dot struct {
	planck *big.Int
}

const (
	// PlanckPerDot number of plancks per Dot.
	PlanckPerDot = 1e10
	// PrintPrecision is the precision with which floats are printed
	// as defined by big.Float.Text.
	PrintPrecision = 3
)

// NewDotFromPlanck creates a new Dot from the given amount of Plancks.
func NewDotFromPlanck(planck *big.Int) *Dot {
	return &Dot{planck}
}

// NewDotsFromPlancks creates new Dots from the given amounts of Plancks.
func NewDotsFromPlancks(planck ...*big.Int) []*Dot {
	ret := make([]*Dot, len(planck))
	for i, p := range planck {
		ret[i] = NewDotFromPlanck(p)
	}
	return ret
}

// String formats a Dot with the correct unit.
// Works for positive and negative values.
func (d *Dot) String() string {
	prefices := []struct {
		thresh *big.Float
		prefix string
	}{
		{big.NewFloat(PlanckPerDot * 1e6), "MDot"},
		{big.NewFloat(PlanckPerDot * 1e3), "KDot"},
		{big.NewFloat(PlanckPerDot), "Dot"},
		{big.NewFloat(PlanckPerDot / 1e3), "mDot"},
		{big.NewFloat(PlanckPerDot / 1e6), "uDot"},
		{big.NewFloat(1), "Planck"},
	}

	planck := new(big.Float).SetInt(d.planck)
	planckAbs := new(big.Float).Abs(planck)
	for _, prefix := range prefices {
		if planckAbs.Cmp(prefix.thresh) >= 0 {
			value := new(big.Float).Quo(planck, prefix.thresh)
			return fmt.Sprintf("%s %s", formatFloat(value), prefix.prefix)
		}
	}

	return "0 Planck"
}

// Planck converts a Dot to Plancks.
func (d *Dot) Planck() *big.Int {
	return new(big.Int).Set(d.planck)
}

// Abs returns the absolute value.
func (d *Dot) Abs() *Dot {
	return NewDotFromPlanck(new(big.Int).Abs(d.planck))
}

func formatFloat(f *big.Float) string {
	return f.Text('f', PrintPrecision)
}
