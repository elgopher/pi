# pi

The retro game development engine for Go, inspired by [Pico-8](https://www.lexaloffle.com/pico-8.php) and powered by [Ebitengine](https://ebiten.org/).

## FAQ

### Is this a new fantasy console?

No, it's not. It's rather a game development library with some additional tools (like a console) which make it simple (and fun!) to write retro games in Go.

### What is a retro game?

It's a game that resembles old 8-bit/16-bit games. This usually means:

* (extremely) Low resolution (like 128x128 pixels)
* Limited number of colors (like 16)
* Very small number of assets (like 256 sprites, map having up to 8K tiles)
* Simple rules (opposite to Paradox grand strategy games)
* Simplified music and SFX

### What similarities does Pi have with Pico-8?

* Most API function names are similar and behave the same way.
* Screen resolution is small, and the number of colors is limited. Although in Pi you can change the resolution and palette.
* You have one small sprite sheet and a map.

### Why would I use it?

Because it's the easiest way to write a game in Go. IMHO ;)

### Is Pi ready to use?

Pi is under development. Only limited functionality is provided. API is not stable. See [roadmap](#roadmap) for details.

## Roadmap

* [x] Present game on the screen
  * [x] add the ability to change the resolution and palette
  * [x] add sprite-sheet loader
  * [ ] add more options: full screen, specifying tps and scale
  * [x] Time function
  * [ ] add a programmatic way to stop the game
* [ ] Implement Graphics API
  * [x] drawing sprites and pixels with camera and clipping support
  * [x] add the ability to directly access pixels on the screen and sprite-sheet
  * [ ] palette transparency
  * [ ] palette swapping
  * [ ] printing text on the screen
  * [ ] drawing shapes
    * [ ] lines, rectangles, circles, ovals
    * [ ] add support for fill patterns
  * [ ] math API
    * [x] Cos, Sin, Atan2
    * [ ] Min, Max, Mid
  * [ ] stretching sprites
* [ ] Add keyboard support
* [ ] Add gamepad/joystick support
* [ ] Add mouse support
* [ ] Development console
  * [ ] pausing, resuming the game
  * [ ] running public functions
  * [ ] sprite-sheet editor
  * [ ] map editor
  * [ ] sound editor
  * [ ] music editor
* [ ] Documentation
  * [ ] Go docs
* [ ] Examples
  * [ ] simple programs for beginners
  * [ ] interactive programs describing how functions work
    * [ ] Sin,Cos,Atan2 visualization
