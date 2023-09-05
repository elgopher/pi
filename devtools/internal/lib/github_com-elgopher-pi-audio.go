// Code generated by 'yaegi extract github.com/elgopher/pi/audio'. DO NOT EDIT.

package lib

import (
	"github.com/elgopher/pi/audio"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["github.com/elgopher/pi/audio/audio"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"EffectArpFast":       reflect.ValueOf(audio.EffectArpFast),
		"EffectArpSlow":       reflect.ValueOf(audio.EffectArpSlow),
		"EffectDrop":          reflect.ValueOf(audio.EffectDrop),
		"EffectFadeIn":        reflect.ValueOf(audio.EffectFadeIn),
		"EffectFadeOut":       reflect.ValueOf(audio.EffectFadeOut),
		"EffectNoEffect":      reflect.ValueOf(audio.EffectNoEffect),
		"EffectSlide":         reflect.ValueOf(audio.EffectSlide),
		"EffectVibrato":       reflect.ValueOf(audio.EffectVibrato),
		"GetStat":             reflect.ValueOf(audio.GetStat),
		"InstrumentNoise":     reflect.ValueOf(audio.InstrumentNoise),
		"InstrumentOrgan":     reflect.ValueOf(audio.InstrumentOrgan),
		"InstrumentPhaser":    reflect.ValueOf(audio.InstrumentPhaser),
		"InstrumentPulse":     reflect.ValueOf(audio.InstrumentPulse),
		"InstrumentSaw":       reflect.ValueOf(audio.InstrumentSaw),
		"InstrumentSfx0":      reflect.ValueOf(audio.InstrumentSfx0),
		"InstrumentSfx1":      reflect.ValueOf(audio.InstrumentSfx1),
		"InstrumentSfx2":      reflect.ValueOf(audio.InstrumentSfx2),
		"InstrumentSfx3":      reflect.ValueOf(audio.InstrumentSfx3),
		"InstrumentSfx4":      reflect.ValueOf(audio.InstrumentSfx4),
		"InstrumentSfx5":      reflect.ValueOf(audio.InstrumentSfx5),
		"InstrumentSfx6":      reflect.ValueOf(audio.InstrumentSfx6),
		"InstrumentSfx7":      reflect.ValueOf(audio.InstrumentSfx7),
		"InstrumentSquare":    reflect.ValueOf(audio.InstrumentSquare),
		"InstrumentTiltedSaw": reflect.ValueOf(audio.InstrumentTiltedSaw),
		"InstrumentTriangle":  reflect.ValueOf(audio.InstrumentTriangle),
		"LoadAudio":           reflect.ValueOf(audio.LoadAudio),
		"Music":               reflect.ValueOf(audio.Music),
		"Pat":                 reflect.ValueOf(&audio.Pat).Elem(),
		"PitchA0":             reflect.ValueOf(audio.PitchA0),
		"PitchA1":             reflect.ValueOf(audio.PitchA1),
		"PitchA2":             reflect.ValueOf(audio.PitchA2),
		"PitchA3":             reflect.ValueOf(audio.PitchA3),
		"PitchA4":             reflect.ValueOf(audio.PitchA4),
		"PitchAs0":            reflect.ValueOf(audio.PitchAs0),
		"PitchAs1":            reflect.ValueOf(audio.PitchAs1),
		"PitchAs2":            reflect.ValueOf(audio.PitchAs2),
		"PitchAs3":            reflect.ValueOf(audio.PitchAs3),
		"PitchAs4":            reflect.ValueOf(audio.PitchAs4),
		"PitchB0":             reflect.ValueOf(audio.PitchB0),
		"PitchB1":             reflect.ValueOf(audio.PitchB1),
		"PitchB2":             reflect.ValueOf(audio.PitchB2),
		"PitchB3":             reflect.ValueOf(audio.PitchB3),
		"PitchB4":             reflect.ValueOf(audio.PitchB4),
		"PitchC0":             reflect.ValueOf(audio.PitchC0),
		"PitchC1":             reflect.ValueOf(audio.PitchC1),
		"PitchC2":             reflect.ValueOf(audio.PitchC2),
		"PitchC3":             reflect.ValueOf(audio.PitchC3),
		"PitchC4":             reflect.ValueOf(audio.PitchC4),
		"PitchC5":             reflect.ValueOf(audio.PitchC5),
		"PitchCs0":            reflect.ValueOf(audio.PitchCs0),
		"PitchCs1":            reflect.ValueOf(audio.PitchCs1),
		"PitchCs2":            reflect.ValueOf(audio.PitchCs2),
		"PitchCs3":            reflect.ValueOf(audio.PitchCs3),
		"PitchCs4":            reflect.ValueOf(audio.PitchCs4),
		"PitchCs5":            reflect.ValueOf(audio.PitchCs5),
		"PitchD0":             reflect.ValueOf(audio.PitchD0),
		"PitchD1":             reflect.ValueOf(audio.PitchD1),
		"PitchD2":             reflect.ValueOf(audio.PitchD2),
		"PitchD3":             reflect.ValueOf(audio.PitchD3),
		"PitchD4":             reflect.ValueOf(audio.PitchD4),
		"PitchD5":             reflect.ValueOf(audio.PitchD5),
		"PitchDs0":            reflect.ValueOf(audio.PitchDs0),
		"PitchDs1":            reflect.ValueOf(audio.PitchDs1),
		"PitchDs2":            reflect.ValueOf(audio.PitchDs2),
		"PitchDs3":            reflect.ValueOf(audio.PitchDs3),
		"PitchDs4":            reflect.ValueOf(audio.PitchDs4),
		"PitchDs5":            reflect.ValueOf(audio.PitchDs5),
		"PitchE0":             reflect.ValueOf(audio.PitchE0),
		"PitchE1":             reflect.ValueOf(audio.PitchE1),
		"PitchE2":             reflect.ValueOf(audio.PitchE2),
		"PitchE3":             reflect.ValueOf(audio.PitchE3),
		"PitchE4":             reflect.ValueOf(audio.PitchE4),
		"PitchF0":             reflect.ValueOf(audio.PitchF0),
		"PitchF1":             reflect.ValueOf(audio.PitchF1),
		"PitchF2":             reflect.ValueOf(audio.PitchF2),
		"PitchF3":             reflect.ValueOf(audio.PitchF3),
		"PitchF4":             reflect.ValueOf(audio.PitchF4),
		"PitchFs0":            reflect.ValueOf(audio.PitchFs0),
		"PitchFs1":            reflect.ValueOf(audio.PitchFs1),
		"PitchFs2":            reflect.ValueOf(audio.PitchFs2),
		"PitchFs3":            reflect.ValueOf(audio.PitchFs3),
		"PitchFs4":            reflect.ValueOf(audio.PitchFs4),
		"PitchG0":             reflect.ValueOf(audio.PitchG0),
		"PitchG1":             reflect.ValueOf(audio.PitchG1),
		"PitchG2":             reflect.ValueOf(audio.PitchG2),
		"PitchG3":             reflect.ValueOf(audio.PitchG3),
		"PitchG4":             reflect.ValueOf(audio.PitchG4),
		"PitchGs0":            reflect.ValueOf(audio.PitchGs0),
		"PitchGs1":            reflect.ValueOf(audio.PitchGs1),
		"PitchGs2":            reflect.ValueOf(audio.PitchGs2),
		"PitchGs3":            reflect.ValueOf(audio.PitchGs3),
		"PitchGs4":            reflect.ValueOf(audio.PitchGs4),
		"Play":                reflect.ValueOf(audio.Play),
		"SFX":                 reflect.ValueOf(&audio.SFX).Elem(),
		"SampleRate":          reflect.ValueOf(constant.MakeFromLiteral("22050", token.INT, 0)),
		"SaveAudio":           reflect.ValueOf(audio.SaveAudio),
		"SetSystem":           reflect.ValueOf(audio.SetSystem),
		"Stop":                reflect.ValueOf(audio.Stop),
		"StopChan":            reflect.ValueOf(audio.StopChan),
		"StopLoop":            reflect.ValueOf(audio.StopLoop),
		"Sync":                reflect.ValueOf(audio.Sync),
		"VolumeLoudest":       reflect.ValueOf(audio.VolumeLoudest),
		"VolumeSilence":       reflect.ValueOf(audio.VolumeSilence),

		// type definitions
		"Effect":      reflect.ValueOf((*audio.Effect)(nil)),
		"Instrument":  reflect.ValueOf((*audio.Instrument)(nil)),
		"LiveReader":  reflect.ValueOf((*audio.LiveReader)(nil)),
		"Note":        reflect.ValueOf((*audio.Note)(nil)),
		"Pattern":     reflect.ValueOf((*audio.Pattern)(nil)),
		"PatternSfx":  reflect.ValueOf((*audio.PatternSfx)(nil)),
		"Pitch":       reflect.ValueOf((*audio.Pitch)(nil)),
		"SoundEffect": reflect.ValueOf((*audio.SoundEffect)(nil)),
		"Stat":        reflect.ValueOf((*audio.Stat)(nil)),
		"Synthesizer": reflect.ValueOf((*audio.Synthesizer)(nil)),
		"System":      reflect.ValueOf((*audio.System)(nil)),
		"Volume":      reflect.ValueOf((*audio.Volume)(nil)),

		// interface wrapper definitions
		"_System": reflect.ValueOf((*_github_com_elgopher_pi_audio_System)(nil)),
	}
}

