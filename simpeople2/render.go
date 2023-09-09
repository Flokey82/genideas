package simpeople2

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/mazznoer/colorgrad"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 480
	tileSize     = 16
)

type Render struct {
	levelWidth   int
	levelHeight  int
	screenWidth  int
	screenHeight int
	camX         float64
	camY         float64
	camScale     float64
	camScaleTo   float64
	mousePanX    int
	mousePanY    int
	Font         font.Face
	BigFont      font.Face
	offscreen    *ebiten.Image
	target       *ebiten.Image
	sprites      *SpriteSheet
	runner       *SpriteSheet
	count        int // animation frame counter
}

func newRender(width, height int) (*Render, error) {
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
	mplusBigFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		return nil, err
	}
	// Load tiles.
	sTiles, err := LoadSpriteSheet(16, tileSize, Spritesheet_png)
	if err != nil {
		return nil, err
	}

	// Load runner animation and convert the frames from 32x32 to 16x16.
	// TODO: Move to TileSet.
	sRunner, err := LoadSpriteSheet(32, tileSize, images.Runner_png)
	if err != nil {
		log.Fatal(err)
	}
	return &Render{
		levelWidth:  width,
		levelHeight: height,
		camX:        float64(width) * tileSize / 2,
		camY:        -float64(height) * tileSize / 2,
		camScale:    1,
		camScaleTo:  1,
		mousePanX:   math.MinInt32,
		mousePanY:   math.MinInt32,
		Font:        mplusNormalFont,
		BigFont:     mplusBigFont,
		sprites:     sTiles,
		runner:      sRunner,
	}, nil
}

// Layout is called when the Game's layout changes.
func (g *Render) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth, g.screenHeight = outsideWidth, outsideHeight
	return g.screenWidth, g.screenHeight
}

func (g *Render) handleInput() {
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
		x, y := g.getTileXY()
		fmt.Printf("Clicked on tile %d,%d\n", x, y)
		// TODO: Add pathfinding from player to tile.
	}
}

func (g *Render) getTileXY() (int, int) {
	// Get mouse position.
	x, y := ebiten.CursorPosition()

	// Convert to world coordinates.
	x, y = g.screenToWorld(x, y)

	// Convert to tile coordinates.
	x, y = worldToTile(x, y)

	return x, y
}

