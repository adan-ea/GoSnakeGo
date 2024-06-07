package game

import (
	"math/rand"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

// Food represents the food
type Food struct {
	position Point
	Sprite   *ebiten.Image
}

// spawnFood creates a new instance of Food
func spawnFood(snakeBody []Point) *Food {
	food := &Food{}

	food.Sprite = images.LoadImage(images.FoodSpritePath)
	for {
		x := rand.Intn(constants.GameWidth / constants.TileSize)
		y := rand.Intn(constants.GameHeight / constants.TileSize)
		occupied := false
		for _, p := range snakeBody {
			if p.X == x && p.Y == y {
				occupied = true
				break
			}
		}
		if !occupied {
			food.position = Point{X: x, Y: y}
			break
		}
	}
	return food
}

// changePosition repositions the food to a new location avoiding the snake's body
func (f *Food) changePosition(snakeBody []Point) {
	for {
		x := rand.Intn(constants.GameWidth / constants.TileSize)
		y := rand.Intn(constants.GameHeight / constants.TileSize)
		occupied := false
		for _, p := range snakeBody {
			if p.X == x && p.Y == y {
				occupied = true
				break
			}
		}
		if !occupied {
			f.position = Point{X: x, Y: y}
			break
		}
	}
}

// Draw renders the food on the screen
func (f *Food) Draw(screen *ebiten.Image) {
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2
	sx := float64(offsetX + f.position.X*constants.TileSize)
	sy := float64(offsetY + f.position.Y*constants.TileSize)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(f.Sprite, op)
}

// getFoodPosition returns the position of the food
func (f *Food) getFoodPosition() Point {
	return f.position
}
