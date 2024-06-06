package game

import (
	"fmt"
	"image"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/food"
	"github.com/adan-ea/GoSnakeGo/snake"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BaseSpeed     = 10 // Base speed (higher means slower)
	SpeedIncrease = 1  // Speed increase factor
	ScoreInterval = 5  // Points interval to increase speed
)

// Game represents the game state and logic
type Game struct {
	Snake            *snake.Snake
	Food             *food.Food
	Score            int
	GameOver         bool
	updateTick       int
	backgroundSprite *ebiten.Image
	numberSprite     *ebiten.Image
}

func NewGame() *Game {
	return &Game{}
}

// Init initializes the game state
func (g *Game) Init() {
	g.Snake = snake.InitSnake()
	g.Food = food.SpawnFood(g.Snake.Body)
	g.backgroundSprite = constants.LoadImage(constants.BackgroundImagePath)
	g.numberSprite = constants.LoadImage(constants.NumbersSpritePath)
	g.GameOver = false
	g.Score = 0
	g.updateTick = 0
}

// Update updates the game state
func (g *Game) Update() error {
	if g.GameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Init()
		}
		return nil
	}

	// Handle key presses to change direction
	HandleKeyPressed(g.Snake)

	// Calculate the current speed based on the score
	currentSpeed := BaseSpeed - (g.Score/ScoreInterval)*SpeedIncrease
	if currentSpeed < 1 {
		currentSpeed = 1
	}

	// Perform a snake movement tick
	g.updateTick++
	if g.updateTick%currentSpeed != 0 {
		return nil
	}
	g.Snake.Update()
	g.HandleCollision()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Calculate the offset to center the game area
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2

	// Draw the background in the game area
	op := &ebiten.DrawImageOptions{}
	for y := 0; y < constants.GameHeight; y += constants.TileSize * 2 {
		for x := 0; x < constants.GameWidth; x += constants.TileSize * 2 {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x+offsetX), float64(y+offsetY))
			screen.DrawImage(g.backgroundSprite, op)
		}
	}

	g.Snake.Draw(screen)
	g.Food.Draw(screen)
	g.DrawScore(screen, g.Score, 40, 7)
}

func (g *Game) DrawScore(screen *ebiten.Image, score int, x, y int) {
	var digitWidths = []int{22, 18, 21, 22, 24, 22, 23, 21, 23, 22}
	digitHeight := 33

	scoreStr := fmt.Sprintf("%d", score)
	currentX := x

	for _, char := range scoreStr {
		digit := int(char - '0')
		digitWidth := digitWidths[digit]
		sx := 0
		for j := 0; j < digit; j++ {
			sx += digitWidths[j]
		}
		sy := 0
		numImage := g.numberSprite.SubImage(image.Rect(sx, sy, sx+digitWidth, sy+digitHeight)).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(currentX), float64(y))
		screen.DrawImage(numImage, op)

		currentX += digitWidth
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func HandleKeyPressed(s *snake.Snake) bool {
	// Control snake
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		s.ChangeDirection(constants.Up)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		s.ChangeDirection(constants.Down)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		s.ChangeDirection(constants.Left)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		s.ChangeDirection(constants.Right)
	}
	return true
}

func (g *Game) HandleCollision() {
	// Check if the snake has collided with the walls
	head := g.Snake.GetHead()
	if head.X < 0 || head.X >= constants.GameWidth/constants.TileSize ||
		head.Y < 0 || head.Y >= constants.GameHeight/constants.TileSize {
		g.GameOver = true
	}

	// Check if the snake has collided with itself
	for i := 1; i < len(g.Snake.Body); i++ {
		if head == g.Snake.Body[i] {
			g.GameOver = true
			break
		}
	}

	// Eating food
	if g.Snake.GetHead() == g.Food.GetPosition() {
		g.Snake.Body = append(g.Snake.Body, g.Food.GetPosition())
		g.Food.Respawn(g.Snake.Body)
		g.Score++
	}
}
