package game

import (
	"bytes"
	"log"

	"github.com/adan-ea/GoSnakeGo/constants"
	raudio "github.com/adan-ea/GoSnakeGo/resources/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
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

/*
func NewGame() *Game {

}
*/

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
		g.board.Update(g.input)

	case ModeGameOver:
		g.themePlayer.Pause()

		if Space() {
			g.board = NewBoard(boardRows, boardCols)
			g.mode = ModeGame
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.mode {
	case ModeGame:
		g.board.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return constants.ScreenWidth, constants.ScreenHeight
}
