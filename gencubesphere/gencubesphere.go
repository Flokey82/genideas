package gencubesphere

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type CubeSphere struct {
	NumPoints     int
	PointsPerFace int
	PointsPerSide int
	HalfCellSize  float64 // The half size of a cell in the grid (used for offsetting points).
}

func NewCubeSphere(numPoints int) *CubeSphere {
	ptsPerFace := numPoints / 6
	ptsPerSideOnFace := int(math.Sqrt(float64(ptsPerFace)))
	cellSize := 2.0 / float64(ptsPerSideOnFace)
	return &CubeSphere{
		NumPoints:     numPoints,
		PointsPerFace: ptsPerFace,
		PointsPerSide: ptsPerSideOnFace,
		HalfCellSize:  cellSize / 2.0,
	}
}

func (sphere *CubeSphere) IndexToFace(index int) (face int) {
	return index / sphere.PointsPerFace
}

func (sphere *CubeSphere) IndexToCoordinates(index int) (x, y, z float64) {
	face := sphere.IndexToFace(index)

	// Calculate the index of the point on the face.
	indexOnFace := index % sphere.PointsPerFace

	// Calculate the coordinates on the face.
	x, y = sphere.IndexToCoordinatesOnFace(indexOnFace)

	// Calculate the coordinates on the cube.
	x, y, z = sphere.FaceCoordinatesToCubeCoordinates(face, x, y)

	// Normalize the coordinates to a unit sphere.
	vLen := math.Sqrt(x*x + y*y + z*z)
	x /= vLen
	y /= vLen
	z /= vLen

	return x, y, z
}

func (sphere *CubeSphere) ExportWavefrontOBJ(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// Vertices
	for i := 0; i < sphere.NumPoints; i++ {
		x, y, z := sphere.IndexToCoordinates(i)
		// Convert to wavefront obj coordinate system.
		x, y, z = x, z, -y
		w.WriteString(fmt.Sprintf("v %f %f %f \n", x, y, z))
	}
	w.Flush()
	return nil
}

func (sphere *CubeSphere) IndexToCoordinatesOnFace(index int) (x, y float64) {
	// Calculate the coordinates on the face.
	x = float64(index%sphere.PointsPerSide) / float64(sphere.PointsPerSide)
	y = float64(index/sphere.PointsPerSide) / float64(sphere.PointsPerSide)

	return x, y
}

func (sphere *CubeSphere) IndexOnFaceToCubeIndex(indexOnFace, face int) (index int) {
	return indexOnFace + face*sphere.PointsPerFace
}

/**
 * The traditional cubemap numbering of faces is [+x, -x, +y, -y, +z, -z].
 * They are arranged visually:
 *
 *       +y                 up
 *    -x +z +x -z     left front right back
 *       -y                down
 *
 * However, I don't actually use GPU cubemaps and I want to make navigation
 * easier, so I'm using a different arrangement: [+y +x -y -x +z -z]
 *
 *   +z                 north
 *   +y +x -y -x        front right back left
 *   -z                 south
 */
func (sphere *CubeSphere) FaceCoordinatesToCubeCoordinates(face int, x, y float64) (cx, cy, cz float64) {
	// TODO: Add support for grid distortion compensation.
	// Either tangent adjustment from:
	// https://www.redblobgames.com/x/1938-square-tiling-of-sphere/
	// or a different approach from:
	// https://catlikecoding.com/unity/tutorials/procedural-meshes/cube-sphere/
	// static float3 CubeToSphere (float3 p) => p * sqrt(
	//  1f - ((p * p).yxx + (p * p).zzy) / 2f + (p * p).yxx * (p * p).zzy / 3f
	// );

	// Our coordinate system looks like this:
	//
	//           +z
	// (back) -y  |
	//          \ |
	//           \|
	//  -x -------\------- +x
	//            |\
	//            | \
	//            |  +y (front)
	//           -z

	// Calculate the coordinates on the cube that will house the unit sphere.
	// The cube is centered around the origin.
	//
	// The cube is 2 units wide and tall, so we need to scale the coordinates to be
	// between -1 and 1.
	//
	//
	// |------2-----|
	// +------------+
	// |\           |\
	// | \          | \
	// |  \         |  \
	// |   +------------+
	// |   |        |   |
	// +---|--------+   |
	//  \  |         \  |
	//   \ |          \ |
	//    \|           \|
	//     +------------+
	switch face {
	case faceFront: // Front
		cx = x*2.0 - 1.0
		cy = 1.0
		cz = 1.0 - y*2.0
	case faceRight: // Right
		cx = 1.0
		cy = 1.0 - x*2.0
		cz = 1.0 - y*2.0
	case faceBack: // Back
		cx = 1.0 - x*2.0
		cy = -1.0
		cz = 1.0 - y*2.0
	case faceLeft: // Left
		cx = -1.0
		cy = x*2.0 - 1.0
		cz = 1.0 - y*2.0
	case faceNorth: // North (Up)
		cx = x*2.0 - 1.0
		cy = y*2.0 - 1.0
		cz = 1.0
	case faceSouth: // South (Down)
		cx = x*2.0 - 1.0
		cy = 1.0 - y*2.0
		cz = -1.0
	}

	return cx, cy, cz
}

