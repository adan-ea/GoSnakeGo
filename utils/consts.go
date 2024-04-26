package utils

const (
	ScreenWidth  = 800
	ScreenHeight = 1024
	GridSize     = 20
)

type Point struct {
	X, Y int
}

var Up = Point{X: 0, Y: -1}
var Down = Point{X: 0, Y: 1}
var Left = Point{X: -1, Y: 0}
var Right = Point{X: 1, Y: 0}
