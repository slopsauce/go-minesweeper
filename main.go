package main

import (
	"fmt"
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

const (
	CELL_SIZE    = 16
	BOARD_WIDTH  = 16
	BOARD_HEIGHT = 16
	MINE_COUNT   = 40
	WINDOW_WIDTH = BOARD_WIDTH*CELL_SIZE + 20
	WINDOW_HEIGHT = BOARD_HEIGHT*CELL_SIZE + 80
)

type CellState int

const (
	HIDDEN CellState = iota
	REVEALED
	FLAGGED
)

type Cell struct {
	IsMine     bool
	State      CellState
	AdjacentMines int
}

type GameState int

const (
	PLAYING GameState = iota
	WON
	LOST
)

type Game struct {
	board     [][]Cell
	gameState GameState
	startTime time.Time
	gameTime  int
	minesLeft int
	firstClick bool
}

func NewGame() *Game {
	game := &Game{
		board:     make([][]Cell, BOARD_HEIGHT),
		gameState: PLAYING,
		minesLeft: MINE_COUNT,
		firstClick: true,
	}
	
	for i := range game.board {
		game.board[i] = make([]Cell, BOARD_WIDTH)
	}
	
	return game
}

func (g *Game) placeMines(avoidX, avoidY int) {
	rand.Seed(time.Now().UnixNano())
	minesPlaced := 0
	
	for minesPlaced < MINE_COUNT {
		x := rand.Intn(BOARD_WIDTH)
		y := rand.Intn(BOARD_HEIGHT)
		
		if !g.board[y][x].IsMine && !(x == avoidX && y == avoidY) {
			g.board[y][x].IsMine = true
			minesPlaced++
		}
	}
	
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
			if nx >= 0 && nx < BOARD_WIDTH && ny >= 0 && ny < BOARD_HEIGHT && g.board[ny][nx].IsMine {
				count++
			}
		}
	}
	return count
}