const (
	faceFront = iota
	faceRight
	faceBack
	faceLeft
	faceNorth
	faceSouth
)

const (
	alongTopSide = iota
	alongRightSide
	alongBottomSide
	alongLeftSide
)

// getAdjacentNeighborSides returns the adjacent faces and which of their sides
// the given face's side is adjacent to (and in which direction) and return it in order:
// [0] = top
// [1] = right
// [2] = bottom
// [3] = left
// The positive or negative direction is indicated by a positive or negative
// value for the side.
func (sphere *CubeSphere) getAdjacentNeighborSides(face int) [4][2]int {
	switch face {
	case faceFront:
		//        _____              _____
		//       |0 1 2|            |     |
		//       |3 4 5|            |  N  |
		//  _____|6_7_8|_____  _____|_____|_____
		// |0 1 2|0 1 2|0 1 2||     |     |     |
		// |3 4 5|3 4 5|3 4 5||  L  |  F  |  R  |
		// |6_7_8|6_7_8|6_7_8||_____|_____|_____|
		//       |0 1 2|  	        |     |
		//       |3 4 5|            |  S  |
		//       |6_7_8| 		    |_____|
		return [4][2]int{
			{faceNorth, alongBottomSide},
			{faceRight, alongLeftSide},
			{faceSouth, alongTopSide},
			{faceLeft, alongRightSide},
		}
	case faceRight:
		//        _____              _____
		//       |6 3 0|            |     |
		//       |7 4 1|            |  N  |
		//  _____|8_5_2|_____  _____|_____|_____
		// |0 1 2|0 1 2|0 1 2||     |     |     |
		// |3 4 5|3 4 5|3 4 5||  F  |  R  |  B  |
		// |6_7_8|6_7_8|6_7_8||_____|_____|_____|
		//       |2 5 8|  	        |     |
		//       |1 4 7|            |  S  |
		//       |0_3_6| 		    |_____|
		return [4][2]int{
			{faceNorth, -alongRightSide},
			{faceBack, alongLeftSide},
			{faceSouth, alongRightSide},
			{faceFront, alongRightSide},
		}
	case faceBack:
		//        _____              _____
		//       |8 7 6|            |     |
		//       |5 4 3|            |  N  |
		//  _____|2_1_0|_____  _____|_____|_____
		// |0 1 2|0 1 2|0 1 2||     |     |     |
		// |3 4 5|3 4 5|3 4 5||  R  |  B  |  L  |
		// |6_7_8|6_7_8|6_7_8||_____|_____|_____|
		//       |8 7 6|  	        |     |
		//       |5 4 3|            |  S  |
		//       |2_1_0| 		    |_____|
		return [4][2]int{
			{faceNorth, -alongTopSide},
			{faceLeft, alongLeftSide},
			{faceSouth, -alongTopSide},
			{faceRight, alongRightSide},
		}
	case faceLeft:
		//        _____              _____
		//       |2 5 8|            |     |
		//       |1 4 7|            |  N  |
		//  _____|0_3_6|_____  _____|_____|_____
		// |0 1 2|0 1 2|0 1 2||     |     |     |
		// |3 4 5|3 4 5|3 4 5||  B  |  L  |  F  |
		// |6_7_8|6_7_8|6_7_8||_____|_____|_____|
		//       |6 3 0|  	        |     |
		//       |7 4 1|            |  S  |
		//       |8_5_2| 		    |_____|
		return [4][2]int{
			{faceNorth, alongLeftSide},
			{faceFront, alongLeftSide},
			{faceSouth, -alongRightSide},
			{faceBack, alongRightSide},
		}
	case faceNorth:
		//	     _____              _____
		//	    |8 7 6|            |     |
		//	    |5 4 3|            |  B  |
		// _____|2_1_0|_____  _____|_____|_____
		//|6 3 0|0 1 2|2 5 8||     |     |     |
		//|7 4 1|3 4 5|1 4 7||  L  |  N  |  R  |
		//|8_5_2|6_7_8|0_3_6||_____|_____|_____|
		//	    |0 1 2|  	       |     |
		//	    |3 4 5|            |  F  |
		//	    |6_7_8| 	       |_____|
		return [4][2]int{
			{faceBack, -alongTopSide},
			{faceRight, -alongTopSide},
			{faceFront, alongTopSide},
			{faceLeft, alongTopSide},
		}
	case faceSouth:
		//	     _____              _____
		//	    |0 1 2|            |     |
		//	    |3 4 5|            |  F  |
		// _____|6_7_8|_____  _____|_____|_____
		//|2 5 8|0 1 2|6 3 0||     |     |     |
		//|1 4 7|3 4 5|7 4 1||  L  |  S  |  R  |
		//|0_3_6|6_7_8|8_5_2||_____|_____|_____|
		//	    |8 7 6|  	       |     |
		//	    |5 4 3|            |  B  |
		//	    |2_1_0| 	       |_____|
		return [4][2]int{
			{faceFront, alongBottomSide},
			{faceRight, alongBottomSide},
			{faceBack, -alongBottomSide},
			{faceRight, -alongBottomSide},
		}
	}

	return [4][2]int{}
}

