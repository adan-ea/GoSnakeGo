package constants

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game Sprites
const (
	BackgroundImagePath = "resources/images/background_brown.png"
	GameOverImagePath   = "resources/images/adriensexyy.png"
)

// Snake Sprites
const (
	HeadSpriteLeftPath = "resources/images/snake/head_sprite.png"
	BodySpritePath     = "resources/images/snake/body_sprite.png"
	TailSpritePath     = "resources/images/snake/tail_sprite.png"
)

// Food Sprites
const (
	FoodSpritePath = "resources/images/food/apple.png"
)

// UI Sprites
const (
	NumbersSpritePath = "resources/images/ui/numbers.png"
)

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}
