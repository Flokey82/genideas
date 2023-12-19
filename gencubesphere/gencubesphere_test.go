package gencubesphere

import (
	"sort"
	"testing"
)

func TestNeighbors(t *testing.T) {
	sphere := NewCubeSphere(3 * 3 * 6) // 3 points per side, 6 faces

	// Face 0 on the front.
	gotNeighbors := sphere.FindDirectNeighbors(sphere.IndexOnFaceToCubeIndex(0, faceFront))
	wantNeighbors := []int{
		sphere.IndexOnFaceToCubeIndex(1, faceFront),
		sphere.IndexOnFaceToCubeIndex(3, faceFront),
		sphere.IndexOnFaceToCubeIndex(6, faceNorth),
		sphere.IndexOnFaceToCubeIndex(2, faceLeft),
	}

	sort.Ints(gotNeighbors)
	sort.Ints(wantNeighbors)

	if len(gotNeighbors) != len(wantNeighbors) {
		t.Errorf("got %v, want %v", gotNeighbors, wantNeighbors)
	}
	for i := range gotNeighbors {
		if gotNeighbors[i] != wantNeighbors[i] {
			t.Errorf("got %v, want %v", gotNeighbors, wantNeighbors)
		}
	}
	sphere = NewCubeSphere(30 * 30 * 6) // 100 points per side, 6 faces
	sphere.ExportWavefrontOBJ("test.obj")

	numPerFace := 4 * 4

	sphere = NewCubeSphere(numPerFace * 6) // 4*4 points per side, 6 faces

	// Check if we can reconstruct the face from the index's lat/lon.
	for _, wantFace := range []int{faceFront, faceBack, faceLeft, faceRight, faceNorth, faceSouth} {
		for i := 0; i < numPerFace; i++ {
			index := sphere.IndexOnFaceToCubeIndex(i, wantFace)
			lat, lon := sphere.IndexToLatLonDeg(index)
			gotFace := sphere.LatLonDegToFace(lat, lon)
			if wantFace != gotFace {
				t.Errorf("%d: Face issue; got %v, want %v", i, gotFace, wantFace)
			}
		}
	}

	// Check if we can convert between indices and lat/lon.
	/*
		for _, wantFace := range []int{faceFront, faceBack, faceLeft, faceRight, faceNorth, faceSouth} {
			for i := 0; i < numPerFace; i++ {
				log.Println("Face", wantFace, "index", i)
				index := sphere.IndexOnFaceToCubeIndex(i, wantFace)
				lat, lon := sphere.IndexToLatLonDeg(index)
				gotIndex := sphere.LatLonDegToIndex(lat, lon)
				if index != gotIndex {
					t.Errorf("%d: Index issue; got %v, want %v", i, gotIndex, index)
				}
			}
		}
	*/
}
