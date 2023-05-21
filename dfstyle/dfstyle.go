// Package dfstyle provides a way to generate a world map in the style of Dwarf Fortress.
// This code is a crude port of https://github.com/Dozed12/df-style-worldgen
package dfstyle

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	WORLD_WIDTH  = 200
	WORLD_HEIGHT = 80

	SCREEN_WIDTH  = 200
	SCREEN_HEIGHT = 80

	CIVILIZED_CIVS = 3
	TRIBAL_CIVS    = 3

	MIN_RIVER_LENGHT = 3

	CIV_MAX_SITES      = 20
	EXPANSION_DISTANCE = 10
	WAR_DISTANCE       = 8
)

type Tile struct {
	temp       float64
	height     float64
	precip     float64
	drainage   float64
	biome      string
	hasRiver   bool
	isCiv      bool
	biomeID    int
	prosperity float64
}

type Race struct {
	Name              string
	PrefBiome         []int
	Strenght          float64
	Size              float64
	ReproductionSpeed float64
	Aggressiveness    float64
	Form              string
}

type CivSite struct {
	x          int
	y          int
	category   string
	suitable   bool
	popcap     int
	Population int
	isCapital  bool
}

type Army struct {
	x    int
	y    int
	Civ  *Civ
	Size int
}

type Civ struct {
	Name            string
	Race            *Race
	Government      GovernmentType
	Color           color.RGBA
	Flag            [][]color.RGBA
	Aggression      float64 // Race.Aggressiveness + Government.Aggressiveness
	Sites           []CivSite
	SuitableSites   []CivSite
	atWar           bool
	Army            *Army
	TotalPopulation int
}

func NewCiv(name string, race *Race, gov GovernmentType, col color.RGBA, flag [][]color.RGBA, aggression float64) *Civ {
	return &Civ{
		Name:       name,
		Race:       race,
		Government: gov,
		Color:      col,
		Flag:       flag,
		Aggression: aggression + race.Aggressiveness + gov.Aggressiveness,
	}
}

func (c *Civ) PrintInfo() {
	fmt.Println(c.Name)
	fmt.Println(c.Race.Name)
	fmt.Println(c.Government.Name)
	fmt.Println("Aggression:", c.Aggression)
	fmt.Println("Suitable Sites:", len(c.SuitableSites), "\n")
}

type GovernmentType struct {
	Name           string
	Description    string
	Aggressiveness float64
	Militarization float64
	TechBonus      float64
}

type War struct {
	Side1 *Civ
	Side2 *Civ
}

// - General Functions -

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func PointDistRound(pt1x, pt1y, pt2x, pt2y int) int {
	distance := abs(pt2x-pt1x) + abs(pt2y-pt1y)
	//distance = round(distance)
	return distance
}

func FlagGenerator(Color color.RGBA) [][]color.RGBA {
	Flag := make([][]color.RGBA, 12)
	for i := range Flag {
		Flag[i] = make([]color.RGBA, 4)
	}

	BackColor1 := Color
	BackColor2 := Palette[rand.Intn(len(Palette))]

	OverColor1 := Palette[rand.Intn(len(Palette))]
	OverColor2 := Palette[rand.Intn(len(Palette))]

	BackFile, _ := os.Open("Background.txt")
	OverlayFile, _ := os.Open("Overlay.txt")

	BTypes := (countLines("Background.txt") + 1) / 5
	OTypes := (countLines("Overlay.txt") + 1) / 5

	Back := rand.Intn(BTypes) + 1
	Overlay := rand.Intn(OTypes) + 1

	for a := 0; a < 53*(Back-1); a++ {
		BackFile.Read(make([]byte, 1))
	}

	for a := 0; a < 53*(Overlay-1); a++ {
		OverlayFile.Read(make([]byte, 1))
	}

	for y := 0; y < 4; y++ {
		for x := 0; x < 12; x++ {

			C := make([]byte, 1)
			BackFile.Read(C)
			for C[0] == '\n' {
				BackFile.Read(C)
			}

			if C[0] == '#' {
				Flag[x][y] = BackColor1
			} else if C[0] == '"' {
				Flag[x][y] = BackColor2
			}

			OverlayFile.Read(C)
			for C[0] == '\n' {
				OverlayFile.Read(C)
			}

			if C[0] == '#' {
				Flag[x][y] = OverColor1
			} else if C[0] == '"' {
				Flag[x][y] = OverColor2
			}
		}
	}

	BackFile.Close()
	OverlayFile.Close()

	return Flag
}

