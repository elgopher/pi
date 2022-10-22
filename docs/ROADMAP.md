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
    * [ ] map API
* [ ] math API
    * [x] Cos, Sin, Atan2
    * [ ] Min, Max, Mid
* [x] Game controller support: gamepad and keyboard
* [x] Mouse support
  * [ ] Add mouse wheel support
* [x] Full keyboard support
* [x] Storing game state like savegames, hall of fame and user preferences
* [ ] Menu screen
  * [ ] controller mapping editor
  * [ ] keyboard mapping editor 
* [ ] Development console
    * [ ] stopping, resuming the game
        * [x] add a programmatic way to stop the game
        * [ ] resume the game using console command
    * [ ] scripting (running π functions)
    * [ ] screen inspector
    * [ ] sprite-sheet editor
    * [ ] map editor
    * [ ] sound editor
    * [ ] music editor
* [ ] Documentation
    * [ ] Go docs
* [ ] Support for different platforms
  * [x] Windows, Linux, macOS
  * [x] Web browsers (WASM)
  * [ ] Android, IOS, Switch
* [ ] Examples
    * [ ] simple programs for beginners
    * [ ] interactive programs describing how functions work
    * [ ] simple working game