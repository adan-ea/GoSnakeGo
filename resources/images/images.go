package images

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game Sprites
const (
	BackgroundImagePath = "resources/images/board/background_brown.png"
	GameOverImagePath   = "resources/images/ui/adriensexyy.png"
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
	DigitHeight       = 33
	StarSpritePath    = "resources/images/ui/star.png"
)

var (
	DigitWidths = []int{22, 18, 21, 22, 24, 22, 23, 21, 23, 22}
)

var (
	BackgroundSprite *ebiten.Image
	GameOverSprite   *ebiten.Image
	HeadSprite       *ebiten.Image
	BodySprite       *ebiten.Image
	TailSprite       *ebiten.Image
	FoodSprite       *ebiten.Image
	NumbersSprite    *ebiten.Image
	StarSprite       *ebiten.Image
)

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

func LoadImages() {
	BackgroundSprite = LoadImage(BackgroundImagePath)
	GameOverSprite = LoadImage(GameOverImagePath)
	HeadSprite = LoadImage(HeadSpriteLeftPath)
	BodySprite = LoadImage(BodySpritePath)
	TailSprite = LoadImage(TailSpritePath)
	FoodSprite = LoadImage(FoodSpritePath)
	NumbersSprite = LoadImage(NumbersSpritePath)
	StarSprite = LoadImage(StarSpritePath)
}
