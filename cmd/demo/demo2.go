package main

import (
    "image/color"
    "github.com/Z6dev/GoTuga"
)


func main() {
    // Create a turtle on a 400x400 white canvas
    t := gotuga.New(400, 400, color.White)

    // --- Draw a filled triangle ---
    t.BeginFill()
    t.FillColor(color.RGBA{255, 0, 0, 255}) // red fill
    t.Forward(100)
    t.Left(120)
    t.Forward(100)
    t.Left(120)
    t.Forward(100)
    t.EndFill()

    // --- Move and draw a filled rectangle ---
    t.PenUp()
    t.Right(90)
    t.Forward(150)
    t.PenDown()
    t.BeginFill()
    t.FillColor(color.RGBA{0, 128, 255, 255}) // blue fill
    t.Rect(120, 80)
    t.EndFill()

    // --- Draw a filled circle ---
    t.PenUp()
    t.Home()
    t.Forward(150)
    t.PenDown()
    t.BeginFill()
    t.FillColor(color.RGBA{0, 200, 0, 255}) // green fill
    t.Circle(60)
    t.EndFill()

    // Save result
    t.SavePNG("gotuga_demo2.png")
}
