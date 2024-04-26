package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"gosnakego/resources/images"
	"gosnakego/utils"
)

type Game struct {
	snake        []utils.Point
	dir          utils.Point
	food         utils.Point
	width        int
	height       int
	updateTick   int
	score        int
	bestScore    int
	gameOver     bool
	count        int
	audioContext *audio.Context
	bgmPlayer    *audio.Player
}

var (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 29
	frameCount  = 6
)

var (
	img         *ebiten.Image
	runnerImage *ebiten.Image
)

func createAudioContext() (*audio.Context, error) {
	context, err := audio.NewContext(44100) // Create audio context
	if err != nil {
		return nil, fmt.Errorf("Failed to create audio context: %w", err)
	}
	return context, nil
}

func main() {
	game := &Game{}

	// Ensure audio context is created before using game
	var err error
	game.audioContext, err = createAudioContext()
	if err != nil {
		log.Fatal(err)
	}

	game.initGame()

	imgData, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imgData)

	ebiten.RunGame(game)
}

func (g *Game) initGame() error {
	g.snake = []utils.Point{{X: 5, Y: 5}, {X: 4, Y: 5}, {X: 3, Y: 5}}
	g.dir = utils.Right
	g.width = utils.ScreenWidth / utils.GridSize
	g.height = utils.ScreenHeight / utils.GridSize
	g.spawnFood()
	g.score = 0
	g.gameOver = false

	f, err := os.ReadFile("assets/music.mp3")
	if err != nil {
		log.Fatalf("Failed to read audio file: %v", err)
	}
	d, err := mp3.Decode(g.audioContext, bytes.NewReader(f)) // Use assigned audioContext
	if err != nil {
		log.Fatalf("Failed to decode MP3: %v", err)
	}
	g.bgmPlayer, err = audio.NewPlayer(g.audioContext, d)
	if err != nil {
		log.Fatalf("Failed to create audio player: %v", err)
	}
	g.bgmPlayer.SetVolume(0.5)
	g.bgmPlayer.Play()

	return nil
}

func HandleKeyPressed(g *Game) bool {
	if (ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)) && g.dir.Y == 0 {
		g.dir = utils.Up
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)) && g.dir.Y == 0 {
		g.dir = utils.Down
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)) && g.dir.X == 0 {
		g.dir = utils.Left
	} else if (ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)) && g.dir.X == 0 {
		g.dir = utils.Right
	}
	return true
}

func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.initGame()
		}
		return nil
	}

	HandleKeyPressed(g)
	g.updateTick++
	g.count++
	if g.updateTick%4 != 0 {
		return nil
	}

	newHead := utils.Point{
		X: g.snake[0].X + g.dir.X,
		Y: g.snake[0].Y + g.dir.Y,
	}
	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		g.gameOver = true
		return nil
	}
	for _, v := range g.snake[1:] {
		if v == newHead {
			g.gameOver = true
			return nil
		}
	}
	g.snake = append([]utils.Point{newHead}, g.snake[:len(g.snake)-1]...)
	if newHead == g.food {
		g.snake = append(g.snake, g.food)
		g.spawnFood()
		g.score++
		if g.score > g.bestScore {
			g.bestScore = g.score
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(utils.GetDarkGreenColor())

	if g.gameOver {
		if img != nil {
			screen.DrawImage(img, nil)
		} else {
			log.Println("Game over image is nil")
			text.Draw(screen, "Game Over - Image not loaded", basicFont(), utils.ScreenWidth/2-150, utils.ScreenHeight/2, utils.GetWhiteColor())
		}
		oneFifthHeight := utils.ScreenHeight / 5
		text.Draw(screen, "Game Over", basicFont(), utils.ScreenWidth/2-50, oneFifthHeight, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicFont(), utils.ScreenWidth/2-40, oneFifthHeight+20, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Best Score: %d", g.bestScore), basicFont(), utils.ScreenWidth/2-50, oneFifthHeight+40, utils.GetBlackColor())
		text.Draw(screen, "Press SPACE to Restart", basicFont(), utils.ScreenWidth/2-90, oneFifthHeight+60, utils.GetBlackColor())
		return
	}

	// Draw food
	foodRect := image.Rect(g.food.X*utils.GridSize, g.food.Y*utils.GridSize, (g.food.X+1)*utils.GridSize, (g.food.Y+1)*utils.GridSize)
	ebitenutil.DrawRect(screen, float64(foodRect.Min.X), float64(foodRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetRedColor())

	// Draw each segment of the snake
	for _, p := range g.snake {
		snakeRect := image.Rect(p.X*utils.GridSize, p.Y*utils.GridSize, (p.X+1)*utils.GridSize, (p.Y+1)*utils.GridSize)
		ebitenutil.DrawRect(screen, float64(snakeRect.Min.X), float64(snakeRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetGreenColor())
	}

	// Draw score and best score at the top of the screen
	text.Draw(screen, fmt.Sprintf("Score: %d, Best Score: %d", g.score, g.bestScore), basicFont(), 10, 20, utils.GetWhiteColor())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return utils.ScreenWidth, utils.ScreenHeight
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
	return basicfont.Face7x13
}
