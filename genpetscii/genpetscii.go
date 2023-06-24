// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	_ "embed"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

var (
	//go:embed vendor/bescii.ttf
	bescii_ttf []byte
)

// https://github.com/damianvila/font-bescii
const (
	screenWidth  = 640
	screenHeight = 520
)

// https://github.com/blackchip-org/retro-cs/blob/master/rcs/cbm/petscii/petscii_table.go
// https://style64.org/petscii/ direct PETSCII to Unicode mapping
const (
	vertical            = '│'
	verticalSlightLeft  = '\U0000e0c7'
	verticalSlightRight = '\U0000e0c8'

	horizontal           = '─'
	horizontalSlightUp   = '\U0000e0c4'
	horizontalSlightDown = '\U0000e0c6'

	treeTopA  = '\U0000f00f'
	treeTopB  = '\U0000f03a'
	treeTrunk = '\U0000f010'

	towerTop = '\U0000f037'

	roundedCorners  = `╭╮╰╯`
	pointyCorners   = `┌┐└┘`
	connectingLines = `─│┴├┬┤┼`
)

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	verticalTiles   = []rune{}
	horizontalTiles = []rune{}
)

func init() {
	verticalTiles = append(verticalTiles, vertical)
	verticalTiles = append(verticalTiles, verticalSlightLeft)  // Slight left vertical
	verticalTiles = append(verticalTiles, verticalSlightRight) // Slight right vertical

	horizontalTiles = append(horizontalTiles, horizontal)
	horizontalTiles = append(horizontalTiles, horizontalSlightUp)   // Slight up horizontal
	horizontalTiles = append(horizontalTiles, horizontalSlightDown) // Slight down horizontal
}

func init() {
	tt, err := opentype.Parse(bescii_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull, // Use quantization to save glyph cache images.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Adjust the line height.
	mplusBigFont = text.FaceWithLineHeight(mplusBigFont, 48)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	riverColor = color.RGBA{0x10, 0x10, 0xff, 0xff}
	roadColor  = color.RGBA{0xcc, 0xcc, 0xcc, 0xff}
	treeColor  = color.RGBA{0, 0xbb, 0, 0xff}
	trunkColor = color.RGBA{0x80, 0x40, 0, 0xff}
	towerColor = color.RGBA{0x80, 0x80, 0x80, 0xff}
)

type Game struct {
	counter       int
	kanjiText     string
	roadText      string
	treeText      string
	treeTrunkText string
	towerText     string
	mapHeight     int
	mapWidth      int
	riverMap      [10][10]bool
	roadMap       [10][10]bool
	treeMap       [10][10]bool
	treeTrunkMap  [10][10]bool
	towerMap      [10][10]bool
}

const (
	t = true
	o = false
)

func NewGame() (*Game, error) {
	g := &Game{
		mapHeight: 10,
		mapWidth:  10,
		riverMap: [10][10]bool{
			{t, o, o, o, o, o, o, t, o, o},
			{t, t, t, o, o, o, t, t, o, o},
			{o, o, t, t, t, o, t, o, o, o},
			{o, o, o, o, t, t, t, o, o, o},
			{o, o, o, o, o, o, t, t, t, o},
			{o, o, o, o, o, o, o, o, t, t},
			{o, o, o, o, o, o, o, o, o, t},
			{o, o, o, o, o, o, o, o, o, t},
			{o, o, o, o, o, o, o, o, o, t},
			{o, o, o, o, o, o, o, o, o, t},
		},
		roadMap: [10][10]bool{
			{o, o, o, o, t, o, o, o, o, o},
			{o, o, o, t, t, o, o, o, o, o},
			{o, o, o, t, o, o, o, o, o, o},
			{o, o, o, t, o, o, o, o, o, o},
			{o, o, o, t, t, t, o, o, o, o},
			{o, o, o, t, o, t, t, t, o, o},
			{o, o, o, t, t, o, o, t, t, o},
			{o, o, o, o, t, o, o, o, t, o},
			{o, o, o, t, t, o, o, o, t, o},
			{o, o, o, t, o, o, o, o, t, o},
		},
		treeMap: [10][10]bool{
			{o, o, o, o, o, o, o, o, t, t},
			{o, o, o, o, o, o, o, t, t, t},
			{o, o, o, o, o, o, o, o, t, t},
			{o, o, o, o, o, o, o, o, o, t},
			{t, o, o, o, o, o, o, o, o, o},
			{t, t, o, o, o, o, o, o, o, o},
			{t, t, t, o, o, t, t, o, o, o},
			{t, t, t, o, o, t, t, t, o, o},
			{t, t, o, o, o, o, t, t, o, o},
			{t, o, o, o, o, o, o, o, o, o},
		},
		treeTrunkMap: [10][10]bool{
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, t, o, o},
			{o, o, o, o, o, o, o, o, t, o},
			{o, o, o, o, o, o, o, o, o, t},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, t, o, o, t, o, o, o, o},
			{o, t, o, o, o, o, t, t, o, o},
		},
		towerMap: [10][10]bool{
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, t, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
			{o, o, o, o, o, o, o, o, o, o},
		},
	}
	g.init()
	return g, nil
}

