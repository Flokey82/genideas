package dfstyle

import (
	"image/color"
	"math"
	"math/rand"
	"strconv"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
)

type Libtcod struct {
	con *console.Console
}

func newLibTcod(width, height int) *Libtcod {
	con, err := console.New(width, height, font.DefaultFont, "ramen example")
	if err != nil {
		panic(err)
	}
	return &Libtcod{
		con: con,
	}
}

func (l *Libtcod) ConsoleSetCustomFont(filename string, flags int) {
}

func (l *Libtcod) ConsolePutCharEx(con int, x, y int, ch rune, fg, bg color.RGBA) {
	if l.con == nil {
		panic("console not initialized")
	}
	l.con.Transform(x, y, t.CharByte(byte(ch)), t.Foreground(concolor.RGB(fg.R, fg.G, fg.B)), t.Background(concolor.RGB(bg.R, bg.G, bg.B)))
}

func (l *Libtcod) ConsoleFlush() {
}

func (l *Libtcod) ConsoleClear() {
	//l.con.ClearAll()                                           // clear console
	//l.con.TransformAll(t.Background(concolor.RGB(50, 50, 50))) // set the background
}

func (l *Libtcod) ConsoleCheckForKeypress(w bool) bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

var libtcod = &Libtcod{}

func clampF(min, max, v float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func fabsmodf(x, y float64) float64 {
	return float64(math.Mod(float64(math.Abs(float64(x))), float64(y)))
}

func colorSetHSV(color *color.RGBA, hue, saturation, value float64) {
	var hueSection int
	var hueFraction, p, q, t float64

	saturation = clampF(0.0, 1.0, saturation)
	value = clampF(0.0, 1.0, value)
	if saturation == 0.0 { // achromatic (grey)
		c := uint8(value*255.0 + 0.5)
		color.R = c
		color.G = c
		color.B = c
		return
	}

	hue = fabsmodf(hue, 360.0)
	hue /= 60.0 // sector 0 to 5
	hueSection = int(math.Floor(float64(hue)))
	hueFraction = hue - float64(hueSection) // fraction between sections
	p = value * (1 - saturation)
	q = value * (1 - saturation*hueFraction)
	t = value * (1 - saturation*(1-hueFraction))

	switch hueSection {
	default:
		fallthrough
	case 0: // red/yellow
		color.R = uint8(value*255.0 + 0.5)
		color.G = uint8(t*255.0 + 0.5)
		color.B = uint8(p*255.0 + 0.5)
	case 1: // yellow/green
		color.R = uint8(q*255.0 + 0.5)
		color.G = uint8(value*255.0 + 0.5)
		color.B = uint8(p*255.0 + 0.5)
	case 2: // green/cyan
		color.R = uint8(p*255.0 + 0.5)
		color.G = uint8(value*255.0 + 0.5)
		color.B = uint8(t*255.0 + 0.5)
	case 3: // cyan/blue
		color.R = uint8(p*255.0 + 0.5)
		color.G = uint8(q*255.0 + 0.5)
		color.B = uint8(value*255.0 + 0.5)
	case 4: // blue/purple
		color.R = uint8(t*255.0 + 0.5)
		color.G = uint8(p*255.0 + 0.5)
		color.B = uint8(value*255.0 + 0.5)
	case 5: // purple/red
		color.R = uint8(value*255.0 + 0.5)
		color.G = uint8(p*255.0 + 0.5)
		color.B = uint8(q*255.0 + 0.5)
	}
}

func colorLerp(c1, c2 color.RGBA, coef float64) color.RGBA {
	return color.RGBA{
		c1.R + uint8(float64(c2.R-c1.R)*coef),
		c1.G + uint8(float64(c2.G-c1.G)*coef),
		c1.B + uint8(float64(c2.B-c1.B)*coef),
		255,
	}
}

var (
	libtcodWhite         = color.RGBA{255, 255, 255, 255}
	libtcodBlack         = color.RGBA{0, 0, 0, 255}
	libtcodRed           = color.RGBA{255, 0, 0, 255}
	libtcodBlue          = color.RGBA{0, 0, 255, 255}
	libtcodLightBlue     = color.RGBA{128, 192, 255, 255}
	libtcodDarkestOrange = color.RGBA{128, 64, 0, 255}
	libtcodDarkerGreen   = color.RGBA{0, 128, 0, 255}
	libtcodLightGray     = color.RGBA{192, 192, 192, 255}
	libtcodDarkerGray    = color.RGBA{128, 128, 128, 255}
	libtcodDarkSepia     = color.RGBA{128, 96, 64, 255}
)

func NamegenGenerate(bla string) string {
	return bla + strconv.Itoa(rand.Intn(100))
}
