## genshops: Name and Slogan Generator

genshops is a procedural generator implemented in Golang that enables you to effortlessly create captivating names and catchy slogans for shops and factions in your fantasy setting. With this tool, you can generate random, humorous names and slogans tailored to your preferences, adding a touch of whimsy to your world. At least that how ChatGPT would describe it.

Anyway, this is currently just an experiment using my glorious genstory generator thingy and is just a draft for now.

### Features

* **Name Generation**: Generate random, funny names using alliteration to make your shops and factions stand out. Examples include "Marv's Magical Miscellanea" and "Percy's Perfect Potions."

* **Slogan Generation**: (TODO) Create memorable slogans for your establishments. Whether you want something witty, clever, or descriptive, the generator has you covered. For instance, you might get a slogan like "Our armor is pricy, but cheap armor costs you an arm and a leg."

* **Token Customization**: Influence the generator's outcome by setting specific tokens. You can specify, for example, the shop owner's name, using the [OWNER_NAME] token, and even influence the goods being sold by setting the [PRODUCT] token.

Feel free to explore the codebase, make modifications, and adapt it to suit your specific needs.

### Usage

```go
package main

import (
    "fmt"

    "github.com/Flokey82/genideas/genshops"
    "github.com/Flokey82/go_gens/genstory"
)

func main() {
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
```

Which prints for example:

```shell
$ Glorbnorb's glorious gear gallery
```

### Contributing

We welcome contributions to enhance genshops. If you have ideas for improvements or would like to fix any issues, please submit a pull request. 
