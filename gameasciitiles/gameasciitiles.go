package gameasciitiles

import (
	"bytes"
	"image"
	"image/color"
	"math"

	"github.com/disintegration/imaging"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/images"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	_ "embed"
)

var (
	//go:embed vendor/bescii.ttf
	bescii_ttf []byte
)

const (
	screenWidth  = 640
	screenHeight = 520
)

type Game struct {
	Font         font.Face // Font used to render the tiles
	FontSmall    font.Face // Font used to render the tiles
	FontTiny     font.Face // Font used to render the tiles
	TileSize     int       // Size of a tile in pixels
	ScreenHeight int       // Height of the screen in tiles
	ScreenWidth  int       // Width of the screen in tiles
	SpriteSheet  *SpriteSheet
	runnerImage  *ebiten.Image
	count        int
}

func NewGame(tileSize int) (*Game, error) {
	tt, err := opentype.Parse(bescii_ttf)
	if err != nil {
		return nil, err
	}

	fontBig, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(tileSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	fontSmall, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(tileSize) / 2,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	fontTiny, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(tileSize) / 4,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}

	ss, err := LoadSpriteSheet(16, tileSize)
	if err != nil {
		return nil, err
	}
	// Decode an image from the image file's byte slice.
	runnerImg, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		return nil, err
	}

	return &Game{
		Font:         fontBig,
		FontSmall:    fontSmall,
		FontTiny:     fontTiny,
		TileSize:     tileSize,
		ScreenHeight: 35,
		ScreenWidth:  60,
		SpriteSheet:  ss,
		runnerImage:  ebiten.NewImageFromImage(runnerImg),
	}, nil
}

func (g *Game) Update() error {
	g.count++
	return nil
}

type colorScale struct {
	R float64
	G float64
	B float64
	A float64
}

