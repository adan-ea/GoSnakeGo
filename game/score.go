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

func (b *Board) DrawScoreWithSprite(screen *ebiten.Image, sprite *ebiten.Image, score int, x, y int) {
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
		numImage := images.NumbersSprite.SubImage(image.Rect(sx, sy, sx+digitWidth, sy+images.DigitHeight)).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(currentX), float64(y))
		screen.DrawImage(numImage, op)

		currentX += digitWidth
	}
}

func (b *Board) drawScore(screen *ebiten.Image, score int, x, y int) {
	b.DrawScoreWithSprite(screen, images.FoodSprite, score, x, y)
}

func (b *Board) drawHighScore(screen *ebiten.Image, score int, x, y int) {
	b.DrawScoreWithSprite(screen, images.StarSprite, score, x, y)
}

// saveHighScore saves the score along with the current date, size text, and time to the scoreboard file
func saveHighScore(score int, size Size) {
	if score == 0 {
		return
	}

	// Read the file contents
	f, err := os.ReadFile(bestScorePath)
	if err != nil {
		panic(err)
	}

	fileContent := string(f)

	// Add the new score with date, size, and time to the file content
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fileContent += fmt.Sprintf("%s;%s;%d\n", currentTime, getSizeText(size), score)
	scoresStr := strings.Split(strings.TrimSpace(fileContent), "\n")

	// Parse the scores, date-time, and size text information
	var scores []struct {
		time  string
		size  string
		score int
	}
	for _, scoreStr := range scoresStr {
		if scoreStr != "" {
			parts := strings.Split(scoreStr, ";")
			if len(parts) != 3 {
				continue // Skip invalid entries
			}
			score, err := strconv.Atoi(parts[2])
			if err != nil {
				panic(err)
			}
			scores = append(scores, struct {
				time  string
				size  string
				score int
			}{parts[0], parts[1], score})
		}
	}

	// Sort the scores by size and then by score in descending order
	sort.SliceStable(scores, func(i, j int) bool {
		if scores[i].size != scores[j].size {
			return scores[i].size < scores[j].size
		}
		return scores[i].score > scores[j].score
	})

	// Keep only the top 5 scores for each size
	seen := make(map[string]int)
	var topScores []string
	for _, s := range scores {
		key := s.size
		if seen[key] < nbScoreSaved {
			topScores = append(topScores, fmt.Sprintf("%s;%s;%d", s.time, s.size, s.score))
			seen[key]++
		}
	}

	// Join the top scores into a single string with newlines
	fileContent = strings.Join(topScores, "\n") + "\n"

	// Write the top scores back to the file
	err = os.WriteFile(bestScorePath, []byte(fileContent), 0644)
	if err != nil {
		panic(err)
	}
}

// getHighestScore returns the highest score from the scoreboard file for the specified size
func getHighestScore(size Size) int {
	// Read the file contents
	f, err := os.ReadFile(bestScorePath)
	if err != nil {
		return 0
	}

	fileContent := string(f)
	scoresStr := strings.Split(strings.TrimSpace(fileContent), "\n")

	// Parse the scores, size, and date-time information
	var highestScore int
	for _, scoreStr := range scoresStr {
		if scoreStr != "" {
			parts := strings.Split(scoreStr, ";")
			if len(parts) != 3 {
				continue // Skip invalid entries
			}
			score, err := strconv.Atoi(parts[2])
			if err != nil {
				return 0
			}
			scoreSize := getTextToSize(parts[1])
			if scoreSize == size && score > highestScore {
				highestScore = score
			}
		}
	}

	return highestScore
}

// updateScore updates the score and high score based on the current game state
func (b *Board) updateScore() {
	b.score++
	if b.score > b.highScore {
		b.highScore = b.score
	}
}
