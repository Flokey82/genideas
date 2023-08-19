package main

import (
	"github.com/Flokey82/genideas/gameasciitiles"
)

func main() {
	g, err := gameasciitiles.NewGame(32)
	if err != nil {
		panic(err)
	}
	if err := g.Run(); err != nil {
		panic(err)
	}
}
