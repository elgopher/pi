// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"fmt"
	"syscall/js"
)

var localStorage = js.Global().Get("localStorage")

func Load(name string) (string, error) {
	val := localStorage.Call("getItem", name)
	if val.IsNull() {
		return "", ErrNotFound
	}
	return val.String(), nil
}

func Save(name string, data string) (err error) {
	defer func() {
		// recover from quota exceeded exception
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	localStorage.Call("setItem", name, data)
	return err
}

func Delete(name string) error {
	localStorage.Call("removeItem", name)
	return nil
}

func Names() ([]string, error) {
	length := localStorage.Get("length").Int()
	names := make([]string, length)
	for i := 0; i < length; i++ {
		names[i] = localStorage.Call("key", i).String()
	}
	return names, nil
}

func Cleanup() {
	localStorage.Call("clear")
}
