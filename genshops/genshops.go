package genshops

import "github.com/Flokey82/go_gens/genstory"

const (
	TokenOwnerName = "[OWNER_NAME]"
	TokenAdj       = "[ADJ]"
	TokenProduct   = "[PRODUCT]"
	TokenShopType  = "[SHOP_TYPE]"
)

var ShopTitleConfig = &genstory.TextConfig{
	TokenPools: map[string][]string{
		TokenOwnerName: owners,
		TokenAdj:       adjectives,
		TokenProduct:   products,
		TokenShopType:  shopType,
	},
	TokenIsMandatory: map[string]bool{},
	Tokens:           []string{TokenOwnerName, TokenAdj, TokenProduct, TokenShopType},
	Templates:        templates,
	UseAllProvided:   true,
	UseAlliteration:  true,
}

var templates = []string{
	"[OWNER_NAME]'s [ADJ] [PRODUCT]",
	"[OWNER_NAME]'s [ADJ] [PRODUCT] [SHOP_TYPE]",
	"[OWNER_NAME]'s [SHOP_TYPE] of [ADJ] [PRODUCT]",
	"[ADJ] [PRODUCT] [SHOP_TYPE]",
	"[ADJ] [PRODUCT]",
}

var owners = []string{
	"Barry",
	"Ben",
	"Bill",
	"Bob",
	"Charlotte",
	"Chloe",
	"Daisy",
	"David",
	"Ella",
	"Emily",
	"Eric",
	"Felicity",
	"Fred",
	"George",
	"Grace",
	"Harry",
	"Henry",
	"Isabella",
	"Jack",
	"Jacob",
	"James",
	"Jessica",
	"Joshua",
	"Leo",
	"Lewis",
	"Liam",
	"Lily",
	"Logan",
	"Lord Lard",
	"Lucas",
	"Lucy",
	"Mara",
	"Marlene",
	"Martha",
	"Mary",
	"Mason",
	"Matthew",
	"Max",
	"Mia",
	"Michael",
	"Nathan",
	"Noah",
	"Oliver",
	"Olivia",
	"Oscar",
	"Poppy",
	"Robert",
	"Ruby",
	"Samuel",
	"Scarlett",
	"Sebastian",
	"Shirley",
	"Sir Loin",
	"Sir Longbottom",
	"Sophia",
	"Sophie",
	"Thomas",
	"Tom",
	"Tommy",
	"Tyler",
	"Ulrich",
	"Ulysses",
	"Uma",
	"Victoria",
	"William",
	"Winston",
	"Wyatt",
	"Xander",
	"Xavier",
	"Yvonne",
	"Zachariah",
	"Zachary",
	"Zoe",
}

var adjectives = []string{
	"ancient",
	"arcane",
	"assorted",
	"battered",
	"beautiful",
	"best",
	"big",
	"bizarre",
	"brand-new",
	"broken",
	"celestial",
	"clean",
	"cobbled",
	"curious",
	"cute",
	"delicate",
	"delightful",
	"devious",
	"dirty",
	"divine",
	"eerie",
	"elegant",
	"enchanted",
	"enchanting",
	"ethereal",
	"exquisite",
	"extraordinary",
	"fancy",
	"fantastic",
	"fine",
	"flawless",
	"fresh",
	"glorious",
	"gorgeous",
	"grand",
	"great",
	"handcrafted",
	"handmade",
	"handy",
	"heavenly",
	"impressive",
	"lovely",
	"magical",
	"magnificent",
	"majestic",
	"marvellous",
	"marvelous",
	"metaphysical",
	"mint",
	"mysterious",
	"mystic",
	"mystical",
	"mystical",
	"new",
	"occult",
	"old",
	"otherworldly",
	"perfect",
	"polished",
	"precious",
	"priceless",
	"pristine",
	"rare",
	"remarkable",
	"resplendent",
	"rusty",
	"scintillating",
	"second-hand",
	"shimmering",
	"shiny",
	"sparkling",
	"special",
	"spiritual",
	"splendid",
	"strange",
	"stunning",
	"sublime",
	"superb",
	"superior",
	"supernatural",
	"supreme",
	"tarnished",
	"terrific",
	"transcendent",
	"transcendental",
	"transmundane",
	"unblemished",
	"uncommon",
	"unimpaired",
	"unique",
	"unmarked",
	"unmarred",
	"unmatched",
	"unparalleled",
	"unrivaled",
	"unspoiled",
	"untouched",
	"unusual",
	"unusual",
	"used",
	"weird",
	"wonderful",
}

var products = []string{
	"accoutrements",
	"apples",
	"armour",
	"belongings",
	"books",
	"boots",
	"clothes",
	"coats",
	"debris",
	"detritus",
	"dresses",
	"equipment",
	"food",
	"garbage",
	"gear",
	"goods",
	"hardware",
	"hodgepodge",
	"inventory",
	"items",
	"jewellery",
	"junk",
	"litter",
	"merchandise",
	"miscellanea",
	"offerings",
	"opportunities",
	"pants",
	"paraphernalia",
	"possessions",
	"potions",
	"produce",
	"provisions",
	"rags",
	"remains",
	"rubbish",
	"scrap",
	"scraps",
	"shoes",
	"skirts",
	"socks",
	"stock",
	"stuff",
	"supplies",
	"things",
	"tools",
	"trappings",
	"trash",
	"treasures",
	"trinkets",
	"trousers",
	"underwear",
	"wares",
	"weapons",
}

var shopType = []string{
	"arcade",
	"atelier",
	"auctionhouse",
	"bazaar",
	"bodega",
	"booth",
	"boutique",
	"corner",
	"cove",
	"craftstore",
	"den",
	"depot",
	"dispensary",
	"distributor",
	"emporium",
	"establishment",
	"exchange",
	"factory",
	"franchise",
	"gallery",
	"hideout",
	"hypermarket",
	"lair",
	"lodge",
	"mall",
	"market",
	"marketplace",
	"outlet",
	"parlour",
	"reseller",
	"salon",
	"shop",
	"shoppe",
	"stall",
	"stand",
	"store",
	"storefront",
	"studio",
	"supermarket",
	"superstore",
	"surplus",
	"trader",
	"tradingpost",
	"upseller",
	"vendor",
	"warehouse",
	"wholesale",
	"workshop",
}
