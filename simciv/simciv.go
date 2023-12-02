// Package simciv is a playground for simulating the spread of civilization.
package simciv

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/vector"
	"github.com/mazznoer/colorgrad"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 480
	tileSize     = 16
)

type Game struct {
	Map        *Map
	width      int
	height     int
	camX       float64
	camY       float64
	camScale   float64
	camScaleTo float64
	mousePanX  int
	mousePanY  int
	offscreen  *ebiten.Image
}

func NewGame(lvlWidth, lvlHeight int) *Game {
	return &Game{
		Map:        NewMap(lvlWidth, lvlHeight, 0),
		camX:       float64(lvlWidth) * tileSize / 2,
		camY:       -float64(lvlHeight) * tileSize / 2,
		camScale:   1,
		camScaleTo: 1,
		mousePanX:  math.MinInt32,
		mousePanY:  math.MinInt32,
	}
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
	worldWidth := float64(g.Map.levelWidth * tileSize)
	worldHeight := float64(g.Map.levelHeight * tileSize)
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
	g.Map.Tick(16)

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

	// Create a gradient that looks like terrain.
	// From white mountain tops to grey mountainsides
	// to dark green plains.
	grad, err := colorgrad.NewGradient().
		Colors(
			color.RGBA{255, 255, 255, 255}, // white
			color.RGBA{128, 128, 128, 255}, // grey
			color.RGBA{128, 128, 128, 255}, // grey
			color.RGBA{0, 128, 0, 255},     // dark green
			color.RGBA{0, 128, 0, 255},     // dark green
		).Interpolation(colorgrad.InterpolationBasis).
		Build()
	if err != nil {
		panic(err)
	}

	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	xCount := g.Map.levelWidth
	for i, t := range g.Map.score {
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
		// Move to current tile position.
		op.GeoM.Translate(float64(x), float64(y))
		// Translate camera position.
		op.GeoM.Translate(-g.camX, g.camY)
		// Zoom.
		op.GeoM.Scale(scale, scale)
		// Center.
		op.GeoM.Translate(cx, cy)

		// Draw tile score to reflect the terrain.
		col := grad.At(t)

		ebitenutil.DrawRect(target, drawX, drawY, tileSize*scale, tileSize*scale, col)
	}

	var maxPop int
	for _, s := range g.Map.Settlements {
		if s.pop > maxPop {
			maxPop = s.pop
		}
	}

	// Draw trade radius.
	for _, s := range g.Map.Settlements {
		if s.pop == 0 {
			continue
		}
		// Get actual world position of tile.
		x := s.x * tileSize
		y := s.y * tileSize

		// Calculate the on-screen position of the tile.
		drawX := (float64(x)-(g.camX))*(g.camScale) + (cx)
		drawY := (float64(y)+(g.camY))*(g.camScale) + (cy)

		// Skip tiles that are not visible.
		if drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(g.width) || drawY-padding > float64(g.height) {
			continue
		}

		op.GeoM.Reset()
		// Move to current tile position.
		op.GeoM.Translate(float64(x), float64(y))
		// Translate camera position.
		op.GeoM.Translate(-g.camX, g.camY)
		// Zoom.
		op.GeoM.Scale(scale, scale)
		// Center.
		op.GeoM.Translate(cx, cy)

		// Depending on the population, the color changes.
		colVal := uint8(255 * float64(s.pop) / float64(maxPop))

		// Draw the trade radius.
		vector.DrawFilledCircle(target, float32(drawX), float32(drawY), float32(s.getTradeRadius()*tileSize*scale), color.RGBA{0, 0, colVal, 10}, true)

		// Draw the protection radius.
		vector.DrawFilledCircle(target, float32(drawX), float32(drawY), float32(s.getProtectionRadius()*tileSize*scale), color.RGBA{colVal, 0, 0, 10}, true)
	}

	// Draw settlements.
	for _, s := range g.Map.Settlements {
		// Get actual world position of tile.
		x := s.x * tileSize
		y := s.y * tileSize

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

		// Draw tile in grayscale.
		// Depending on the population, the color changes.
		colVal := uint8(255 * float64(s.pop) / float64(maxPop))
		vector.DrawFilledRect(target, float32(drawX), float32(drawY), float32(tileSize*scale), float32(tileSize*scale), color.RGBA{0, colVal, 0, 255}, true)
	}

	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-cx, -cy)
		op.GeoM.Scale(float64(g.camScale), float64(g.camScale))
		op.GeoM.Translate(cx, cy)
		screen.DrawImage(target, op)
	}

	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f\nDATE %d/%d", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY, g.Map.day, g.Map.year))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return g.width, g.height
}
