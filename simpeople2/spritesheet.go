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
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Grass         []*ebiten.Image
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
	img := s.Src.SubImage(image.Rect(x*s.TileSize, y*s.TileSize, (x+1)*s.TileSize, (y+1)*s.TileSize)).(*ebiten.Image)
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
func LoadSpriteSheet(tileSize, targetSize int, data []byte) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y, tileSize int, sheet *ebiten.Image) *ebiten.Image {
		return sheet.SubImage(image.Rect(x*tileSize, y*tileSize, (x+1)*tileSize, (y+1)*tileSize)).(*ebiten.Image)
	}

	// scaleSprite returns the sprite scaled to the provided size.
	scaleSprite := func(img *ebiten.Image) *ebiten.Image {
		newImage := ebiten.NewImage(targetSize, targetSize)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(targetSize)/float64(tileSize), float64(targetSize)/float64(tileSize))
		newImage.DrawImage(img, op)
		return newImage
	}

	// Populate SpriteSheet.
	s := &SpriteSheet{
		TileSize:      tileSize,
		TargetSize:    targetSize,
		Width:         img.Bounds().Dx() / tileSize,
		Height:        img.Bounds().Dy() / tileSize,
		Src:           sheet,
		subimageCache: make(map[int]*ebiten.Image),
	}
	s.Grass = []*ebiten.Image{
		scaleSprite(spriteAt(18, 8, tileSize, sheet)),
		scaleSprite(spriteAt(18, 9, tileSize, sheet)),
		scaleSprite(spriteAt(19, 8, tileSize, sheet)),
		scaleSprite(spriteAt(19, 9, tileSize, sheet)),
	}

	return s, nil
}
