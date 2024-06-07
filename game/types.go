package game

// Mode represents the game mode
type Mode int

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

// Point represents a point in 2D space
type Point struct {
	x, y int
}

// Direction represents the direction of the snake's movement
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)
