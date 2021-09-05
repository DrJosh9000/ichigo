package engine

var (
	// Oblique projections
	CabinetProjection  = ParallelProjection{0.5, 0.5}
	CavalierProjection = ParallelProjection{1, 1}

	// Axonometric projections
	DimetricProjection  = ParallelProjection{0, 0.5}
	ElevationProjection = ParallelProjection{0, 0}
	IsometricProjection = ParallelProjection{0, 0.707106781186548} // 1 ÷ √2
	TrimetricProjection = ParallelProjection{0, 1}
	HexPrismProjection  = ParallelProjection{0, 0.816496580927726} // √2 ÷ √3
)

type ParallelProjection struct {
	ZX, ZY float64
}

func (π ParallelProjection) Project(p Point3) (px, py float64) {
	px = float64(p.X) + π.ZX*float64(p.Z)
	py = float64(p.Y) + π.ZY*float64(p.Z)
	return px, py
}
