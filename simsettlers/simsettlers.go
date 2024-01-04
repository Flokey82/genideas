// Package simsettlers contains the simulation of settlers settling in / founding a new village.
// This simulation should include the following:
// - Picking a suitable location based on multiple factors
// - Harvesting resources if needed
// - Building new buildings
// - Expanding individual plots (farm plots, buildings, etc.)
// - Demolishing buildings
// - Abandoning buildings / re-using abandoned buildings
// - Growing the village
package simsettlers

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/Flokey82/go_gens/genheightmap"
	"github.com/Flokey82/go_gens/vectors"
	"github.com/s0rg/fantasyname"
)

// Map represents the world map, which is a simple grid of tiles.
type Map struct {
	Day          uint16
	Year         int
	Height       int
	Width        int
	Elevation    []float64
	Flux         []float64
	TileType     []int
	Dungeons     []*Building // The dungeons.
	Root         *Building   // The root building, which the settlers will build around.
	Cemetery     *Building   // The cemetery.
	Buildings    []*Building
	Construction []*Building
	Resources    int
	Population   int
	RealPop      []*Person
	firstGen     [2]fmt.Stringer // First name generators (male/female).
	lastGen      fmt.Stringer    // Last name generators.
}

// first name prefixes for fantasyname generator.
const firstNamePrefix = "!(bil|bal|ban|hil|ham|hal|hol|hob|wil|me|or|ol|od|gor|for|fos|tol|ar|fin|ere|leo|vi|bi|bren|thor)"

// NewMap creates a new map with the given height and width.
func NewMap(height, width int) *Map {
	m := &Map{
		Height:     height,
		Width:      width,
		Elevation:  make([]float64, height*width),
		Flux:       make([]float64, height*width),
		TileType:   make([]int, height*width),
		Resources:  100,
		Population: 15,
		Cemetery:   NewBuilding(0, 0, BuildingTypeCemetery),
	}

	// Initialize name generation.

	// Female first names.
	genFirstF, err := fantasyname.Compile(firstNamePrefix+"(|ga|orbise|apola|adure|mosi|ri|i|na|olea|ne)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	m.firstGen[0] = genFirstF

	// Male first names.
	genFirstM, err := fantasyname.Compile(firstNamePrefix+"(|go|orbis|apol|adur|mos|ole|n)", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	m.firstGen[1] = genFirstM

	// Last names.
	genLast, err := fantasyname.Compile("!BsVc", fantasyname.Collapse(true), fantasyname.RandFn(rand.Intn))
	if err != nil {
		log.Fatal(err)
	}
	m.lastGen = genLast

	m.addNRandomPeople(m.Population)

	// Normalize the elevation.
	m.genElevation()

	// Calculate the flux.
	m.calcFlux()

	return m
}

func (m *Map) genElevation() {
	// Generate a slope.
	genSlope := genheightmap.GenSlope(vectors.Vec2{X: float64(-m.Width), Y: float64(-m.Height)})
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Elevation[x+y*m.Width] = 0.0001 * genSlope(float64(x), float64(y))
		}
	}

	// Generate two mountain ranges that form a diagonal valley.
	genMountain1 := genheightmap.GenMountainRange(
		vectors.Vec2{X: float64(m.Width) * 3 / 3, Y: float64(m.Height) * 1 / 3},
		vectors.Vec2{X: 0, Y: float64(m.Height) * 3 / 3},
		15, 30.0, 0.1, 5.0, true)
	genMountain2 := genheightmap.GenMountainRange(
		vectors.Vec2{X: float64(m.Width) * 2 / 3, Y: 0},
		vectors.Vec2{X: 0, Y: float64(m.Height) * 2 / 3},
		15, 30.0, 0.1, 5.0, true)

	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Elevation[x+y*m.Width] += genMountain1(float64(x), float64(y)) + genMountain2(float64(x), float64(y))
		}
	}

	// Generate a random noise map.
	genNoise := genheightmap.GenNoise(0, 10.4)
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			m.Elevation[x+y*m.Width] += 4.0 * genNoise(float64(x)*0.005, float64(y)*0.005)
		}
	}

	// Generate the elevation.
	normalize(m.Elevation)

	// Fill sinks.
	m.fillSinks()
	normalize(m.Elevation)

	down := m.calcDownhill()
	// check if there are any sinks left
	for i := range m.Elevation {
		if down[i] == -1 {
			log.Printf("Sink at %d,%d", i%m.Width, i/m.Width)
		}
	}
}

