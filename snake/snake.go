package snake

import (
	"image"
	"time"

	"github.com/adan-ea/GoSnakeGo/constants"
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
	Body          []constants.Point
	Direction     constants.Direction
	currentFrame  int
	lastFrameTime time.Time
	headSprite    *ebiten.Image
	bodySprite    *ebiten.Image
	tailSprite    *ebiten.Image
}

// InitSnake creates a new instance of Snake
func InitSnake() *Snake {
	snake := &Snake{
		Body: []constants.Point{
			{X: 3, Y: 1},
			{X: 2, Y: 1},
			{X: 1, Y: 1}},
		Direction:     constants.Right,
		currentFrame:  0,
		lastFrameTime: time.Now(),
	}
	snake.headSprite = constants.LoadImage(constants.HeadSpriteLeftPath)
	snake.bodySprite = constants.LoadImage(constants.BodySpritePath)
	snake.tailSprite = constants.LoadImage(constants.TailSpritePath)

	return snake
}

func (s *Snake) Update() {
	// Update the animation frame
	if time.Since(s.lastFrameTime) >= frameDelay {
		s.currentFrame = (s.currentFrame + 1) % frameCount
		s.lastFrameTime = time.Now()
	}

	s.MoveSnake()
}

// MoveSnake moves the snake one step in its current direction
func (s *Snake) MoveSnake() {
	// Create a copy of the head position to avoid modifying the original directly
	newHead := s.Body[0]

	// Calculate the new position of the head based on the direction
	switch s.Direction {
	case constants.Up:
		newHead.Y--
	case constants.Down:
		newHead.Y++
	case constants.Left:
		newHead.X--
	case constants.Right:
		newHead.X++
	}
	// Move each segment of the snake's body to the position of the segment in front of it
	for i := len(s.Body) - 1; i > 0; i-- {
		s.Body[i] = s.Body[i-1]
	}

	// Update the position of the head
	s.Body[0] = newHead
}

// ChangeDirection changes the direction of the snake
func (s *Snake) ChangeDirection(newDir constants.Direction) {
	// Prevent the snake from turning back on itself
	if (s.Direction == constants.Up && newDir == constants.Down) ||
		(s.Direction == constants.Down && newDir == constants.Up) ||
		(s.Direction == constants.Left && newDir == constants.Right) ||
		(s.Direction == constants.Right && newDir == constants.Left) {
		return
	}
	s.Direction = newDir
}

// Draw draws the snake onto the screen
func (s *Snake) Draw(screen *ebiten.Image) {
	// Calculate the offset to center the game area
	offsetX := (constants.ScreenWidth - constants.GameWidth) / 2
	offsetY := (constants.ScreenHeight - constants.GameHeight) / 2

	// Draw the snake's body and tail first
	for i := len(s.Body) - 1; i >= 0; i-- {
		part := s.Body[i]
		sx := float64(offsetX + part.X*constants.TileSize)
		sy := float64(offsetY + part.Y*constants.TileSize)

		switch {
		case i == 0:
			s.HandleHead(screen, sx, sy)
		case i == len(s.Body)-1:
			s.HandleTail(screen, sx, sy, i)
		default:
			s.HandleBody(screen, sx, sy, i)
		}
	}
}

// HandleHead draws the snake's head
func (s *Snake) HandleHead(screen *ebiten.Image, sx, sy float64) {
	frameX := s.currentFrame * headWidth
	var frameY int

	switch s.Direction {
	case constants.Up:
		frameY = 1 * headHeight
	case constants.Down:
		frameY = 2 * headHeight
	case constants.Left:
		frameY = 3 * headHeight
	case constants.Right:
		frameY = 0 * headHeight
	}

	headImage := s.headSprite.SubImage(image.Rect(frameX, frameY, frameX+headWidth, frameY+headHeight)).(*ebiten.Image)
	headOp := &ebiten.DrawImageOptions{}
	headOp.GeoM.Translate(sx, sy)
	screen.DrawImage(headImage, headOp)
}

// HandleBody draws the snake's body
func (s *Snake) HandleBody(screen *ebiten.Image, sx, sy float64, i int) {
	prev := s.Body[i-1]
	curr := s.Body[i]
	next := s.Body[i+1]

	var bodyImage *ebiten.Image
	switch {
	// Vertical
	case prev.X == next.X:
		bodyImage = s.bodySprite.SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	// Horizontal
	case prev.Y == next.Y:
		bodyImage = s.bodySprite.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)

	// Top left corner
	case (prev.X > curr.X && next.Y < curr.Y) || (next.X > curr.X && prev.Y < curr.Y):
		bodyImage = s.bodySprite.SubImage(image.Rect(4*frameWidth, 0, 5*frameWidth, frameHeight)).(*ebiten.Image)

	// Top right corner
	case (prev.X < curr.X && next.Y < curr.Y) || (next.X < curr.X && prev.Y < curr.Y):
		bodyImage = s.bodySprite.SubImage(image.Rect(5*frameWidth, 0, 6*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom right corner
	case (prev.X < curr.X && next.Y > curr.Y) || (next.X < curr.X && prev.Y > curr.Y):
		bodyImage = s.bodySprite.SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)

	// Bottom left corner
	case (prev.X > curr.X && next.Y > curr.Y) || (next.X > curr.X && prev.Y > curr.Y):
		bodyImage = s.bodySprite.SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	}

	bodyOp := &ebiten.DrawImageOptions{}
	bodyOp.GeoM.Translate(sx, sy)
	screen.DrawImage(bodyImage, bodyOp)
}

// HandleTail draws the snake's tail
func (s *Snake) HandleTail(screen *ebiten.Image, sx, sy float64, i int) {
	tail := s.Body[i]
	prev := s.Body[i-1]

	var tailImage *ebiten.Image
	switch {
	case tail.X > prev.X: // Going left
		tailImage = s.tailSprite.SubImage(image.Rect(3*frameWidth, 0, 4*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.X < prev.X: // Going right
		tailImage = s.tailSprite.SubImage(image.Rect(0, 0, frameWidth, frameHeight)).(*ebiten.Image)
	case tail.Y > prev.Y: // Going up
		tailImage = s.tailSprite.SubImage(image.Rect(2*frameWidth, 0, 3*frameWidth, frameHeight)).(*ebiten.Image)
	case tail.Y < prev.Y: // Going down
		tailImage = s.tailSprite.SubImage(image.Rect(frameWidth, 0, 2*frameWidth, frameHeight)).(*ebiten.Image)
	}

	tailOp := &ebiten.DrawImageOptions{}
	tailOp.GeoM.Translate(sx, sy)
	screen.DrawImage(tailImage, tailOp)
}

// Returns the position of the snake's head
func (s *Snake) GetHead() constants.Point {
	return s.Body[0]
}
