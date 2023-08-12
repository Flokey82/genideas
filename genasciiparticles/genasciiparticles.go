package genasciiparticles

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"strings"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/Flokey82/go_gens/vectors"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

const (
	labelWindow     = "ramen - ASCII particles example"
	labelWorldView  = "World View"
	labelPlayerInfo = "Player Info"
)

var player = struct {
	X int
	Y int
}{3, 3}

type Game struct {
	rootView       *console.Console
	worldView      *console.Console
	playerInfoView *console.Console
	world          [][]byte
	currentEffect  int
}

func New() (*Game, error) {
	rootView, err := console.New(60, 35, font.DefaultFont, labelWindow)
	if err != nil {
		return nil, err
	}

	worldView, err := rootView.CreateSubConsole(0, 1, rootView.Width-20, rootView.Height-1)
	if err != nil {
		return nil, err
	}

	playerInfoView, err := rootView.CreateSubConsole(worldView.Width, 1, 20, rootView.Height-1)
	if err != nil {
		return nil, err
	}

	// converts levelLayout to world
	var world [][]byte
	lines := strings.Split(levelLayout, "\n")
	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}
		world = append(world, []byte(lines[i]))
	}
	return &Game{
		rootView:       rootView,
		worldView:      worldView,
		playerInfoView: playerInfoView,
		world:          world,
	}, nil
}

var levelLayout = `
#####################
#         #    #    #
#    #    #         #
#    ######    #    #
#              #    #
##  #############  ##
#    #    #    #    #
#    #         #    #
#    ######         #
#              #    #
#              #    #
#              #    #
#              #    #
#              #    #
#              #    #
#####################
`

type particle struct {
	age    float64
	maxAge float64
	pos    vectors.Vec2
	speed  vectors.Vec2
	preset *ParticlePreset
}

var particles []particle

type Effect struct {
	Name   string
	Preset *ParticlePreset
}

var effects = []Effect{
	{"fairy fire", presetFire},
	{"ice blast", presetIce},
}

// checks if a tile is solid (tile content is not a space ' ' character)
func (g *Game) isSolid(x int, y int) bool {
	if y < 0 || y >= len(g.world) {
		return true
	}
	if x < 0 || x >= len(g.world[y]) {
		return true
	}
	return g.world[y][x] != ' '
}

func (g *Game) Run() {
	g.rootView.SetTickHook(g.Tick)
	g.rootView.SetPreRenderHook(g.PreRender)
	g.rootView.Start(2)
}

func (g *Game) Tick(timeElapsed float64) error {
	// Move player
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && !g.isSolid(player.X, player.Y-1) {
		player.Y -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) && !g.isSolid(player.X, player.Y+1) {
		player.Y += 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) && !g.isSolid(player.X-1, player.Y) {
		player.X -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && !g.isSolid(player.X+1, player.Y) {
		player.X += 1
	}

	// Move particles and remove them if they are too old.
	// TODO: Add collision detection for particles and other behavior (like bouncing of walls).
	for i := range particles {
		particles[i].age += timeElapsed
		particles[i].pos = particles[i].pos.Add(particles[i].speed.Mul(timeElapsed))
	}

	// Remove particles that are too old or collide with a solid tile.
	for i := len(particles) - 1; i >= 0; i-- {
		if particles[i].age > particles[i].maxAge || g.isSolid(int(particles[i].pos.X), int(particles[i].pos.Y)) {
			particles = append(particles[:i], particles[i+1:]...)
		}
	}

	// Switch effect.
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.currentEffect = (g.currentEffect + 1) % len(effects)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Get the current effect and its preset.
		effect := effects[g.currentEffect]
		preset := effect.Preset

		// Get base values from the preset.
		mag := preset.SpeedMag
		maxAge := preset.MaxAge

		// Calculate the origin of the effect.
		origin := vectors.Vec2{X: float64(player.X) + 0.5, Y: float64(player.Y) + 0.5}

		// Spawn particles.
		for i := 0; i < preset.NumParticles; i++ {
			curMaxAge := maxAge * (1 + (rand.Float64()*2-1)*preset.Variance)

			// Calculate speed vector (evenly distributed in a circle)
			curSpeedMag := mag * (1 + (rand.Float64()*2-1)*preset.Variance)
			angle := float64(i) * math.Pi * 2 / float64(preset.NumParticles)
			speedVec := vectors.Vec2{X: math.Cos(angle), Y: math.Sin(angle)}.Mul(curSpeedMag)

			// Add particle
			particles = append(particles, particle{
				age:    0,
				maxAge: curMaxAge,
				pos:    origin,
				speed:  speedVec,
				preset: preset,
			})
		}
	}

	return nil
}

