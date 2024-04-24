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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDot_String tests the String() method in a non-exhaustive manner.
func TestDot_String(t *testing.T) {
	str := NewDotFromPlanck(big.NewInt(0)).String()
	assert.Equal(t, "0 Planck", str)

	str = NewDotFromPlanck(big.NewInt(999)).String()
	assert.Equal(t, "999.000 Planck", str)
	str = NewDotFromPlanck(big.NewInt(-999)).String()
	assert.Equal(t, "-999.000 Planck", str)

	str = NewDotFromPlanck(big.NewInt(PlanckPerDot)).String()
	assert.Equal(t, "1.000 Dot", str)
	str = NewDotFromPlanck(big.NewInt(-PlanckPerDot)).String()
	assert.Equal(t, "-1.000 Dot", str)

	str = NewDotFromPlanck(big.NewInt(PlanckPerDot / 2)).String()
	assert.Equal(t, "500.000 mDot", str)
	str = NewDotFromPlanck(big.NewInt(PlanckPerDot / -2)).String()
	assert.Equal(t, "-500.000 mDot", str)

	str = NewDotFromPlanck(big.NewInt(PlanckPerDot * 2000000)).String()
	assert.Equal(t, "2.000 MDot", str)
	str = NewDotFromPlanck(big.NewInt(PlanckPerDot * -2000000)).String()
	assert.Equal(t, "-2.000 MDot", str)
}

func TestDot_Abs(t *testing.T) {
	dot := NewDotFromPlanck(big.NewInt(-10))
	abs := dot.Abs()

	assert.NotEqual(t, abs, dot)
	assert.Equal(t, abs, NewDotFromPlanck(big.NewInt(10)))
}
