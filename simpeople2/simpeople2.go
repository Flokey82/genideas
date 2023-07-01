package simpeople2

// World represents the world of the simulation.
type World struct {
	Objects []*Object
	People  []*Person
}

// NewWorld creates a new world.
func NewWorld() *World {
	return &World{}
}

// Tick ticks the world.
func (w *World) Tick() {
	for _, p := range w.People {
		p.Tick()
		p.Log()
	}
}
