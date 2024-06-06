package constants

// Point represents a point in 2D space
type Point struct {
	X, Y int
}

// Direction represents the direction of the snake's movement
type Direction int

// Direction constants
const (
	Up Direction = iota
	Down
	Left
	Right
)