func (g *Game) init() {
	/*
		// Random walk the river map, starting from 0, 0
		x := 0
		y := 0
		g.riverMap[y][x] = true
		for i := 0; i < 10; i++ {
			// Pick a random direction and check if we can move there.
			// If not, try again.
			for {
				dx := rand.Intn(3) - 1
				dy := rand.Intn(3) - 1
				if dx == 0 && dy == 0 {
					continue
				}
				if x+dx < 0 || x+dx >= g.mapWidth {
					continue
				}
				if y+dy < 0 || y+dy >= g.mapHeight {
					continue
				}
				if g.riverMap[y+dy][x+dx] {
					continue
				}
				x += dx
				y += dy
				g.riverMap[y][x] = true
				break
			}
			if x == g.mapWidth-1 && y == g.mapHeight-1 {
				break
			}
		}
	*/
}

func (g *Game) Update() error {
	getMapChar := func(x, y int, m [10][10]bool) rune {
		if x < 0 || x >= g.mapWidth {
			return ' '
		}
		if y < 0 || y >= g.mapHeight {
			return ' '
		}
		if !m[y][x] {
			return ' '
		}

		// Check which neighbors are set and pick the correct character.
		// We encode the set neighbors as a 4-bit number.
		// 0b0000 = no neighbors
		// 0b0001 = left neighbor
		// 0b0010 = right neighbor
		// 0b0100 = top neighbor
		// 0b1000 = bottom neighbor
		// 0b0101 = left and top neighbor
		// etc.
		var neighbors int
		if x > 0 && m[y][x-1] {
			neighbors |= 0b0001
		}
		if x < g.mapWidth-1 && m[y][x+1] {
			neighbors |= 0b0010
		}
		if y > 0 && m[y-1][x] {
			neighbors |= 0b0100
		}
		if y < g.mapHeight-1 && m[y+1][x] {
			neighbors |= 0b1000
		}

		// Pick the correct character.
		// roundedCorners  = `╭╮╰╯`
		// pointyCorners   = `┌┐└┘`
		// connectingLines = `─│┴├┬┤┼`
		var c rune
		switch neighbors {
		case 0b0000:
			c = ' '
		case 0b0001: // left
			c = horizontal
		case 0b0010: // right
			c = horizontal
		case 0b0011: // left and right
			c = horizontal
		case 0b0100: // top
			c = vertical
		case 0b0101: // left and top
			c = rune('╯')
		case 0b0110: // right and top
			c = rune('╰')
		case 0b0111: // left, right and top
			c = rune('┴')
		case 0b1000: // bottom
			c = vertical
		case 0b1001: // left and bottom
			c = rune('╮')
		case 0b1010: // right and bottom
			c = rune('╭')
		case 0b1011: // left, right and bottom
			c = rune('┬')
		case 0b1100: // top and bottom
			c = vertical
		case 0b1101: // left, top and bottom
			c = rune('┤')
		case 0b1110: // right, top and bottom
			c = rune('├')
		case 0b1111: // left, right, top and bottom
			c = rune('┼')
		}
		return c
	}

	// Change the text color for each second.
	if g.counter%ebiten.TPS() == 0 {
		g.kanjiText = ""
		// We will render the river map as a series of characters.
		for y := 0; y < g.mapHeight; y++ {
			for x := 0; x < g.mapWidth; x++ {
				g.kanjiText += string(getMapChar(x, y, g.riverMap))
			}
			g.kanjiText += "\n"
		}

		// Render roads.
		g.roadText = ""
		for y := 0; y < g.mapHeight; y++ {
			for x := 0; x < g.mapWidth; x++ {
				g.roadText += string(getMapChar(x, y, g.roadMap))
			}
			g.roadText += "\n"
		}

		// Render trees.
		g.treeText = ""
		for y := 0; y < g.mapHeight; y++ {
			for x := 0; x < g.mapWidth; x++ {
				if !g.treeMap[y][x] {
					g.treeText += string(' ')
					continue
				}
				g.treeText += string(treeTopA)
			}
			g.treeText += "\n"
		}

		// Render tree trunks.
		g.treeTrunkText = ""
		for y := 0; y < g.mapHeight; y++ {
			for x := 0; x < g.mapWidth; x++ {
				if !g.treeTrunkMap[y][x] {
					g.treeTrunkText += string(' ')
					continue
				}
				g.treeTrunkText += string(treeTrunk)
			}
			g.treeTrunkText += "\n"
		}

		// Render towers.
		g.towerText = ""
		for y := 0; y < g.mapHeight; y++ {
			for x := 0; x < g.mapWidth; x++ {
				if !g.towerMap[y][x] {
					g.towerText += string(' ')
					continue
				}
				g.towerText += string(towerTop)
			}
			g.towerText += "\n"
		}
	}
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const x = 20

	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 30, color.White)

	// Draw river text lines
	text.Draw(screen, g.kanjiText, mplusBigFont, x, 80, riverColor)

	// Draw road text lines
	text.Draw(screen, g.roadText, mplusBigFont, x, 80, roadColor)

	// Draw tower text lines
	text.Draw(screen, g.towerText, mplusBigFont, x, 80, towerColor)

	// Draw tree text lines
	text.Draw(screen, g.treeText, mplusBigFont, x, 80, treeColor)

	// Draw tree trunk text lines
	text.Draw(screen, g.treeTrunkText, mplusBigFont, x+5, 80, trunkColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("PETSCII (Ebitengine Demo)")
	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
