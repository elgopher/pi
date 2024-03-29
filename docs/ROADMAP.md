## Roadmap

* [x] Present game on the screen
    * [x] ability to change the resolution and palette
    * [x] sprite-sheet loader
    * [ ] more options: full screen, specifying tps and scale
    * [x] Game loop
* [ ] Implement Graphics API
    * [x] drawing sprites and pixels with camera and clipping support
    * [x] add the ability to directly access pixels on the screen and sprite-sheet
    * [x] palette transparency
    * [x] palette swapping
    * [ ] printing text on the screen
        * [x] system font
        * [x] custom font 
        * [ ] additional features: escape characters, offsets
    * [ ] drawing shapes
        * [x] rectangles, lines, circles
        * [ ] fill patterns
    * [ ] stretching sprites
* [ ] Audio
    * [x] playing sound effects
    * [x] updating sound effects and music when the game is run
    * [ ] SFX notes effects
    * [ ] SFX modifiers
    * [ ] ADSR (envelopes)
    * [ ] playing music
* [ ] map API
* [ ] math API
    * [x] Cos, Sin, Atan2
    * [x] Min, Max, Mid for integers
    * [x] Mid for float64
* [x] Game controller support: gamepad and keyboard
* [x] Mouse support
  * [ ] Add mouse wheel support
* [x] Full keyboard support
* [x] Storing game state like savegames, hall of fame and user preferences
* [ ] Menu screen
  * [ ] controller mapping editor
  * [ ] keyboard mapping editor 
* [ ] Development tools
    * [x] controlling the game
        * [x] pausing, resuming
        * [x] add a programmatic way to stop the game
    * [x] taking screenshots
    * [x] screen inspector
      * [x] presenting pixel coords and color
      * [ ] zoom-in, zoom-out
      * [ ] drawing on the screen using Pi functions
        * [x] Set, shapes
        * [x] Spr
        * [ ] Print
    * [x] **scripting/REPL** - write Go code **live** when the game is running
    * [ ] palette inspector
      * [ ] display, draw palette
    * [ ] sprite-sheet editor
    * [ ] map editor
    * [ ] sound editor
    * [ ] music editor
* [ ] Documentation
    * [x] Go docs
    * [ ] More detailed Go docs
    * [ ] Tutorials
* [ ] Support for different platforms
  * [x] Windows, Linux, macOS
  * [x] Web browsers (WASM)
  * [ ] Android, IOS, Switch
* [ ] Examples
    * [ ] simple programs for beginners
    * [ ] interactive programs describing how functions work
    * [ ] simple working game