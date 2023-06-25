// Package gametiles implements a simple game demo based on square tiles using ebiten for rendering.
package gametiles

import (
	"fmt"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 480
	tileSize     = 16
)

type Game struct {
	currentLevel *Level
	levelWidth   int
	levelHeight  int
	width        int
	height       int
	camX         float64
	camY         float64
	camScale     float64
	camScaleTo   float64
	mousePanX    int
	mousePanY    int
	offscreen    *ebiten.Image
	sprites      *SpriteSheet
}

func (g *Game) Update() error {
	// Update target zoom level.
	var scrollY float64
	if ebiten.IsKeyPressed(ebiten.KeyC) || ebiten.IsKeyPressed(ebiten.KeyPageDown) {
		scrollY = -0.25
	} else if ebiten.IsKeyPressed(ebiten.KeyE) || ebiten.IsKeyPressed(ebiten.KeyPageUp) {
		scrollY = 0.25
	} else {
		_, scrollY = ebiten.Wheel()
		if scrollY < -1 {
			scrollY = -1
		} else if scrollY > 1 {
			scrollY = 1
		}
	}
	// TODO: Fix camera position.
	g.camScaleTo += scrollY * (g.camScaleTo / 7)

	// Clamp target zoom level.
	if g.camScaleTo < 0.01 {
		g.camScaleTo = 0.01
	} else if g.camScaleTo > 100 {
		g.camScaleTo = 100
	}

	// Smooth zoom transition.
	div := 10.0
	if g.camScaleTo > g.camScale {
		g.camScale += (g.camScaleTo - g.camScale) / div
	} else if g.camScaleTo < g.camScale {
		g.camScale -= (g.camScale - g.camScaleTo) / div
	}

	// Pan camera via keyboard.
	pan := 7.0 / g.camScale
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.camX -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.camX += pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.camY -= pan
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.camY += pan
	}

	// Pan camera via mouse.
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if g.mousePanX == math.MinInt32 && g.mousePanY == math.MinInt32 {
			g.mousePanX, g.mousePanY = ebiten.CursorPosition()
		} else {
			x, y := ebiten.CursorPosition()
			dx, dy := float64(g.mousePanX-x)*(pan/100), float64(g.mousePanY-y)*(pan/100)
			g.camX, g.camY = g.camX-dx, g.camY+dy
		}
	} else if g.mousePanX != math.MinInt32 || g.mousePanY != math.MinInt32 {
		g.mousePanX, g.mousePanY = math.MinInt32, math.MinInt32
	}

	// Clamp camera position.
	// TODO: Fix camera clamping when zoomed in.
	worldWidth := float64(g.levelWidth * tileSize)
	worldHeight := float64(g.levelHeight * tileSize)
	if g.camX < 0 {
		g.camX = 0
	} else if g.camX > worldWidth {
		g.camX = worldWidth
	}
	if g.camY < -worldHeight {
		g.camY = -worldHeight
	} else if g.camY > 0 {
		g.camY = 0
	}

	// If we click, print the tile we clicked on.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := getTileXY(g)
		fmt.Printf("Clicked on tile %d,%d\n", x, y)
		// TODO: Add pathfinding from player to tile.
	}

	return nil
}

func getTileXY(g *Game) (int, int) {
	// Get mouse position.
	x, y := ebiten.CursorPosition()

	// Convert to world coordinates.
	x, y = g.screenToWorld(x, y)

	// Convert to tile coordinates.
	x, y = worldToTile(x, y)

	return x, y
}

func (g *Game) screenToWorld(x, y int) (int, int) {
	// Convert to world coordinates.
	x = int(float64(x)/g.camScale - g.camX)
	y = int(float64(y)/g.camScale - g.camY)

	return x, y
}

func worldToTile(x, y int) (int, int) {
	// Convert to tile coordinates.
	x = x / tileSize
	y = y / tileSize

	return x, y
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	padding := float64(tileSize) * g.camScale
	cx, cy := float64(g.width/2), float64(g.height/2)

	scaleLater := g.camScale > 1
	target := screen
	scale := g.camScale

	// When zooming in, tiles can have slight bleeding edges.
	// To avoid them, render the result on an offscreen first and then scale it later.
	if scaleLater {
		if g.offscreen != nil {
			w, h := g.offscreen.Size()
			sw, sh := screen.Size()
			if w != sw || h != sh {
				g.offscreen.Dispose()
				g.offscreen = nil
			}
		}
		if g.offscreen == nil {
			g.offscreen = ebiten.NewImage(screen.Size())
		}
		target = g.offscreen
		target.Clear()
		scale = 1
	}

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	xCount := g.levelWidth
	for i, t := range g.currentLevel.Tiles {
		// Get actual world position of tile.
		x := (i % xCount) * tileSize
		y := (i / xCount) * tileSize

		// Calculate the on-screen position of the tile.
		drawX := (float64(x)-(g.camX))*(g.camScale) + (cx)
		drawY := (float64(y)+(g.camY))*(g.camScale) + (cy)

		// Skip tiles that are not visible.
		if drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(g.width) || drawY-padding > float64(g.height) {
			continue
		}

		op.GeoM.Reset()
		// Move to current isometric position.
		op.GeoM.Translate(float64(x), float64(y))
		// Translate camera position.
		op.GeoM.Translate(-g.camX, g.camY)
		// Zoom.
		op.GeoM.Scale(scale, scale)
		// Center.
		op.GeoM.Translate(cx, cy)

		// Draw tile.
		val := t
		// Draw terrain.
		switch val.Type() {
		case TileTypeWater:
			target.DrawImage(g.sprites.Water, op)
		case TileTypeDirt:
			target.DrawImage(g.sprites.Dirt, op)
		case TileTypeGrass:
			target.DrawImage(g.sprites.Grass, op)
		case TileTypeRock:
			target.DrawImage(g.sprites.Rock, op)
		default:
			target.DrawImage(g.sprites.Snow, op)
		}

		// Draw trees if present.
		//if val.HasTrees() {
		//	target.DrawImage(g.sprites.Trees, op)
		//}
	}

	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.camScale), float64(g.camScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}

	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return g.width, g.height
}

func NewGame(lvlWidth, lvlHeight int) *Game {
	// Load sprites and convert them from 32x32 to 16x16.
	s, err := LoadSpriteSheet(32, tileSize)
	if err != nil {
		log.Fatal(err)
	}
	l, err := NewLevel(lvlWidth, lvlHeight)
	if err != nil {
		log.Fatal(err)
	}
	return &Game{
		currentLevel: l,
		levelWidth:   lvlWidth,
		levelHeight:  lvlHeight,
		camX:         float64(lvlWidth) * tileSize / 2,
		camY:         -float64(lvlHeight) * tileSize / 2,
		camScale:     1,
		camScaleTo:   1,
		mousePanX:    math.MinInt32,
		mousePanY:    math.MinInt32,
		sprites:      s,
	}
}
