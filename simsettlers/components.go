package simsettlers

import (
	"math"

	"github.com/Flokey82/go_gens/vectors"
)

type ObjectType struct {
	Name    string
	Actions []*Action
}

func (o *ObjectType) New(l Location) *Object {
	return &Object{
		ObjectType: o,
		Location:   l,
	}
}

type Object struct {
	*ObjectType
	Owner    *Person // The owner might be unknown?
	Location         // Location might be a reference to a person, or a fixed location in the world.
}

// UpdateLocation updates the location of the object.
func (o *Object) UpdateLocation(l Location) {
	o.Location = l
}

type Location interface {
	Position() vectors.Vec2 // Current position of the location
}

type LocationFixed struct {
	Pos vectors.Vec2
}

func (l *LocationFixed) Position() vectors.Vec2 {
	return l.Pos
}

type LocationPerson struct {
	X     float64
	Y     float64
	Speed vectors.Vec2
}

// Position implements the Location interface.
func (p *LocationPerson) Position() vectors.Vec2 {
	return vectors.Vec2{X: p.X, Y: p.Y}
}

const walkingSpeed = 0.5

// SetDirection sets the speed of the person to the given direction.
func (p *LocationPerson) SetDirection(x, y float64) {
	// Calculate the speed vector from our position and the target position.
	p.Speed = vectors.Vec2{X: x - p.X, Y: y - p.Y}

	// Check the length of the vector. If it is smaller than the walking speed,
	// we do not need to adjust the magnitude of the vector.
	if p.Speed.Len() > walkingSpeed {
		p.Speed = p.Speed.Normalize().Mul(walkingSpeed)
	}
}

func (p *LocationPerson) distanceTo(x, y float64) float64 {
	dX, dY := p.X-x, p.Y-y
	return math.Sqrt(dX*dX + dY*dY)
}

// Move moves the person by their speed.
func (p *LocationPerson) Move(elapsed float64) {
	p.X += p.Speed.X * elapsed
	p.Y += p.Speed.Y * elapsed
}
