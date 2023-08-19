# gameasciitiles

GameAsciiTiles is (will be) a simple library that demonstrates mixed ASCII and raster tile based roguelike game rendering.

Tiles, entities, and effects can be sprites, tiles (from a tileset), or ASCII characters. The library supports (will support) loading and rendering of multiple fonts at the same time, allowing to mix and match fonts for different purposes, or provide a transition between visually distinct areas.

![alt text](/gameasciitiles/images/rgb.png "Screenshot")

## Tile source

https://opengameart.org/content/orthographic-outdoor-tiles

## Font source

This project uses the font from https://github.com/damianvila/font-bescii Thank you for making this available!

## TODO

- [ ] Sprites / tiles
    - [X] Simple ASCII / glyph
        - [ ] Caching of pre-rendered glyphs
        - [X] Superscript
        - [ ] Subscript
    - [X] Simple sprite tiles
    - [ ] Common tile interface
- [ ] Tilesets
    - [X] Loading
    - [X] Rendering
    - [ ] Caching
    - [X] Animation
    - [ ] Manipulation
    - [ ] Tile selection based on neighbors
        - [ ] Connecting wall tiles, roads, etc.
