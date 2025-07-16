// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// A minimal snake (worm) game implementation.
package main

import (
	_ "embed"
	"math/rand"
	"slices"
	"strconv"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/picofont"
	"github.com/elgopher/pi/piebiten"
	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/pipad"
)

var snake []pi.Position           // snake body segments
var fruit pi.Position             // fruit location
var direction pi.Position         // current snake heading
var possibleDirection pi.Position // next possible snake direction

var frame = 0
var speed int
var gameOver = false

const gridSize = 8
const width = 16
const height = 16

var leftDirection = pi.Position{X: -1}
var rightDirection = pi.Position{X: 1}
var upDirection = pi.Position{Y: -1}
var downDirection = pi.Position{Y: 1}

func startNewGame() {
	gameOver = false

	speed = 5
	direction = pi.Position{X: 1, Y: 0}
	possibleDirection = direction
	fruit = pi.Position{X: 8, Y: 8}
	snake = []pi.Position{
		{X: 4, Y: 4},
		{X: 3, Y: 4},
		{X: 2, Y: 4},
	}
}

func handleUserInput() {
	if (pikey.Duration(pikey.Left) > 0 || pipad.Duration(pipad.Left) > 0) && direction.X == 0 {
		possibleDirection = leftDirection
	}
	if (pikey.Duration(pikey.Right) > 0 || pipad.Duration(pipad.Right) > 0) && direction.X == 0 {
		possibleDirection = rightDirection
	}
	if (pikey.Duration(pikey.Up) > 0 || pipad.Duration(pipad.Top) > 0) && direction.Y == 0 {
		possibleDirection = upDirection
	}
	if (pikey.Duration(pikey.Down) > 0 || pipad.Duration(pipad.Bottom) > 0) && direction.Y == 0 {
		possibleDirection = downDirection
	}
}

func spawnFruit() {
	fruit.X = rand.Intn(width)
	fruit.Y = rand.Intn(height)
}

func update() {
	if gameOver {
		if pikey.Duration(pikey.Enter) > 0 || pipad.Duration(pipad.A) > 0 {
			startNewGame()
		}
		return
	}

	handleUserInput()

	frame += 1
	if frame%speed == 0 {
		direction = possibleDirection
		// create new head position
		newPos := snake[0].Add(direction)

		// collisions
		// check collision with wall
		if newPos.X < 0 || newPos.X >= width || newPos.Y < 0 || newPos.Y >= height {
			gameOver = true
			return
		}
		// check collision with the snake itself
		for i := 0; i < len(snake); i++ {
			if snake[i] == newPos {
				gameOver = true
				return
			}
		}

		// move the snake body
		snake = slices.Insert(snake, 0, newPos)
		// check if it eats the apple
		if newPos == fruit {
			spawnFruit()
			if len(snake)%10 == 0 && speed > 0 {
				speed -= 1 // increase speed
			}
		} else {
			snake = snake[:len(snake)-1] // remove tail
		}
	}
}

func draw() {
	pi.Screen().Clear(0)

	drawGrid()
	drawFruit()
	drawSnake()

	if gameOver {
		score := "SCORE: " + strconv.Itoa(len(snake)-3)
		picofont.Sheet.PrintStroked(score, 54, 58, 7, 5)
		pi.SetColor(7)
		picofont.Sheet.Print("HIT ENTER TO START", 33, 74)
	}
}

func drawGrid() {
	pi.SetColor(1)
	for i := 0; i < width; i++ {
		pi.Line(i*gridSize, 0, i*gridSize, height*gridSize)
		pi.Line(0, i*gridSize, width*gridSize, i*gridSize)
	}
}

func drawFruit() {
	verticalShift := frame % 10 / 5 // simple animation
	pi.Spr(fruitSprite, fruit.X*gridSize, fruit.Y*gridSize+verticalShift)
}

func drawSnake() {
	var headSprite pi.Sprite
	switch direction {
	case leftDirection:
		headSprite = headHorizontal.WithFlipX(true) // reuse sprite
	case rightDirection:
		headSprite = headHorizontal
	case upDirection:
		headSprite = headVertical
	case downDirection:
		headSprite = headVertical.WithFlipY(true) // reuse sprite
	}
	pi.Spr(headSprite, snake[0].X*gridSize, snake[0].Y*gridSize)
	for i := 1; i < len(snake); i++ {
		bodySegment := snake[i]
		pi.Spr(bodySprite, bodySegment.X*gridSize, bodySegment.Y*gridSize)
	}
}

//go:embed "sprites.png"
var spritesPNG []byte

var fruitSprite, headVertical, headHorizontal, bodySprite pi.Sprite

func main() {
	pi.Palette = pi.DecodePalette(spritesPNG)
	sprites := pi.DecodeCanvas(spritesPNG)
	fruitSprite = pi.SpriteFrom(sprites, 0, 0, 8, 8)
	headVertical = pi.SpriteFrom(sprites, 8, 0, 8, 8)
	headHorizontal = pi.SpriteFrom(sprites, 16, 0, 8, 8)
	bodySprite = pi.SpriteFrom(sprites, 24, 0, 8, 8)

	pi.SetTPS(30) // 60 is for hardcore players!
	pi.SetScreenSize(gridSize*width, gridSize*height)
	pi.Update = update
	pi.Draw = draw

	startNewGame()

	piebiten.Run()
}
