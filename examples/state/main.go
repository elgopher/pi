// Example showing how to save and load the state.
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/elgopher/pi/state"
)

// Savegame is a struct which will be stored permanently.
// Please note that all fields to be stored must be public.
type Savegame struct {
	// headers
	Version int
	Date    time.Time

	// payload
	PlayerPosX int
	PlayerPosY int
}

func main() {
	var err error
	var lang string

	// start with loading a single value
	if err = state.Load("language", &lang); errors.Is(err, state.ErrNotFound) {
		// language was not found, use default language
		lang = "EN"
		// save default language
		if err = state.Save("language", lang); err != nil {
			panic(err)
		}
	}
	// there could be some other error returned by Load
	if err != nil {
		panic(err)
	}

	fmt.Println(lang) // will print "EN"

	// Now store a savegame.
	// Savegame is a struct, having multiple fields. All public fields will be stored.
	saveGame := Savegame{
		Version: 1,
		Date:    time.Now(),
	}

	// Store game state for slot-0.
	if err = state.Save("slot-0", saveGame); err != nil {
		panic(err)
	}

	// load existing save game
	if err = state.Load("slot-0", &saveGame); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", saveGame)

	// return all state names
	names, err := state.All()
	if err != nil {
		panic(err)
	}

	// remove all saved states
	for _, name := range names {
		if err = state.Delete(name); err != nil {
			panic(err)
		}
	}
}
