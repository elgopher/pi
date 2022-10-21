// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/elgopher/pi/state/internal"
)

var (
	ErrNotFound             = errors.New("state not found")    // ErrNotFound is an expected error which is returned when state is not stored.
	ErrInvalidStateName     = errors.New("invalid state name") // ErrInvalidStateName is a programmer error which is returned when state name is invalid.
	ErrNilStateOutput       = errors.New("nil state output")   // ErrNilStateOutput is a programmer error which is returned when output is nil.
	ErrStateUnmarshalFailed = errors.New("state unmarshal failed")
	ErrStateMarshalFailed   = errors.New("state marshal failed")
)

// Load reads the persistent game data with specified name. Data will be stored in out param.
// The type of out should be compatible with type used during Save. For example, trying to load
// stored string into an int will return the ErrStateUnmarshalFailed.
//
// ErrNotFound error is returned when state does not exist. Please check the error with
// following code:
//
//	if errors.Is(err, pi.ErrNotFound) { ... }
//
// ErrNilStateOutput is returned when out is nil.
func Load[T any](name string, out *T) error {
	if err := validateStateName(name); err != nil {
		return err
	}

	if out == nil {
		return ErrNilStateOutput
	}

	str, err := internal.Load(name)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("error loading persistent data for name %s: %w", name, err)
	}

	err = json.Unmarshal([]byte(str), out)
	if err != nil {
		return fmt.Errorf("%w for state %s: %s", ErrStateUnmarshalFailed, name, err.Error())
	}

	return nil
}

// Save permanently stores the data with given name. Data could be loaded using Load after restarting the game.
// Data can be of any type: string, int, time.Time, slice, map or a struct. Only struct public fields will be stored.
//
// Please note that on some platforms there are limits for how much data could be stored.
// Web browsers for example have 5MB total data limit. So, if you want your game to be portable
// across different platforms, please limit the size of data below that number (in fact a game which stores
// more than hundreds of KBs per savegame does not look retro anymore).
// Similarly to data size limit, there is a limit for the number of states. It should be no more than tens of states.
// Generally using structs instead of primitive types like string or int could overcome this limit. For example,
// a single struct could have all the user preferences, or the entire hall of fame.
//
// Saving a state is an atomic operation. Either the entire operation is successful and state is stored, or error is reported
// and no data is stored (previous data is not updated). You can leverage this feature to store similar data in a consistent
// manner - either all changes to data are applied or no changes at all. For example, you could design a struct having all user
// preferences, or a struct dedicated for an entire savegame. Such design will highly decrease the chance of consistency problems
// in case the game/OS crashes or during power loss.
//
// Name cannot be empty, have characters "/", "\", or by longer than 32 characters
func Save(name string, data any) error {
	if err := validateStateName(name); err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("%w for state %s: %s", ErrStateMarshalFailed, name, err.Error())
	}

	if err = internal.Save(name, string(bytes)); err != nil {
		return fmt.Errorf("error saving state for %s: %w", name, err)
	}

	return nil
}

func validateStateName(name string) error {
	if name == "" {
		return fmt.Errorf("%w: empty state name", ErrInvalidStateName)
	}
	if len(name) > 32 {
		return fmt.Errorf("%w: state name longer than 32 chars", ErrInvalidStateName)
	}
	if strings.Contains(name, "\\") {
		return fmt.Errorf("%w: state name cannot have \\", ErrInvalidStateName)
	}
	if strings.Contains(name, "/") {
		return fmt.Errorf("%w: state name cannot have /", ErrInvalidStateName)
	}

	return nil
}

// Delete permanently deletes data with given name.
func Delete(name string) error {
	if err := validateStateName(name); err != nil {
		return err
	}

	return internal.Delete(name)
}

// All returns names of all states.
func All() ([]string, error) {
	return internal.Names()
}