func (m *Map) fillSinks() {
	// Planchon-Darboux algorithm for filling sinks
	const epsilon = 1e-5
	const infinity = 999999
	newh := make([]float64, len(m.Elevation))
	for i := range m.Elevation {
		if m.isNearEdge(i) {
			newh[i] = m.Elevation[i]
		} else {
			newh[i] = infinity
		}
	}
	for {
		changed := false
		for i := range m.Elevation {
			if newh[i] == m.Elevation[i] {
				continue
			}
			nbs := m.Neighbors(i%m.Width, i/m.Width)
			for _, j := range nbs {
				if m.Elevation[i] >= newh[j]+epsilon {
					newh[i] = m.Elevation[i]
					changed = true
					break
				}
				oh := newh[j] + epsilon
				if newh[i] > oh && oh > m.Elevation[i] {
					newh[i] = oh
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}
	m.Elevation = newh
}

func (m *Map) isNearEdge(i int) bool {
	x := i % m.Width
	y := i / m.Width
	return x < 2 || x > m.Width-3 || y < 2 || y > m.Height-3
}

func (m *Map) calcDownhill() []int {
	downhill := make([]int, len(m.Flux))
	for i := range m.Flux {
		downhill[i] = -1
	}

	// Populate the downhill array, which will contain the lowest neighbor for each point (or -1 if there is none).
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			// Find the lowest neighbor.
			lowest := -1
			lowestElevation := m.Elevation[x+y*m.Width]

			// Avoid artifacts at the edges of the map, so don't go downhill there.
			if !m.isNearEdge(x + y*m.Width) {
				for _, neighbor := range m.Neighbors(x, y) {
					if m.Elevation[neighbor] < lowestElevation {
						lowest = neighbor
						lowestElevation = m.Elevation[neighbor]
					}
				}
			}
			downhill[x+y*m.Width] = lowest
		}
	}
	return downhill
}

func (m *Map) calcFlux() {
	// Initialize all flux values to 1/numpoints.
	byElevation := make([]int, len(m.Flux)) // This will contain the indices of the points, sorted by elevation.
	for i := range m.Flux {
		m.Flux[i] = 1.0 / float64(len(m.Flux))
		byElevation[i] = i
	}

	// Sort the points by elevation.
	sort.Slice(byElevation, func(i, j int) bool {
		return m.Elevation[byElevation[i]] > m.Elevation[byElevation[j]]
	})

	// Calculate the flux.
	downhill := m.calcDownhill()
	for _, i := range byElevation {
		if downhill[i] != -1 {
			m.Flux[downhill[i]] += m.Flux[i]
		}
	}

	// Draw a max-flux line through the map.
	// This is for debugging purposes only and will create a "river" through the map.
	/*
		y := m.Height / 2
		for x := 0; x < m.Width; x++ {
			m.Flux[x+y*m.Width] = 1.0
		}
	*/
}

func (m *Map) calcFitnessScore() []float64 {
	// Calculate the fitness score for each point.
	fitness := make([]float64, len(m.Flux))
	for i := range fitness {
		// We want the lowest flux in the cell, but the highest flux in the neighbors.
		if m.Flux[i] > fluxRiverThreshold {
			// This cell is not suitable.
			fitness[i] = math.Inf(-1)
			continue
		}
		fitness[i] = -m.Flux[i]
		for _, neighbor := range m.Neighbors(i%m.Width, i/m.Width) {
			fitness[i] += m.Flux[neighbor]
		}

		// Add distance to the border of the map.
		// NOTE: Yuk, this is a hack. We should use a proper distance function.
		x := i % m.Width
		y := i / m.Width
		x -= m.Width / 2
		y -= m.Height / 2

		// Now invert the distance, so that the center of the map has the highest score.
		fitness[i] += 1.0 / (1.0 + math.Sqrt(float64(x*x+y*y)))
	}
	return fitness
}

func normalize(values []float64) {
	min := 0.0
	max := 0.0
	for _, e := range values {
		if e < min {
			min = e
		}
		if e > max {
			max = e
		}
	}
	if min == max {
		return
	}
	for i := range values {
		values[i] = (values[i] - min) / (max - min)
	}
}

const fluxRiverThreshold = 0.01

// ExportPNG exports the map as a PNG file.
func (m *Map) ExportPNG(filename string) error {
	// We will draw the elevation as a grayscale image and
	// the flux above a certain threshold in blue.
	img := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))

	// Calculate the fitness score for each point.
	fs := m.calcFitnessScore()
	fs = m.calcFitnessScoreHouse(false)
	// normalize(fs)
	fs = m.Elevation
	//fs = m.calcFitnessScoreDungeon()
	//fs = m.Flux
	//
	normalize(m.Flux)

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

	// Encode the image as PNG.
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// Neighbors returns the neighbors of the given point.
func (m *Map) Neighbors(x, y int) []int {
	var neighbors []int
	if x > 0 {
		neighbors = append(neighbors, x-1+y*m.Width)
	}
	if x < m.Width-1 {
		neighbors = append(neighbors, x+1+y*m.Width)
	}
	if y > 0 {
		neighbors = append(neighbors, x+(y-1)*m.Width)
	}
	if y < m.Height-1 {
		neighbors = append(neighbors, x+(y+1)*m.Width)
	}
	allowCross := true
	if allowCross {
		if x > 0 && y > 0 {
			neighbors = append(neighbors, x-1+(y-1)*m.Width)
		}
		if x < m.Width-1 && y > 0 {
			neighbors = append(neighbors, x+1+(y-1)*m.Width)
		}
		if x > 0 && y < m.Height-1 {
			neighbors = append(neighbors, x-1+(y+1)*m.Width)
		}
		if x < m.Width-1 && y < m.Height-1 {
			neighbors = append(neighbors, x+1+(y+1)*m.Width)
		}
	}
	return neighbors
}

