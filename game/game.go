package game

import (
	"bytes"
	"fmt"
	"log"

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
	BaseSpeed     = 10 // Base speed (higher means slower)
	SpeedIncrease = 1  // Speed increase factor
	ScoreInterval = 5  // Points interval to increase speed

	bestScorePath = "resources/scoreboard.txt"
	nbScoreSaved  = 5
)

// Game represents the game state and logic
type Game struct {
	mode       Mode
	snake      *Snake
	food       *Food
	score      int
	highScore  int
	updateTick int

	backgroundSprite *ebiten.Image
	numberSprite     *ebiten.Image
	starSprite       *ebiten.Image
	audioContext     *audio.Context
	hitPlayer        *audio.Player
	eatPlayer        *audio.Player
	gameOverPlayer   *audio.Player
	themePlayer      *audio.Player
}

func NewGame() *Game {
	return &Game{}
}

// Init initializes the game state
func (g *Game) Init() {

	g.snake = initSnake()
	g.food = spawnFood(g.snake.Body)
	g.score = 0
	g.highScore = getHighestScore()
	g.updateTick = 0

	g.backgroundSprite = images.LoadImage(images.BackgroundImagePath)
	g.numberSprite = images.LoadImage(images.NumbersSpritePath)
	g.starSprite = images.LoadImage(images.StarSpritePath)

	g.initAudio()

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

		if newDir, ok := Dir(); ok {
			g.snake.ChangeDirection(newDir)
		}

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

		// Check for game over state
		if g.isGameOver() {
			g.themePlayer.Pause()
			g.gameOverPlayer.Rewind()

			g.mode = ModeGameOver
			g.hitPlayer.Play()
			saveHighScore(g.score)
		}

		g.eatApple()

	case ModeGameOver:
		g.gameOverPlayer.Play()
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.mode = ModeGame
			g.Init()
			g.themePlayer.Rewind()
			g.themePlayer.Play()
		}
	}

	return nil
}

// looks if the snake ate an apple and updates the score
func (g *Game) eatApple() {
	// Eating food
	if g.snake.getHead() == g.food.getFoodPosition() {
		if err := g.eatPlayer.Rewind(); err != nil {
			return
		}
		g.eatPlayer.Play()
		g.snake.Body = append(g.snake.Body, g.food.getFoodPosition())
		g.food.changePosition(g.snake.Body)
		g.updateScore()
	}
}

// Add a method to check if the game is over
func (g *Game) isGameOver() bool {
	head := g.snake.getHead()
	// Check if the snake has collided with the walls
	if head.X < 0 || head.X >= constants.GameWidth/constants.TileSize ||
		head.Y < 0 || head.Y >= constants.GameHeight/constants.TileSize {
		return true
	}

	// Check if the snake has collided with itself
	for i := 1; i < len(g.snake.Body); i++ {
		if head == g.snake.Body[i] {
			return true
		}
	}

	return false
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
				screen.DrawImage(g.backgroundSprite, op)
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
