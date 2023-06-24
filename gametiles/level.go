package gametiles

import "github.com/ojrac/opensimplex-go"

type Level struct {
	Width    int
	Height   int
	TileSize int
	Tiles    []Tile
}

func NewLevel(width, height int) (*Level, error) {
	l := &Level{
		Width:    width,
		Height:   height,
		TileSize: tileSize,
		Tiles:    make([]Tile, width*height),
	}

	noise := opensimplex.New(0)
	for y := 0; y < l.Height; y++ {
		for x := 0; x < l.Width; x++ {
			// Calculate the noise value for the current tile.
			tx, ty := l.TileXYToPixelPos(x, y)
			l.Tiles[y*l.Width+x] = Tile(noise.Eval2(float64(tx)/100, float64(ty)/100) * 255)
		}
	}

	return l, nil
}

// TileXYToPixelPos returns the center coordinates in pixel for the given tile in the grid.
func (l *Level) TileXYToPixelPos(x, y int) (int, int) {
	return x*l.TileSize + l.TileSize/2, y*l.TileSize + l.TileSize/2
}

type Tile byte

// TileType represents the type of a tile.
type TileType int

const (
	TileTypeGrass TileType = iota
	TileTypeDirt
	TileTypeRock
	TileTypeWater
	TileTypeSnow
	TileTypeTrees
)

func (t Tile) Type() TileType {
	switch {
	case t < 40:
		return TileTypeWater
	case t < 60:
		return TileTypeDirt
	case t < 120:
		return TileTypeGrass
	case t < 220:
		return TileTypeRock
	default:
		return TileTypeSnow
	}
}

func (t Tile) HasTrees() bool {
	return t > 50 && t < 80
}
