package game

import (
	"image"
	"time"

	"github.com/adan-ea/GoSnakeGo/constants"
	"github.com/adan-ea/GoSnakeGo/resources/images"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	spriteWidth  = 240
	spriteHeight = 240
	frameWidth   = 32
	frameHeight  = 32
	headWidth    = 40
	headHeight   = 40
	frameCount   = 6
	frameDelay   = 100 // Delay between frames
)

// Snake represents the snake
type Snake struct {
	Body          []Point
	Direction     Direction
	justAte       bool
	currentFrame  int
	lastFrameTime time.Time
}

// NewSnake creates a new snake with the given body and direction
func NewSnake() *Snake {
	return &Snake{
		Body: []Point{
			{x: 3, y: 1},
			{x: 2, y: 1},
			{x: 1, y: 1}},
		Direction:     Right,
		currentFrame:  0,
		lastFrameTime: time.Now(),
	}
}

// Returns the position of the snake's head
func (s *Snake) Head() Point {
	return s.Body[0]
}

// HeadHits checks if the snake's head is at the given position
func (s *Snake) HeadHits(x, y int) bool {
	h := s.Head()

	return h.x == x && h.y == y
}

// MoveSnake moves the snake one step in its current direction
func (s *Snake) MoveSnake() error {
	// Create a copy of the head position to avoid modifying the original directly
	newHead := s.Body[0]

	// Calculate the new position of the head based on the direction
	switch s.Direction {
	case Up:
		newHead.y--
	case Down:
		newHead.y++
	case Left:
		newHead.x--
	case Right:
		newHead.x++
	}
	// Move each segment of the snake's body to the position of the segment in front of it
	for i := len(s.Body) - 1; i > 0; i-- {
		s.Body[i] = s.Body[i-1]
	}

	// Update the position of the head
	s.Body[0] = newHead

	return nil
}

func (s *Snake) ChangeDirection(newDir Direction) {
	opposites := map[Direction]Direction{
		Up:    Down,
		Right: Left,
		Down:  Up,
		Left:  Right,
	}

	// don't allow changing direction to opposite
	if o, ok := opposites[newDir]; ok && o != s.Direction {
		s.Direction = newDir
	}
}

// Draw draws the snake onto the screen
func (s *Snake) Draw(screen *ebiten.Image) {
	// Calculate the offset to center the game area
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2

	// Draw the snake's body and tail first
	for i := len(s.Body) - 1; i >= 0; i-- {
		part := s.Body[i]
		sx := float64(offsetX + part.x*constants.TileSize)
		sy := float64(offsetY + part.y*constants.TileSize)

		switch {
		case i == 0:
			s.handleHead(screen, sx, sy)
		case i == len(s.Body)-1:
			s.handleTail(screen, sx, sy, i)
		default:
			s.handleBody(screen, sx, sy, i)
		}
	}
}

// handleHead draws the snake's head
func (s *Snake) handleHead(screen *ebiten.Image, sx, sy float64) {
	frameX := s.currentFrame * headWidth
	var frameY int

	switch s.Direction {
	case Up:
		frameY = 1 * headHeight
	case Down:
		frameY = 2 * headHeight
	case Left:
		frameY = 3 * headHeight
	case Right:
		frameY = 0 * headHeight
	}

	headImage := images.HeadSprite.SubImage(image.Rect(frameX, frameY, frameX+headWidth, frameY+headHeight)).(*ebiten.Image)
	headOp := &ebiten.DrawImageOptions{}
	headOp.GeoM.Translate(sx, sy)
	screen.DrawImage(headImage, headOp)
}

// handleBody draws the snake's body
func (s *Snake) handleBody(screen *ebiten.Image, sx, sy float64, i int) {
	prev := s.Body[i-1]
	curr := s.Body[i]
	next := s.Body[i+1]

	var bodyImage *ebiten.Image
	switch {
	// Vertical
	case prev.x == next.x:
		bodyImage = images.BodySprite.SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	// Horizontal
	case prev.y == next.y:
		bodyImage = images.BodySprite.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)

	// Top left corner
	case (prev.x > curr.x && next.y < curr.y) || (next.x > curr.x && prev.y < curr.y):
		bodyImage = images.BodySprite.SubImage(image.Rect(4*frameWidth, 0, 5*frameWidth, frameHeight)).(*ebiten.Image)

	// Top right corner
	case (prev.x < curr.x && next.y < curr.y) || (next.x < curr.x && prev.y < curr.y):
		bodyImage = images.BodySprite.SubImage(image.Rect(5*frameWidth, 0, 6*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom right corner
	case (prev.x < curr.x && next.y > curr.y) || (next.x < curr.x && prev.y > curr.y):
		bodyImage = images.BodySprite.SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom left corner
	case (prev.x > curr.x && next.y > curr.y) || (next.x > curr.x && prev.y > curr.y):
		bodyImage = images.BodySprite.SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	}

	bodyOp := &ebiten.DrawImageOptions{}
	bodyOp.GeoM.Translate(sx, sy)
	screen.DrawImage(bodyImage, bodyOp)
}

// handleTail draws the snake's tail
func (s *Snake) handleTail(screen *ebiten.Image, sx, sy float64, i int) {
	tail := s.Body[i]
	prev := s.Body[i-1]

	var tailImage *ebiten.Image
	switch {
	case tail.x > prev.x: // Going left
		tailImage = images.TailSprite.SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.x < prev.x: // Going right
		tailImage = images.TailSprite.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)
	case tail.y > prev.y: // Going up
		tailImage = images.TailSprite.SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.y < prev.y: // Going down
		tailImage = images.TailSprite.SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	}

	tailOp := &ebiten.DrawImageOptions{}
	tailOp.GeoM.Translate(sx, sy)
	screen.DrawImage(tailImage, tailOp)
}
