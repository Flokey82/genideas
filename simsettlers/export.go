package simsettlers

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/sizeofint/webpanimation"
)

type webpExport struct {
	anim     *webpanimation.WebpAnimation
	config   webpanimation.WebPConfig
	timeline int
	timestep int
}

func newWebPExport(width, height int) *webpExport {
	anim := webpanimation.NewWebpAnimation(width, height, 0)
	anim.WebPAnimEncoderOptions.SetKmin(9)
	anim.WebPAnimEncoderOptions.SetKmax(17)

	config := webpanimation.NewWebpConfig()
	config.SetLossless(1)
	return &webpExport{
		anim:     anim,
		config:   config,
		timeline: 0,
		timestep: 5,
	}
}

func (m *webpExport) ExportWebp(name string) error {
	// Write the final frame.
	m.timeline += m.timestep
	if err := m.anim.AddFrame(nil, m.timeline, m.config); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	// Encode animation and write result bytes in buffer.
	if err = m.anim.Encode(f); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	m.anim.ReleaseMemory() // TODO: This doesn't really prevent crashes?

	return nil
}

func (m *Map) storeWebPFrame() error {
	// Write the current map to the animation.
	// Create a colored image of the given width and height.
	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	m.drawFrame(img)
	if err := m.Export.anim.AddFrame(img, m.Export.timeline, m.Export.config); err != nil {
		return err
	}
	m.Export.timeline += m.Export.timestep
	return nil
}

// ExportPNG exports the map as a PNG file.
func (m *Map) ExportPNG(filename string) error {
	// We will draw the elevation as a grayscale image and
	// the flux above a certain threshold in blue.
	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	m.drawFrame(img)

	// Encode the image as PNG.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func (m *Map) drawFrame(img *image.RGBA) {
	fs := m.Elevation
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			// Draw the elevation as grayscale.
			// We will use the full range of grayscale values from 0 to 255.
			// The lowest point will be black, the highest point white.
			// We will use the elevation as a percentage of the full range.
			// This means that the lowest point will be black, the highest point white.
			// The lowest point will be black, the highest point white.
			img.Set(x, y, color.Gray{uint8(fs[x+y*m.Width] * 255)})
		}
	}

	// Draw the flux.
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			if m.Flux[x+y*m.Width] > fluxRiverThreshold {
				// Scale the flux to the range 0-255.
				fluxVal := uint8(m.Flux[x+y*m.Width] * 255)

				img.Set(x, y, color.RGBA{fluxVal, fluxVal, 255, 255})
			}
		}
	}

	// Draw the buildings.
	for _, b := range m.Buildings {
		img.Set(b.X, b.Y, color.RGBA{255, 0, 0, 255})
	}

	// Draw the construction sites.
	for _, b := range m.Construction {
		img.Set(b.X, b.Y, color.RGBA{255, 255, 0, 255})
	}

	// Draw the root building.
	img.Set(m.Root.X, m.Root.Y, color.RGBA{0, 255, 0, 255})

	// Draw the dungeon.
	for _, d := range m.Dungeons {
		img.Set(d.X, d.Y, color.RGBA{233, 128, 0, 255})
	}

	// Draw the people.
	for _, p := range m.RealPop {
		img.Set(int(p.X), int(p.Y), color.RGBA{0, 100, 255, 255})
	}
}
