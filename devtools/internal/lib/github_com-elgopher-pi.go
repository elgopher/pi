// Code generated by 'yaegi extract github.com/elgopher/pi'. DO NOT EDIT.

package lib

import (
	"go/constant"
	"go/token"
	"reflect"

	"github.com/elgopher/pi"
)

func init() {
	Symbols["github.com/elgopher/pi/pi"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Atan2":                     reflect.ValueOf(pi.Atan2),
		"Audio":                     reflect.ValueOf(pi.Audio),
		"Btn":                       reflect.ValueOf(pi.Btn),
		"BtnBits":                   reflect.ValueOf(pi.BtnBits),
		"BtnPlayer":                 reflect.ValueOf(pi.BtnPlayer),
		"Btnp":                      reflect.ValueOf(pi.Btnp),
		"BtnpBits":                  reflect.ValueOf(pi.BtnpBits),
		"BtnpPlayer":                reflect.ValueOf(pi.BtnpPlayer),
		"Camera":                    reflect.ValueOf(pi.Camera),
		"CameraReset":               reflect.ValueOf(pi.CameraReset),
		"Circ":                      reflect.ValueOf(pi.Circ),
		"CircFill":                  reflect.ValueOf(pi.CircFill),
		"Clip":                      reflect.ValueOf(pi.Clip),
		"ClipReset":                 reflect.ValueOf(pi.ClipReset),
		"Cls":                       reflect.ValueOf(pi.Cls),
		"ClsCol":                    reflect.ValueOf(pi.ClsCol),
		"ColorTransparency":         reflect.ValueOf(&pi.ColorTransparency).Elem(),
		"Controllers":               reflect.ValueOf(&pi.Controllers).Elem(),
		"Cos":                       reflect.ValueOf(pi.Cos),
		"CustomFont":                reflect.ValueOf(pi.CustomFont),
		"DisplayPalette":            reflect.ValueOf(&pi.DisplayPalette).Elem(),
		"Down":                      reflect.ValueOf(pi.Down),
		"Draw":                      reflect.ValueOf(&pi.Draw).Elem(),
		"DrawPalette":               reflect.ValueOf(&pi.DrawPalette).Elem(),
		"GameLoopStopped":           reflect.ValueOf(&pi.GameLoopStopped).Elem(),
		"Left":                      reflect.ValueOf(pi.Left),
		"Line":                      reflect.ValueOf(pi.Line),
		"Load":                      reflect.ValueOf(pi.Load),
		"MaxInt":                    reflect.ValueOf(pi.MaxInt[int]), // TODO Generic functions not supported by Yaegi yet
		"Mid":                       reflect.ValueOf(pi.Mid),
		"MidInt":                    reflect.ValueOf(pi.MidInt[int]), // TODO Generic functions not supported by Yaegi yet
		"MinInt":                    reflect.ValueOf(pi.MinInt[int]), // TODO Generic functions not supported by Yaegi yet
		"MouseBtn":                  reflect.ValueOf(pi.MouseBtn),
		"MouseBtnDuration":          reflect.ValueOf(&pi.MouseBtnDuration).Elem(),
		"MouseBtnp":                 reflect.ValueOf(pi.MouseBtnp),
		"MouseLeft":                 reflect.ValueOf(pi.MouseLeft),
		"MouseMiddle":               reflect.ValueOf(pi.MouseMiddle),
		"MousePos":                  reflect.ValueOf(pi.MousePos),
		"MousePosition":             reflect.ValueOf(&pi.MousePosition).Elem(),
		"MouseRight":                reflect.ValueOf(pi.MouseRight),
		"NewPixMap":                 reflect.ValueOf(pi.NewPixMap),
		"NewPixMapWithPixels":       reflect.ValueOf(pi.NewPixMapWithPixels),
		"O":                         reflect.ValueOf(pi.O),
		"Pal":                       reflect.ValueOf(pi.Pal),
		"PalDisplay":                reflect.ValueOf(pi.PalDisplay),
		"PalReset":                  reflect.ValueOf(pi.PalReset),
		"Palette":                   reflect.ValueOf(&pi.Palette).Elem(),
		"Palt":                      reflect.ValueOf(pi.Palt),
		"PaltReset":                 reflect.ValueOf(pi.PaltReset),
		"Pget":                      reflect.ValueOf(pi.Pget),
		"Print":                     reflect.ValueOf(pi.Print),
		"Pset":                      reflect.ValueOf(pi.Pset),
		"Rect":                      reflect.ValueOf(pi.Rect),
		"RectFill":                  reflect.ValueOf(pi.RectFill),
		"Reset":                     reflect.ValueOf(pi.Reset),
		"Right":                     reflect.ValueOf(pi.Right),
		"Scr":                       reflect.ValueOf(pi.Scr),
		"ScreenCamera":              reflect.ValueOf(&pi.ScreenCamera).Elem(),
		"SetCustomFontHeight":       reflect.ValueOf(pi.SetCustomFontHeight),
		"SetCustomFontSpecialWidth": reflect.ValueOf(pi.SetCustomFontSpecialWidth),
		"SetCustomFontWidth":        reflect.ValueOf(pi.SetCustomFontWidth),
		"SetScreenSize":             reflect.ValueOf(pi.SetScreenSize),
		"Sget":                      reflect.ValueOf(pi.Sget),
		"Sin":                       reflect.ValueOf(pi.Sin),
		"Spr":                       reflect.ValueOf(pi.Spr),
		"SprSheet":                  reflect.ValueOf(pi.SprSheet),
		"SprSize":                   reflect.ValueOf(pi.SprSize),
		"SprSizeFlip":               reflect.ValueOf(pi.SprSizeFlip),
		"SpriteHeight":              reflect.ValueOf(constant.MakeFromLiteral("8", token.INT, 0)),
		"SpriteWidth":               reflect.ValueOf(constant.MakeFromLiteral("8", token.INT, 0)),
		"Sset":                      reflect.ValueOf(pi.Sset),
		"Stop":                      reflect.ValueOf(pi.Stop),
		"SystemFont":                reflect.ValueOf(pi.SystemFont),
		"Time":                      reflect.ValueOf(&pi.Time).Elem(),
		"Up":                        reflect.ValueOf(pi.Up),
		"Update":                    reflect.ValueOf(&pi.Update).Elem(),
		"UseEmptySpriteSheet":       reflect.ValueOf(pi.UseEmptySpriteSheet),
		"X":                         reflect.ValueOf(pi.X),

		// type definitions
		"AudioSystem": reflect.ValueOf((*pi.AudioSystem)(nil)),
		"Button":      reflect.ValueOf((*pi.Button)(nil)),
		"Controller":  reflect.ValueOf((*pi.Controller)(nil)),
		"Font":        reflect.ValueOf((*pi.Font)(nil)),
		//"Int":         reflect.ValueOf((*pi.Int)(nil)),  // TODO Generic constraints not supported by Yaegi yet
		"MouseButton": reflect.ValueOf((*pi.MouseButton)(nil)),
		"PixMap":      reflect.ValueOf((*pi.PixMap)(nil)),
		"Pointer":     reflect.ValueOf((*pi.Pointer)(nil)),
		"Position":    reflect.ValueOf((*pi.Position)(nil)),
		"Region":      reflect.ValueOf((*pi.Region)(nil)),

		// interface wrapper definitions
		"_Int": reflect.ValueOf((*_github_com_elgopher_pi_Int)(nil)),
	}
}

// _github_com_elgopher_pi_Int is an interface wrapper for Int type
type _github_com_elgopher_pi_Int struct {
	IValue interface{}
}
