package main

import (
	"math/rand"

	"github.com/Flokey82/genideas/simpeople2"
	"github.com/Flokey82/go_gens/vectors"
)

func main() {
	// Test decay of sleep.
	mSleep := simpeople2.MotiveTypeSleep.New()

	for i := 0; i < 20; i++ {
		mSleep.Log()
		mSleep.Tick()
	}

	// Test decay of food.
	mFood := simpeople2.MotiveTypeFood.New()

	for i := 0; i < 20; i++ {
		mFood.Log()
		mFood.Tick()
	}

	// New world.
	w := simpeople2.NewWorld()

	// Add some objects.
	bed := simpeople2.ObjectTypeBed.New(vectors.Vec2{
		X: rand.Float64() * 50,
		Y: rand.Float64() * 50,
	})
	w.Objects = append(w.Objects, bed)

	fridge := simpeople2.ObjectTypeFridge.New(vectors.Vec2{
		X: rand.Float64() * 50,
		Y: rand.Float64() * 50,
	})
	w.Objects = append(w.Objects, fridge)

	tv := simpeople2.ObjectTypeCouch.New(vectors.Vec2{
		X: rand.Float64() * 50,
		Y: rand.Float64() * 50,
	})
	w.Objects = append(w.Objects, tv)

	toilet := simpeople2.ObjectTypeToilet.New(vectors.Vec2{
		X: rand.Float64() * 50,
		Y: rand.Float64() * 50,
	})
	w.Objects = append(w.Objects, toilet)

	shower := simpeople2.ObjectTypeShower.New(vectors.Vec2{
		X: rand.Float64() * 50,
		Y: rand.Float64() * 50,
	})
	w.Objects = append(w.Objects, shower)

	// Add some people.
	p1 := w.NewPerson("Person 1")
	w.People = append(w.People, p1)

	//p2 := w.NewPerson("Person 2")
	//w.People = append(w.People, p2)

	// Tick the world.
	for i := 0; i < 40; i++ {
		w.Tick()
	}
}