// _github_com_elgopher_pi_audio_System is an interface wrapper for System type
type _github_com_elgopher_pi_audio_System struct {
	IValue    interface{}
	WGetMusic func(patterNo int) audio.Pattern
	WGetSfx   func(sfxNo int) audio.SoundEffect
	WLoad     func(a0 []byte) error
	WMusic    func(patterNo int, fadeMs int, channelMask byte)
	WPlay     func(sfxNo int, channel int, offset int, length int)
	WSave     func() ([]byte, error)
	WSetMusic func(patternNo int, _ audio.Pattern)
	WSetSfx   func(sfxNo int, e audio.SoundEffect)
	WStat     func() audio.Stat
	WStop     func(sfxNo int)
	WStopChan func(channel int)
	WStopLoop func(channel int)
}

func (W _github_com_elgopher_pi_audio_System) GetMusic(patterNo int) audio.Pattern {
	return W.WGetMusic(patterNo)
}
func (W _github_com_elgopher_pi_audio_System) GetSfx(sfxNo int) audio.SoundEffect {
	return W.WGetSfx(sfxNo)
}
func (W _github_com_elgopher_pi_audio_System) Load(a0 []byte) error {
	return W.WLoad(a0)
}
func (W _github_com_elgopher_pi_audio_System) Music(patterNo int, fadeMs int, channelMask byte) {
	W.WMusic(patterNo, fadeMs, channelMask)
}
func (W _github_com_elgopher_pi_audio_System) Play(sfxNo int, channel int, offset int, length int) {
	W.WPlay(sfxNo, channel, offset, length)
}
func (W _github_com_elgopher_pi_audio_System) Save() ([]byte, error) {
	return W.WSave()
}
func (W _github_com_elgopher_pi_audio_System) SetMusic(patternNo int, p audio.Pattern) {
	W.WSetMusic(patternNo, p)
}
func (W _github_com_elgopher_pi_audio_System) SetSfx(sfxNo int, e audio.SoundEffect) {
	W.WSetSfx(sfxNo, e)
}
func (W _github_com_elgopher_pi_audio_System) Stat() audio.Stat {
	return W.WStat()
}
func (W _github_com_elgopher_pi_audio_System) Stop(sfxNo int) {
	W.WStop(sfxNo)
}
func (W _github_com_elgopher_pi_audio_System) StopChan(channel int) {
	W.WStopChan(channel)
}
func (W _github_com_elgopher_pi_audio_System) StopLoop(channel int) {
	W.WStopLoop(channel)
}
