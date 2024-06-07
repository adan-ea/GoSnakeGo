package game

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/adan-ea/GoSnakeGo/constants"
	raudio "github.com/adan-ea/GoSnakeGo/resources/audio"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	bestScorePath = "resources/scoreboard.txt"
	nbScoreSaved  = 5
)

// Game represents the game state and logic
type Game struct {
	mode      Mode
	rows      int
	cols      int
	snake     *Snake
	food      *Food
	score     int
	highScore int
	timer     time.Time

	audioContext   *audio.Context
	hitPlayer      *audio.Player
	eatPlayer      *audio.Player
	gameOverPlayer *audio.Player
	themePlayer    *audio.Player
}

func NewGame(rows int, cols int) *Game {
	game := &Game{
		rows:      rows,
		cols:      cols,
		timer:     time.Now(),
		highScore: getHighestScore(),
		snake:     NewSnake(),
	}
	game.initAudio()
	game.placeFood()

	return game
}

// initAudio initializes the audio context and players
func (g *Game) initAudio() {
	if g.audioContext == nil {
		g.audioContext = audio.NewContext(48000)
	}

	hitD, err := vorbis.DecodeWithoutResampling(bytes.NewReader(raudio.Hit_ogg))
	if err != nil {
		log.Fatal(err)
	}

	g.hitPlayer, err = g.audioContext.NewPlayer(hitD)
	if err != nil {
		log.Fatal(err)
	}

	eatD, err := vorbis.DecodeWithoutResampling(bytes.NewReader(raudio.Eat_ogg))
	if err != nil {
		log.Fatal(err)
	}

	g.eatPlayer, err = g.audioContext.NewPlayer(eatD)
	if err != nil {
		log.Fatal(err)
	}

	gameOverD, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.GameOver_wav))
	if err != nil {
		log.Fatal(err)
	}

	themeD, err := wav.DecodeWithoutResampling(bytes.NewReader(raudio.TetrisTheme_wav))
	if err != nil {
		log.Fatal(err)
	}

	g.themePlayer, err = g.audioContext.NewPlayer(themeD)
	if err != nil {
		log.Fatal(err)
	}

	g.gameOverPlayer, err = g.audioContext.NewPlayer(gameOverD)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	BaseInterval  = time.Millisecond * 200
	MinInterval   = time.Millisecond * 50
	SpeedIncrease = time.Millisecond * 5
)

func CalculateInterval(score int) time.Duration {
	// Calculate the new interval by decreasing it linearly with the score
	newInterval := BaseInterval - time.Duration(score)*SpeedIncrease

	// Ensure the interval does not go below the minimum interval
	if newInterval < MinInterval {
		return MinInterval
	}
	return newInterval
}

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
			g.themePlayer.Rewind()
		}
	case ModeGame:
		if !g.themePlayer.IsPlaying() {
			g.themePlayer.Rewind()
			g.themePlayer.Play()
		}

		// snake goes faster when there are more points
		interval := CalculateInterval(g.score)
		if newDir, ok := Dir(); ok {
			g.snake.ChangeDirection(newDir)
		}

		if time.Since(g.timer) >= interval {
			if err := g.moveSnake(); err != nil {
				return err
			}

			g.timer = time.Now()
		}

	case ModeGameOver:
		g.themePlayer.Pause()
		g.gameOverPlayer.Play()
		
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			NewGame(g.rows, g.cols)
			g.mode = ModeGame
			g.themePlayer.Rewind()
			g.themePlayer.Play()
		}
	}

	return nil
}

func (g *Game) placeFood() {
	var x, y int
	validPosition := false

	for !validPosition {
		x = rand.Intn(g.rows)
		y = rand.Intn(g.cols)

		validPosition = true
		for _, p := range g.snake.Body {
			if p.x == x && p.y == y {
				validPosition = false
				break
			}
		}
	}

	g.food = NewFood(x, y)
}

func (g *Game) moveSnake() error {
	// remove tail first, add 1 in front
	g.snake.Move()
	if g.snakeLeftBoard() || g.snake.HeadHitsBody() {
		g.mode = ModeGameOver
		saveHighScore(g.score)
		return nil
	}

	if g.snake.HeadHits(g.food.x, g.food.y) {
		// the snake grows on the next move
		g.snake.justAte = true

		g.placeFood()
		g.updateScore()
	}

	return nil
}

func (g *Game) snakeLeftBoard() bool {
	head := g.snake.Head()
	return head.x > g.cols-1 || head.y > g.rows-1 || head.x < 0 || head.y < 0
}

func (g *Game) Draw(screen *ebiten.Image) {

	switch g.mode {
	case ModeGame:
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

		g.snake.Draw(screen)
		g.food.Draw(screen)
		g.drawScore(screen, g.score, 40, 7)
		g.drawHighScore(screen, g.highScore, 550, 7)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
