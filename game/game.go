package game

import (
	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/food"
	"github.com/adan-ea/GoSnakeGo/snake"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BaseSpeed     = 10 // Base speed (higher means slower)
	SpeedIncrease = 1  // Speed increase factor
	ScoreInterval = 5  // Points interval to increase speed

	bestScorePath = "resources/scoreboard.txt"
	nbScoreSaved  = 5
)

// Game represents the game state and logic
type Game struct {
	snake            *snake.Snake
	food             *food.Food
	score            int
	highScore        int
	gameOver         bool
	updateTick       int
	backgroundSprite *ebiten.Image
	numberSprite     *ebiten.Image
	starSprite       *ebiten.Image
}

func NewGame() *Game {
	return &Game{}
}

// Init initializes the game state
func (g *Game) Init() {
	g.snake = snake.InitSnake()
	g.food = food.SpawnFood(g.snake.Body)
	g.backgroundSprite = constants.LoadImage(constants.BackgroundImagePath)
	g.numberSprite = constants.LoadImage(constants.NumbersSpritePath)
	g.starSprite = constants.LoadImage(constants.StarSpritePath)
	g.gameOver = false
	g.score = 0
	g.highScore = getHighestScore()
	g.updateTick = 0
}

// Update updates the game state
func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.Init()
		}
		return nil
	}

	// Handle key presses to change direction
	HandleKeyPressed(g.snake)

	// Calculate the current speed based on the score
	currentSpeed := BaseSpeed - (g.score/ScoreInterval)*SpeedIncrease
	if currentSpeed < 1 {
		currentSpeed = 1
	}

	// Perform a snake movement tick
	g.updateTick++
	if g.updateTick%currentSpeed != 0 {
		return nil
	}
	g.snake.Update()
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

	g.snake.Draw(screen)
	g.food.Draw(screen)
	g.drawScore(screen, g.score, 40, 7)
	g.drawHighScore(screen, g.highScore, 550, 7)
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
	head := g.snake.GetHead()
	if head.X < 0 || head.X >= constants.GameWidth/constants.TileSize ||
		head.Y < 0 || head.Y >= constants.GameHeight/constants.TileSize {
		g.gameOver = true
		saveHighScore(g.score)
	}

	// Check if the snake has collided with itself
	for i := 1; i < len(g.snake.Body); i++ {
		if head == g.snake.Body[i] {
			g.gameOver = true
			saveHighScore(g.score)
			break
		}
	}

	// Eating food
	if g.snake.GetHead() == g.food.GetPosition() {
		g.snake.Body = append(g.snake.Body, g.food.GetPosition())
		g.food.Respawn(g.snake.Body)
		g.updateScore()
	}
}
