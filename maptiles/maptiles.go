// Package maptiles implements a map tiles skeleton package implementing the tiles á la Google maps.
// See: https://www.maptiler.com/google-maps-coordinates-tile-bounds-projection
package maptiles

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hajimehoshi/ebiten/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 480
	tileSize     = 256
)

type Game struct {
	currentZoom     int // Current zoom level.
	numTilesPerAxis int // Number of tiles on each axis at the current zoom level.
	width           int
	height          int
	camX            float64
	camY            float64
	camScale        float64
	camScaleTo      float64
	mousePanX       int
	mousePanY       int
	Font            font.Face
}

func NewGame() (*Game, error) {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		return nil, err
	}
	const dpi = 72
	mplusNormalFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		return nil, err
	}
	return &Game{
		currentZoom:     0,
		numTilesPerAxis: 1,
		Font:            mplusNormalFont,
		camX:            tileSize / 2,
		camY:            -tileSize / 2,
		camScale:        1,
		camScaleTo:      1,
		mousePanX:       math.MinInt32,
		mousePanY:       math.MinInt32,
	}, nil
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
	const maxZoom = 28 // Same as Google maps.
	if g.camScaleTo < 0.01 {
		g.camScaleTo = 0.01
	} else if g.camScaleTo > maxZoom {
		g.camScaleTo = maxZoom
	}

	// Smooth zoom transition.
	div := 10.0
	if g.camScaleTo > g.camScale {
		g.camScale += (g.camScaleTo - g.camScale) / div
	} else if g.camScaleTo < g.camScale {
		g.camScale -= (g.camScale - g.camScaleTo) / div
	}

	// Clamp zoom level.
	if g.camScale < 0.01 {
		g.camScale = 0.01
	} else if g.camScale > maxZoom {
		g.camScale = maxZoom
	}

	// Update the number of tiles at the current zoom level.
	if g.currentZoom != int(g.camScale) {
		g.currentZoom = int(g.camScale)
		g.numTilesPerAxis = zoomToNumTiles(g.currentZoom)
	}

	// Calculate the pan and zoom speed based on the current zoom level.
	panScale := 1 / math.Pow(2, g.camScale)

	// Pan camera via keyboard.
	pan := 7.0 * panScale
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
	worldWidth := float64(tileSize)
	worldHeight := float64(tileSize)
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
		fmt.Printf("Clicked on tile %d,%d,%d\n", x, y, g.currentZoom)
	}

	return nil
}

func getTileXY(g *Game) (int, int) {
	// Get mouse position.
	x, y := ebiten.CursorPosition()

	// Convert to world coordinates.
	x, y = g.screenToWorld(x, y)

	// Convert to tile coordinates.
	x, y = g.worldToTile(x, y)

	return x, y
}

func (g *Game) screenToWorld(x, y int) (int, int) {
	// Convert to world coordinates.
	scale := math.Pow(2, g.camScale)

	cx, cy := float64(g.width)/2, float64(g.height)/2

	wx := (((float64(x) - cx) / scale) + g.camX)
	wy := (((float64(y) - cy) / scale) - g.camY)
	return int(wx), int(wy)
}

func (g *Game) worldToTile(x, y int) (int, int) {
	// Convert to tile coordinates.
	currTileSize := float64(tileSize) / float64(g.numTilesPerAxis)
	xf := float64(x) / currTileSize
	yf := float64(y) / currTileSize
	return int(xf), int(yf)
}

var lastSize float64

func (g *Game) Draw(screen *ebiten.Image) {
	zoom := int(g.camScale)

	// Calculate the tile size at the current zoom level.
	numTiles := float64(g.numTilesPerAxis)
	currTileSize := float64(tileSize) / numTiles

	// Scaling factor for the tiles.
	// This is to scale the tiles to the current zoom level (g.camScale).
	// Since the tiles shrink exponentially as the zoom level increases,
	// we need to scale the tiles by 2^zoom to scale them at the same rate.
	scale := math.Pow(2, g.camScale)

	cx, cy := float64(g.width)/2, float64(g.height)/2

	// Determine bounding box of visible tiles and only draw those instead
	// of testing all tiles.
	tx0 := ((-1-cx)/scale + g.camX) / currTileSize
	ty0 := ((-1-cy)/scale - g.camY) / currTileSize
	tx1 := ((float64(g.width)+1-cx)/scale + g.camX) / currTileSize
	ty1 := ((float64(g.height)+1-cy)/scale - g.camY) / currTileSize

	// Calculate the range of tiles to draw.
	minTileX := int(math.Max(0, math.Floor(tx0)))
	minTileY := int(math.Max(0, math.Floor(ty0)))
	maxTileX := int(math.Min(numTiles, math.Ceil(tx1)))
	maxTileY := int(math.Min(numTiles, math.Ceil(ty1)))

	for xt := minTileX; xt < maxTileX; xt++ {
		for yt := minTileY; yt < maxTileY; yt++ {
			x, y := float64(xt)*currTileSize, float64(yt)*currTileSize

			// Calculate the on-screen position of the tile.
			drawX := (float64(x)-(g.camX))*(scale) + (cx)
			drawY := (float64(y)+(g.camY))*(scale) + (cy)

			// Draw the tile label at the tile's position.
			g.drawLabel(screen, drawX, drawY, fmt.Sprintf("%d,%d,%d", xt, yt, zoom), color.White)

			// Draw an empty rectangle around the tile.
			vector.StrokeRect(screen, float32(drawX), float32(drawY), float32(currTileSize*(scale)), float32(currTileSize*(scale)), 1, color.White, false)
		}
	}
	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY))
}

func tileToScreenPixels(x, y, zoom int) (float64, float64) {
	n := 1 << uint(zoom)
	return float64(x) * tileSize / float64(n), float64(y) * tileSize / float64(n)
}

func zoomToNumTiles(zoom int) int {
	return 1 << uint(zoom)
}

func tileSizeAtZoom(zoom int) float64 {
	return tileSize / float64(zoomToNumTiles(zoom))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.width = outsideWidth
	g.height = outsideHeight
	return g.width, g.height
}

func (g *Game) drawLabel(screen *ebiten.Image, x, y float64, label string, col color.Color) {
	// Calculate the on-screen position of the label.
	drawX, drawY := x, y

	screen.Set(int(drawX), int(drawY), col)
	text.Draw(screen, label, g.Font, int(drawX), int(drawY), col)
}
