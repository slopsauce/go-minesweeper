package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// Game constants
const (
	CELL_SIZE     = 16
	BOARD_WIDTH   = 16
	BOARD_HEIGHT  = 16
	MINE_COUNT    = 40
	WINDOW_WIDTH  = BOARD_WIDTH*CELL_SIZE + 20
	WINDOW_HEIGHT = BOARD_HEIGHT*CELL_SIZE + 80
)

// UI layout constants
const (
	BOARD_OFFSET_X = 13
	BOARD_OFFSET_Y = 58
	SMILEY_SIZE    = 24
)

// Color constants
var (
	ColorGray     = color.RGBA{192, 192, 192, 255}
	ColorWhite    = color.RGBA{255, 255, 255, 255}
	ColorDarkGray = color.RGBA{128, 128, 128, 255}
	ColorBlack    = color.RGBA{0, 0, 0, 255}
	ColorRed      = color.RGBA{255, 0, 0, 255}
	ColorYellow   = color.RGBA{255, 255, 0, 255}
)

// Number colors for mine counts
var NumberColors = []color.RGBA{
	{0, 0, 255, 255},     // 1 - blue
	{0, 128, 0, 255},     // 2 - green
	{255, 0, 0, 255},     // 3 - red
	{0, 0, 128, 255},     // 4 - dark blue
	{128, 0, 0, 255},     // 5 - maroon
	{0, 128, 128, 255},   // 6 - teal
	{0, 0, 0, 255},       // 7 - black
	{128, 128, 128, 255}, // 8 - gray
}

type CellState int

const (
	HIDDEN CellState = iota
	REVEALED
	FLAGGED
)

type Cell struct {
	IsMine        bool
	State         CellState
	AdjacentMines int
}

type GameState int

const (
	PLAYING GameState = iota
	WON
	LOST
)

type Game struct {
	board      [][]Cell
	gameState  GameState
	startTime  time.Time
	gameTime   int
	minesLeft  int
	firstClick bool
	rng        *rand.Rand
}

func NewGame() *Game {
	game := &Game{
		board:      make([][]Cell, BOARD_HEIGHT),
		gameState:  PLAYING,
		minesLeft:  MINE_COUNT,
		firstClick: true,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for i := range game.board {
		game.board[i] = make([]Cell, BOARD_WIDTH)
	}

	return game
}

// Fisher-Yates shuffle algorithm for efficient mine placement
func (g *Game) placeMines(avoidX, avoidY int) {
	// Create list of all valid positions
	positions := make([]image.Point, 0, BOARD_WIDTH*BOARD_HEIGHT-1)
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			if x != avoidX || y != avoidY {
				positions = append(positions, image.Point{x, y})
			}
		}
	}

	// Fisher-Yates shuffle and pick first MINE_COUNT positions
	for i := 0; i < MINE_COUNT && i < len(positions); i++ {
		j := g.rng.Intn(len(positions)-i) + i
		positions[i], positions[j] = positions[j], positions[i]

		pos := positions[i]
		g.board[pos.Y][pos.X].IsMine = true
	}

	// Calculate adjacent mine counts
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			if !g.board[y][x].IsMine {
				g.board[y][x].AdjacentMines = g.countAdjacentMines(x, y)
			}
		}
	}
}

func (g *Game) countAdjacentMines(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if g.isValidPosition(nx, ny) && g.board[ny][nx].IsMine {
				count++
			}
		}
	}
	return count
}

func (g *Game) isValidPosition(x, y int) bool {
	return x >= 0 && x < BOARD_WIDTH && y >= 0 && y < BOARD_HEIGHT
}

func (g *Game) revealCell(x, y int) {
	if !g.isValidPosition(x, y) {
		return
	}

	cell := &g.board[y][x]
	if cell.State != HIDDEN {
		return
	}

	cell.State = REVEALED

	if cell.IsMine {
		g.gameState = LOST
		// Reveal all mines when game is lost
		for my := 0; my < BOARD_HEIGHT; my++ {
			for mx := 0; mx < BOARD_WIDTH; mx++ {
				if g.board[my][mx].IsMine {
					g.board[my][mx].State = REVEALED
				}
			}
		}
		return
	}

	// Flood fill for empty cells
	if cell.AdjacentMines == 0 {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				g.revealCell(x+dx, y+dy)
			}
		}
	}

	g.checkWin()
}