func countLines(filename string) int {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}

func LowestNeighbour(X, Y int, World [][]Tile) (int, int, int) {
	minval := 1.0
	x := 0
	y := 0

	if World[X+1][Y].height < minval && X+1 < WORLD_WIDTH {
		minval = World[X+1][Y].height
		x = X + 1
		y = Y
	}

	if World[X][Y+1].height < minval && Y+1 < WORLD_HEIGHT {
		minval = World[X][Y+1].height
		x = X
		y = Y + 1
	}

	if World[X-1][Y].height < minval && X-1 > 0 {
		minval = World[X-1][Y].height
		x = X - 1
		y = Y
	}

	if World[X][Y-1].height < minval && Y-1 > 0 {
		minval = World[X][Y-1].height
		x = X
		y = Y - 1
	}

	error := 0

	if x == 0 && y == 0 {
		error = 1
	}

	return x, y, error
}

// - MapGen Functions -
func PoleGen(hm *Heightmap, NS int) {
	if NS == 0 {
		rng := rand.Intn(4) + 2
		for i := 0; i < WORLD_WIDTH; i++ {
			for j := 0; j < rng; j++ {
				hm.SetValue(i, WORLD_HEIGHT-1-j, 0.31)
			}
			rng += rand.Intn(3) - 1
			if rng > 6 {
				rng = 5
			}
			if rng < 2 {
				rng = 2
			}
		}
	}

	if NS == 1 {
		rng := rand.Intn(4) + 2
		for i := 0; i < WORLD_WIDTH; i++ {
			for j := 0; j < rng; j++ {
				hm.SetValue(i, j, 0.31)
			}
			rng += rand.Intn(3) - 1
			if rng > 6 {
				rng = 5
			}
			if rng < 2 {
				rng = 2
			}
		}
	}
}

func TectonicGen(hm *Heightmap, hor int) {
	TecTiles := make([][]int, WORLD_WIDTH)
	for i := range TecTiles {
		TecTiles[i] = make([]int, WORLD_HEIGHT)
	}

	// Define Tectonic Borders
	if hor == 1 {
		pos := rand.Intn(WORLD_HEIGHT/10) + WORLD_HEIGHT/10
		for x := 0; x < WORLD_WIDTH; x++ {
			TecTiles[x][pos] = 1
			pos += rand.Intn(5) - 3
			if pos < 0 {
				pos = 0
			}
			if pos > WORLD_HEIGHT-1 {
				pos = WORLD_HEIGHT - 1
			}
		}
	}

	if hor == 0 {
		pos := rand.Intn(WORLD_WIDTH/10) + WORLD_WIDTH/10
		for y := 0; y < WORLD_HEIGHT; y++ {
			TecTiles[pos][y] = 1
			pos += rand.Intn(5) - 3
			if pos < 0 {
				pos = 0
			}
			if pos > WORLD_WIDTH-1 {
				pos = WORLD_WIDTH - 1
			}
		}
	}

	// Apply elevation to borders
	for x := WORLD_WIDTH / 10; x < WORLD_WIDTH-WORLD_WIDTH/10; x++ {
		for y := WORLD_HEIGHT / 10; y < WORLD_HEIGHT-WORLD_HEIGHT/10; y++ {
			if TecTiles[x][y] == 1 && hm.GetValue(x, y) > 0.3 {
				hm.AddHill(x, y, float64(rand.Intn(3)+2), rand.Float64()*0.03+0.15)
			}
		}
	}
}

func Temperature(temp *Heightmap, hm *Heightmap) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			heighteffect := 0.0
			if y > WORLD_HEIGHT/2 {
				temp.SetValue(x, y, float64(WORLD_HEIGHT-y)-heighteffect)
			} else {
				temp.SetValue(x, y, float64(y)-heighteffect)
			}
			heighteffect = hm.GetValue(x, y)
			if heighteffect > 0.8 {
				heighteffect = heighteffect * 5
				if y > WORLD_HEIGHT/2 {
					temp.SetValue(x, y, float64(WORLD_HEIGHT-y)-heighteffect)
				} else {
					temp.SetValue(x, y, float64(y)-heighteffect)
				}
			}
			if heighteffect < 0.25 {
				heighteffect = heighteffect * 10
				if y > WORLD_HEIGHT/2 {
					temp.SetValue(x, y, float64(WORLD_HEIGHT-y)-heighteffect)
				} else {
					temp.SetValue(x, y, float64(y)-heighteffect)
				}
			}
		}
	}
}

