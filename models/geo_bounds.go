package models

type Point struct {
	Lat float32
	Lon float32
}

type SquareField struct {
	HighLeftPoint  *Point
	DownRightPoint *Point
}

type CircleField struct {
	Center *Point
	Radius int
}
