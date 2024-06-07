package game

import (
	"math/rand"
	"time"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

type Board struct {
	rows      int
	cols      int
	food      *Food
	snake     *Snake
	score     int
	highScore int
	gameOver  bool
	timer     time.Time
}

func NewBoard(rows int, cols int) *Board {
	game := &Board{
		rows:      rows,
		cols:      cols,
		timer:     time.Now(),
		highScore: getHighestScore(),
		snake:     NewSnake(),
	}
	game.placeFood()

	return game
}

func (b *Board) Update(input *Input) error {
	if b.gameOver {
		return nil
	}

	// snake goes faster when there are more points
	interval := calculateInterval(b.score)
	if newDir, ok := Dir(); ok {
		b.snake.ChangeDirection(newDir)
	}

	if time.Since(b.timer) >= interval {
		if err := b.moveSnake(); err != nil {
			return err
		}

		b.timer = time.Now()
	}

	return nil
}

func (b *Board) moveSnake() error {
	// remove tail first, add 1 in front
	b.snake.Move()
	if b.snakeLeftBoard() || b.snake.HeadHitsBody() {
		b.gameOver = true
		saveHighScore(b.score)
		return nil
	}

	if b.snake.HeadHits(b.food.x, b.food.y) {
		// the snake grows on the next move
		b.snake.justAte = true

		b.placeFood()
		b.updateScore()
	}

	return nil
}

func (b *Board) snakeLeftBoard() bool {
	head := b.snake.Head()
	return head.x > b.cols-1 || head.y > b.rows-1 || head.x < 0 || head.y < 0
}

func (b *Board) placeFood() {
	var x, y int
	validPosition := false

	for !validPosition {
		x = rand.Intn(b.rows)
		y = rand.Intn(b.cols)

		validPosition = true
		for _, p := range b.snake.body {
			if p.x == x && p.y == y {
				validPosition = false
				break
			}
		}
	}

	b.food = NewFood(x, y)
}

const (
	BaseInterval  = time.Millisecond * 200
	MinInterval   = time.Millisecond * 50
	SpeedIncrease = time.Millisecond * 5
)

func calculateInterval(score int) time.Duration {
	// Calculate the new interval by decreasing it linearly with the score
	newInterval := BaseInterval - time.Duration(score)*SpeedIncrease

	// Ensure the interval does not go below the minimum interval
	if newInterval < MinInterval {
		return MinInterval
	}
	return newInterval
}

func (b *Board) Draw(screen *ebiten.Image) {
	// Calculate the offset to center the game area
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2

	// Draw the background in the game area
	op := &ebiten.DrawImageOptions{}
	for y := 0; y < constants.GameHeight; y += constants.TileSize * 2 {
		for x := 0; x < constants.GameWidth; x += constants.TileSize * 2 {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(x+offsetX), float64(y+offsetY))
			screen.DrawImage(images.BackgroundSprite, op)
		}
	}

	b.snake.Draw(screen)
	b.food.Draw(screen)
	b.drawScore(screen, b.score, 40, 7)
	b.drawHighScore(screen, b.highScore, 550, 7)
}
