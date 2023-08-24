package main

import (
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

var iterators pi.Iterators

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	pi.Update = func() {
		//if pi.MouseBtn(pi.MouseLeft) {
		//	iterators = append(iterators, movePixel(pi.MousePos))
		//}
		if pi.MouseBtnp(pi.MouseLeft) {
			for j := 0; j < 8000; j++ { // 250bytes per coroutine
				iterators = append(iterators, complexIteratorAlternative()) // 30x faster than coroutines
			}
		}
	}

	pi.Draw = func() {
		pi.Cls()
		iterators = iterators.Next()
		fmt.Println(len(iterators))
	}

	ebitengine.MustRun()
}

func movePixel(pos pi.Position) pi.Iterator {
	i := 0
	return func() bool {
		if i == 128 {
			return false
		}

		if i%2 == 0 { // draw pixel every 2 frames
			pi.Set(pos.X+i, pos.Y+i, byte(rand.Intn(16)))
		}

		i++

		return true
	}
}

func moveHero(startX, stopX, minSpeed, maxSpeed int) pi.Iterator {
	anim := randomMove(startX, stopX, minSpeed, maxSpeed)
	finished := false

	return func() bool {
		if finished {
			return false
		}

		x, hasNext := anim()
		if !hasNext {
			finished = true
		}
		pi.Set(x, 20, 7)
		return hasNext
	}
}

// Reusable iterator which returns int
func randomMove(start, stop, minSpeed, maxSpeed int) func() (int, bool) {
	pos := start

	return func() (int, bool) {
		speed := rand.Intn(maxSpeed - minSpeed)
		if stop > start {
			pos = pi.MinInt(stop, pos+speed) // move pos in stop direction by random speed
		} else {
			pos = pi.MaxInt(stop, pos-speed)
		}

		return pos, pos != stop
	}
}

func complexIterator() pi.Iterator {
	return pi.Sequence(
		sleep(10),
		moveHero(10, 120, 5, 10),
		sleep(20),
		moveHero(120, 10, 2, 10),
	)
}

func complexIteratorAlternative() pi.Iterator {
	sleep10 := sleep(10)                      // + 2 allocations
	move := moveHero(10, 120, 5, 10)          // + 4 allocations
	sleep20 := sleep(20)                      // + 2 allocations
	moveBackwards := moveHero(120, 10, 2, 10) // + 4 allocations

	return func() bool {
		if sleep10() {
			return true
		}
		if move() {
			return true
		}
		if sleep20() {
			return true
		}
		if moveBackwards() {
			return true
		}
		return false
	}
}

func complexIterator2() pi.Iterator {
	return pi.Sequence(
		sleep(90),
		func() bool {
			fmt.Println("After 90 frames")
			return false
		},
	)
}

// this is better than complexIterator2
func complexIterator3() pi.Iterator {
	sleep := sleep(90)

	return func() bool {
		if sleep() {
			return true
		}

		fmt.Println("After 90 frames")
		return false
	}
}

// this is better event better than complexIterator3
func complexIterator4() pi.Iterator {
	i := 0

	return func() bool {
		switch {
		case i < 90:
			i++
			return true
		case i == 90:
			fmt.Println("After 90 frames")
		}

		return false
	}
}

func finishOnNextCall() bool {
	return false
}

func sleep(iterations int) pi.Iterator {
	if iterations <= 0 {
		return finishOnNextCall
	}

	i := 0

	return func() bool {
		if i == iterations {
			return false
		}

		i++

		return i != iterations
	}
}
