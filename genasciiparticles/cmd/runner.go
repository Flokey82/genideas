package main

import "github.com/Flokey82/genideas/genasciiparticles"

func main() {
	g, err := genasciiparticles.New()
	if err != nil {
		panic(err)
	}
	g.Run()
}
