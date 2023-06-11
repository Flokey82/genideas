// Package gamehex implements an example hexagonal game board and uses ebiten for rendering.
package gamehex

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	currentLevel *Level
	width        int
	height       int
	offscreen    *ebiten.Image
	clickedTileX int
	clickedTileY int
}

// NewGame returns a new isometric demo Game.
func NewGame() (*Game, error) {
	l, err := NewLevel(32, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to create new level: %s", err)
	}

	return &Game{
		currentLevel: l,
	}, nil
}

// Update reads current user input and updates the Game state.
func (g *Game) Update() error {
	// If we have a mouse click, we calculate the tile we clicked on and store it
	// for rendering later.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.clickedTileX, g.clickedTileY = g.currentLevel.PosToHexTileXY(x, y)
	}
	return nil
}

// Draw draws the Game on the screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// Render level.
	g.renderLevel(screen)
}

// Layout is called when the Game's layout changes.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width, g.height = outsideWidth, outsideHeight
	return g.width, g.height
}

// renderLevel draws the current Level on the screen.
func (g *Game) renderLevel(screen *ebiten.Image) {
	for y := 0; y < g.currentLevel.Height; y++ {
		for x := 0; x < g.currentLevel.Width; x++ {
			if g.clickedTileX == x && g.clickedTileY == y {
				continue
			}
			px, py := g.currentLevel.HexTileXYToPixelPos(x, y)
			g.currentLevel.drawHex(screen, float64(px), float64(py), float64(g.currentLevel.tileSize/2))
		}
	}
}
