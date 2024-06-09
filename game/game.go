package game

import (
	"bytes"
	"image/color"
	"log"
	"strconv"

	"github.com/adan-ea/GoSnakeGo/constants"
	raudio "github.com/adan-ea/GoSnakeGo/resources/audio"
	"github.com/adan-ea/GoSnakeGo/resources/fonts"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	bestScorePath = "resources/scoreboard.txt"
	nbScoreSaved  = 5
	boardRows     = 18
	boardCols     = 18
)

// Game represents the game state and logic
type Game struct {
	input *Input
	board *Board
	mode  Mode

	audioContext   *audio.Context
	hitPlayer      *audio.Player
	eatPlayer      *audio.Player
	gameOverPlayer *audio.Player
	themePlayer    *audio.Player
}

func NewGame(rows int, cols int) *Game {
	game := &Game{
		input: NewInput(),
		board: NewBoard(boardRows, boardCols),
	}
	game.initAudio()

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

func (g *Game) Update() error {
	switch g.mode {
	case ModeTitle:
		if Space() {
			g.mode = ModeGame
			g.themePlayer.Rewind()
		}
	case ModeGame:
		if g.board.gameOver {
			g.gameOverPlayer.Rewind()
			g.gameOverPlayer.Play()
			g.mode = ModeGameOver
		}

		if !g.themePlayer.IsPlaying() {
			g.themePlayer.Rewind()
			g.themePlayer.Play()
		}

		if g.board.snake.justAte {
			g.eatPlayer.Rewind()
			g.eatPlayer.Play()
		}
		g.board.Update(g.input)

	case ModeGameOver:
		g.themePlayer.Pause()

		if Space() {
			g.board = NewBoard(boardRows, boardCols)
			g.mode = ModeGame
		}

		if Escape() {
			g.mode = ModeTitle
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case ModeGame:
		g.board.Draw(screen)
	case ModeGameOver:
		g.DrawGameOver(screen)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func (g *Game) DrawGameOver(screen *ebiten.Image) {
	sx := float64(constants.ScreenWidth / 2)
	sy := float64(constants.ScreenHeight / 2)

	// Create the options and set the scale
	op := &ebiten.DrawImageOptions{}
	scaleFactor := 0.8
	op.GeoM.Scale(scaleFactor, scaleFactor)

	// Adjust the translation to center the scaled image
	imageWidth := images.GameOverSprite.Bounds().Dx()
	imageHeight := images.GameOverSprite.Bounds().Dy()
	scaledWidth := float64(imageWidth) * scaleFactor
	scaledHeight := float64(imageHeight) * scaleFactor
	op.GeoM.Translate(sx-scaledWidth/2, sy-scaledHeight/2)

	// Draw the scaled image
	screen.DrawImage(images.GameOverSprite, op)

	// Set the positions for the text
	gameOverText := "Game Over"
	scoreText := "Score: " + strconv.Itoa(g.board.score)
	pressSpaceText := "Press space to play again"
	pressEscapeText := "Press escape to return to the title screen"

	gameOverX := (constants.ScreenWidth - font.MeasureString(fonts.BigFont, gameOverText).Round()) / 2
	gameOverY := (constants.ScreenHeight / 2) - 50

	scoreX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, scoreText).Round()) / 2
	scoreY := gameOverY + 50

	pressStartX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, pressSpaceText).Round()) / 2
	pressStartY := scoreY + 30

	pressEscapeX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, pressEscapeText).Round()) / 2
	pressEscapeY := pressStartY + 30

	// Draw the text
	text.Draw(screen, gameOverText, fonts.BigFont, gameOverX, gameOverY, color.White)
	text.Draw(screen, scoreText, fonts.RegularFont, scoreX, scoreY, color.White)
	text.Draw(screen, pressSpaceText, fonts.RegularFont, pressStartX, pressStartY, color.White)
	text.Draw(screen, pressEscapeText, fonts.RegularFont, pressEscapeX, pressEscapeY, color.White)
}
