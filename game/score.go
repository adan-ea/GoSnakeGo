package game

import (
	"fmt"
	"image"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) DrawScoreWithSprite(screen *ebiten.Image, sprite *ebiten.Image, score int, x, y int) {
	spriteWidth := sprite.Bounds().Dx()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(sprite, op)

	scoreStr := fmt.Sprintf("%d", score)
	currentX := x + spriteWidth + 5

	for _, char := range scoreStr {
		digit := int(char - '0')
		digitWidth := images.DigitWidths[digit]
		sx := 0
		for j := 0; j < digit; j++ {
			sx += images.DigitWidths[j]
		}
		sy := 0
		numImage := g.numberSprite.SubImage(image.Rect(sx, sy, sx+digitWidth, sy+images.DigitHeight)).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(currentX), float64(y))
		screen.DrawImage(numImage, op)

		currentX += digitWidth
	}
}

func (g *Game) drawScore(screen *ebiten.Image, score int, x, y int) {
	g.DrawScoreWithSprite(screen, g.food.Sprite, score, x, y)
}

func (g *Game) drawHighScore(screen *ebiten.Image, score int, x, y int) {
	g.DrawScoreWithSprite(screen, g.starSprite, score, x, y)
}

// saveHighScore saves the score along with the current date and time to the scoreboard file
func saveHighScore(score int) {
	if score == 0 {
		return
	}

	// Read the file contents
	f, err := os.ReadFile(bestScorePath)
	if err != nil {
		panic(err)
	}

	fileContent := string(f)

	// Add the new score with date and time to the file content
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fileContent += fmt.Sprintf("%s;%d\n", currentTime, score)
	scoresStr := strings.Split(strings.TrimSpace(fileContent), "\n")

	// Parse the scores and date-time information
	var scores []struct {
		time  string
		score int
	}
	for _, scoreStr := range scoresStr {
		if scoreStr != "" {
			parts := strings.Split(scoreStr, ";")
			if len(parts) != 2 {
				continue // Skip invalid entries
			}
			score, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			scores = append(scores, struct {
				time  string
				score int
			}{parts[0], score})
		}
	}

	// Sort the scores in descending order
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Keep only the top 5 scores
	var topScores []string
	for i := 0; i < len(scores) && i < nbScoreSaved; i++ {
		topScores = append(topScores, fmt.Sprintf("%s;%d", scores[i].time, scores[i].score))
	}

	// Join the top scores into a single string with newlines
	fileContent = strings.Join(topScores, "\n") + "\n"

	// Write the top scores back to the file
	err = os.WriteFile(bestScorePath, []byte(fileContent), 0644)
	if err != nil {
		panic(err)
	}
}

// getHighestScore returns the highest score from the scoreboard file
func getHighestScore() int {
	// Read the file contents
	f, err := os.ReadFile(bestScorePath)
	if err != nil {
		return 0
	}

	fileContent := string(f)
	scoresStr := strings.Split(strings.TrimSpace(fileContent), "\n")

	// Parse the scores and date-time information
	var highestScore int
	for _, scoreStr := range scoresStr {
		if scoreStr != "" {
			parts := strings.Split(scoreStr, ";")
			if len(parts) != 2 {
				continue // Skip invalid entries
			}
			score, err := strconv.Atoi(parts[1])
			if err != nil {
				return 0
			}
			if score > highestScore {
				highestScore = score
			}
		}
	}

	return highestScore
}

// updateScore updates the score and high score based on the current game state
func (g *Game) updateScore() {
	g.score++
	if g.score > g.highScore {
		g.highScore = g.score
	}
}
