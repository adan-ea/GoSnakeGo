package images

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game Sprites
const (
	backgroundImagePath = "resources/images/board/background_brown.png"
	gameOverImagePath   = "resources/images/ui/adriensexyy.png"
)

// Snake Sprites
const (
	headSpriteLeftPath = "resources/images/snake/head_sprite.png"

	bodySpriteBluePath   = "resources/images/snake/body_sprite_blue.png"
	bodySpritePurplePath = "resources/images/snake/body_sprite_purple.png"
	bodySpriteRedPath = "resources/images/snake/body_sprite_red.png"

	tailSpriteBluePath   = "resources/images/snake/tail_sprite_blue.png"
	tailSpritePurplePath = "resources/images/snake/tail_sprite_purple.png"
	tailSpriteRedPath = "resources/images/snake/tail_sprite_red.png"
)

// Food Sprites
const (
	foodSpritePath = "resources/images/food/apple.png"
)

// UI Sprites
const (
	numbersSpritePath = "resources/images/ui/numbers.png"
	starSpritePath    = "resources/images/ui/star.png"
	iconSpritePath    = "resources/images/ui/icon.png"
)

// Size of the numbers in the numbers sprite
var (
	DigitWidths = []int{22, 18, 21, 22, 24, 22, 23, 21, 23, 22}
	DigitHeight = 33
)

// Actual loaded images
var (
	BackgroundSprite *ebiten.Image
	GameOverSprite   *ebiten.Image
	HeadSprite       *ebiten.Image
	BodySprite       map[int]*ebiten.Image
	TailSprite       map[int]*ebiten.Image
	FoodSprite       *ebiten.Image
	NumbersSprite    *ebiten.Image
	StarSprite       *ebiten.Image
	IconSprite       *ebiten.Image
)

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

func LoadImages() {
	BackgroundSprite = LoadImage(backgroundImagePath)
	GameOverSprite = LoadImage(gameOverImagePath)

	HeadSprite = LoadImage(headSpriteLeftPath)

	BodySprite = map[int]*ebiten.Image{
		0: LoadImage(bodySpriteBluePath),
		1: LoadImage(bodySpritePurplePath),
		2: LoadImage(bodySpriteRedPath),
	}

	TailSprite = map[int]*ebiten.Image{
		0: LoadImage(tailSpriteBluePath),
		1: LoadImage(tailSpritePurplePath),
		2: LoadImage(tailSpriteRedPath),
	}

	FoodSprite = LoadImage(foodSpritePath)
	NumbersSprite = LoadImage(numbersSpritePath)
	StarSprite = LoadImage(starSpritePath)
	IconSprite = LoadImage(iconSpritePath)
}
