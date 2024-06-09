package game

import (
	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

// Food represents the food
type Food struct {
	x, y int
}

func NewFood(x, y int) *Food {
	return &Food{
		x: x,
		y: y,
	}
}

// Draw renders the food on the screen
func (f *Food) Draw(screen *ebiten.Image, offsetX, offsetY int) {
	sx := float64(offsetX + f.x*constants.TileSize)
	sy := float64(offsetY + f.y*constants.TileSize)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(sx, sy)
	screen.DrawImage(images.FoodSprite, op)
}
