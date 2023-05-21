package dfstyle

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/BigJk/ramen/concolor"
	"github.com/BigJk/ramen/t"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var Palette = []color.RGBA{
	{255, 45, 33, 255},   // Red
	{254, 80, 0, 255},    // Orange
	{0, 35, 156, 255},    // Blue
	{71, 45, 96, 255},    // Purple
	{0, 135, 199, 255},   // Ocean Blue
	{254, 221, 0, 255},   // Yellow
	{255, 255, 255, 255}, // White
	{99, 102, 106, 255},  // Gray
}
var Wars []War
var isRunning bool
var needUpdate bool

func Main() {
	// ###################################################################################### - Startup - ######################################################################################
	// Start Console and set custom font
	libtcod = newLibTcod(SCREEN_WIDTH, SCREEN_HEIGHT)
	//libtcod.ConsoleSetCustomFont("Andux_cp866ish.png", libtcod.FONT_LAYOUT_ASCII_INROW)
	//libtcod.ConsoleInitRoot(SCREEN_WIDTH, SCREEN_HEIGHT, "pyWorld", false, libtcod.RENDERER_SDL) // Set true for fullscreen

	// Palette

	// libtcod.SysSetFPS(30)
	// libtcod.ConsoleSetFullscreen(true)
	isRunning = false
	needUpdate = false

	// World Gen
	World := make([][]Tile, WORLD_WIDTH)
	for i := range World {
		World[i] = make([]Tile, WORLD_HEIGHT)
	}
	World = MasterWorldGen()
	log.Println("done world gen")

	// Normal Map Initialization
	Chars, Colors := NormalMap(World)
	log.Println("done normal map")

	// Read Races
	Races := ReadRaces()
	log.Println("done read races")

	// Read Governments
	Govern := ReadGovern()
	log.Println("done read govern")

	// Civ Gen
	Civs := make([]Civ, CIVILIZED_CIVS+TRIBAL_CIVS)
	Civs = CivGen(Races, Govern)
	log.Println("done civ gen")

	// Setup Civs
	Civs = SetupCivs(Civs, World, Chars, Colors)
	log.Println("done setup civs")

	// Print Map
	BiomeMap(Chars, Colors)
	log.Println("done biome map")

	// Month 0
	Month := 0

	// Reset Wars
	Wars = make([]War, 0)
	// Select Map Mode
	con := libtcod.con

	con.SetTickHook(func(t float64) error {
		// Update the console
		//log.Println("Tick")
		return nil
	})

	con.SetPreRenderHook(func(screen *ebiten.Image, timeDelta float64) error {
		// Simulation
		if isRunning {
			ProcessCivs(World, Civs, Chars, Colors, Month)

			// DEBUG Print Month
			Month++
			fmt.Println("Month: ", Month)

			// End Simulation
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
				//timer = 0
				isRunning = false
				fmt.Println("*PAUSED*")
				time.Sleep(time.Second)
				log.Println("done simulation")
			}

			// Flush Console
			if needUpdate {
				con.ClearAll()                                           // clear console
				con.TransformAll(t.Background(concolor.RGB(50, 50, 50))) // set the background
				BiomeMap(Chars, Colors)
				needUpdate = false
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Start Simulation
			isRunning = true
			needUpdate = true
			fmt.Println("*RUNNING*")
			time.Sleep(time.Second)
			log.Println("start simulation")
		}

		keys := inpututil.PressedKeys()
		for _, key := range keys {
			if key == ebiten.KeyT {
				TerrainMap(World)
			} else if key == ebiten.KeyH {
				HeightGradMap(World)
			} else if key == ebiten.KeyW {
				TempGradMap(World)
			} else if key == ebiten.KeyP {
				PrecipGradMap(World)
			} else if key == ebiten.KeyD {
				DrainageGradMap(World)
			} else if key == ebiten.KeyF {
				ProsperityGradMap(World)
			} else if key == ebiten.KeyB {
				BiomeMap(Chars, Colors)
			} else if key == ebiten.KeyR {
				fmt.Print("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
				fmt.Println(" * NEW WORLD *")
				Month = 0
				Wars = []War{}
				World = MasterWorldGen()
				Races := ReadRaces()
				Govern := ReadGovern()
				Civs = CivGen(Races, Govern)
				Chars, Colors = NormalMap(World)
				Civs = SetupCivs(Civs, World, Chars, Colors)
				BiomeMap(Chars, Colors)
			}
		}
		return nil
	})

	// Start the main loop
	con.Start(1)
}