func (g *Game) toggleFlag(x, y int) {
	if !g.isValidPosition(x, y) {
		return
	}

	cell := &g.board[y][x]
	if cell.State == REVEALED {
		return
	}

	if cell.State == HIDDEN {
		cell.State = FLAGGED
		g.minesLeft--
	} else {
		cell.State = HIDDEN
		g.minesLeft++
	}
}

func (g *Game) checkWin() {
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			cell := &g.board[y][x]
			if !cell.IsMine && cell.State != REVEALED {
				return
			}
		}
	}
	g.gameState = WON
}

func (g *Game) resetGame() {
	// Reset board
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			g.board[y][x] = Cell{}
		}
	}

	g.gameState = PLAYING
	g.minesLeft = MINE_COUNT
	g.firstClick = true
	g.gameTime = 0
}

func (g *Game) screenToBoard(screenX, screenY int) (int, int) {
	boardX := (screenX - BOARD_OFFSET_X) / CELL_SIZE
	boardY := (screenY - BOARD_OFFSET_Y) / CELL_SIZE
	return boardX, boardY
}

func (g *Game) isClickOnSmiley(x, y int) bool {
	smileyX := WINDOW_WIDTH/2 - SMILEY_SIZE/2
	smileyY := 16
	return x >= smileyX && x <= smileyX+SMILEY_SIZE &&
		y >= smileyY && y <= smileyY+SMILEY_SIZE
}