func Percipitation(preciphm, temphm *Heightmap) {
	preciphm.Add(2)

	/*
		for x := 0; x < WORLD_WIDTH; x++ {
			for y := 0; y < WORLD_HEIGHT; y++ {
				temp := temphm.GetValue(x, y)
			}
		}
	*/

	precip := newNoise(2, TCOD_NOISE_DEFAULT_HURST, TCOD_NOISE_DEFAULT_LACUNARITY)
	preciphm.AddFBM(precip, 2, 2, 0, 0, 32, 1, 1)
	preciphm.Normalize(0.0, 1.0)
}

func RiverGen(World [][]Tile) {
	log.Println("Generating Rivers")
	X := rand.Intn(WORLD_WIDTH)
	Y := rand.Intn(WORLD_HEIGHT)

	var XCoor []int
	var YCoor []int

	tries := 0

	//prev := ""

	for World[X][Y].height < 0.8 {
		tries += 1
		X = rand.Intn(WORLD_WIDTH)
		Y = rand.Intn(WORLD_HEIGHT)

		if tries > 2000 {
			return
		}
	}

	XCoor = append(XCoor, X)
	YCoor = append(YCoor, Y)

	for World[X][Y].height >= 0.2 {
		X, Y, error := LowestNeighbour(X, Y, World)
		if error == 1 {
			return
		}

		if World[X][Y].hasRiver || World[X+1][Y].hasRiver || World[X-1][Y].hasRiver || World[X][Y+1].hasRiver || World[X][Y-1].hasRiver {
			break
		}

		if Contains(XCoor, X) && Contains(YCoor, Y) {
			break
		}

		XCoor = append(XCoor, X)
		YCoor = append(YCoor, Y)
	}

	if len(XCoor) <= MIN_RIVER_LENGHT {
		return
	}

	for x := 0; x < len(XCoor); x++ {
		if World[XCoor[x]][YCoor[x]].height < 0.2 {
			break
		}
		World[XCoor[x]][YCoor[x]].hasRiver = true
		if World[XCoor[x]][YCoor[x]].height >= 0.2 && x == len(XCoor)-1 {
			World[XCoor[x]][YCoor[x]].hasRiver = true // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~Change to Lake later
		}
	}
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Prosperity(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			World[x][y].prosperity = (1.0 - math.Abs(World[x][y].precip-0.6) + 1.0 - math.Abs(World[x][y].temp-0.5) + World[x][y].drainage) / 3
		}
	}
}

