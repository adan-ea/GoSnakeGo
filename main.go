package main

import (
	"image"
	"log"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/game"
	"github.com/adan-ea/GoSnakeGo/resources/audio"
	"github.com/adan-ea/GoSnakeGo/resources/fonts"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	images.InitImages()
	audio.InitAudio()
	fonts.InitFonts()

	g := game.NewGame()

	ebiten.SetWindowSize(constants.ScreenWidth, constants.ScreenHeight)
	ebiten.SetWindowIcon([]image.Image{images.IconSprite})
	ebiten.SetWindowTitle("Go Snake Go!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
