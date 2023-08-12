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

type Pos struct {
	X int
	Y int
}

type Game struct {
	rootView       *console.Console
	worldView      *console.Console
	playerInfoView *console.Console
	world          [][]byte
	player         Pos
	playerTarget   Pos
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
		player:         Pos{X: 3, Y: 3},
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
	{"magic missile", presetMagicMissile},
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
	// If clicked, update the player target.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		cx, cy := ebiten.CursorPosition()
		relX := cx / font.DefaultFont.TileWidth
		relY := cy / font.DefaultFont.TileHeight
		relY -= 1 // Subtract the header height.

		// Get the screen offset depending on the player position.
		midX := g.worldView.Width / 2
		midY := g.worldView.Height / 2

		// Calculate the target position.
		g.playerTarget.X = relX + g.player.X - midX
		g.playerTarget.Y = relY + g.player.Y - midY
	}

	// Move player
	if inpututil.IsKeyJustPressed(ebiten.KeyW) && !g.isSolid(g.player.X, g.player.Y-1) {
		g.player.Y -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) && !g.isSolid(g.player.X, g.player.Y+1) {
		g.player.Y += 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) && !g.isSolid(g.player.X-1, g.player.Y) {
		g.player.X -= 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && !g.isSolid(g.player.X+1, g.player.Y) {
		g.player.X += 1
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
		origin := vectors.Vec2{X: float64(g.player.X) + 0.5, Y: float64(g.player.Y) + 0.5}

		// Spawn particles.
		for i := 0; i < preset.NumParticles; i++ {
			curMaxAge := maxAge * (1 + (rand.Float64()*2-1)*preset.Variance)

			// Calculate speed vector (evenly distributed in a circle)
			curSpeedMag := mag * (1 + (rand.Float64()*2-1)*preset.Variance)
			var angle float64
			if preset.Pattern == PatternDirectional {
				// Directional pattern, angle is the angle between the origin and the target.
				angle = math.Atan2(float64(g.playerTarget.Y)+0.5-origin.Y, float64(g.playerTarget.X)+0.5-origin.X)
			} else {
				// Circle pattern, angle is evenly distributed.
				angle = float64(i) * math.Pi * 2 / float64(preset.NumParticles)
			}
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
			g.worldView.Transform(midX-g.player.X+x, midY-g.player.Y+y, t.CharByte(g.world[y][x]))
		}
	}

	// Draw particles.
	for _, part := range particles {
		drawTile := func(deltas [][2]int, progress float64) {
			// Get the color from the palette based on the progress
			// which we will use as the foreground color.
			partCol := part.preset.GetConColor(progress)

			// Get the character from the palette based on the progress.
			char := part.preset.GetChar(progress)

			// Get the color from the palette based on the progress
			// which we will use as the background color.
			backCol := part.preset.GetConColor(math.Pow(progress, 2))

			// Draw the particle at the given deltas.
			for _, delta := range deltas {
				// Check if the delta is within a solid tile. If so, skip it.
				if g.isSolid(int(part.pos.X)+delta[0], int(part.pos.Y)+delta[1]) {
					continue
				}
				// TODO: Color of the tile should depend on how close the particle is to the center of the tile.
				g.worldView.Transform(int(part.pos.X)+midX-g.player.X+delta[0], int(part.pos.Y)+midY-g.player.Y+delta[1], t.CharRune(char), t.Foreground(partCol), t.Background(backCol))
			}
		}

		// Draw the particle.
		drawTile([][2]int{{0, 0}}, 1-part.age/part.maxAge)

		// Draw the bloom.
		// TODO: Instead, bloom should be averaged by the surrounding tiles.
		for i := 0; i < part.preset.Bloom; i++ {
			progress := 1 - part.age/part.maxAge - float64(i)*0.1
			if progress < 0 {
				continue
			}

			// Draw the particle's direct neighbors as part of the bloom.
			drawTile([][2]int{{i, 0}, {-i, 0}, {0, i}, {0, -i}}, progress)

			// Draw the cross bloom.
			progress = 1 - part.age/part.maxAge - float64(i)*0.2
			if progress < 0 {
				continue
			}

			// Draw the bloom corners.
			drawTile([][2]int{{i, i}, {-i, i}, {i, -i}, {-i, -i}}, progress)
		}
	}

	// draw player in the middle
	g.worldView.Transform(midX, midY, t.CharByte('@'), t.Foreground(concolor.Green))

	// draw player info
	g.playerInfoView.PrintBounded(1, 1, g.playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", g.player.X, g.player.Y))
	g.playerInfoView.PrintBounded(1, 3, g.playerInfoView.Width-2, 4, "SPACE: effect")
	g.playerInfoView.PrintBounded(1, 4, g.playerInfoView.Width-2, 6, "TAB: next effect")
	g.playerInfoView.PrintBounded(1, 6, g.playerInfoView.Width-2, 8, fmt.Sprintf("FX: %s", effects[g.currentEffect].Name))
	g.playerInfoView.PrintBounded(1, 8, g.playerInfoView.Width-2, 10, "WASD: move")
	g.playerInfoView.PrintBounded(1, 9, g.playerInfoView.Width-2, 12, "LMB: set target")
	g.playerInfoView.PrintBounded(1, 11, g.playerInfoView.Width-2, 14, fmt.Sprintf("TGT: X=%d Y=%d", g.playerTarget.X, g.playerTarget.Y))

	return nil
}

type Pattern int

const (
	PatternCircle Pattern = iota
	PatternDirectional
)

type ParticlePreset struct {
	CharPalette  []rune       // Character palette
	ColorPalette []color.RGBA // RGBA color palette
	NumParticles int          // Number of particles to spawn
	SpeedMag     float64      // Speed magnitude
	MaxAge       float64      // Max age of a particle
	Variance     float64      // Variance of the speed magnitude and max age
	Bloom        int          // Bloom intensity of the particle
	Pattern      Pattern      // Pattern of the particles
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
	Bloom:        2,
	Pattern:      PatternCircle,
}

var presetIce = &ParticlePreset{
	CharPalette:  iceCharPalette,
	ColorPalette: iceColorPalette,
	NumParticles: 20,
	SpeedMag:     2,
	MaxAge:       2.5,
	Variance:     0,
	Bloom:        2,
	Pattern:      PatternCircle,
}

var presetMagicMissile = &ParticlePreset{
	CharPalette:  []rune{'|'},
	ColorPalette: magicMissileColorPalette,
	NumParticles: 10,
	SpeedMag:     20,
	MaxAge:       0.5,
	Variance:     0.2,
	Bloom:        1,
	Pattern:      PatternDirectional,
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

// laser palette :)
// Creates a nice laser effect ranging from white to orange to red.
var magicMissileColorPalette = []color.RGBA{
	{R: 255, G: 255, B: 255}, //  0 - white
	{R: 255, G: 255, B: 191}, //  1
	{R: 255, G: 255, B: 127}, //  2
	{R: 255, G: 255, B: 63},  //  3
	{R: 255, G: 255, B: 0},   //  4 - yellow
	{R: 255, G: 191, B: 0},   //  5
	{R: 255, G: 127, B: 0},   //  6
	{R: 255, G: 63, B: 0},    //  7 - orange
	{R: 255, G: 0, B: 0},     //  8 - red
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
