// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mem

// KeyDuration has info how many frames in a row a keyboard button was pressed:
// Index of array is equal to key button constant, for example
// KeyDuration[1] has button duration for Shift key.
var KeyDuration [97]uint

const (
	KeyShift        = 1
	KeyCtrl         = 3
	KeyAlt          = 5 // Please note that on some keyboard layouts on Windows the right alt is a combination of Ctrl+Alt
	KeyCap          = 7
	KeyBack         = '\b' // 8
	KeyTab          = '\t' // 9
	KeyEnter        = '\n' // 10
	KeyF1           = 11
	KeyF2           = 12
	KeyF3           = 13
	KeyF4           = 14
	KeyF5           = 15
	KeyF6           = 16
	KeyF7           = 17
	KeyF8           = 18
	KeyF9           = 19
	KeyF10          = 20
	KeyF11          = 21
	KeyF12          = 22
	KeyLeft         = 23
	KeyRight        = 24
	KeyUp           = 25
	KeyDown         = 26
	KeyEsc          = 27
	KeySpace        = ' '  // 32
	KeyApostrophe   = '\'' // 39
	KeyComma        = ','  // 44
	KeyMinus        = '-'  // 45
	KeyPeriod       = '.'  // 46
	KeySlash        = '/'  // 47
	KeyDigit0       = '0'  // 48
	KeyDigit1       = '1'  // 49
	KeyDigit2       = '2'  // 50
	KeyDigit3       = '3'  // 51
	KeyDigit4       = '4'  // 52
	KeyDigit5       = '5'  // 53
	KeyDigit6       = '6'  // 54
	KeyDigit7       = '7'  // 55
	KeyDigit8       = '8'  // 56
	KeyDigit9       = '9'  // 57
	KeySemicolon    = ';'  // 59
	KeyEqual        = '='  // 61
	KeyA            = 'A'  // 65
	KeyB            = 'B'  // 66
	KeyC            = 'C'  // 67
	KeyD            = 'D'  // 68
	KeyE            = 'E'  // 69
	KeyF            = 'F'  // 70
	KeyG            = 'G'  // 71
	KeyH            = 'H'  // 72
	KeyI            = 'I'  // 73
	KeyJ            = 'J'  // 74
	KeyK            = 'K'  // 75
	KeyL            = 'L'  // 76
	KeyM            = 'M'  // 77
	KeyN            = 'N'  // 78
	KeyO            = 'O'  // 79
	KeyP            = 'P'  // 80
	KeyQ            = 'Q'  // 81
	KeyR            = 'R'  // 82
	KeyS            = 'S'  // 83
	KeyT            = 'T'  // 84
	KeyU            = 'U'  // 85
	KeyV            = 'V'  // 86
	KeyW            = 'W'  // 87
	KeyX            = 'X'  // 88
	KeyY            = 'Y'  // 89
	KeyZ            = 'Z'  // 90
	KeyBracketLeft  = '['  // 91
	KeyBackslash    = '\\' // 92
	KeyBracketRight = ']'  // 93
	KeyBackquote    = '`'  // 96
)