func MasterWorldGen() [][]Tile {
	fmt.Println(" * World Gen START * ")
	starttime := time.Now()

	//Heightmap
	hm := HeightmapNew(WORLD_WIDTH, WORLD_HEIGHT)

	for i := 0; i < 250; i++ {
		hm.AddHill(rand.Intn(WORLD_WIDTH), rand.Intn(WORLD_HEIGHT), float64(rand.Intn(16-12)+12), float64(rand.Intn(10-6)+6))
	}
	fmt.Println("- Main Hills -")

	for i := 0; i < 1000; i++ {
		hm.AddHill(rand.Intn(WORLD_WIDTH), rand.Intn(WORLD_HEIGHT), float64(rand.Intn(4-2)+2), float64(rand.Intn(10-6)+6))
	}
	fmt.Println("- Small Hills -")

	hm.Normalize(0.0, 1.0)

	noisehm := HeightmapNew(WORLD_WIDTH, WORLD_HEIGHT)
	noise2d := newNoise(2, TCOD_NOISE_DEFAULT_HURST, TCOD_NOISE_DEFAULT_LACUNARITY)
	noisehm.AddFBM(noise2d, 6, 6, 0, 0, 32, 1, 1)
	noisehm.Normalize(0.0, 1.0)
	HeightmapMultiplyHm(hm, noisehm, hm)
	fmt.Println("- Apply Simplex -")

	PoleGen(hm, 0)
	fmt.Println("- South Pole -")

	PoleGen(hm, 1)
	fmt.Println("- North Pole -")

	TectonicGen(hm, 0)
	TectonicGen(hm, 1)
	fmt.Println("- Tectonic Gen -")

	hm.RainErosion(WORLD_WIDTH*WORLD_HEIGHT, 0.07, 0)
	fmt.Println("- Erosion -")

	hm.Clamp(0.0, 1.0)

	//Temperature
	temp := HeightmapNew(WORLD_WIDTH, WORLD_HEIGHT)
	Temperature(temp, hm)
	temp.Normalize(0.0, 1.0)
	fmt.Println("- Temperature Calculation -")

	//Precipitation

	preciphm := HeightmapNew(WORLD_WIDTH, WORLD_HEIGHT)
	Percipitation(preciphm, temp)
	preciphm.Normalize(0.0, 1.0)
	fmt.Println("- Percipitation Calculation -")

	//Drainage

	drainhm := HeightmapNew(WORLD_WIDTH, WORLD_HEIGHT)
	drain := newNoise(2, TCOD_NOISE_DEFAULT_HURST, TCOD_NOISE_DEFAULT_LACUNARITY)
	drainhm.AddFBM(drain, 2, 2, 0, 0, 32, 1, 1)
	drainhm.Normalize(0.0, 1.0)
	fmt.Println("- Drainage Calculation -")

	// VOLCANISM - RARE AT SEA FOR NEW ISLANDS (?) RARE AT MOUNTAINS > 0.9 (?) RARE AT TECTONIC BORDERS (?)

	elapsed_time := time.Since(starttime)
	fmt.Println(" * World Gen DONE *    in: ", elapsed_time.Seconds(), " seconds")

	//Initialize Tiles with Map values
	World := make([][]Tile, WORLD_WIDTH)
	for x := range World {
		World[x] = make([]Tile, WORLD_HEIGHT)
		for y := range World[x] {
			World[x][y] = Tile{
				height:   hm.GetValue(x, y),
				temp:     temp.GetValue(x, y),
				precip:   preciphm.GetValue(x, y),
				drainage: drainhm.GetValue(x, y),
			}
		}
	}

	fmt.Println("- Tiles Initialized -")

	//Prosperity
	Prosperity(World)
	fmt.Println("- Prosperity Calculation -")

	biomeCount := make([]int, 20)
	//Biome info to Tile
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {

			if World[x][y].precip >= 0.10 && World[x][y].precip < 0.33 && World[x][y].drainage < 0.5 {
				World[x][y].biomeID = 3
				if rand.Intn(2)+1 == 2 {
					World[x][y].biomeID = 16
				}
			}

			if World[x][y].precip >= 0.10 && World[x][y].precip > 0.33 {
				World[x][y].biomeID = 2
				if World[x][y].precip >= 0.66 {
					World[x][y].biomeID = 1
				}
			}

			if World[x][y].precip >= 0.33 && World[x][y].precip < 0.66 && World[x][y].drainage >= 0.33 {
				World[x][y].biomeID = 15
				if rand.Intn(5)+1 == 5 {
					World[x][y].biomeID = 5
				}
			}

			if World[x][y].temp > 0.2 && World[x][y].precip >= 0.66 && World[x][y].drainage > 0.33 {
				World[x][y].biomeID = 5
				if World[x][y].precip >= 0.75 {
					World[x][y].biomeID = 6
				}
				if rand.Intn(5)+1 == 5 {
					World[x][y].biomeID = 15
				}
			}

			if World[x][y].precip >= 0.10 && World[x][y].precip < 0.33 && World[x][y].drainage >= 0.5 {
				World[x][y].biomeID = 16
				if rand.Intn(2)+1 == 2 {
					World[x][y].biomeID = 14
				}
			}

			if World[x][y].precip < 0.10 {
				World[x][y].biomeID = 4
				if World[x][y].drainage > 0.5 {
					World[x][y].biomeID = 16
					if rand.Intn(2)+1 == 2 {
						World[x][y].biomeID = 14
					}
				}
				if World[x][y].drainage >= 0.66 {
					World[x][y].biomeID = 8
				}
			}

			if World[x][y].height <= 0.2 {
				World[x][y].biomeID = 0
			}

			if World[x][y].temp <= 0.2 && World[x][y].height > 0.15 {
				World[x][y].biomeID = rand.Intn(13-11) + 11
			}

			if World[x][y].height > 0.6 {
				World[x][y].biomeID = 9
			}
			if World[x][y].height > 0.9 {
				World[x][y].biomeID = 10
			}
			//spew.Dump(World[x][y])
			biomeCount[World[x][y].biomeID]++
		}
	}

	spew.Dump(biomeCount)
	fmt.Println("- BiomeIDs Atributed -")

	//River Gen

	for x := 0; x < 1; x++ {
		RiverGen(World)
	}
	fmt.Println("- River Gen -")

	//Free Heightmaps
	//libtcod.HeightmapDelete(hm)
	//libtcod.HeightmapDelete(temp)
	//libtcod.HeightmapDelete(noisehm)

	fmt.Println(" * Biomes/Rivers Sorted *")

	return World
}

