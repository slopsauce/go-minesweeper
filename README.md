# Go Minesweeper

A faithful recreation of the classic Windows 3.1 Minesweeper game built in Go using the Ebiten game engine.

## Play Online

🎮 **[Play Minesweeper in your browser!](https://slopsauce.github.io/go-minesweeper/)**

No download required - the game runs directly in your web browser using WebAssembly.

## Features

- **Authentic Windows 3.1 Graphics**: Pixel-perfect recreation of the original game's appearance
- **Classic Gameplay**: 16x16 grid with 40 mines (Expert difficulty)
- **Proper Game Mechanics**:
  - First click is always safe (mines are placed after first click)
  - Efficient Fisher-Yates algorithm for mine placement
  - Flood fill algorithm for revealing adjacent empty cells
  - Right-click flagging system
  - Timer and mine counter with 999-second cap
  - Win/lose detection with proper smiley face states
- **3D Visual Effects**: Authentic raised/sunken button effects
- **LED-style Displays**: Classic red digital counter appearance
- **Professional Code Quality**: Clean, maintainable Go code following best practices
- **Web Support**: Playable in browser via WebAssembly

## Controls

- **Left Click**: Reveal cell
- **Right Click**: Flag/unflag cell
- **Smiley Face Button**: Reset game

## Installation

### Option 1: Download Pre-built Binary (Recommended)

1. Go to [Releases](https://github.com/slopsauce/go-minesweeper/releases)
2. Download the appropriate binary for your platform:
   - **Windows**: `minesweeper-windows-amd64.exe`
   - **macOS (Intel)**: `minesweeper-macos-amd64`
   - **macOS (Apple Silicon)**: `minesweeper-macos-arm64`
   - **Linux**: `minesweeper-linux-amd64`
3. Run the executable directly - no installation required!

### Option 2: Build from Source

**Prerequisites**: Go 1.21 or later, Git

```bash
git clone https://github.com/slopsauce/go-minesweeper.git
cd go-minesweeper
go run main.go
```

**Build executable:**
```bash
go build -o minesweeper main.go
./minesweeper
```


## Game Rules

1. **Objective**: Clear all cells that don't contain mines
2. **Numbers**: Indicate how many mines are adjacent to that cell
3. **Flags**: Right-click to flag suspected mines
4. **First Click**: Always safe - mines are placed after your first move
5. **Winning**: Reveal all non-mine cells
6. **Losing**: Click on a mine

## Technical Details

- **Engine**: [Ebiten](https://ebiten.org/) - A dead simple 2D game library for Go
- **Graphics**: Custom pixel-perfect rendering matching Windows 3.1 style
- **Window Size**: 276x336 pixels (authentic size)
- **Performance**: 60 FPS, optimized for smooth gameplay
- **Algorithm**: Fisher-Yates shuffle for O(n) mine placement
- **Code Quality**: Professional Go code with proper error handling and bounds checking
- **Web Support**: Compiled to WebAssembly for browser deployment
- **CI/CD**: GitHub Actions automatically builds and deploys to GitHub Pages

## Code Structure

- `main.go` - Complete game implementation (497 lines)
- Clean separation between game logic and rendering
- Authentic Windows 3.1 UI rendering with 3D effects
- Efficient mouse input handling with coordinate conversion
- Optimized mine placement and flood fill algorithms
- Color constants and helper methods for maintainability

## Smiley Face States

- **😊 Normal**: Default playing state
- **😵 Dead**: Game over (hit a mine)
- **😎 Cool**: Victory (all mines found)

## Dependencies

```go
github.com/hajimehoshi/ebiten/v2
golang.org/x/image/font/basicfont
```

## Screenshots

The game faithfully recreates the Windows 3.1 Minesweeper experience:

- Classic gray 3D interface
- Authentic mine graphics with spikes and highlight
- Proper flag triangular design
- LED-style red digital displays
- Pixel-perfect smiley face expressions

## Contributing

Feel free to submit issues and pull requests! This project aims to maintain authenticity to the original Windows 3.1 Minesweeper while being written in idiomatic Go.

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Original Minesweeper game by Microsoft
- Ebiten game engine by Hajime Hoshi
- Windows 3.1 design inspiration

---

🤖 **Generated with [Claude Code](https://claude.ai/code)**

**Co-Authored-By:** Claude <noreply@anthropic.com>