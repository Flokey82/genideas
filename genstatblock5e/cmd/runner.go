package main

import (
	"log"

	"github.com/Flokey82/genideas/genstatblock5e"
)

func main() {
	attrs := genstatblock5e.GenerateAttributes(5, 10, 19, 1.0)
	log.Printf("Attributes: %+v", attrs)
}