func (g *Game) Update() error {
	// Update timer
	if g.gameState == PLAYING && !g.firstClick {
		g.gameTime = int(time.Since(g.startTime).Seconds())
		if g.gameTime > 999 {
			g.gameTime = 999 // Cap at 999 like original
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Check smiley face click
		if g.isClickOnSmiley(x, y) {
			g.resetGame()
			return nil
		}

		// Check board click
		boardX, boardY := g.screenToBoard(x, y)
		if g.isValidPosition(boardX, boardY) {
			if g.firstClick {
				g.placeMines(boardX, boardY)
				g.startTime = time.Now()
				g.firstClick = false
			}
			g.revealCell(boardX, boardY)
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		boardX, boardY := g.screenToBoard(x, y)
		if g.isValidPosition(boardX, boardY) {
			g.toggleFlag(boardX, boardY)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ColorGray)
	g.drawHeader(screen)
	g.drawBoard(screen)
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	// Main header panel
	g.draw3DPanel(screen, 10, 10, float32(WINDOW_WIDTH-20), 40, true)

	// Mine counter panel
	g.draw3DPanel(screen, 16, 16, 40, 28, false)
	g.drawLEDDisplay(screen, 18, 18, fmt.Sprintf("%03d", g.minesLeft))

	// Timer panel
	timerX := float32(WINDOW_WIDTH - 56)
	g.draw3DPanel(screen, timerX, 16, 40, 28, false)
	g.drawLEDDisplay(screen, timerX+2, 18, fmt.Sprintf("%03d", g.gameTime))

	// Smiley face button
	smileyX := float32(WINDOW_WIDTH/2 - SMILEY_SIZE/2)
	smileyY := float32(16)
	g.draw3DPanel(screen, smileyX, smileyY, SMILEY_SIZE, SMILEY_SIZE, true)
	g.drawSmileyFace(screen, smileyX+2, smileyY+2)

	// Board panel
	g.draw3DPanel(screen, 10, 55, float32(BOARD_WIDTH*CELL_SIZE+6), float32(BOARD_HEIGHT*CELL_SIZE+6), false)
}

func (g *Game) draw3DPanel(screen *ebiten.Image, x, y, width, height float32, raised bool) {
	// Background
	vector.DrawFilledRect(screen, x, y, width, height, ColorGray, false)

	if raised {
		// Raised effect (buttons)
		vector.DrawFilledRect(screen, x, y, width-1, height-1, ColorWhite, false)
		vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, ColorGray, false)
		vector.DrawFilledRect(screen, x+width-1, y, 1, height, ColorDarkGray, false)
		vector.DrawFilledRect(screen, x, y+height-1, width, 1, ColorDarkGray, false)
	} else {
		// Sunken effect (displays)
		vector.DrawFilledRect(screen, x, y, width, height, ColorDarkGray, false)
		vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, ColorGray, false)
		vector.DrawFilledRect(screen, x+width-1, y+1, 1, height-2, ColorWhite, false)
		vector.DrawFilledRect(screen, x+1, y+height-1, width-2, 1, ColorWhite, false)
	}
}

func (g *Game) drawLEDDisplay(screen *ebiten.Image, x, y float32, displayText string) {
	// Black background for LED display
	vector.DrawFilledRect(screen, x, y, 36, 24, ColorBlack, false)

	// LED-style text in red
	bounds := text.BoundString(basicfont.Face7x13, displayText)
	textX := int(x) + (36-bounds.Dx())/2
	textY := int(y) + (24+bounds.Dy())/2 - 2

	text.Draw(screen, displayText, basicfont.Face7x13, textX, textY, ColorRed)
}

func (g *Game) drawSmileyFace(screen *ebiten.Image, x, y float32) {
	// Yellow circular face with black outline
	vector.DrawFilledCircle(screen, x+10, y+10, 9, ColorYellow, false)
	vector.StrokeCircle(screen, x+10, y+10, 9, 1, ColorBlack, false)

	if g.gameState == LOST {
		// Dead face - X eyes and O mouth
		g.drawXEyes(screen, x, y)
		vector.StrokeCircle(screen, x+10, y+13, 2, 1, ColorBlack, false)
	} else if g.gameState == WON {
		// Happy face with sunglasses
		vector.DrawFilledRect(screen, x+5, y+6, 10, 3, ColorBlack, false)
		vector.DrawFilledRect(screen, x+6, y+7, 8, 1, ColorYellow, false)
		// Big smile
		vector.DrawFilledRect(screen, x+7, y+12, 6, 1, ColorBlack, false)
		vector.DrawFilledRect(screen, x+8, y+13, 4, 1, ColorBlack, false)
		vector.DrawFilledRect(screen, x+9, y+14, 2, 1, ColorBlack, false)
	} else {
		// Normal eyes and smile
		vector.DrawFilledRect(screen, x+7, y+7, 1, 1, ColorBlack, false)
		vector.DrawFilledRect(screen, x+13, y+7, 1, 1, ColorBlack, false)
		vector.DrawFilledRect(screen, x+8, y+12, 4, 1, ColorBlack, false)
		vector.DrawFilledRect(screen, x+9, y+13, 2, 1, ColorBlack, false)
	}
}

func (g *Game) drawXEyes(screen *ebiten.Image, x, y float32) {
	// Left X eye
	vector.DrawFilledRect(screen, x+6, y+6, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+7, y+7, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+8, y+8, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+8, y+6, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+7, y+7, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+6, y+8, 1, 1, ColorBlack, false)

	// Right X eye
	vector.DrawFilledRect(screen, x+12, y+6, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+13, y+7, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+14, y+8, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+14, y+6, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+13, y+7, 1, 1, ColorBlack, false)
	vector.DrawFilledRect(screen, x+12, y+8, 1, 1, ColorBlack, false)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			g.drawCell(screen, x, y)
		}
	}
}

func (g *Game) drawCell(screen *ebiten.Image, x, y int) {
	cell := &g.board[y][x]
	screenX := float32(x*CELL_SIZE + BOARD_OFFSET_X)
	screenY := float32(y*CELL_SIZE + BOARD_OFFSET_Y)

	if cell.State == HIDDEN {
		g.drawRaisedCell(screen, screenX, screenY)
	} else if cell.State == REVEALED {
		g.drawRevealedCell(screen, screenX, screenY, cell)
	} else if cell.State == FLAGGED {
		g.drawRaisedCell(screen, screenX, screenY)
		g.drawFlag(screen, screenX, screenY)
	}
}

