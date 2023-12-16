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
}
