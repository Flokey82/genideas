package main

import (
	"github.com/Flokey82/genideas/genlandmarknames"
)

func main() {
	gen := genlandmarknames.NewNameGenerators(0)
	println("Desert:")
	for i := 0; i < 100; i++ {
		println(gen.Desert.Generate(int64(i), i%2 == 0))
	}
	println("Mountain:")
	for i := 0; i < 100; i++ {
		println(gen.Mountain.Generate(int64(i), i%2 == 0))
	}
	println("MountainRange:")
	for i := 0; i < 100; i++ {
		println(gen.MountainRange.Generate(int64(i), i%2 == 0))
	}
	println("Forest:")
	for i := 0; i < 100; i++ {
		println(gen.Forest.Generate(int64(i), i%2 == 0))
	}
	println("Swamp:")
	for i := 0; i < 100; i++ {
		println(gen.Swamp.Generate(int64(i), i%2 == 0))
	}
	println("River:")
	for i := 0; i < 100; i++ {
		println(gen.River.Generate(int64(i), i%2 == 0, i%2 == 0))
	}
	println("Lake:")
	for i := 0; i < 100; i++ {
		println(gen.Lake.Generate(int64(i), i%2 == 0))
	}
	println("Plains:")
	for i := 0; i < 100; i++ {
		println(gen.Plains.Generate(int64(i), i%2 == 0))
	}

	println("Lake (alternative):")
	for i := 0; i < 100; i++ {
		st, err := genlandmarknames.LakeNameConfig.Generate(nil)
		if err != nil {
			println(err.Error())
		} else {
			println(st.Text)
		}
	}

	println("River (alternative):")
	for i := 0; i < 100; i++ {
		st, err := genlandmarknames.RiverNameConfig.Generate(nil)
		if err != nil {
			println(err.Error())
		} else {
			println(st.Text)
		}
	}
}
