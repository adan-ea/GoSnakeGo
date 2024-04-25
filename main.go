package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridSize     = 20
)

type Game struct {
	snake      []Point
	dir        Point
	food       Point
	width      int
	height     int
	updateTick int
	score      int
}

type Point struct {
	X, Y int
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (g *Game) Update() error {
	g.updateTick++
	if g.updateTick%10 != 0 {
		return nil
	}

	newHead := Point{
		X: g.snake[0].X + g.dir.X,
		Y: g.snake[0].Y + g.dir.Y,
	}

	// Collision with boundaries
	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		return g.initGame()
	}

	// Collision with itself
	for _, v := range g.snake[1:] {
		if v == newHead {
			return g.initGame()
		}
	}

	// Move snake
	g.snake = append([]Point{newHead}, g.snake[:len(g.snake)-1]...)

	// Eating food
	if newHead == g.food {
		g.snake = append(g.snake, g.food)
		g.spawnFood()
		g.score++
		fmt.Println(g.score)
	}

	// Control snake
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.dir.Y == 0 {
		g.dir = Point{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.dir.Y == 0 {
		g.dir = Point{X: 0, Y: 1}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && g.dir.X == 0 {
		g.dir = Point{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && g.dir.X == 0 {
		g.dir = Point{X: 1, Y: 0}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 10, G: 30, B: 10, A: 255})

	// Draw food
	foodRect := image.Rect(g.food.X*gridSize, g.food.Y*gridSize, (g.food.X+1)*gridSize, (g.food.Y+1)*gridSize)
	ebitenutil.DrawRect(screen, float64(foodRect.Min.X), float64(foodRect.Min.Y), float64(gridSize), float64(gridSize), color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw snake
	for _, p := range g.snake {
		snakeRect := image.Rect(p.X*gridSize, p.Y*gridSize, (p.X+1)*gridSize, (p.Y+1)*gridSize)
		ebitenutil.DrawRect(screen, float64(snakeRect.Min.X), float64(snakeRect.Min.Y), float64(gridSize), float64(gridSize), color.RGBA{R: 0, G: 255, B: 0, A: 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) initGame() error {
	g.snake = []Point{{X: 5, Y: 5}, {X: 4, Y: 5}, {X: 3, Y: 5}}
	g.dir = Point{X: 1, Y: 0}
	g.width = screenWidth / gridSize
	g.height = screenHeight / gridSize
	g.spawnFood()
	g.score = 0
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
			g.food = Point{X: x, Y: y}
			break
		}
	}
}

func main() {
	game := &Game{}
	game.initGame()

	ebiten.RunGame(game)
}