func ReadRaces() []Race {
	RacesFile := "Races.txt"

	file, err := os.Open(RacesFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var Races []Race

	for scanner.Scan() {
		Info := make([]string, 7)
		for y := 0; y < 7; y++ {
			data := scanner.Text()
			start := strings.Index(data, "]") + 1
			end := len(data)
			Info[y] = data[start:end]
			if y < 6 {
				scanner.Scan()
			}
		}
		spew.Dump(Info)
		PreferedBiomesStr := strings.Fields(Info[1])
		PreferedBiomes := make([]int, len(PreferedBiomesStr))
		for i, v := range PreferedBiomesStr {
			PreferedBiomes[i], _ = strconv.Atoi(v)
		}
		log.Println("Prefered Biomes:")
		log.Println(PreferedBiomesStr)
		log.Println(PreferedBiomes)
		race := Race{Info[0], PreferedBiomes, atoi(Info[2]), atoi(Info[3]), atoi(Info[4]), atoi(Info[5]), Info[6]}
		spew.Dump(race)
		Races = append(Races, race)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("- Races Read -")
	spew.Dump(Races)
	return Races
}

func ReadGovern() []GovernmentType {
	GovernFile := "CivilizedGovernment.txt"
	NLines := countLines(GovernFile)
	NGovern := NLines / 5

	f, err := os.Open(GovernFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var Governs []GovernmentType
	for x := 0; x < NGovern; x++ {
		Info := make([]string, 5)
		for y := 0; y < 5; y++ {
			scanner.Scan()
			data := scanner.Text()
			start := strings.Index(data, "]") + 1
			end := len(data)
			Info[y] = data[start:end]
		}
		Governs = append(Governs, GovernmentType{Info[0], Info[1], atoi(Info[2]), atoi(Info[3]), atoi(Info[4])})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("- Government Types Read -")
	return Governs
}

func atoi(s string) float64 {
	i, _ := strconv.Atoi(s)
	return float64(i)
}

/*
func countLines(filename string) int {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}*/

func CivGen(Races []Race, Govern []GovernmentType) []Civ {
	Civs := []Civ{}

	for x := 0; x < CIVILIZED_CIVS; x++ {
		log.Println("Generating Civs - ", x, "/", CIVILIZED_CIVS, " - ", float64(x)/float64(CIVILIZED_CIVS)*100, "%")
		Name := NamegenGenerate("Fantasy male") + " Civilization"

		var race Race
		for _, i := range rand.Perm(len(Races)) {
			race = Races[i]
			if race.Form == "civilized" {
				break
			}
		}
		if race.Form == "" {
			log.Println("No civilized Found")
			break
		}

		Government := Govern[rand.Intn(len(Govern))]

		Color := Palette[rand.Intn(len(Palette))]

		Flag := FlagGenerator(Color)

		//Initialize Civ
		civ := NewCiv(Name, &race, Government, Color, Flag, 0)
		Civs = append(Civs, *civ)
	}

	for a := 0; a < TRIBAL_CIVS; a++ {
		log.Println("Generating Civs - ", a+CIVILIZED_CIVS, "/", CIVILIZED_CIVS+TRIBAL_CIVS, " - ", float64(a+CIVILIZED_CIVS)/float64(CIVILIZED_CIVS+TRIBAL_CIVS)*100, "%")

		Name := NamegenGenerate("Fantasy male") + " Tribe"

		var race Race
		for _, i := range rand.Perm(len(Races)) {
			race = Races[i]
			if race.Form == "tribal" {
				break
			}
		}
		if race.Form == "" {
			log.Println("No Tribes Found")
			break
		}

		Government := GovernmentType{"Tribal", "*PLACE HOLDER*", 2, 50, 0}

		Color := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256))}

		Flag := FlagGenerator(Color)

		//Initialize Civ
		civ := NewCiv(Name, &race, Government, Color, Flag, 0)
		Civs = append(Civs, *civ)
	}

	fmt.Println("- Civs Generated -")

	return Civs
}

func SetupCivs(Civs []Civ, World [][]Tile, Chars [][]int, Colors [][]color.RGBA) []Civ {
	for x := 0; x < len(Civs); x++ {
		civ := &Civs[x]
		civ.Sites = []CivSite{}
		civ.SuitableSites = []CivSite{}

		for i := 0; i < WORLD_WIDTH; i++ {
			for j := 0; j < WORLD_HEIGHT; j++ {
				for _, biome := range civ.Race.PrefBiome {
					if World[i][j].biomeID == biome {
						civ.SuitableSites = append(civ.SuitableSites, CivSite{x: i, y: j, category: "", suitable: true, popcap: 0})
					}
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		if len(civ.SuitableSites) == 0 {
			fmt.Println("No Suitable Sites for Civ ", x)
			continue
		}
		randIndex := rand.Intn(len(civ.SuitableSites))
		for World[civ.SuitableSites[randIndex].x][civ.SuitableSites[randIndex].y].isCiv == true {
			civ.SuitableSites = append(civ.SuitableSites[:randIndex], civ.SuitableSites[randIndex+1:]...)
			randIndex = rand.Intn(len(civ.SuitableSites))
		}

		X := civ.SuitableSites[randIndex].x
		Y := civ.SuitableSites[randIndex].y

		World[X][Y].isCiv = true

		FinalProsperity := World[X][Y].prosperity * 150
		if World[X][Y].hasRiver {
			FinalProsperity = FinalProsperity * 1.5
		}
		PopCap := 4*civ.Race.ReproductionSpeed + FinalProsperity
		PopCap = PopCap * 2 //Capital Bonus
		PopCap = math.Round(float64(PopCap))

		civ.Sites = append(civ.Sites, CivSite{x: X, y: Y, category: "Village", suitable: false, popcap: int(PopCap)})

		civ.Sites[0].isCapital = true

		civ.Sites[0].Population = 20

		Chars[X][Y] = 31
		Colors[X][Y] = civ.Color

		civ.PrintInfo()
	}

	fmt.Println("- Civs Setup -")
	fmt.Println(" * Civ Gen DONE *")

	return Civs
}

// ##################################################################################### - PROCESS CIVS - ##################################################################################
func NewSite(Civ Civ, Origin CivSite, World [][]Tile, Chars [][]int, Colors [][]color.RGBA) Civ {
	rand.Seed(time.Now().UnixNano())
	randIndex := rand.Intn(len(Civ.SuitableSites))

	Tries := 0

	for PointDistRound(Origin.x, Origin.y, Civ.SuitableSites[randIndex].x, Civ.SuitableSites[randIndex].y) > EXPANSION_DISTANCE || World[Civ.SuitableSites[randIndex].x][Civ.SuitableSites[randIndex].y].isCiv == true {
		if Tries > 200 {
			return Civ
		}
		Tries += 1
		randIndex = rand.Intn(len(Civ.SuitableSites))
	}

	X := Civ.SuitableSites[randIndex].x
	Y := Civ.SuitableSites[randIndex].y

	World[X][Y].isCiv = true

	FinalProsperity := World[X][Y].prosperity * 150
	if World[X][Y].hasRiver {
		FinalProsperity = FinalProsperity * 1.5
	}
	PopCap := 3*Civ.Race.ReproductionSpeed + FinalProsperity
	PopCap = math.Round(float64(PopCap))

	Civ.Sites = append(Civ.Sites, CivSite{x: X, y: Y, category: "Village", suitable: false, popcap: int(PopCap)})

	Civ.Sites[len(Civ.Sites)-1].Population = 20

	Chars[X][Y] = 31
	Colors[X][Y] = Civ.Color

	needUpdate = true

	return Civ
}

func ProcessCivs(World [][]Tile, Civs []Civ, Chars [][]int, Colors [][]color.RGBA, Month int) {
	fmt.Println("------------------------------------------")
	for x := 0; x < len(Civs); x++ {
		fmt.Println(Civs[x].Name)
		fmt.Println(Civs[x].Race.Name)

		Civs[x].TotalPopulation = 0

		// Site
		for y := 0; y < len(Civs[x].Sites); y++ {

			// Population
			NewPop := int(math.Round(float64(Civs[x].Sites[y].Population) * Civs[x].Race.ReproductionSpeed / 1500))

			if Civs[x].Sites[y].Population > Civs[x].Sites[y].popcap/2 {
				NewPop /= 6
			}

			Civs[x].Sites[y].Population += NewPop

			// Expand
			if Civs[x].Sites[y].Population > Civs[x].Sites[y].popcap {
				Civs[x].Sites[y].Population = int(math.Round(float64(Civs[x].Sites[y].popcap)))
				if len(Civs[x].Sites) < CIV_MAX_SITES {
					Civs[x].Sites[y].Population = int(math.Round(float64(Civs[x].Sites[y].popcap) / 2))
					Civs[x] = NewSite(Civs[x], Civs[x].Sites[y], World, Chars, Colors)
				}
			}

			Civs[x].TotalPopulation += Civs[x].Sites[y].Population

			// Diplomacy
			for a := 0; a < CIVILIZED_CIVS+TRIBAL_CIVS; a++ {
				for b := 0; b < len(Civs[a].Sites); b++ {
					if x == a {
						break
					}
					if PointDistRound(Civs[x].Sites[y].x, Civs[x].Sites[y].y, Civs[a].Sites[b].x, Civs[a].Sites[b].y) < WAR_DISTANCE {
						AlreadyWar := false
						for c := 0; c < len(Wars); c++ {
							if (Wars[c].Side1 == &Civs[x] && Wars[c].Side2 == &Civs[a]) || (Wars[c].Side1 == &Civs[a] && Wars[c].Side2 == &Civs[x]) {
								// Already at War
								AlreadyWar = true
							}
						}
						if AlreadyWar == false {
							// Start War and form armies if dot have army yet
							Wars = append(Wars, War{Side1: &Civs[x], Side2: &Civs[a]})
							if Civs[a].atWar == false { // if not already at war form new army
								Civs[a].Army = &Army{
									x:    Civs[a].Sites[0].x,
									y:    Civs[a].Sites[0].y,
									Civ:  &Civs[a],
									Size: int(float64(Civs[a].TotalPopulation) * Civs[a].Government.Militarization / 100)}
								Civs[a].atWar = true
							}
							if Civs[x].atWar == false { // if not already at war form new army
								Civs[x].Army = &Army{
									x:    Civs[x].Sites[0].x,
									y:    Civs[x].Sites[0].y,
									Civ:  &Civs[x],
									Size: int(float64(Civs[x].TotalPopulation) * Civs[x].Government.Militarization / 100)}
								Civs[x].atWar = true
							}
						}
					}
				}
			}

			fmt.Println("X:", Civs[x].Sites[y].x, "Y:", Civs[x].Sites[y].y, "Population:", Civs[x].Sites[y].Population)
		}

		if Civs[x].Army != nil {
			fmt.Println(Civs[x].Army.x, Civs[x].Army.y, Civs[x].Army.Size, "\n")
		}
	}
}

//####################################################################################### - MAP MODES - ####################################################################################

// # --------------------------------------------------------------------------------- Print Map (Terrain) --------------------------------------------------------------------------------
func TerrainMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			hm_v := World[x][y].height
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '0', libtcodBlue, libtcodBlack)
			if hm_v > 0.1 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '1', libtcodBlue, libtcodBlack)
			}
			if hm_v > 0.2 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '2', Palette[0], libtcodBlack)
			}
			if hm_v > 0.3 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '3', Palette[0], libtcodBlack)
			}
			if hm_v > 0.4 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '4', Palette[0], libtcodBlack)
			}
			if hm_v > 0.5 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '5', Palette[0], libtcodBlack)
			}
			if hm_v > 0.6 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '6', Palette[0], libtcodBlack)
			}
			if hm_v > 0.7 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '7', Palette[0], libtcodBlack)
			}
			if hm_v > 0.8 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '8', libtcodDarkSepia, libtcodBlack)
			}
			if hm_v > 0.9 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '9', libtcodLightGray, libtcodBlack)
			}
			if hm_v > 0.99 {
				libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '^', libtcodDarkerGray, libtcodBlack)
			}
		}
	}
	libtcod.ConsoleFlush()
}

