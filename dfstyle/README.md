# DFStyle

NOTE: This is a crude port of https://github.com/Dozed12/df-style-worldgen to Go :)

All credit goes to the original author!

This modified version uses github.com/BigJk/ramen and github.com/hajimehoshi/ebiten for rendering. There are also some really ugly ports of parts of libtcod in there somewhere... Which I probably should tidy up.

This is in really rough shape right now, but I will tidy things up as I go.

## TODO

- Fix river generation
- Add wars and war resolution
- Add record of history

![alt text](/dfstyle/images/screen.png "Screenshot")

# What is it?

df-style-worldgen is a 2D fantasy world generator inspired by Dwarf Fortress. It generates 2D worlds with multiple map modes and eventually simulate civilizations, gods and beasts and their history.


# Instructions

After opening pyWorld it generates a new map for you and displays the Biome Map Mode

Keys:

- r - Generate brand new world
- b - Display Biome Map Mode
- p - Display Precipitation Map Mode
- d - Display Drainage Map Mode
- w - Display Temperature Map Mode
- h - Display Altitude Map Mode
- t - Display Obsolete Terrain Map Mode
- f - Display Prosperity Map Mode

Simulation:

Currently only Civ expansion is generated. It's visible in the map by ▼ chars. More information is visible on the side console that displays Civ name, Civ race, Civ government form and each month display all sites and their population.

- SPACE to start/pause simulation