func (g *Game) drawRaisedCell(screen *ebiten.Image, x, y float32) {
	vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, ColorGray, false)
	vector.DrawFilledRect(screen, x, y, CELL_SIZE-1, CELL_SIZE-1, ColorWhite, false)
	vector.DrawFilledRect(screen, x+1, y+1, CELL_SIZE-2, CELL_SIZE-2, ColorGray, false)
	vector.DrawFilledRect(screen, x+CELL_SIZE-1, y, 1, CELL_SIZE, ColorDarkGray, false)
	vector.DrawFilledRect(screen, x, y+CELL_SIZE-1, CELL_SIZE, 1, ColorDarkGray, false)
}

func (g *Game) drawRevealedCell(screen *ebiten.Image, x, y float32, cell *Cell) {
	vector.DrawFilledRect(screen, x, y, CELL_SIZE, CELL_SIZE, ColorDarkGray, false)
	vector.DrawFilledRect(screen, x+1, y+1, CELL_SIZE-2, CELL_SIZE-2, ColorGray, false)

	if cell.IsMine {
		g.drawMine(screen, x, y)
	} else if cell.AdjacentMines > 0 {
		g.drawNumber(screen, x, y, cell.AdjacentMines)
	}
}

func (g *Game) drawMine(screen *ebiten.Image, x, y float32) {
	centerX := x + 8
	centerY := y + 8

	// Main mine body
	vector.DrawFilledCircle(screen, centerX, centerY, 4, ColorBlack, false)

	// 8 spikes radiating from center
	vector.DrawFilledRect(screen, centerX-0.5, y+1, 1, 14, ColorBlack, false) // Vertical
	vector.DrawFilledRect(screen, x+1, centerY-0.5, 14, 1, ColorBlack, false) // Horizontal
	vector.DrawFilledRect(screen, x+2, y+2, 12, 1, ColorBlack, false)         // Diagonal
	vector.DrawFilledRect(screen, x+2, y+13, 12, 1, ColorBlack, false)        // Diagonal
	vector.DrawFilledRect(screen, x+2, y+2, 1, 12, ColorBlack, false)         // Diagonal
	vector.DrawFilledRect(screen, x+13, y+2, 1, 12, ColorBlack, false)        // Diagonal

	// White highlight
	vector.DrawFilledRect(screen, centerX-1, centerY-1, 2, 2, ColorWhite, false)
}

func (g *Game) drawFlag(screen *ebiten.Image, x, y float32) {
	// Flag pole
	vector.DrawFilledRect(screen, x+8, y+2, 1, 12, ColorBlack, false)

	// Flag triangle
	vector.DrawFilledRect(screen, x+9, y+2, 5, 3, ColorRed, false)
	vector.DrawFilledRect(screen, x+9, y+5, 4, 1, ColorRed, false)
	vector.DrawFilledRect(screen, x+9, y+6, 3, 1, ColorRed, false)
	vector.DrawFilledRect(screen, x+9, y+7, 2, 1, ColorRed, false)
	vector.DrawFilledRect(screen, x+9, y+8, 1, 1, ColorRed, false)
}

func (g *Game) drawNumber(screen *ebiten.Image, x, y float32, number int) {
	if number < 1 || number > 8 {
		return
	}

	numberColor := NumberColors[number-1]
	numberText := fmt.Sprintf("%d", number)
	bounds := text.BoundString(basicfont.Face7x13, numberText)
	textX := int(x) + (CELL_SIZE-bounds.Dx())/2
	textY := int(y) + (CELL_SIZE+bounds.Dy())/2 - 2

	text.Draw(screen, numberText, basicfont.Face7x13, textX, textY, numberColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WINDOW_WIDTH, WINDOW_HEIGHT
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
	ebiten.SetWindowTitle("Minesweeper")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
