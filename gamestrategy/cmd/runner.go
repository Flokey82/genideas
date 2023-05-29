package main

import (
	"fmt"

	"github.com/Flokey82/genideas/gamestrategy"
)

func main() {
	fmt.Println("Hello world!")
	g := gamestrategy.NewGrid(100, 100)
	g.AddPlayer(gamestrategy.NewPlayer("Player 1"))
	g.AddPlayer(gamestrategy.NewPlayer("Player 2"))
	g.AddPlayer(gamestrategy.NewPlayer("Player 3"))

	for i := 0; i < 1000; i++ {
		g.Tick()
	}

	g.ExportWebp("test.webp")
}
