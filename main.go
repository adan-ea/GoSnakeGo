package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	//"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"gosnakego/resources/images"
	"gosnakego/utils"
)

type Game struct {
	snake         []utils.Point
	dir           utils.Point
	food          utils.Point
	width         int
	height        int
	updateTick    int
	score         int
	bestScore     int
	gameOver      bool
	count         int
	scoreregister bool
}

var (
	frameOX     = 0
	frameOY     = 0
	frameWidth  = 32
	frameHeight = 29
	frameCount  = 6
)

var (
	img          *ebiten.Image
	runnerImage  *ebiten.Image
	file         *os.File
	file_content string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	imgb, _, err := image.Decode(bytes.NewReader(images.AdrienSexyy_png))
	if err != nil {
		log.Fatal(err)
	}
	img = ebiten.NewImageFromImage(imgb)

	f, err := ioutil.ReadFile("resources/scoreboard.txt")
	file_content = string(f)

}

func HandleKeyPressed(g *Game) bool {
	// Control snake
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
		if g.scoreregister == false {
			if g.score != 0 {
				file_content += strconv.Itoa(g.score) + "\n"

				scoresStr := strings.Split(file_content, "\n")

				scores := make([]int, len(scoresStr))
				for i, scoreStr := range scoresStr {
					fmt.Println(scoreStr)
					if scoreStr != "" {
						score, err := strconv.Atoi(scoreStr)
						if err != nil {
							panic(err)
						}
						scores[i] = score
					}
				}

				// Trier les scores en ordre décroissant
				sort.Slice(scores, func(i, j int) bool {
					return scores[i] > scores[j]
				})

				// Convertir les scores triés en chaînes
				for i, score := range scores {
					if score != 0 {
						scoresStr[i] = strconv.Itoa(score)
					}
				}

				// Joindre les scores avec des sauts de ligne
				file_content = strings.Join(scoresStr, "\n")
				file_content += "\n"

				ioutil.WriteFile("resources/scoreboard.txt", []byte(file_content), 0644)
			}
			g.scoreregister = true
		}

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

	// Collision with boundaries
	if newHead.X < 0 || newHead.X >= g.width || newHead.Y < 0 || newHead.Y >= g.height {
		g.gameOver = true
		return nil
	}

	// Collision with itself
	for _, v := range g.snake[1:] {
		if v == newHead {
			g.gameOver = true
			return nil
		}
	}

	// Move snake
	g.snake = append([]utils.Point{newHead}, g.snake[:len(g.snake)-1]...)

	// Eating food
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
		screen.DrawImage(img, nil)

		oneFifthHeight := utils.ScreenHeight / 8
		// Dessiner le texte à la position calculée
		text.Draw(screen, "Game Over", basicFont(), utils.ScreenWidth/2-25, oneFifthHeight, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Score: %d", g.score), basicFont(), utils.ScreenWidth/2-20, oneFifthHeight+20, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Best Score: %d", g.bestScore), basicFont(), utils.ScreenWidth/2-30, oneFifthHeight+40, utils.GetBlackColor())
		text.Draw(screen, fmt.Sprintf("Scoreboard: %s", file_content), basicFont(), utils.ScreenWidth/2-30, oneFifthHeight+40, utils.GetBlackColor())
		text.Draw(screen, "Press SPACE to Restart", basicFont(), utils.ScreenWidth/2-80, oneFifthHeight+60, utils.GetBlackColor())

		return
	}

	// Draw food
	foodRect := image.Rect(g.food.X*utils.GridSize, g.food.Y*utils.GridSize, (g.food.X+1)*utils.GridSize, (g.food.Y+1)*utils.GridSize)
	ebitenutil.DrawRect(screen, float64(foodRect.Min.X), float64(foodRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetRedColor())

	// Draw snake
	for _, p := range g.snake[1:] {
		snakeRect := image.Rect(p.X*utils.GridSize, p.Y*utils.GridSize, (p.X+1)*utils.GridSize, (p.Y+1)*utils.GridSize)
		ebitenutil.DrawRect(screen, float64(snakeRect.Min.X), float64(snakeRect.Min.Y), float64(utils.GridSize), float64(utils.GridSize), utils.GetGreenColor())
	}

	// Draw score
	text.Draw(screen, fmt.Sprintf("Score: %d, Best Score: %d", g.score, g.bestScore), basicFont(), 10, 20, utils.GetWhiteColor())

	// Draw head
	op := &ebiten.DrawImageOptions{}

	// Translate to the center of the grid cell occupied by the snake head
	op.GeoM.Translate(float64(g.snake[0].X*utils.GridSize)+utils.GridSize/2, float64(g.snake[0].Y*utils.GridSize)+utils.GridSize/2)

	// Rotate the snake head based on its direction
	switch g.dir {
	case utils.Up:
		op.GeoM.Rotate(2)
	case utils.Down:
		op.GeoM.Rotate(1)
	case utils.Left:
		op.GeoM.Invert()
	}

	// Translate to the center of the snake head image
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)

	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return utils.ScreenWidth, utils.ScreenHeight
}

func (g *Game) initGame() error {
	g.snake = []utils.Point{{X: 5, Y: 5}, {X: 4, Y: 5}, {X: 3, Y: 5}}
	g.dir = utils.Right
	g.width = utils.ScreenWidth / utils.GridSize
	g.height = utils.ScreenHeight / utils.GridSize
	g.spawnFood()
	g.score = 0
	g.gameOver = false
	g.scoreregister = false
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
			g.food = utils.Point{X: x, Y: y}
			break
		}
	}
}

func basicFont() font.Face {
	const dpi = 72
	return basicfont.Face7x13
}

func main() {
	game := &Game{}
	game.initGame()
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)
	ebiten.RunGame(game)
}