func BiomeMap(Chars [][]int, Colors [][]color.RGBA) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, rune(Chars[x][y]), Colors[x][y], libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func HeightGradMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			hm_v := World[x][y].height
			HeightColor := color.RGBA{255, 255, 255, 255}
			colorSetHSV(&HeightColor, 0, 0, hm_v) // Set lightness to hm_v so higher heightmap value -> "whiter"
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '\u00DE', HeightColor, libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func TempGradMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			tempv := World[x][y].temp
			tempcolor := colorLerp(libtcodWhite, libtcodRed, tempv)
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '\u00DE', tempcolor, libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func PrecipGradMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			precipv := World[x][y].precip
			precipcolor := colorLerp(libtcodWhite, libtcodLightBlue, precipv)
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '\u00DE', precipcolor, libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func DrainageGradMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			drainv := World[x][y].drainage
			draincolor := colorLerp(libtcodDarkestOrange, libtcodWhite, drainv)
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '\u00DE', draincolor, libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func ProsperityGradMap(World [][]Tile) {
	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			prosperitynv := World[x][y].prosperity
			if World[x][y].hasRiver {
				prosperitynv *= 1.5
			}
			prosperitycolor := colorLerp(libtcodBlack, libtcodDarkerGreen, prosperitynv)
			libtcod.ConsolePutCharEx(0, x, y+SCREEN_HEIGHT/2-WORLD_HEIGHT/2, '\u0333', prosperitycolor, libtcodBlack)
		}
	}
	libtcod.ConsoleFlush()
}

