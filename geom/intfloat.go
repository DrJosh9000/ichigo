/*
Copyright 2021 Josh Deprez

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package geom

import (
	"fmt"
	"math"
)

// IntFloat represents a number as an integer part plus a fractional part.
// This can represent reals in the int range with decent precision.
type IntFloat struct {
	I int
	F float64
}

// ToIntFloat converts a float64 directly into an IntFloat.
func ToIntFloat(f float64) IntFloat {
	return IntFloat{I: 0, F: f}.Canon()
}

// Canon returns a value equal to x, in canonical form (0 ≤ F < 1).
// Each possible value only has one canonical form.
func (x IntFloat) Canon() IntFloat {
	i, f := math.Modf(x.F)
	if f < 0 {
		i--
		f = 1 + f
	}
	return IntFloat{I: x.I + int(i), F: f}
}

// Float converts the value into a float64.
func (x IntFloat) Float() float64 {
	return float64(x.I) + x.F
}

func (x IntFloat) String() string {
	return fmt.Sprintf("%d + %f", x.I, x.F)
}

// Lt reports x < y. x and y must be in canonical form for the
// comparison to be meaningful.
func (x IntFloat) Lt(y IntFloat) bool {
	if x.I == y.I {
		return x.F < y.F
	}
	return x.I < y.I
}

// Gt reports x > y. x and y must be in canonical form for the
// comparison to be meaningful.
func (x IntFloat) Gt(y IntFloat) bool {
	if x.I == y.I {
		return x.F > y.F
	}
	return x.I > y.I
}

// Add returns x+y (not canonicalised).
func (x IntFloat) Add(y IntFloat) IntFloat {
	return IntFloat{I: x.I + y.I, F: x.F + y.F}
}

// Neg returns -x (not canonicalised).
func (x IntFloat) Neg() IntFloat {
	return IntFloat{I: -x.I, F: -x.F}
}

// Sub returns x-y (not canonicalised).
func (x IntFloat) Sub(y IntFloat) IntFloat {
	return IntFloat{I: x.I - y.I, F: x.F - y.F}
}

// Mul returns x*y, canonicalised.
func (x IntFloat) Mul(y IntFloat) IntFloat {
	return IntFloat{
		I: x.I * y.I,
		F: float64(x.I)*y.F + x.F*float64(y.I) + x.F*y.F,
	}.Canon()
}

// Inv returns 1/x, canonicalised.
func (x IntFloat) Inv() IntFloat {
	return ToIntFloat(1 / x.Float())
}

// Div returns x/y, canonicalised.
func (x IntFloat) Div(y IntFloat) IntFloat {
	return ToIntFloat(x.Float() / y.Float())
}
