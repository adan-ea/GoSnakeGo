package game

import "math/rand"

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

const nbSize = 5
const (
	Small Size = iota
	Medium
	Large
	ExtraLarge
	RandomSize
)

func getGridSize(size Size) (int, int) {
	if size == RandomSize {
		size = Size(rand.Intn(nbSize - 1))
	}
	switch size {
	case Small:
		return 14, 14
	case Medium:
		return 16, 16
	case Large:
		return 18, 18
	case ExtraLarge:
		return 20, 20
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
	case ExtraLarge:
		return "Extra Large"
	case RandomSize:
		return "Random"
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
	case "Extra Large":
		return ExtraLarge
	}
	return Large
}

func getSizeFromRowsCols(rows, cols int) Size {
	if rows == 14 && cols == 14 {
		return Small
	} else if rows == 16 && cols == 16 {
		return Medium
	} else if rows == 18 && cols == 18 {
		return Large
	}
	return ExtraLarge
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

// Number of colors available including Random
const nbColors = 4
const (
	Blue Color = iota
	Purple
	Red
	RandomColor
)

func getColorText(color Color) string {
	switch color {
	case Blue:
		return "Blue"
	case Purple:
		return "Purple"
	case Red:
		return "Red"
	case RandomColor:
		return "Random"
	}
	return "Blue"
}
