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
	frameDelay   = 50 // Delay between frames
)

// Snake represents the snake
type Snake struct {
	body             []Point
	direction        Direction
	color            Color
	justAte          bool
	changedDirection bool
	currentFrame     int
	lastFrameTime    time.Time
}

// NewSnake creates a new snake with the given body and direction
func NewSnake(color Color) *Snake {
	return &Snake{
		body: []Point{
			{x: 1, y: 1},
			{x: 2, y: 1},
			{x: 3, y: 1}},
		color:         color,
		lastFrameTime: time.Now(),
	}
}

// Returns the position of the snake's head
func (s *Snake) Head() Point {
	return s.body[len(s.body)-1]
}

func (s *Snake) ChangeDirection(newDir Direction) {
	opposites := map[Direction]Direction{
		Up:    Down,
		Right: Left,
		Down:  Up,
		Left:  Right,
	}
	if !s.changedDirection {
		// Prevent the snake from reversing direction
		if o, ok := opposites[newDir]; ok && o != s.direction {
			s.direction = newDir
			s.changedDirection = true
		}
	}
}

// HeadHits checks if the snake's head is at the given position
func (s *Snake) HeadHits(x, y int) bool {
	h := s.Head()

	return h.x == x && h.y == y
}

func (s *Snake) HeadHitsBody() bool {
	h := s.Head()
	bodyWithoutHead := s.body[:len(s.body)-1]

	for _, b := range bodyWithoutHead {
		if b.x == h.x && b.y == h.y {
			return true
		}
	}

	return false
}

// Move moves the snake one step in its current direction
func (s *Snake) Move() {
	s.changedDirection = false
	// Create a copy of the head position to avoid modifying the original directly
	h := s.Head()
	newHead := Point{x: h.x, y: h.y}
	// Calculate the new position of the head based on the direction
	switch s.direction {
	case Up:
		newHead.y--
	case Down:
		newHead.y++
	case Left:
		newHead.x--
	case Right:
		newHead.x++
	}

	if s.justAte {
		s.body = append(s.body, newHead)
		s.justAte = false
	} else {
		s.body = append(s.body[1:], newHead)
	}
}

// Draw draws the snake onto the screen
func (s *Snake) Draw(screen *ebiten.Image) {
	// Calculate the offset to center the game area
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2
	s.UpdateAnimation()
	// Draw the snake's body and tail first
	for i := 0; i < len(s.body); i++ {
		part := s.body[i]
		sx := float64(offsetX + part.x*constants.TileSize)
		sy := float64(offsetY + part.y*constants.TileSize)

		switch {
		case i == 0:
			s.handleTail(screen, sx, sy, i)
		case i == len(s.body)-1:
			s.handleHead(screen, sx, sy)
		default:
			s.handleBody(screen, sx, sy, i)
		}
	}
}

// handleHead draws the snake's head
func (s *Snake) handleHead(screen *ebiten.Image, sx, sy float64) {
	frameX := s.currentFrame * headWidth
	var frameY int

	switch s.direction {
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
	prev := s.body[i-1]
	curr := s.body[i]
	next := s.body[i+1]

	var bodyImage *ebiten.Image
	switch {
	// Vertical
	case prev.x == next.x:
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	// Horizontal
	case prev.y == next.y:
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)

	// Top left corner
	case (prev.x > curr.x && next.y < curr.y) || (next.x > curr.x && prev.y < curr.y):
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(4*frameWidth, 0, 5*frameWidth, frameHeight)).(*ebiten.Image)

	// Top right corner
	case (prev.x < curr.x && next.y < curr.y) || (next.x < curr.x && prev.y < curr.y):
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(5*frameWidth, 0, 6*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom right corner
	case (prev.x < curr.x && next.y > curr.y) || (next.x < curr.x && prev.y > curr.y):
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom left corner
	case (prev.x > curr.x && next.y > curr.y) || (next.x > curr.x && prev.y > curr.y):
		bodyImage = images.BodySprite[int(s.color)].SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	}

	bodyOp := &ebiten.DrawImageOptions{}
	bodyOp.GeoM.Translate(sx, sy)
	screen.DrawImage(bodyImage, bodyOp)
}

// handleTail draws the snake's tail
func (s *Snake) handleTail(screen *ebiten.Image, sx, sy float64, i int) {
	tail := s.body[i]
	prev := s.body[i+1]

	var tailImage *ebiten.Image
	switch {
	case tail.x > prev.x: // Going left
		tailImage = images.TailSprite[int(s.color)].SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.x < prev.x: // Going right
		tailImage = images.TailSprite[int(s.color)].SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)
	case tail.y > prev.y: // Going up
		tailImage = images.TailSprite[int(s.color)].SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.y < prev.y: // Going down
		tailImage = images.TailSprite[int(s.color)].SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	}

	tailOp := &ebiten.DrawImageOptions{}
	tailOp.GeoM.Translate(sx, sy)
	screen.DrawImage(tailImage, tailOp)
}

// UpdateAnimation updates the snake's animation frame
func (s *Snake) UpdateAnimation() {
	if time.Since(s.lastFrameTime) >= time.Millisecond*frameDelay {
		s.currentFrame = (s.currentFrame + 1) % frameCount
		s.lastFrameTime = time.Now()
	}
}