func NormalMap(World [][]Tile) ([][]int, [][]color.RGBA) {
	Chars := make([][]int, WORLD_WIDTH)
	Colors := make([][]color.RGBA, WORLD_WIDTH)

	for i := range Chars {
		Chars[i] = make([]int, WORLD_HEIGHT)
		Colors[i] = make([]color.RGBA, WORLD_HEIGHT)
	}

	SymbolDictionary := func(x int) rune {
		char := rune(0)
		if x == 15 || x == 8 {
			if rand.Intn(2)+1 == 2 {
				char = 251
			} else {
				char = ','
			}
		}
		if x == 1 {
			if rand.Intn(2)+1 == 2 {
				char = 244
			} else {
				char = 131
			}
		}
		if x == 2 {
			if rand.Intn(2)+1 == 2 {
				char = '"'
			} else {
				char = 163
			}
		}
		return map[int]rune{
			0:  '\u00FB',
			1:  char,
			2:  char,
			3:  'n',
			4:  '\u00FB',
			5:  24,
			6:  rune(6 - rand.Intn(2)),
			8:  char,
			9:  127,
			10: 30,
			11: 176,
			12: 177,
			13: 178,
			14: 'n',
			15: char,
			16: 139,
		}[x]
	}

	ColorDictionary := func(x int) color.RGBA {
		badlands := color.RGBA{204, 159, 81, 255}
		icecolor := color.RGBA{176, 223, 215, 255}
		darkgreen := color.RGBA{68, 158, 53, 255}
		lightgreen := color.RGBA{131, 212, 82, 255}
		water := color.RGBA{13, 103, 196, 255}
		mountain := color.RGBA{185, 192, 162, 255}
		desert := color.RGBA{255, 218, 90, 255}
		return map[int]color.RGBA{
			0:  water,
			1:  darkgreen,
			2:  lightgreen,
			3:  lightgreen,
			4:  desert,
			5:  darkgreen,
			6:  darkgreen,
			8:  badlands,
			9:  mountain,
			10: mountain,
			11: icecolor,
			12: icecolor,
			13: icecolor,
			14: darkgreen,
			15: lightgreen,
			16: darkgreen,
		}[x]
	}

	for x := 0; x < WORLD_WIDTH; x++ {
		for y := 0; y < WORLD_HEIGHT; y++ {
			Chars[x][y] = int(SymbolDictionary(World[x][y].biomeID))
			Colors[x][y] = ColorDictionary(World[x][y].biomeID)
			if World[x][y].hasRiver {
				Chars[x][y] = 'o'
				Colors[x][y] = libtcodLightBlue
			}
		}
	}

	return Chars, Colors
}
