package main

import (
	"log"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/game"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	images.LoadImages()
	g := game.NewGame(18, 18)

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowTitle("Snake Game")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
