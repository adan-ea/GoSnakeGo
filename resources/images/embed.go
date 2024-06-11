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
	bodySpriteRedPath    = "resources/images/snake/body_sprite_red.png"

	tailSpriteBluePath   = "resources/images/snake/tail_sprite_blue.png"
	tailSpritePurplePath = "resources/images/snake/tail_sprite_purple.png"
	tailSpriteRedPath    = "resources/images/snake/tail_sprite_red.png"
)

// Food Sprites
const (
	foodSpritePath = "resources/images/food/apple.png"
)

// UI Sprites
const (
	numbersSpritePath = "resources/images/ui/numbers.png"
	trophySpritePath  = "resources/images/ui/trophy.png"
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
	TrophySprite     *ebiten.Image
	IconSprite       *ebiten.Image
)

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(err)
	}
	return img
}

func InitImages() {
	BackgroundSprite = loadImage(backgroundImagePath)
	GameOverSprite = loadImage(gameOverImagePath)

	HeadSprite = loadImage(headSpriteLeftPath)

	BodySprite = map[int]*ebiten.Image{
		0: loadImage(bodySpriteBluePath),
		1: loadImage(bodySpritePurplePath),
		2: loadImage(bodySpriteRedPath),
	}

	TailSprite = map[int]*ebiten.Image{
		0: loadImage(tailSpriteBluePath),
		1: loadImage(tailSpritePurplePath),
		2: loadImage(tailSpriteRedPath),
	}

	FoodSprite = loadImage(foodSpritePath)
	NumbersSprite = loadImage(numbersSpritePath)
	TrophySprite = loadImage(trophySpritePath)
	IconSprite = loadImage(iconSpritePath)
}
