package main

import (
	"math"
	"math/rand"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/ebitengine"
)

var iterators pi.Iterators

func main() {
	for i := 0; i < 32; i++ {
		pos := pi.Position{X: rand.Intn(128), Y: rand.Intn(128)}
		radius := rand.Intn(9) + 1
		color := byte(rand.Intn(15)) + 1

		iterators = append(iterators, animateBall(pos, radius, color))
	}

	pi.Draw = func() {
		pi.Cls()
		iterators = iterators.Next()
	}

	ebitengine.MustRun()
}

func animateBall(currentPos pi.Position, radius int, color byte) pi.Iterator {
	var moveBall func() (pi.Position, bool) // moveBall will have move iterator which calculates new position on each call
	framesStoodStill := 0

	return func() bool {
		pi.CircFill(currentPos.X, currentPos.Y, radius, color)

		notMoving := moveBall == nil
		if notMoving {
			framesStoodStill++
			if framesStoodStill == 90 {
				// It's time to move the ball to a new position
				newPos := pi.Position{X: rand.Intn(128), Y: rand.Intn(128)}
				moveBall = move(currentPos, newPos)
				framesStoodStill = 0
			}

			return true // animateBall iterator never ends
		}

		var hasNext bool
		currentPos, hasNext = moveBall() // run iterator which returns new position
		if !hasNext {                    // hasNext = false means that iterator was finished
			moveBall = nil
		}

		return true // animateBall iterator never ends
	}
}

const speed = 2

func move(from, to pi.Position) func() (pi.Position, bool) {
	dy := float64(to.Y - from.Y)
	dx := float64(to.X - from.X)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
	x, y := float64(from.X), float64(from.Y)

	stepX := speed * dx / distance
	stepY := speed * dy / distance
	steps := int(distance / speed)
	step := 0

	return func() (pi.Position, bool) {
		x += stepX
		y += stepY
		step++

		newPos := pi.Position{X: int(x), Y: int(y)}
		return newPos, steps != step // iterator will finish if steps == step
	}
}
