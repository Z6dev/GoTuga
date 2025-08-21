# GoTuga

`gotuga` is a Go package that implements something similar to a turtle graphics system, Like the one found in Python!  
It allows drawing on an image canvas using simple movement and drawing commands.

---

## Features

- Create and manage a drawable canvas.
- Control a "turtle" with position, heading, pen state, and pen properties.
- Draw lines, rectangles, polygons, and circles.
- Support for filled shapes with customizable fill color.
- Export the final drawing as a PNG image.

---

## Documentation

## Installation

```bash
go get github.com/Z6dev/GoTuga
```

## Usage Example

```go
package main

import (
    "image/color"
    gotuga "github.com/Z6dev/GoTuga"
)

func main() {
    // Create a new canvas 500×500 with white background
    t := gotuga.New(500, 500, color.White)

    // Set pen color and width
    t.SetColor(color.RGBA{255, 0, 0, 255}) // Red
    t.SetWidth(3)

    // Draw a square
    for i := 0; i < 4; i++ {
        t.Forward(100)
        t.Left(90)
    }

    // Save to PNG
    t.SavePNG("square.png")
}

```
