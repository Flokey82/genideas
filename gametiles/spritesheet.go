package gametiles

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

var (
	//go:embed tiles/default_tiles.png
	Spritesheet_png []byte
)

// SpriteSheet represents a collection of sprite images.
type SpriteSheet struct {
	Grass *ebiten.Image
	Dirt  *ebiten.Image
	Water *ebiten.Image
	Rock  *ebiten.Image
	Snow  *ebiten.Image
}

// LoadSpriteSheet loads the embedded SpriteSheet.
func LoadSpriteSheet(tileSize, targetSize int) (*SpriteSheet, error) {
	img, _, err := image.Decode(bytes.NewReader(Spritesheet_png))
	if err != nil {
		return nil, err
	}

	sheet := ebiten.NewImageFromImage(img)

	// spriteAt returns a sprite at the provided coordinates.
	spriteAt := func(x, y int) *ebiten.Image {
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
	s := &SpriteSheet{}
	s.Grass = scaleSprite(spriteAt(6, 0))
	s.Dirt = scaleSprite(spriteAt(1, 0))
	s.Water = scaleSprite(spriteAt(8, 0))
	s.Rock = scaleSprite(spriteAt(7, 0))
	s.Snow = scaleSprite(spriteAt(5, 1))

	return s, nil
}
