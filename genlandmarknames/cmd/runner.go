package main

import (
	"github.com/Flokey82/genideas/genlandmarknames"
)

func main() {
	gen := genlandmarknames.NewNameGenerators(0)
	for i := 0; i < 100; i++ {
		println(gen.Desert.Generate(int64(i), i%2 == 0))
	}
	for i := 0; i < 100; i++ {
		println(gen.Mountain.Generate(int64(i), i%2 == 0))
	}
	for i := 0; i < 100; i++ {
		println(gen.MountainRange.Generate(int64(i), i%2 == 0))
	}
	for i := 0; i < 100; i++ {
		println(gen.Forest.Generate(int64(i), i%2 == 0))
	}
	for i := 0; i < 100; i++ {
		println(gen.Swamp.Generate(int64(i), i%2 == 0))
	}
	for i := 0; i < 100; i++ {
		println(gen.River.Generate(int64(i), i%2 == 0, i%2 == 0))
	}
}
