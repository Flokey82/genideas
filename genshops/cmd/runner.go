package main

import (
	"log"

	"github.com/Flokey82/genideas/genshops"
	"github.com/Flokey82/go_gens/genstory"
)

func main() {
	for i := 0; i < 100; i++ {
		shop, err := genshops.ShopTitleConfig.Generate(nil)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(shop.Text)
		}
	}

	for i := 0; i < 100; i++ {
		slogan, err := genshops.SloganConfig.Generate(nil)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(slogan.Text)
		}
	}

	// Specify a shop owner.
	shop, err := genshops.ShopTitleConfig.Generate([]genstory.TokenReplacement{{
		Token:       genshops.TokenOwnerName,
		Replacement: "Glorbnorb",
	}})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(shop.Text)
	}
}
