package geom

import "strconv"

// Rat is a (small) rational number implementation. Overflow can happen.
type Rat struct{ N, D int }

// IntRat returns the rational representation of n.
func IntRat(n int) Rat { return Rat{N: n, D: 1} }

// String returns a nice string representation like "-3/5".
func (r Rat) String() string {
	if r.D == 1 {
		return strconv.Itoa(r.N)
	}
	return strconv.Itoa(r.N) + "/" + strconv.Itoa(r.D)
}

// Int returns r.N / r.D.
func (r Rat) Int() int { return r.N / r.D }

// Rem returns r.N % r.D.
func (r Rat) Rem() int { return r.N % r.D }

// Canon puts the rational number into reduced form.
func (r Rat) Canon() Rat {
	if r.D == 0 {
		panic("division by zero")
	}
	if r.N == 0 {
		r.D = 1
		return r
	}
	if r.D < 0 {
		r.N, r.D = -r.N, -r.D
	}
	if d := gcd(abs(r.N), r.D); d > 1 {
		r.N /= d
		r.D /= d
	}
	return r
}

// Neg returns -r.
func (r Rat) Neg() Rat {
	r.N = -r.N
	return r
}

// Add returns r + q.
func (r Rat) Add(q Rat) Rat {
	return Rat{
		N: r.N*q.D + q.N*r.D,
		D: r.D * q.D,
	}.Canon()
}

// Sub returns r - q.
func (r Rat) Sub(q Rat) Rat {
	return Rat{
		N: r.N*q.D - q.N*r.D,
		D: r.D * q.D,
	}.Canon()
}

// Mul returns r * q.
func (r Rat) Mul(q Rat) Rat {
	return Rat{
		N: r.N * q.N,
		D: r.D * q.D,
	}.Canon()
}

// Invert returns 1/r if it exists, otherwise panics.
func (r Rat) Invert() Rat {
	return Rat{N: r.D, D: r.N}.Canon()
}

// Div returns r/q.
func (r Rat) Div(q Rat) Rat {
	return Rat{
		N: r.N * q.D,
		D: r.D * q.N,
	}.Canon()
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func gcd(a, b int) int {
	if a < b {
		a, b = b, a
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
