package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"gosnakego/utils"
)

type Game struct {
	snake      []utils.Point
	dir        utils.Point
	food       utils.Point
	width      int
	height     int
	updateTick int
	score      int
	bestScore  int
	gameOver   bool
}

var img *ebiten.Image

func init() {
	rand.Seed(time.Now().UnixNano())
	var err error
	img, _, err = ebitenutil.NewImageFromFile("assets/adriensexyy.png")
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.initGame()
		}
		return nil
	}

	g.updateTick++
	if g.updateTick%10 != 0 {
		return nil
	}

	newHead := utils.Point{
		X: g.snake[0].X + g.dir.X,
		Y: g.snake[0].Y + g.dir.Y,
	}

	// Collision with boundaries
	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		g.gameOver = true
		return nil
	}

	// Collision with itself
	for _, v := range g.snake[1:] {
		if v == newHead {
			g.gameOver = true
			return nil
		}
	}

	// Move snake
	g.snake = append([]utils.Point{newHead}, g.snake[:len(g.snake)-1]...)

	// Eating food
	if newHead == g.food {
		g.snake = append(g.snake, g.food)
		g.spawnFood()
		g.score++

		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}

	// Control snake
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && g.dir.Y == 0 {
		g.dir = utils.Up
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && g.dir.Y == 0 {
		g.dir = utils.Down
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)) && g.dir.X == 0 {
		g.dir = utils.Left
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)) && g.dir.X == 0 {
		g.dir = utils.Right
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(utils.GetDarkGreenColor())

	if g.gameOver {
		screen.DrawImage(img, nil)

		oneFifthHeight := utils.ScreenHeight / 8
		// Dessiner le texte à la position calculée
		text.Draw(screen, "Game Over", basicFont(), utils.ScreenWidth/2-25, oneFifthHeight, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicFont(), utils.ScreenWidth/2-20, oneFifthHeight+20, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Best Score: %d", g.bestScore), basicFont(), utils.ScreenWidth/2-30, oneFifthHeight+40, utils.GetBlackColor())
		text.Draw(screen, "Press SPACE to Restart", basicFont(), utils.ScreenWidth/2-80, oneFifthHeight+60, utils.GetBlackColor())

		return
	}

	// Draw food
	foodRect := image.Rect(g.food.X*utils.GridSize, g.food.Y*utils.GridSize, (g.food.X+1)*utils.GridSize, (g.food.Y+1)*utils.GridSize)
	ebitenutil.DrawRect(screen, float64(foodRect.Min.X), float64(foodRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetRedColor())

	// Draw snake
	for _, p := range g.snake {
		snakeRect := image.Rect(p.X*utils.GridSize, p.Y*utils.GridSize, (p.X+1)*utils.GridSize, (p.Y+1)*utils.GridSize)
		ebitenutil.DrawRect(screen, float64(snakeRect.Min.X), float64(snakeRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetGreenColor())
	}

	// Draw score
	text.Draw(screen, fmt.Sprintf("Score: %d, Best Score: %d", g.score, g.bestScore), basicFont(), 10, 20, utils.GetWhiteColor())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return utils.ScreenWidth, utils.ScreenHeight
}

func (g *Game) initGame() error {
	g.snake = []utils.Point{{X: 5, Y: 5}, {X: 4, Y: 5}, {X: 3, Y: 5}}
	g.dir = utils.Right
	g.width = utils.ScreenWidth / utils.GridSize
	g.height = utils.ScreenHeight / utils.GridSize
	g.spawnFood()
	g.score = 0
	g.gameOver = false
	return nil
}

func (g *Game) spawnFood() {
	for {
		x := rand.Intn(g.width)
		y := rand.Intn(g.height)
		occupied := false
		for _, p := range g.snake {
			if p.X == x && p.Y == y {
				occupied = true
				break
			}
		}
		if !occupied {
			g.food = utils.Point{X: x, Y: y}
			break
		}
	}
}

func basicFont() font.Face {
	const dpi = 72
	return basicfont.Face7x13
}

func main() {
	game := &Game{}
	game.initGame()

	ebiten.RunGame(game)
}
