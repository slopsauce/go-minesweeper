# Claude Development Log

This project was developed using Claude Code, an AI assistant for software development.

## Project Overview

**Task**: Implement a Windows 3.1-style Minesweeper game in Go using Ebiten with complete production-ready setup
**Development Date**: July 18-19, 22, 2025

## Development Memories and Best Practices

- Always ensure CLAUDE.md and README.md are up-to-date before committing, tagging, or releasing.
- The game is now playable online at https://slopsauce.github.io/go-minesweeper/
- WebAssembly builds are automatically updated via GitHub Actions on push to main

## WebAssembly Deployment (July 22, 2025)

Successfully deployed the game to GitHub Pages using WebAssembly:

1. **Created Web Build Infrastructure**:
   - Added `build-web.sh` script for WebAssembly compilation
   - Created HTML wrapper files (`index.html` and `game.html`)
   - Styled with Windows 3.1 desktop theme (teal background #008080)

2. **GitHub Pages Setup**:
   - Deployed to `docs/` directory
   - Enabled GitHub Pages via API: `gh api repos/slopsauce/go-minesweeper/pages -X POST --field 'source[branch]=main' --field 'source[path]=/docs'`
   - Game accessible at: https://slopsauce.github.io/go-minesweeper/

3. **Continuous Deployment**:
   - Added GitHub Actions workflow (`.github/workflows/deploy-web.yml`)
   - Automatically rebuilds WASM on push to main
   - Updates `docs/` folder with latest build
   - No manual intervention needed for updates

4. **Technical Details**:
   - Ebiten v2.8.8 has excellent WebAssembly support
   - No code changes required - game works identically in browser
   - WASM file size: ~9.4MB (compresses well with gzip)
   - Uses Go's standard `wasm_exec.js` runtime

## Complete Development Process

The initial development created a pixel-perfect Windows 3.1 Minesweeper clone with:
- Authentic graphics and 3D effects
- Proper game mechanics (first click safety, flood fill, etc.)
- Professional code quality
- Cross-platform desktop builds
- GitHub releases with pre-built binaries

The web deployment phase added:
- Browser-based gameplay via WebAssembly
- Automatic deployment pipeline
- Zero-maintenance updates

## Key Commands

```bash
# Deploy to GitHub Pages (automatic via push to main)
git push origin main

# The game is automatically built and deployed via GitHub Actions
```