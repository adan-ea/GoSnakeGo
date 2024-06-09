package game

import (
	"image/color"
	"strconv"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/resources/audio"
	"github.com/adan-ea/GoSnakeGo/resources/fonts"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	bestScorePath = "resources/scoreboard.txt"
	nbScoreSaved  = 5
)

var sizeChanged bool
var colorChanged bool

// Game represents the game state and logic
type Game struct {
	input *Input
	board *Board
	size  Size
	color Color
	mode  Mode
}

func NewGame() *Game {
	game := &Game{
		input: newInput(),
	}

	return game
}

func (g *Game) Update() error {

	switch g.mode {
	case ModeTitle:
		handleSizeOption(g)
		handleColorOption(g)
		if Space() {
			g.board = newBoard(g.size, g.color)
			g.mode = ModeGame
		}
	case ModeGame:
		if g.board.gameOver {
			audio.PlayOnce(audio.GameOverPlayer)
			g.mode = ModeGameOver
		}

		audio.PlayLoop(audio.ThemePlayer)
		g.board.Update(g.input)

	case ModeGameOver:
		audio.ThemePlayer.Pause()

		if Space() {
			g.board = newBoard(g.size, g.color)
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
	case ModeTitle:
		g.DrawMainPage(screen)
	case ModeGame:
		g.board.Draw(screen)
	case ModeGameOver:
		g.DrawGameOver(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}

func handleColorOption(g *Game) {
	if KeyC() && !colorChanged {
		colorChanged = true
		g.color = (g.color + 1) % nbColors
	} else if !KeyC() {
		colorChanged = false
	}
}

func handleSizeOption(g *Game) {
	if KeyS() && !sizeChanged {
		sizeChanged = true
		g.size = (g.size + 1) % nbSize
	} else if !KeyS() {
		sizeChanged = false
	}
}

func (g *Game) DrawMainPage(screen *ebiten.Image) {
	title := "Go Snake Go!"
	sizeText := "Size: " + getSizeText(g.size)
	colorText := "Color: " + getColorText(g.color)
	startText := "Space to start"

	// Set the positions for the text
	titleX := (constants.ScreenWidth - font.MeasureString(fonts.BigFont, title).Round()) / 2
	titleY := (constants.ScreenHeight / 2) - 150

	sizeX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, sizeText).Round()) / 2
	sizeY := titleY + 50

	colorX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, colorText).Round()) / 2
	colorY := sizeY + 50

	startX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, startText).Round()) / 2
	startY := constants.ScreenHeight - 50

	// Draw the text
	text.Draw(screen, title, fonts.BigFont, titleX, titleY, color.White)
	text.Draw(screen, sizeText, fonts.RegularFont, sizeX, sizeY, color.White)
	text.Draw(screen, colorText, fonts.RegularFont, colorX, colorY, color.White)
	text.Draw(screen, startText, fonts.RegularFont, startX, startY, color.White)
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
	gameOverY := (constants.ScreenHeight / 2) - 150

	scoreX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, scoreText).Round()) / 2
	scoreY := gameOverY + 50

	pressStartX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, pressSpaceText).Round()) / 2
	pressStartY := constants.ScreenHeight - 50

	pressEscapeX := (constants.ScreenWidth - font.MeasureString(fonts.RegularFont, pressEscapeText).Round()) / 2
	pressEscapeY := pressStartY + 30

	// Draw the text
	text.Draw(screen, gameOverText, fonts.BigFont, gameOverX, gameOverY, color.White)
	text.Draw(screen, scoreText, fonts.RegularFont, scoreX, scoreY, color.White)
	text.Draw(screen, pressSpaceText, fonts.RegularFont, pressStartX, pressStartY, color.White)
	text.Draw(screen, pressEscapeText, fonts.RegularFont, pressEscapeX, pressEscapeY, color.White)
}
