package main

import (
	"math/rand"
	"net/http"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/coro"
	"github.com/elgopher/pi/ebitengine"
)

var coroutines coro.Routines

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	pi.Update = func() {
		if pi.MouseBtnp(pi.MouseLeft) {
			//r := movePixel(pi.MousePos)
			for j := 0; j < 8000; j++ { // (~6-9KB per COROUTINE). Pico-8 has 4000 coroutines limit
				r := coro.New(func(yield coro.Yield) {
					sleep(10, yield)
					moveHero(10, 120, 5, 10, yield)
					sleep(20, yield)
					moveHero(120, 10, 2, 10, yield)
				})
				coroutines = append(coroutines, r) // complexCoroutine is 2 coroutines - 12-18KB in total
			}
		}
	}

	pi.Draw = func() {
		pi.Cls()
		coroutines = coroutines.ResumeAll()
		//devtools.Export("coroutines", coroutines)
	}

	ebitengine.Run()
}

func movePixel(pos pi.Position, yield coro.Yield) {
	for i := 0; i < 64; i++ {
		pi.Set(pos.X+i, pos.Y+i, byte(rand.Intn(16)))
		yield()
		yield()
	}
}

func moveHero(startX, stopX, minSpeed, maxSpeed int, yield coro.Yield) {
	anim := coro.WithReturn(randomMove(startX, stopX, minSpeed, maxSpeed))

	for {
		x, hasMore := anim.Resume()
		pi.Set(x, 20, 7)
		if hasMore {
			yield()
		} else {
			return
		}
	}
}

// Reusable coroutine which returns int.
func randomMove(start, stop, minSpeed, maxSpeed int) func(yield coro.YieldReturn[int]) {
	pos := start

	return func(yield coro.YieldReturn[int]) {
		for {
			speed := rand.Intn(maxSpeed - minSpeed)
			if stop > start {
				pos = pi.MinInt(stop, pos+speed) // move pos in stop direction by random speed
			} else {
				pos = pi.MaxInt(stop, pos-speed)
			}

			if pos == stop {
				return
			} else {
				yield(pos)
			}
		}
	}
}

func sleep(iterations int, yield coro.Yield) {
	for i := 0; i < iterations; i++ {
		yield()
	}
}