func (g *Render) screenToWorld(x, y int) (int, int) {
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

func (g *World) Draw(screen *ebiten.Image) {
	padding := float64(tileSize) * g.camScale
	cx, cy := float64(g.screenWidth/2), float64(g.screenHeight/2)

	// Check if we have zoomed in (scaleLater) or not.
	scaleLater := g.camScale > 1
	g.target = screen

	// When zooming in, tiles can have slight bleeding edges.
	// To avoid them, render the result on an offscreen first and then scale it later.
	if scaleLater {
		// Check if we need to re-initialize the offscreen.
		if g.offscreen != nil {
			// If the screen size has changed, dispose the old offscreen.
			w, h := g.offscreen.Size()
			sw, sh := screen.Size()
			if w != sw || h != sh {
				g.offscreen.Dispose()
				g.offscreen = nil
			}
		}

		// Create a new offscreen if needed.
		if g.offscreen == nil {
			g.offscreen = ebiten.NewImage(screen.Size())
		}

		// Render to the offscreen.
		g.target = g.offscreen
		g.target.Clear()
	}

	isOutOfBounds := func(x, y float64) bool {
		// Calculate the on-screen position of the tile to verify if it is visible.
		drawX := (x-g.camX)*g.camScale + cx
		drawY := (y+g.camY)*g.camScale + cy

		// Skip tiles that are not visible.
		return drawX+padding < 0 || drawY+padding < 0 || drawX-padding > float64(g.screenWidth) || drawY-padding > float64(g.screenHeight)
	}

	op := &ebiten.DrawImageOptions{}
	// Draw each tile with each DrawImage call.
	// As the source images of all DrawImage calls are always the same,
	// this rendering is done very efficiently.
	// For more detail, see https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Image.DrawImage
	getDrawOp := func(x, y float64) *ebiten.DrawImageOptions {
		// Skip tiles that are not visible.
		if isOutOfBounds(x, y) {
			return nil
		}

		op.GeoM.Reset()
		// Move to current position P (in absolute world coordinates).
		// NOTE: I use screen coordinates, so 0, 0 is the top left corner.
		//
		//                              | - y
		//                              |
		//                              |
		//                              |           *P
		//                              |
		//                              |
		//                              |                       *C
		//  ____________________________|______________________________ x
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//                              | y

		// .... and translate the point to be relative to the camera position C. (P - C)
		//
		//                              | - y
		//                              |
		//                              |
		//                              |
		//                  *P          |
		//                              |
		//                              |
		//  ____________________________*C_____________________________ x
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//                              | y
		// ______________________________
		// |                            |
		// |                            |
		// |              ______________|______________
		// |             |   *P         |             |
		// |             |   |          |             |
		// |             |   |          |             |
		// |             |   |----------*C            |
		// |_____________|______________|             |
		//               |                            |
		//               |                            |
		//               |____________________________|
		op.GeoM.Translate(x-g.camX, y+g.camY)

		// Zoom.
		if !scaleLater {
			// Scale the image relative to the camera position.
			// The higher the zoom, the larger the distance in pixels between the camera and the point.
			// ______________________________
			// |                            | - y
			// |                            |
			// |                            |
			// |                 *P         |
			// |                 |    ______|______
			// |                 |    |     |     |
			// |      -----------|----|-----*C----|-------------------- x
			// |______________________|_____|_____| <- Scaled by 2
			//                              |
			//                              |
			//                              |
			//                              |
			//                              | y
			op.GeoM.Scale(g.camScale, g.camScale)
		}
		// Translate the point from "relative to camera" to "relative to screen origin"
		// (with the camera being the center of the screen SC). (P - C) + SC/2
		//
		//                              | - y
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//                              |
		//  ____________________________|______________________________ x
		//                              |   *P
		//                              |
		//                              |
		//                              |              *C
		//                              |
		//                              |
		//                              | y
		// Screen:
		// ______________________________
		// |   *P                       |
		// |   |                        |
		// |   |                        |
		// |   |----------*C            |
		// |                            |
		// |                            |
		// |                            |
		// |____________________________|
		op.GeoM.Translate(cx, cy)

		return op
	}

	render := func(target *ebiten.Image) {
		xCount := g.levelWidth
		for i := 0; i < g.levelWidth*g.levelHeight; i++ {
			// Get actual world position of the tile.
			x := float64((i % xCount) * tileSize)
			y := float64((i / xCount) * tileSize)

			if op := getDrawOp(x, y); op != nil {
				// Draw ground.
				target.DrawImage(g.sprites.GetSubImageID(g.Level.GetGround(i%xCount, i/xCount)), op)

				// Draw tile.
				if t := g.Level.GetTile(i%xCount, i/xCount); t != 0 {
					target.DrawImage(g.sprites.GetSubImageID(t), op)
				}
			}
		}

		// Draw the objects on top.
		for _, o := range g.Objects {
			// Get actual world position of the object.
			x := o.Position.X * tileSize
			y := o.Position.Y * tileSize

			if op := getDrawOp(x, y); op != nil {
				target.DrawImage(g.sprites.GetSubImageID(286), op)
			}
		}

		// Draw the players on top.
		for _, p := range g.People {
			// Get actual world position of the person.
			x := p.Position.X * tileSize
			y := p.Position.Y * tileSize

			if op := getDrawOp(x, y); op != nil {
				// Draw animation frame.
				const frameCount = 8
				target.DrawImage(g.runner.GetSubImageXY((g.count/5)%frameCount, 1), op)
			}
		}
	}

	render(g.target)

	if scaleLater {
		op := &ebiten.DrawImageOptions{}
		// Translate the origin from the top left corner of the screen to the center (where the camera is).
		op.GeoM.Translate(-cx, -cy)
		// Scale the image relative to the camera position.
		op.GeoM.Scale(g.camScale, g.camScale)
		// Translate the point from "relative to camera" to "relative to screen origin"
		op.GeoM.Translate(cx, cy)
		// Draw the offscreen buffer to the screen.
		screen.DrawImage(g.target, op)
	}

	// Render labels
	drawLabel := func(x, y float64, label string, col color.Color) {
		// Calculate the on-screen position of the label.
		drawX := ((x*tileSize)-g.camX)*g.camScale + cx
		drawY := ((y*tileSize)+g.camY)*g.camScale + cy
		screen.Set(int(drawX), int(drawY), col)
		text.Draw(screen, label, g.Font, int(drawX), int(drawY), col)
	}

	// Generate a color gradient for the happiness.
	colorGrad := colorgrad.Rainbow()
	colsHappiness := colorGrad.Colors(20)

	// Generate a unique color for each person and object.
	colsEntity := colorGrad.Colors(uint(len(g.People) + len(g.Objects) + 1))

	// Draw the name of the person with a color based on their happiness.
	for _, p := range g.People {
		drawLabel(p.Position.X, p.Position.Y, p.Name, colsHappiness[int((p.Happiness()+100)/10)])
	}

	// Draw the name of the object.
	for i, p := range g.Objects {
		drawLabel(p.Position.X, p.Position.Y, p.Name, colsEntity[i+len(g.People)])
	}

	// Print the attributes of each person in the sidebar.
	sideBar := 200                      // Width of the sidebar.
	sideBarX := g.screenWidth - sideBar // X offset of the sidebar.

	for i, p := range g.People {
		// Print stats of person.
		startOffset := i*20*8 + 20
		col := colsEntity[i]
		text.Draw(screen, p.Name, g.Font, sideBarX, startOffset, col)

		// Draw the state of the motives, happiness and current action.
		lines := []string{
			p.Motives[0].String(),
			p.Motives[1].String(),
			p.Motives[2].String(),
			p.Motives[3].String(),
			p.Motives[4].String(),
			fmt.Sprintf("Happiness: %.2f", p.Happiness()),
		}
		if p.Action != nil {
			lines = append(lines, fmt.Sprintf("Action: %s", p.Action.Name))
		} else {
			lines = append(lines, "Action: None")
		}
		text.Draw(screen, strings.Join(lines, "\n"), g.Font, sideBarX, startOffset+20, color.White)
	}

	// Print game info.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("KEYS WASD EC\nFPS  %0.0f\nTPS  %0.0f\nSCA  %0.2f\nPOS  %0.0f,%0.0f", ebiten.ActualFPS(), ebiten.ActualTPS(), g.camScale, g.camX, g.camY))
}
