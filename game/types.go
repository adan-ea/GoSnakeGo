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

// Size represents the size of the board
type Size int

const (
	Small Size = iota
	Medium
	Large
)

func getGridSize(size Size) (int, int) {
	switch size {
	case Small:
		return 14, 14
	case Medium:
		return 16, 16
	case Large:
		return 18, 18
	}
	return 18, 18
}

func getSizeText(size Size) string {
	switch size {
	case Small:
		return "Small"
	case Medium:
		return "Medium"
	case Large:
		return "Large"
	}
	return "Large"
}

func getTextToSize(text string) Size {
	switch text {
	case "Small":
		return Small
	case "Medium":
		return Medium
	case "Large":
		return Large
	}
	return Large
}

func getSizeFromRowsCols(rows, cols int) Size {
	if rows == 14 && cols == 14 {
		return Small
	} else if rows == 16 && cols == 16 {
		return Medium
	}
	return Large
}

// Direction represents the direction of the snake's movement
type Direction int

const (
	Right Direction = iota
	Left
	Down
	Up
)

// Color represents possible colors for the snake
type Color int

const (
	Blue Color = iota
	Purple
	Red
)
