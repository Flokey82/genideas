package simpeople2

// findPath finds a path from e to target and returns the path.
func findPath(w *World, e *Person, target *Object) []*Node {
	// New A* pathfinder.
	p := NewPathfinder(w, int(e.Position.X), int(e.Position.Y), int(target.Position.X), int(target.Position.Y))
	// Find the path.
	return p.FindPath()
}

// Node represents a node in the pathfinder.
type Node struct {
	Parent *Node
	X, Y   int
}

func newNode(x, y int, parent *Node) *Node {
	return &Node{
		Parent: parent,
		X:      x,
		Y:      y,
	}
}

// Pathfinder represents a pathfinder.
type Pathfinder struct {
	world *World
	start *Node
	end   *Node
}

// NewPathfinder creates a new pathfinder.
func NewPathfinder(w *World, x, y, tx, ty int) *Pathfinder {
	return &Pathfinder{
		world: w,
		start: newNode(x, y, nil),
		end:   newNode(tx, ty, nil),
	}
}

// FindPath finds a path from the start to the end.
// NOTE: Poor man's A*.
// TODO: Use heuristics to prioritize the search to nodes that are closer to the
// end node.
func (p *Pathfinder) FindPath() []*Node {
	var closed []*Node
	open := []*Node{p.start}
	// While there are still nodes to check.
	for len(open) > 0 {
		// Get the first node from the open list.
		// NOTE: This is pretty inefficient to do it this way due to the
		// constant re-allocation. In theory we could allocate a slice with a
		// reasonable capacity and use an index to keep track of the position
		// and once we've exhausted the capacity, we can copy the last node to
		// the first position and set the index to 0 and truncate the slice.
		n := open[0]
		open = open[1:]
		// If this is the end node, we are done.
		if n.X == p.end.X && n.Y == p.end.Y {
			return p.reconstructPath(n)
		}
		// Add the node to the closed list.
		closed = append(closed, n)
		// Get the neighbors.
		neighbors := p.getNeighbors(n)
		// For each neighbor.
		for _, neighbor := range neighbors {
			// If the neighbor is in the closed list, skip it.
			if isInList(neighbor, closed) {
				continue
			}
			// If the neighbor is not in the open list, add it.
			if !isInList(neighbor, open) {
				open = append(open, neighbor)
			}
		}
	}
	// No path found.
	return nil
}

// reconstructPath reconstructs the path from the end node.
func (p *Pathfinder) reconstructPath(n *Node) []*Node {
	path := []*Node{n}
	for n.Parent != nil {
		n = n.Parent
		path = append([]*Node{n}, path...)
	}
	return path
}

// getNeighbors returns the neighbors of the node.
// TODO: Re-use a pre-allocated slice for the neighbors.
func (p *Pathfinder) getNeighbors(n *Node) []*Node {
	// Check each direction.
	var neighbors []*Node
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			// Skip the center.
			if dx == 0 && dy == 0 {
				continue
			}
			// If we can move to the tile, add it as a neighbor.
			if !p.world.IsSolid(n.X+dx, n.Y+dy) {
				neighbors = append(neighbors, newNode(n.X+dx, n.Y+dy, n))
			}
		}
	}
	return neighbors
}

// isInList returns true if the node is in the list.
func isInList(n *Node, list []*Node) bool {
	for _, node := range list {
		if node.X == n.X && node.Y == n.Y {
			return true
		}
	}
	return false
}