func (g *Game) revealCell(x, y int) {
	if x < 0 || x >= BOARD_WIDTH || y < 0 || y >= BOARD_HEIGHT {
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
	if x < 0 || x >= BOARD_WIDTH || y < 0 || y >= BOARD_HEIGHT {
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
	g.board = make([][]Cell, BOARD_HEIGHT)
	for i := range g.board {
		g.board[i] = make([]Cell, BOARD_WIDTH)
	}
	g.gameState = PLAYING
	g.minesLeft = MINE_COUNT
	g.firstClick = true
	g.gameTime = 0
}

func (g *Game) Update() error {
	if g.gameState == PLAYING && !g.firstClick {
		g.gameTime = int(time.Since(g.startTime).Seconds())
	}
	
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		
		// Check if clicking on smiley face button
		if x >= WINDOW_WIDTH/2-12 && x <= WINDOW_WIDTH/2+12 && y >= 16 && y <= 40 {
			g.resetGame()
			return nil
		}
		
		boardX := (x - 13) / CELL_SIZE
		boardY := (y - 58) / CELL_SIZE
		
		if boardX >= 0 && boardX < BOARD_WIDTH && boardY >= 0 && boardY < BOARD_HEIGHT {
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
		boardX := (x - 13) / CELL_SIZE
		boardY := (y - 58) / CELL_SIZE
		
		if boardX >= 0 && boardX < BOARD_WIDTH && boardY >= 0 && boardY < BOARD_HEIGHT {
			g.toggleFlag(boardX, boardY)
		}
	}
	
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{192, 192, 192, 255})
	
	g.drawHeader(screen)
	g.drawBoard(screen)
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	headerColor := color.RGBA{192, 192, 192, 255}
	screen.Fill(headerColor)
	
	// Main header panel with 3D effect
	g.draw3DPanel(screen, 10, 10, float32(WINDOW_WIDTH-20), 40, true)
	
	// Mine counter panel
	g.draw3DPanel(screen, 16, 16, 40, 28, false)
	g.drawLEDDisplay(screen, 18, 18, fmt.Sprintf("%03d", g.minesLeft))
	
	// Timer panel
	g.draw3DPanel(screen, float32(WINDOW_WIDTH-56), 16, 40, 28, false)
	g.drawLEDDisplay(screen, float32(WINDOW_WIDTH-54), 18, fmt.Sprintf("%03d", g.gameTime))
	
	// Smiley face button
	smileyX := float32(WINDOW_WIDTH/2 - 12)
	smileyY := float32(16)
	
	g.draw3DPanel(screen, smileyX, smileyY, 24, 24, true)
	g.drawSmileyFace(screen, smileyX+2, smileyY+2)
	
	// Board panel
	g.draw3DPanel(screen, 10, 55, float32(WINDOW_WIDTH-20), float32(BOARD_HEIGHT*CELL_SIZE+6), false)
}

func (g *Game) draw3DPanel(screen *ebiten.Image, x, y, width, height float32, raised bool) {
	// Background
	vector.DrawFilledRect(screen, x, y, width, height, color.RGBA{192, 192, 192, 255}, false)
	
	if raised {
		// Raised effect (buttons)
		vector.DrawFilledRect(screen, x, y, width-1, height-1, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, x+width-1, y, 1, height, color.RGBA{128, 128, 128, 255}, false)
		vector.DrawFilledRect(screen, x, y+height-1, width, 1, color.RGBA{128, 128, 128, 255}, false)
	} else {
		// Sunken effect (displays)
		vector.DrawFilledRect(screen, x, y, width, height, color.RGBA{128, 128, 128, 255}, false)
		vector.DrawFilledRect(screen, x+1, y+1, width-2, height-2, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, x+width-1, y+1, 1, height-2, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledRect(screen, x+1, y+height-1, width-2, 1, color.RGBA{255, 255, 255, 255}, false)
	}
}

func (g *Game) drawLEDDisplay(screen *ebiten.Image, x, y float32, displayText string) {
	// Black background for LED display
	vector.DrawFilledRect(screen, x, y, 36, 24, color.RGBA{0, 0, 0, 255}, false)
	
	// LED-style text in red
	bounds := text.BoundString(basicfont.Face7x13, displayText)
	textX := int(x) + (36-bounds.Dx())/2
	textY := int(y) + (24+bounds.Dy())/2 - 2
	
	text.Draw(screen, displayText, basicfont.Face7x13, textX, textY, color.RGBA{255, 0, 0, 255})
}

func (g *Game) drawSmileyFace(screen *ebiten.Image, x, y float32) {
	// Yellow circular face with black outline
	vector.DrawFilledCircle(screen, x+10, y+10, 9, color.RGBA{255, 255, 0, 255}, false)
	vector.StrokeCircle(screen, x+10, y+10, 9, 1, color.RGBA{0, 0, 0, 255}, false)
	
	if g.gameState == LOST {
		// Dead face - X eyes
		vector.DrawFilledRect(screen, x+6, y+6, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+7, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+8, y+8, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+8, y+6, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+7, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+6, y+8, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		
		vector.DrawFilledRect(screen, x+12, y+6, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+13, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+14, y+8, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+14, y+6, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+13, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+12, y+8, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		
		// O mouth
		vector.StrokeCircle(screen, x+10, y+13, 2, 1, color.RGBA{0, 0, 0, 255}, false)
	} else if g.gameState == WON {
		// Sunglasses
		vector.DrawFilledRect(screen, x+5, y+6, 10, 3, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+6, y+7, 8, 1, color.RGBA{255, 255, 0, 255}, false)
		
		// Big smile
		vector.DrawFilledRect(screen, x+7, y+12, 6, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+8, y+13, 4, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+9, y+14, 2, 1, color.RGBA{0, 0, 0, 255}, false)
	} else {
		// Normal eyes - black dots
		vector.DrawFilledRect(screen, x+7, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+13, y+7, 1, 1, color.RGBA{0, 0, 0, 255}, false)
		
		// Normal smile
		vector.DrawFilledRect(screen, x+8, y+12, 4, 1, color.RGBA{0, 0, 0, 255}, false)
		vector.DrawFilledRect(screen, x+9, y+13, 2, 1, color.RGBA{0, 0, 0, 255}, false)
	}
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
	screenX := float32(x*CELL_SIZE + 13)
	screenY := float32(y*CELL_SIZE + 58)
	
	if cell.State == HIDDEN {
		// Hidden cell - raised 3D effect
		vector.DrawFilledRect(screen, screenX, screenY, CELL_SIZE, CELL_SIZE, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, screenX, screenY, CELL_SIZE-1, CELL_SIZE-1, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledRect(screen, screenX+1, screenY+1, CELL_SIZE-2, CELL_SIZE-2, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, screenX+CELL_SIZE-1, screenY, 1, CELL_SIZE, color.RGBA{128, 128, 128, 255}, false)
		vector.DrawFilledRect(screen, screenX, screenY+CELL_SIZE-1, CELL_SIZE, 1, color.RGBA{128, 128, 128, 255}, false)
	} else if cell.State == REVEALED {
		// Revealed cell - flat with border
		vector.DrawFilledRect(screen, screenX, screenY, CELL_SIZE, CELL_SIZE, color.RGBA{128, 128, 128, 255}, false)
		vector.DrawFilledRect(screen, screenX+1, screenY+1, CELL_SIZE-2, CELL_SIZE-2, color.RGBA{192, 192, 192, 255}, false)
		
		if cell.IsMine {
			g.drawMine(screen, screenX, screenY)
		} else if cell.AdjacentMines > 0 {
			g.drawNumber(screen, screenX, screenY, cell.AdjacentMines)
		}
	} else if cell.State == FLAGGED {
		// Flagged cell - raised 3D effect like hidden
		vector.DrawFilledRect(screen, screenX, screenY, CELL_SIZE, CELL_SIZE, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, screenX, screenY, CELL_SIZE-1, CELL_SIZE-1, color.RGBA{255, 255, 255, 255}, false)
		vector.DrawFilledRect(screen, screenX+1, screenY+1, CELL_SIZE-2, CELL_SIZE-2, color.RGBA{192, 192, 192, 255}, false)
		vector.DrawFilledRect(screen, screenX+CELL_SIZE-1, screenY, 1, CELL_SIZE, color.RGBA{128, 128, 128, 255}, false)
		vector.DrawFilledRect(screen, screenX, screenY+CELL_SIZE-1, CELL_SIZE, 1, color.RGBA{128, 128, 128, 255}, false)
		
		g.drawFlag(screen, screenX, screenY)
	}
}

func (g *Game) drawMine(screen *ebiten.Image, x, y float32) {
	centerX := x + 8
	centerY := y + 8
	
	// Main mine body - black circle
	vector.DrawFilledCircle(screen, centerX, centerY, 4, color.RGBA{0, 0, 0, 255}, false)
	
	// 8 spikes radiating from center
	// Vertical spike
	vector.DrawFilledRect(screen, centerX-0.5, y+1, 1, 14, color.RGBA{0, 0, 0, 255}, false)
	// Horizontal spike
	vector.DrawFilledRect(screen, x+1, centerY-0.5, 14, 1, color.RGBA{0, 0, 0, 255}, false)
	// Diagonal spikes (4 directions)
	vector.DrawFilledRect(screen, x+2, y+2, 12, 1, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+2, y+13, 12, 1, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+2, y+2, 1, 12, color.RGBA{0, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+13, y+2, 1, 12, color.RGBA{0, 0, 0, 255}, false)
	
	// White highlight on upper left of mine
	vector.DrawFilledRect(screen, centerX-1, centerY-1, 2, 2, color.RGBA{255, 255, 255, 255}, false)
}

func (g *Game) drawFlag(screen *ebiten.Image, x, y float32) {
	// Flag pole - thin vertical line
	vector.DrawFilledRect(screen, x+8, y+2, 1, 12, color.RGBA{0, 0, 0, 255}, false)
	
	// Flag - red triangle
	vector.DrawFilledRect(screen, x+9, y+2, 5, 3, color.RGBA{255, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+9, y+5, 4, 1, color.RGBA{255, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+9, y+6, 3, 1, color.RGBA{255, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+9, y+7, 2, 1, color.RGBA{255, 0, 0, 255}, false)
	vector.DrawFilledRect(screen, x+9, y+8, 1, 1, color.RGBA{255, 0, 0, 255}, false)
}

func (g *Game) drawNumber(screen *ebiten.Image, x, y float32, number int) {
	colors := []color.RGBA{
		{0, 0, 255, 255},     // 1 - blue
		{0, 128, 0, 255},     // 2 - green
		{255, 0, 0, 255},     // 3 - red
		{0, 0, 128, 255},     // 4 - dark blue
		{128, 0, 0, 255},     // 5 - maroon
		{0, 128, 128, 255},   // 6 - teal
		{0, 0, 0, 255},       // 7 - black
		{128, 128, 128, 255}, // 8 - gray
	}
	
	numberColor := colors[number-1]
	
	// Use text package for better centered text
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
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}