func (g *Game) getTilePos(x, y int) (int, int) {
	return x * g.TileSize, y * g.TileSize
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fill screen with grass.
	for y := 0; y < g.ScreenHeight; y++ {
		for x := 0; x < g.ScreenWidth; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*g.TileSize), float64(y*g.TileSize))
			screen.DrawImage(g.SpriteSheet.Grass[x*y%4], op)
		}
	}
	text.Draw(screen, "Hello, World!", g.Font, 0, g.TileSize, color.White)
	text.Draw(screen, "╭──┬────┬──╮", g.Font, 0, 3*g.TileSize, color.White)
	text.Draw(screen, "│  │    │  │", g.Font, 0, 4*g.TileSize, color.White)
	text.Draw(screen, "│  │   ╭╯  │", g.Font, 0, 5*g.TileSize, color.White)
	text.Draw(screen, "│  │   │   │", g.Font, 0, 6*g.TileSize, color.White)
	text.Draw(screen, "╰──┴───┴───╯", g.Font, 0, 7*g.TileSize, color.White)

	// TODO: Apply drop shadow to layer that contains glyphs.

	drawGlyph := func(glyph rune, x, y int, glyphColor colorScale, superscript *SuperScript, bounce bool) {
		dr, img, _, _, ok := g.FontSmall.Glyph(fixed.Point26_6{}, glyph)
		if ok {
			// Get the tile position.
			tilePosX, tilePosY := g.getTilePos(x, y)

			// Draw the glyph onto the screen.
			op := &ebiten.DrawImageOptions{}

			// Translate the glyph to the correct position (center of tile).
			bounds := dr.Bounds()
			size := bounds.Size()
			dx := (g.TileSize - size.X) / 2
			dy := (g.TileSize - size.Y) / 2

			// Bounce in a sine based on count.
			if bounce {
				dy += int(float64(g.TileSize) * 0.05 * math.Sin(2*math.Pi*float64(g.count%50)/50))
			}

			op.GeoM.Translate(float64(tilePosX)+float64(dx), float64(tilePosY)+float64(dy))
			// Change color of the glyph.
			op.ColorM.Scale(glyphColor.R, glyphColor.G, glyphColor.B, glyphColor.A)
			img = ebiten.NewImageFromImage(DropShadow(img, 3))
			eImg := ebiten.NewImageFromImage(img)
			screen.DrawImage(eImg, op)

			if superscript != nil {
				// Superscript of 'superscript'
				superImg := g.assembleStringImage(g.FontTiny, superscript.Text)
				op = &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tilePosX)+float64(dx)+float64(dr.Dx()), float64(tilePosY)+float64(dy)-float64(g.TileSize)/4)
				// Change color of the superscript.
				op.ColorM.Scale(superscript.Color.R, superscript.Color.G, superscript.Color.B, superscript.Color.A)
				superImg = ebiten.NewImageFromImage(DropShadow(superImg, 2))
				screen.DrawImage(superImg, op)
			}
		}
	}

	drawGlyph('@', 0, 1, colorScale{0, 1, 0, 1}, newSuperScript("!", 1, 0, 0, 1), false)
	drawGlyph('!', 1, 1, colorScale{1, 0, 0, 1}, newSuperScript("?", 0, 0, 1, 1), true)
	drawGlyph('2', 2, 1, colorScale{0, 0, 1, 1}, nil, false)
	drawGlyph('X', 3, 1, colorScale{1, 0, 0, 1}, nil, true)

	const (
		frameOX     = 0
		frameOY     = 32
		frameWidth  = 32
		frameHeight = 32
		frameCount  = 8
	)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(4*float64(g.TileSize), 1*float64(g.TileSize))
	i := (g.count / 5) % frameCount
	sx, sy := frameOX+i*frameWidth, frameOY
	screen.DrawImage(g.runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	drawTileIDXY := func(x, y int, id int) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x*g.TileSize), float64(y*g.TileSize))
		screen.DrawImage(g.SpriteSheet.GetSubImageID(id), op)
	}

	drawTileIDXY(4, 7, 176)
	drawTileIDXY(5, 7, 177)
	drawTileIDXY(6, 7, 184)
	drawTileIDXY(4, 8, 201)
	drawTileIDXY(5, 8, 202)
	drawTileIDXY(6, 8, 209)
	drawTileIDXY(4, 9, 226)
	drawTileIDXY(5, 9, 227)
	drawTileIDXY(6, 9, 234)
	drawTileIDXY(4, 10, 251)
	drawTileIDXY(5, 10, 252)
	drawTileIDXY(6, 10, 259)
	drawTileIDXY(4, 11, 276)
	drawTileIDXY(5, 11, 277)
	drawTileIDXY(6, 11, 284)

}

type SuperScript struct {
	Text  string
	Color colorScale
}

func newSuperScript(text string, r, g, b, a float64) *SuperScript {
	return &SuperScript{
		Text:  text,
		Color: colorScale{r, g, b, a},
	}
}

func (g *Game) assembleStringImage(font font.Face, txt string) *ebiten.Image {
	metrics := font.Metrics()
	height := metrics.Height.Ceil()
	rect := text.BoundString(font, txt)
	bounds := rect.Bounds()
	size := bounds.Size()
	img := ebiten.NewImage(size.X+3, height+2)
	text.Draw(img, txt, font, 0, height, color.White)
	return img
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Run() error {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("GameASCIITiles (Ebitengine Demo)")
	return ebiten.RunGame(g)
}

func DropShadow(img image.Image, size float64) image.Image {
	bounds := img.Bounds()
	sizeInt := int(math.Ceil(size)) * 4
	final := imaging.New(bounds.Dx()+sizeInt, bounds.Dy()+sizeInt, color.Alpha{})
	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			_, _, _, a := img.At(x, y).RGBA()
			final.Set(x+sizeInt/2, y+sizeInt/2, color.RGBA{0x0, 0x0, 0x0, uint8(a / 2)})
		}
	}
	final = imaging.Blur(final, size)
	final = imaging.Overlay(final, img, image.Point{sizeInt / 2, sizeInt / 2}, 1)
	return final
}

const dpi = 72
