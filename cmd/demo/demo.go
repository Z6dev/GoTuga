package main

import (
	"image/color"
	"log"
	"math"

	gotuga "github.com/Z6dev/GoTuga"
)

func main() {
	// Create a 1024x768 white canvas
	t := gotuga.New(1024, 768, color.White)

	// Style
	t.SetColor(color.RGBA{0, 0, 0, 255})
	t.SetWidth(3)

	// Draw axes
	t.PenUp()
	t.GoTo(-480, 0)
	t.PenDown()
	t.Forward(960)
	t.PenUp()
	t.GoTo(0, -340)
	t.SetHeading(90)
	t.PenDown()
	t.Forward(680)

	// A square
	t.PenUp()
	t.GoTo(-200, -200)
	t.SetHeading(0)
	t.PenDown()
	t.Rect(200, 200)

	// A triangle
	t.PenUp()
	t.GoTo(200, -200)
	t.SetHeading(0)
	t.PenDown()
	t.Polygon(3, 180)

	// A circle
	t.PenUp()
	t.GoTo(0, 150)
	t.SetHeading(0)
	t.PenDown()
	t.SetWidth(5)
	t.SetColor(color.RGBA{30, 144, 255, 255}) // dodger blue
	t.Circle(100)

	// A spiral
	t.PenUp()
	t.GoTo(0, 0)
	t.SetHeading(0)
	t.PenDown()
	t.SetColor(color.RGBA{220, 20, 60, 255}) // crimson
	t.SetWidth(2)
	step := 4.0
	for i := 0; i < 120; i++ {
		t.Forward(step)
		t.Left(15)
		step *= 1.02
	}

	// A sine wave
	t.SetColor(color.RGBA{34, 139, 34, 255}) // forest green
	t.SetWidth(2)
	t.PenUp()
	t.GoTo(-480, 250)
	t.PenDown()
	t.SetHeading(0)
	for x := -480.0; x <= 480.0; x += 2 {
		y := 250 + 40*math.Sin(x*math.Pi/120)
		t.GoTo(x, y)
	}

	if err := t.SavePNG("turtle_demo.png"); err != nil {
		log.Fatal(err)
	}
}