// Settle picks a suitable location for the settlers to settle and builds the first building.
func (m *Map) Settle() {
	// Calculate the fitness score for each point.
	fs := m.calcFitnessScore()

	// Find the best point.
	best := 0
	for i := range fs {
		if fs[i] > fs[best] {
			best = i
		}
	}

	// Build the first building.
	m.Root = m.AddBuilding(best%m.Width, best/m.Width, BuildingTypeMarket)

	// Add the dungeon.
	// Find the best spot for the dungeon.
	const numDungeons = 4
	for i := 0; i < numDungeons; i++ {
		fs = m.calcFitnessScoreDungeon()
		best = 0
		normalize(fs)
		for i := range fs {
			if fs[i] > fs[best] {
				best = i
			}
		}
		x := best % m.Width
		y := best / m.Width
		m.Dungeons = append(m.Dungeons, m.AddBuilding(x, y, BuildingTypeDungeon))
	}
}

// Tick advances the simulation by one tick.
func (m *Map) Tick() {
	m.Day++
	if m.Day > 365 {
		m.Day = 1
		m.Year++
	}

	// Age the population.
	m.agePop()

	// Get all yields for this tick.
	// TODO: Only calculate yields for buildings that are
	// inhabited.
	var yields int
	for _, b := range m.Buildings {
		yields += b.Yield()
	}

	// Add the yields to the resources.
	m.Resources += yields

	// Advance unoccupied building decay.
	m.tickBuildings()

	// Advance building construction.
	m.advanceConstruction()

	// Construct more houses if needed.
	m.tickPeople()
	// m.constructMoreHouses()

	// Match singles.
	m.matchSingles()

	// Advance pregnancies.
	m.advancePregnancies()
}