// getNThIndexOnSide returns the index of the point on the given side that is
// the nth point on that side.
func (sphere *CubeSphere) getNThIndexOnSide(side, n int) (index int) {
	switch side {
	case alongTopSide:
		return n
	case alongRightSide:
		return (sphere.PointsPerSide - 1) + n*sphere.PointsPerSide
	case alongBottomSide:
		return sphere.PointsPerFace - sphere.PointsPerSide + n
	case alongLeftSide:
		return n * sphere.PointsPerSide
	}

	return index
}

// getNThIndexOnSideOfFace returns the index of the point on the given side of
// the given face that is the nth point on that side.
func (sphere *CubeSphere) getNThIndexOnSideOfFace(face, side, n int) (index int) {
	return sphere.getNThIndexOnSide(side, n) + face*sphere.PointsPerFace
}

// FindDirectNeighbors returns the indices of the points that are directly
// adjacent to the given index and can either be on the same face or on a
// neighboring face.
func (sphere *CubeSphere) FindDirectNeighbors(index int) (neighbors []int) {
	// Calculate the face of the index.
	face := sphere.IndexToFace(index)
	faceStartIndex := face * sphere.PointsPerFace

	// Calculate the index of the point on the face.
	indexOnFace := index % sphere.PointsPerFace

	// Calculate the coordinates on the face.
	x, y := indexOnFace%sphere.PointsPerSide, indexOnFace/sphere.PointsPerSide

	// Add the neighbors on the same face.
	for _, dxy := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		nx := x + dxy[0]
		if nx < 0 || nx >= sphere.PointsPerSide {
			continue
		}
		ny := y + dxy[1]
		if ny < 0 || ny >= sphere.PointsPerSide {
			continue
		}
		neighbors = append(neighbors, (sphere.PointsPerSide*ny+nx)+faceStartIndex)
	}

	// Get the adjacent faces and which of their sides the given face's side is
	// adjacent to.
	adjacentFaces := sphere.getAdjacentNeighborSides(face)

	getNeighbor := func(side, n int) int {
		adjacentFace := adjacentFaces[side][0]
		adjacentSide := adjacentFaces[side][1]
		if adjacentSide < 0 {
			return sphere.getNThIndexOnSideOfFace(adjacentFace, -adjacentSide, sphere.PointsPerSide-n)
		}
		return sphere.getNThIndexOnSideOfFace(adjacentFace, adjacentSide, n)
	}

	// Figure out which sides the index is adjacent to.
	if x == 0 { // Left side
		neighbors = append(neighbors, getNeighbor(alongLeftSide, y))
	} else if x == sphere.PointsPerSide-1 { // Right side
		neighbors = append(neighbors, getNeighbor(alongRightSide, y))
	}
	if y == 0 { // Top side
		neighbors = append(neighbors, getNeighbor(alongTopSide, x))
	} else if y == sphere.PointsPerSide-1 { // Bottom side
		neighbors = append(neighbors, getNeighbor(alongBottomSide, x))
	}

	return neighbors
}
