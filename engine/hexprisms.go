package engine

type HexPrismMap struct {
	Map map[Point3]*HexPrism
}

type HexPrism struct {
	pos Point3
	hpm *HexPrismMap
}
