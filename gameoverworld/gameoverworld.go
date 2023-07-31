package gameoverworld

import (
	"fmt"
	"math"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/console"
	"github.com/BigJk/ramen/font"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/ojrac/opensimplex-go"
)

const (
	labelWindow     = "game-overworld"
	labelWorldView  = "World View"
	labelPlayerInfo = "Player Info"
)

var player = struct {
	X int
	Y int
}{3, 3}

type Game struct {
	rootView         *console.Console
	worldView        *console.Console
	playerInfoView   *console.Console
	noise            *Noise
	currentElevation int
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

	return &Game{
		rootView:       rootView,
		worldView:      worldView,
		playerInfoView: playerInfoView,
		noise:          NewNoise(5, 0.5, 1234),
	}, nil
}

func (g *Game) Tick(timeElapsed float64) error {
	// Move player
	playerMoved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		player.Y -= 1
		playerMoved = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		player.Y += 1
		playerMoved = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		player.X -= 1
		playerMoved = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		player.X += 1
		playerMoved = true
	}

	// Update current elevation based on player position
	if playerMoved {
		g.currentElevation = g.ElevationLevelAt(player.X, player.Y)
	}

	// Switch between elevation levels
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.currentElevation++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.currentElevation--
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

	minX := player.X - midX
	minY := player.Y - midY
	maxX := player.X + midX
	maxY := player.Y + midY

	verticalRange := 2 // Range per elevation level

	// draw world
	for y := minY; y < maxY; y++ {
		for x := minX; x < maxX; x++ {
			elevation := g.ElevationLevelAt(x, y)
			if elevation < g.currentElevation {
				var con concolor.Color
				if g.currentElevation-elevation <= verticalRange {
					con = concolor.RGB(0, 0, 0)
				} else {
					con = concolor.RGB(55, 55, 55)
				}
				g.worldView.Transform(midX-player.X+x, midY-player.Y+y, t.CharByte('-'), t.Foreground(con))
				continue
			}
			if elevation >= g.currentElevation+verticalRange {
				var con concolor.Color
				if elevation-g.currentElevation <= verticalRange {
					con = concolor.RGB(255, 255, 255)
				} else {
					con = concolor.RGB(128, 128, 128)
				}
				g.worldView.Transform(midX-player.X+x, midY-player.Y+y, t.CharByte('#'), t.Foreground(con))
				continue
			}
			var con concolor.Color
			if elevation < 0 {
				con = concolor.RGB(0, 191, 255)
			} else {
				con = concolor.RGB(34, 139, 34)
			}
			g.worldView.Transform(midX-player.X+x, midY-player.Y+y, t.CharByte('.'), t.Foreground(con))
		}
	}

	// draw player in the middle if we are on the same elevation
	if g.currentElevation == g.ElevationLevelAt(player.X, player.Y) {
		g.worldView.Transform(midX, midY, t.CharByte('@'), t.Foreground(concolor.Green))
	}

	// draw player info
	g.playerInfoView.PrintBounded(1, 1, g.playerInfoView.Width-2, 2, fmt.Sprintf("X=%d Y=%d", player.X, player.Y))
	g.playerInfoView.PrintBounded(1, 3, g.playerInfoView.Width-2, 2, fmt.Sprintf("Elevation=%d", g.currentElevation))

	for i := 0; i < 20; i++ {
		if 10-i == g.currentElevation {
			g.playerInfoView.Transform(g.playerInfoView.Width-1, i+1, t.CharByte('#'), t.Background(concolor.RGB(255, 255, 255)))
		} else {
			g.playerInfoView.Transform(g.playerInfoView.Width-1, i+1, t.CharByte('|'), t.Background(concolor.RGB(128, 128, 128)))
		}
	}
	return nil
}

func (g *Game) ElevationLevelAt(x, y int) int {
	return int(g.noise.Eval2(float64(x)/20, float64(y)/20)*20 - 10)
}

func (g *Game) Run() {
	g.rootView.SetTickHook(g.Tick)
	g.rootView.SetPreRenderHook(g.PreRender)
	g.rootView.Start(2)
}

// Noise is a wrapper for opensimplex.Noise, initialized with
// a given seed, persistence, and number of octaves.
type Noise struct {
	Octaves     int
	Persistence float64
	Amplitudes  []float64
	Seed        int64
	OS          opensimplex.Noise
}

// NewNoise returns a new Noise.
func NewNoise(octaves int, persistence float64, seed int64) *Noise {
	n := &Noise{
		Octaves:     octaves,
		Persistence: persistence,
		Amplitudes:  make([]float64, octaves),
		Seed:        seed,
		OS:          opensimplex.NewNormalized(seed),
	}

	// Initialize the amplitudes.
	for i := range n.Amplitudes {
		n.Amplitudes[i] = math.Pow(persistence, float64(i))
	}

	return n
}

// Eval3 returns the noise value at the given point.
func (n *Noise) Eval3(x, y, z float64) float64 {
	var sum, sumOfAmplitudes float64
	for octave := 0; octave < n.Octaves; octave++ {
		frequency := 1 << octave
		fFreq := float64(frequency)
		sum += n.Amplitudes[octave] * n.OS.Eval3(x*fFreq, y*fFreq, z*fFreq)
		sumOfAmplitudes += n.Amplitudes[octave]
	}
	return sum / sumOfAmplitudes
}

// Eval2 returns the noise value at the given point.
func (n *Noise) Eval2(x, y float64) float64 {
	var sum, sumOfAmplitudes float64
	for octave := 0; octave < n.Octaves; octave++ {
		frequency := 1 << octave
		fFreq := float64(frequency)
		sum += n.Amplitudes[octave] * n.OS.Eval2(x*fFreq, y*fFreq)
		sumOfAmplitudes += n.Amplitudes[octave]
	}
	return sum / sumOfAmplitudes
}

// PlusOneOctave returns a new Noise with one more octave.
func (n *Noise) PlusOneOctave() *Noise {
	return NewNoise(n.Octaves+1, n.Persistence, n.Seed)
}
