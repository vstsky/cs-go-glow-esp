package cs

type Entity struct {
	Health   uint32
	Team     uint32
	Position EntityPosition
	Pointer  uint32
}

type EntityPosition struct {
	X float32
	Y float32
	Z float32
}

type Color struct {
	r float32
	g float32
	b float32
}
