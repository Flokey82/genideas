package main

import (
	"github.com/Flokey82/genideas/gameoverworld"
)

func main() {
	world, err := gameoverworld.New()
	if err != nil {
		panic(err)
	}
	world.Run()
}
