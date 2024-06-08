package game

import "github.com/hajimehoshi/ebiten/v2"

type Input struct{}

func NewInput() *Input {
	return &Input{}
}

func Dir() (Direction, bool) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return Up, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return Left, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return Right, true
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return Down, true
	}

	return 0, false
}

func Space() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}


