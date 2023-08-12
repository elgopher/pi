package lib

import (
	"fmt"
	"reflect"
)

func init() {
	mt := []reflect.Type{
		reflect.TypeOf((*fmt.Formatter)(nil)).Elem(),
		reflect.TypeOf((*fmt.Stringer)(nil)).Elem(),
	}

	MapTypes[reflect.ValueOf(fmt.Errorf)] = mt
	MapTypes[reflect.ValueOf(fmt.Fprint)] = mt
	MapTypes[reflect.ValueOf(fmt.Fprintf)] = mt
	MapTypes[reflect.ValueOf(fmt.Fprintln)] = mt
	MapTypes[reflect.ValueOf(fmt.Print)] = mt
	MapTypes[reflect.ValueOf(fmt.Printf)] = mt
	MapTypes[reflect.ValueOf(fmt.Println)] = mt
	MapTypes[reflect.ValueOf(fmt.Sprint)] = mt
	MapTypes[reflect.ValueOf(fmt.Sprintf)] = mt
	MapTypes[reflect.ValueOf(fmt.Sprintln)] = mt

	mt = []reflect.Type{reflect.TypeOf((*fmt.Scanner)(nil)).Elem()}

	MapTypes[reflect.ValueOf(fmt.Scan)] = mt
	MapTypes[reflect.ValueOf(fmt.Scanf)] = mt
	MapTypes[reflect.ValueOf(fmt.Scanln)] = mt
}
