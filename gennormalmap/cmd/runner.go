package main

import (
	"image"
	"image/png"
	"os"

	"github.com/Flokey82/genideas/gennormalmap"
)

func main() {
	// Open heightmap.
	f, err := os.Open("test.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	// Convert to gray scale.
	// TODO: Skip this step?
	bounds := img.Bounds()
	maxX := bounds.Max.X
	maxY := bounds.Max.Y
	grayImg := image.NewGray(bounds)
	for y := 0; y < maxY; y++ {
		for x := 0; x < maxX; x++ {
			grayImg.Set(x, y, img.At(x, y))
		}
	}

	outImg := gennormalmap.MapNormals(grayImg)

	// Save output.
	out, err := os.Create("test_out.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(out, outImg); err != nil {
		panic(err)
	}

	out.Close()
}
