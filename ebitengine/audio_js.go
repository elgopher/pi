package ebitengine

import (
	"errors"
	"fmt"
	"runtime"
	"syscall/js"
	"unsafe"

	"github.com/elgopher/pi"
)

func startAudio() (stop func(), _ error) {
	if AudioStream == nil {
		AudioStream = pi.Audio()
	}

	audio, err := newAudioContext()
	if err != nil {
		return func() {}, err
	}

	processor := audio.createScriptProcessor(0, 0, 1)
	bufferSize := processor.bufferSize()

	browserBuffer := make([]float32, bufferSize)
	piBuffer := make([]float64, bufferSize)

	processor.setOnaudioprocess(func(_ ScriptProcessorNode, e AudioProcessingEvent) any {
		//fmt.Println(js.Value(e.outputBuffer()).Get("length"))
		start := 0
		for {
			n, err := AudioStream.Read(piBuffer[start:])
			if err != nil {
				panic(err)
			}

			for i := 0; i < n; i++ {
				browserBuffer[i] = float32(piBuffer[i])
			}

			start += n
			if start >= len(piBuffer) {
				break
			}
		}

		e.outputBuffer().copyToChannel(browserBuffer[:start], 0, 1)
		return nil
	})
	js.Value(processor).Call("connect", js.Value(audio).Get("destination"))

	startAudioOnUserInteraction("mouseup", audio)
	startAudioOnUserInteraction("touchend", audio)

	return func() {
		audio.close()
	}, nil
}

func startAudioOnUserInteraction(event string, audio AudioContext) {
	var audioIsReady bool

	var f js.Func
	f = js.FuncOf(func(this js.Value, arguments []js.Value) any {
		if !audioIsReady {
			fmt.Println("RESUMED")
			audio.resume()
			audioIsReady = true
		}
		js.Global().Get("document").Call("removeEventListener", event, f)
		return nil
	})
	js.Global().Get("document").Call("addEventListener", event, f)
}

func newAudioContext() (AudioContext, error) {
	class := js.Global().Get("AudioContext")
	if !class.Truthy() {
		class = js.Global().Get("webkitAudioContext")
	}
	if !class.Truthy() {
		return AudioContext{}, errors.New("oto: AudioContext or webkitAudioContext was not found")
	}

	options := js.Global().Get("Object").New()
	options.Set("sampleRate", audioSampleRate)

	return AudioContext(class.New(options)), nil
}

type AudioContext js.Value

func (a AudioContext) createScriptProcessor(bufferSize, numberOfInputChannels, numberOfOutputChannels int) ScriptProcessorNode {
	v := js.Value(a).Call("createScriptProcessor", bufferSize, numberOfInputChannels, numberOfOutputChannels)
	return ScriptProcessorNode(v)
}

func (a AudioContext) destination() AudioDestinationNode {
	return AudioDestinationNode(js.Value(a).Get("destination"))
}

func (a AudioContext) resume() {
	js.Value(a).Call("resume")
}

func (a AudioContext) close() {
	js.Value(a).Call("close")
}

type ScriptProcessorNode js.Value

func (n ScriptProcessorNode) setOnaudioprocess(f func(this ScriptProcessorNode, e AudioProcessingEvent) any) {
	js.Value(n).Set("onaudioprocess", js.FuncOf(func(this js.Value, args []js.Value) any {
		return f(n, AudioProcessingEvent(args[0]))
	}))
}

func (n ScriptProcessorNode) connectDestination(node AudioDestinationNode) {
	js.Value(n).Call("connect", node)
}

func (n ScriptProcessorNode) bufferSize() int {
	return js.Value(n).Get("bufferSize").Int()
}

type AudioProcessingEvent js.Value

func (o AudioProcessingEvent) outputBuffer() AudioBuffer {
	return AudioBuffer(js.Value(o).Get("outputBuffer"))
}

type AudioBuffer js.Value

func (a AudioBuffer) copyToChannel(source []float32, channelNumber int, startInChannel int) {
	js.Value(a).Call("copyToChannel", float32SliceToTypedArray(source), channelNumber, startInChannel)
}

func float32SliceToTypedArray(s []float32) js.Value {
	bs := unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), len(s)*4)
	a := js.Global().Get("Uint8Array").New(len(bs))
	js.CopyBytesToJS(a, bs)
	runtime.KeepAlive(s)
	buf := a.Get("buffer")
	return js.Global().Get("Float32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
}

type AudioDestinationNode js.Value
