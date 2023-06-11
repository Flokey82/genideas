package gamehex

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Level struct {
	Width    int
	Height   int
	tileSize int
	Tiles    []Tile
}

func NewLevel(width, height int) (*Level, error) {
	return &Level{
		Width:    width,
		Height:   height,
		tileSize: 32,
		Tiles:    make([]Tile, width*height),
	}, nil
}

// HexTilePos returns the tile position for the given tile in the hex grid.
func (l *Level) HexTileXYToPixelPos(x, y int) (int, int) {
	//layoutHorizontal := false
	//if layoutHorizontal {
	// Every odd row is shifted by half a tile.
	px, py := x*l.tileSize*2+l.tileSize/2, (y+1)*l.tileSize/2
	if y%2 == 1 {
		px += l.tileSize
	}
	return px, py
	//} else {
	//	px, py := x*l.tileSize/2+l.tileSize/2, y*l.tileSize*2+l.tileSize/2
	//	if x%2 == 1 {
	//		py += l.tileSize
	//	}
	//	return px, py
	//}
}

// PosToHexTileXY returns the tile position for the given pixel position.
func (l *Level) PosToHexTileXY(px, py int) (int, int) {
	// Reverse HexTileXYToPixelPos.
	//layoutHorizontal := false
	//if layoutHorizontal {
	x, y := px/(l.tileSize*2), py/(l.tileSize/2)-1
	return x, y
	//} else {
	//	x, y := px/(l.tileSize/2), py/(l.tileSize*2)
	//	return x, y
	//}
}

func (l *Level) drawHex(background *ebiten.Image, xCenter float64, yCenter float64, innerRadius float64) {
	lineColor := color.RGBA{uint8(255), uint8(255), uint8(102), 50}
	//layoutHorizontal := false
	//if layoutHorizontal {
	// Draw a hexagon with the given center and inner radius.
	ebitenutil.DrawLine(background, xCenter+innerRadius, yCenter, xCenter+innerRadius/2, yCenter+innerRadius*math.Sqrt(3)/2, lineColor)
	ebitenutil.DrawLine(background, xCenter+innerRadius/2, yCenter+innerRadius*math.Sqrt(3)/2, xCenter-innerRadius/2, yCenter+innerRadius*math.Sqrt(3)/2, lineColor)
	ebitenutil.DrawLine(background, xCenter-innerRadius/2, yCenter+innerRadius*math.Sqrt(3)/2, xCenter-innerRadius, yCenter, lineColor)
	ebitenutil.DrawLine(background, xCenter-innerRadius, yCenter, xCenter-innerRadius/2, yCenter-innerRadius*math.Sqrt(3)/2, lineColor)
	ebitenutil.DrawLine(background, xCenter-innerRadius/2, yCenter-innerRadius*math.Sqrt(3)/2, xCenter+innerRadius/2, yCenter-innerRadius*math.Sqrt(3)/2, lineColor)
	ebitenutil.DrawLine(background, xCenter+innerRadius/2, yCenter-innerRadius*math.Sqrt(3)/2, xCenter+innerRadius, yCenter, lineColor)
	//} else {
	//	ebitenutil.DrawLine(background, xCenter, yCenter+innerRadius, xCenter+innerRadius*math.Sqrt(3)/2, yCenter+innerRadius/2, lineColor)
	//	ebitenutil.DrawLine(background, xCenter+innerRadius*math.Sqrt(3)/2, yCenter+innerRadius/2, xCenter+innerRadius*math.Sqrt(3)/2, yCenter-innerRadius/2, lineColor)
	//	ebitenutil.DrawLine(background, xCenter+innerRadius*math.Sqrt(3)/2, yCenter-innerRadius/2, xCenter, yCenter-innerRadius, lineColor)
	//	ebitenutil.DrawLine(background, xCenter, yCenter-innerRadius, xCenter-innerRadius*math.Sqrt(3)/2, yCenter-innerRadius/2, lineColor)
	//	ebitenutil.DrawLine(background, xCenter-innerRadius*math.Sqrt(3)/2, yCenter-innerRadius/2, xCenter-innerRadius*math.Sqrt(3)/2, yCenter+innerRadius/2, lineColor)
	//	ebitenutil.DrawLine(background, xCenter-innerRadius*math.Sqrt(3)/2, yCenter+innerRadius/2, xCenter, yCenter+innerRadius, lineColor)
	//}
}

type Tile byte
