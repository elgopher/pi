# pi

The retro game engine for Go, inspired by [Pico-8](https://www.lexaloffle.com/pico-8.php) and powered by [Ebiten](https://ebiten.org/).

## FAQ

### Is this a new fantasy console?

No, it's not. It's rather a game library with dev tools console which make it simple (and fun!) to write retro (and pixelated) games in Go. 

### What similarities does it have with Pico-8?

* Most API function names are similar and behave the same way.
* Screen resolution is small, number of colors is limited. Although in Pi you can change the resolution and palette.
* You have one small sprite sheet and a map.

### Why would I use it?

Because it's the easiest way to write a game in Go. IMHO ;) 

### Is Pi ready to use?

Pi is under development. Only limited functionality is provided. API is not stable. See [roadmap](#roadmap) for details.

## Roadmap

* [x] Present game on the screen
  * [x] add ability to change the resolution and palette
  * [x] add sprite-sheet loader
  * [ ] add more options: full screen, specifying tps and scale
  * [x] Time function
  * [ ] add a programmatic way to stop the game
* [ ] Implement Graphics API
  * [x] drawing sprites and pixels with camera and clipping support
  * [x] add ability to directly access pixels and on the screen and sprite-sheet
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
* [ ] Dev tools console
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
  * [ ] interactive programs describing how functions works
    * [ ] Sin,Cos,Atan2 visualization