func (g *Game) PreRender(screen *ebiten.Image, timeDelta float64) error {
	// clear console
	g.rootView.ClearAll()
	g.rootView.TransformAll(t.Background(concolor.RGB(50, 50, 50)))

	g.worldView.ClearAll()
	g.worldView.TransformAll(t.Background(concolor.RGB(55, 55, 55)), t.Char(0))

	g.playerInfoView.ClearAll()

	// draw header
	g.rootView.TransformArea(0, 0, g.rootView.Width, 1, t.Background(concolor.RGB(80, 80, 80)))
	g.rootView.Print(2, 0, labelWorldView, t.Foreground(concolor.White))
	g.rootView.Print(g.worldView.Width+2, 0, labelPlayerInfo, t.Foreground(concolor.White))

	// draw world
	midX := g.worldView.Width / 2
	midY := g.worldView.Height / 2

	for y := range g.world {
		for x := range g.world[y] {
			if g.world[y][x] == ' ' {
				continue
			}

			g.worldView.Transform(midX-player.X+x, midY-player.Y+y, t.CharByte(g.world[y][x]))
		}
	}

	// Draw particles.
	for _, part := range particles {
		// Get the color from the palette based on the progress
		// which we will use as the foreground color.
		partCol := part.preset.GetConColor(1 - part.age/part.maxAge)

		// Get the character from the palette based on the progress.
		char := part.preset.GetChar(1 - part.age/part.maxAge)

		// Get the color from the palette based on the progress
		// which we will use as the background color.
		backCol := part.preset.GetConColor(math.Pow(1-part.age/part.maxAge, 2))

		// Draw the particle.
		g.worldView.Transform(int(part.pos.X)+midX-player.X, int(part.pos.Y)+midY-player.Y, t.CharRune(char), t.Foreground(partCol), t.Background(backCol))
	}

	// draw player in the middle
	g.worldView.Transform(midX, midY, t.CharByte('@'), t.Foreground(concolor.Green))

	// draw player info
	g.playerInfoView.PrintBounded(1, 1, g.playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", player.X, player.Y))
	g.playerInfoView.PrintBounded(1, 3, g.playerInfoView.Width-2, 4, "SPACE: effect")
	g.playerInfoView.PrintBounded(1, 4, g.playerInfoView.Width-2, 6, "TAB: next effect")
	g.playerInfoView.PrintBounded(1, 6, g.playerInfoView.Width-2, 8, fmt.Sprintf("Effect: %s", effects[g.currentEffect].Name))

	return nil
}

type ParticlePreset struct {
	CharPalette  []rune       // Character palette
	ColorPalette []color.RGBA // RGBA color palette
	NumParticles int          // Number of particles to spawn
	SpeedMag     float64      // Speed magnitude
	MaxAge       float64      // Max age of a particle
	Variance     float64      // Variance of the speed magnitude and max age
}

// GetChar returns a character from the char palette based on the progress (0.0 - 1.0).
func (p *ParticlePreset) GetChar(progress float64) rune {
	return p.CharPalette[int(progress*float64(len(p.CharPalette)))%len(p.CharPalette)]
}

// GetColor returns a color from the fire palette based on the progress (0.0 - 1.0).
func (p *ParticlePreset) GetColor(progress float64) color.RGBA {
	return p.ColorPalette[int(progress*float64(len(p.ColorPalette)))%len(p.ColorPalette)]
}

// GetConColor returns a concolor from the fire palette based on the progress (0.0 - 1.0).
func (p *ParticlePreset) GetConColor(progress float64) concolor.Color {
	col := p.GetColor(progress)
	return concolor.RGB(col.R, col.G, col.B)
}

var presetFire = &ParticlePreset{
	CharPalette:  fireCharPalette,
	ColorPalette: fireColorPalette,
	NumParticles: 15,
	SpeedMag:     4,
	MaxAge:       1.5,
	Variance:     0.5,
}

var presetIce = &ParticlePreset{
	CharPalette:  iceCharPalette,
	ColorPalette: iceColorPalette,
	NumParticles: 20,
	SpeedMag:     2,
	MaxAge:       2.5,
	Variance:     0,
}

// doom fire palette :)
// Creates a nice fire effect that starts with white (36), goes to yellow (26), orange (17) then red (11) and finally black (0).
var fireColorPalette = []color.RGBA{
	{R: 7, G: 7, B: 7},       //  0 - black
	{R: 31, G: 7, B: 7},      //  1
	{R: 47, G: 15, B: 7},     //  2
	{R: 71, G: 15, B: 7},     //  3
	{R: 87, G: 23, B: 7},     //  4
	{R: 103, G: 31, B: 7},    //  5
	{R: 119, G: 31, B: 7},    //  6
	{R: 143, G: 39, B: 7},    //  7
	{R: 159, G: 47, B: 7},    //  8
	{R: 175, G: 63, B: 7},    //  9
	{R: 191, G: 71, B: 7},    // 10
	{R: 199, G: 71, B: 7},    // 11 - red
	{R: 223, G: 79, B: 7},    // 12
	{R: 223, G: 87, B: 7},    // 13
	{R: 223, G: 87, B: 7},    // 14
	{R: 215, G: 95, B: 7},    // 15
	{R: 215, G: 95, B: 7},    // 16
	{R: 215, G: 103, B: 15},  // 17 - orange
	{R: 207, G: 111, B: 15},  // 18
	{R: 207, G: 119, B: 15},  // 19
	{R: 207, G: 127, B: 15},  // 20
	{R: 207, G: 135, B: 23},  // 21
	{R: 199, G: 135, B: 23},  // 22
	{R: 199, G: 143, B: 23},  // 23
	{R: 199, G: 151, B: 31},  // 24
	{R: 191, G: 159, B: 31},  // 25
	{R: 191, G: 159, B: 31},  // 26 - yellow
	{R: 191, G: 167, B: 39},  // 27
	{R: 191, G: 167, B: 39},  // 28
	{R: 191, G: 175, B: 47},  // 29
	{R: 183, G: 175, B: 47},  // 30
	{R: 183, G: 183, B: 47},  // 31
	{R: 183, G: 183, B: 55},  // 32
	{R: 207, G: 207, B: 111}, // 33
	{R: 223, G: 223, B: 159}, // 34
	{R: 239, G: 239, B: 199}, // 35
	{R: 255, G: 255, B: 255}, // 36 - white
}

// ice blast palette :)
// Creates a nice ice / frost effect ranging from white (21) to cyan (14), blue (7) and black (0).
var iceColorPalette = []color.RGBA{
	{R: 0, G: 0, B: 0},       //  0 - black
	{R: 0, G: 0, B: 36},      //  1
	{R: 0, G: 0, B: 72},      //  2
	{R: 0, G: 0, B: 109},     //  3
	{R: 0, G: 0, B: 145},     //  4
	{R: 0, G: 0, B: 182},     //  5
	{R: 0, G: 0, B: 218},     //  6
	{R: 0, G: 0, B: 255},     //  7 - blue
	{R: 0, G: 36, B: 255},    //  8
	{R: 0, G: 72, B: 255},    //  9
	{R: 0, G: 109, B: 255},   // 10
	{R: 0, G: 145, B: 255},   // 11
	{R: 0, G: 182, B: 255},   // 12
	{R: 0, G: 218, B: 255},   // 13
	{R: 0, G: 255, B: 255},   // 14 - cyan
	{R: 36, G: 255, B: 255},  // 15
	{R: 72, G: 255, B: 255},  // 16
	{R: 109, G: 255, B: 255}, // 17
	{R: 145, G: 255, B: 255}, // 18
	{R: 182, G: 255, B: 255}, // 19
	{R: 218, G: 255, B: 255}, // 20
	{R: 255, G: 255, B: 255}, // 21 - white
}

// char palette :)
var fireCharPalette = []rune{
	' ', //  0
	'.', //  1
	':', //  2
	'-', //  3
	'=', //  4
	'+', //  5
	'/', //  6
	't', //  7
	'z', //  8
	'U', //  9
	'w', // 10
	'*', // 11
	'o', // 12
	'O', // 13
	'#', // 14
	'@', // 15
}

// char palette :)
var iceCharPalette = []rune{
	' ',  //  0
	'.',  //  1
	':',  //  2
	'-',  //  3
	'|',  //  4
	'+',  //  5
	'x',  //  6
	'/',  //  7
	'-',  //  8
	'\\', //  9
	'X',  //  10
	'=',  //  11
	'H',  //  12
	'*',  //  13
	'#',  //  14
}
