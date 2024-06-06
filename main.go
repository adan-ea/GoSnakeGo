package main

import (
	"log"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()

	g.Init()

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
