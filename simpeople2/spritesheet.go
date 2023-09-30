package simpeople2

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

var (
	//go:embed vendor/tiles.png
	Spritesheet_png []byte
	//go:embed vendor/kenney/roguelikeSheet_transparent.png
	SpritesheetKenney_png []byte
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Grass         []*ebiten.Image
	Spacing       int // Spacing between tiles in pixels
	TileSize      int // Size of a tile in pixels (original size)
	TargetSize    int // Size of a tile in pixels (scaled)
	Width         int // Width of the spritesheet in tiles
	Height        int // Height of the spritesheet in tiles
	Src           *ebiten.Image
	subimageCache map[int]*ebiten.Image
}

func (s *SpriteSheet) GetSubImageXY(x, y int) *ebiten.Image {
	if img, ok := s.subimageCache[x+y*s.Width]; ok {
		return img
	}
	img := s.Src.SubImage(image.Rect(
		x*(s.TileSize+s.Spacing),
		y*(s.TileSize+s.Spacing),
		x*(s.TileSize+s.Spacing)+s.TileSize,
		y*(s.TileSize+s.Spacing)+s.TileSize,
	)).(*ebiten.Image)
	newImage := ebiten.NewImage(s.TargetSize, s.TargetSize)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(s.TargetSize)/float64(s.TileSize), float64(s.TargetSize)/float64(s.TileSize))
	newImage.DrawImage(img, op)
	s.subimageCache[x+y*s.Width] = newImage
	return newImage
}

func (s *SpriteSheet) GetSubImageID(id int) *ebiten.Image {
	if img, ok := s.subimageCache[id]; ok {
		return img
	}
	x := id % s.Width
	y := id / s.Width
	img := s.GetSubImageXY(x, y)
	s.subimageCache[id] = img
	return img
}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(tileSize, targetSize, spacing int, data []byte) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &SpriteSheet{
		Spacing:       spacing,
		TileSize:      tileSize,
		TargetSize:    targetSize,
		Width:         (img.Bounds().Dx() + spacing) / (tileSize + spacing),
		Height:        (img.Bounds().Dy() + spacing) / (tileSize + spacing),
		Src:           ebiten.NewImageFromImage(img),
		subimageCache: make(map[int]*ebiten.Image),
	}, nil
}
