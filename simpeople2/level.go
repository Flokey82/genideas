package simpeople2

import (
	"math/rand"

	"github.com/Flokey82/go_gens/vectors"
)

type Level struct {
	Width   int
	Height  int
	Ground  []int
	Tiles   []int
	Objects []*Object
}

func NewLevel(width, height int) *Level {
	l := &Level{
		Width:  width,
		Height: height,
		Ground: make([]int, width*height),
		Tiles:  make([]int, width*height),
	}
	for i := 0; i < width*height; i++ {
		l.Ground[i] = 915
	}

	// Draw a border around the level.

	// Corners.
	l.Tiles[0] = 721                        // Top left.
	l.Tiles[width-1] = 722                  // Top right.
	l.Tiles[(height-1)*width] = 778         // Bottom left.
	l.Tiles[(height-1)*width+width-1] = 779 // Bottom right.

	// Top and bottom.
	for i := 1; i < width-1; i++ {
		l.Tiles[i] = 719
		l.Tiles[(height-1)*width+i] = 719
	}

	// Left and right.
	for i := 1; i < height-1; i++ {
		l.Tiles[i*width] = 777
		l.Tiles[i*width+width-1] = 777
	}

	// Encode the walls depending on the number of adjacent walls.
	// Each neighbor (N, E, S, W) is encoded as a bit in the following order: N, E, S, W.
	// 0 = no wall, 1 = wall.
	// For example, a wall with a wall to the north and east would be encoded as 1001 = 9.
	// The following table shows the tile ID for each combination of neighbors:
	table := []int{
		0,   // 0000
		713, // 0001
		768, // 0010
		715, // 0011
		711, // 0100
		712, // 0101
		714, // 0110
		774, // 0111
		825, // 1000
		772, // 1001
		770, // 1010
		773, // 1011
		771, // 1100
		717, // 1101
		716, // 1110
		769, // 1111
	}

	// Create a maze by subdividing the level into 2 parts, adding a wall of width 1 between them, and repeating.
	// NOTE: We will remove a single block from the wall to allow for a door.
	// NOTE: This is a hacky implementation, but it works for now.
	// TODO: Implement a proper maze generation algorithm.
	var subDivide func(x1, y1, x2, y2 int)
	subDivide = func(x1, y1, x2, y2 int) {
		blockWidth := x2 - x1
		blockHeight := y2 - y1

		// Now divide the level into 2 parts, leaving a wall of width 1 between them.
		// NOTE: We will remove a single block from the wall to allow for a door.
		xBelow5 := blockWidth < 5
		yBelow5 := blockHeight < 5

		// If both dimensions are below 5 or one below 3, stop dividing.
		if xBelow5 && yBelow5 || blockHeight < 3 || blockWidth < 3 {
			return
		}

		// Select an axis to divide on.
		axis := 0 // 0 = horizontal, 1 = vertical
		if !xBelow5 && !yBelow5 {
			axis = rand.Intn(2) // If both dimensions are above 5, select a random axis.
		} else if yBelow5 {
			axis = 1 // If the y dimension is below 5, divide vertically.
		}

		if axis == 0 {
			// Horizontal axis.
			// Select a random position to divide on.
			dividePos := rand.Intn(blockHeight-2) + 1

			// Add a wall.
			for i := 0; i < blockWidth; i++ {
				l.Tiles[(y1+dividePos)*width+x1+i] = 769
			}

			// Remove a single block from the wall to allow for a door.
			var doorIdx int
			if rand.Intn(2) == 0 {
				doorIdx = (y1+dividePos)*width + x1 + blockWidth - 1
			} else {
				doorIdx = (y1+dividePos)*width + x1
			}
			l.Tiles[doorIdx] = 0

			// Recurse, leaving a gap where the wall is.
			subDivide(x1, y1, x2, y1+dividePos)
			subDivide(x1, y1+dividePos+1, x2, y2)
		} else {
			// Vertical axis.
			// Select a random position to divide on.
			dividePos := rand.Intn(blockWidth-2) + 1

			// Add a wall.
			for i := 0; i < blockHeight; i++ {
				l.Tiles[(y1+i)*width+x1+dividePos] = 769
			}

			// Remove a single block from the wall to allow for a door.
			var doorIdx int
			if rand.Intn(2) == 0 {
				doorIdx = (y1+blockHeight-1)*width + x1 + dividePos
			} else {
				doorIdx = (y1)*width + x1 + dividePos
			}
			l.Tiles[doorIdx] = 0

			// Recurse, leaving a gap where the wall is.
			subDivide(x1, y1, x1+dividePos, y2)
			subDivide(x1+dividePos+1, y1, x2, y2)
		}
	}
	subDivide(1, 1, width-1, height-1)

	tileCopy := make([]int, width*height)
	copy(tileCopy, l.Tiles)

	// Encode the walls... loop through all tiles and encode the walls depending on the number of adjacent walls.
	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			// If this tile is nothing, skip it.
			if l.Tiles[y*width+x] == 0 {
				continue
			}
			// Encode the walls depending on the number of adjacent walls.
			// Each neighbor (N, E, S, W) is encoded as a bit in the following order: N, E, S, W.
			// 0 = no wall, 1 = wall.
			var encoded int

			// Check north.
			if tileCopy[(y-1)*width+x] != 0 {
				encoded |= 1 << 3
			}

			// Check east.
			if tileCopy[y*width+x+1] != 0 {
				encoded |= 1 << 2
			}

			// Check south.
			if tileCopy[(y+1)*width+x] != 0 {
				encoded |= 1 << 1
			}

			// Check west.
			if tileCopy[y*width+x-1] != 0 {
				encoded |= 1 << 0
			}
			l.Tiles[y*width+x] = table[encoded]
		}
	}

	// Add some objects.
	l.placeObjectRandom(ObjectTypeBed)
	l.placeObjectRandom(ObjectTypeFridge)
	l.placeObjectRandom(ObjectTypeCouch)
	l.placeObjectRandom(ObjectTypeToilet)
	l.placeObjectRandom(ObjectTypeShower)
	return l
}

func (l *Level) placeObjectRandom(ot *ObjectType) {
	pos := vectors.Vec2{
		X: float64(rand.Intn(l.Width)),
		Y: float64(rand.Intn(l.Height)),
	}
	for l.IsSolid(int(pos.X), int(pos.Y)) {
		pos.X = float64(rand.Intn(l.Width))
		pos.Y = float64(rand.Intn(l.Height))
	}
	l.Objects = append(l.Objects, ot.New(pos))
}

func (l *Level) GetTile(x, y int) int {
	if x < 0 || y < 0 || x >= l.Width || y >= l.Height {
		return 0
	}
	return l.Tiles[y*l.Width+x]
}

func (l *Level) GetGround(x, y int) int {
	if x < 0 || y < 0 || x >= l.Width || y >= l.Height {
		return 0
	}
	return l.Ground[y*l.Width+x]
}

func (l *Level) IsSolid(x, y int) bool {
	if x < 0 || y < 0 || x >= l.Width || y >= l.Height {
		return true
	}
	return l.Tiles[y*l.Width+x] != 0
}
