package gotuga

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Turtle struct {
	canvas     *image.RGBA
	W, H       int
	bg         color.Color
	x, y       float64
	headingDeg float64
	penDown    bool
	penColor   color.Color
	penWidth   float64

	filling   bool
	fillColor color.Color
	fillPath  []image.Point // collected pixel coords
}

// New creates a new turtle with a W×H canvas and a background color.
// The turtle starts at (0,0) facing 0° (east), pen down, black ink, width 2px.
func New(W, H int, bg color.Color) *Turtle {
	if bg == nil {
		bg = color.RGBA{A: 0}
	}
	t := &Turtle{
		canvas:     image.NewRGBA(image.Rect(0, 0, W, H)),
		W:          W,
		H:          H,
		bg:         bg,
		x:          0,
		y:          0,
		headingDeg: 0,
		penDown:    true,
		penColor:   color.Black,
		penWidth:   2,
	}
	t.fillCanvas(bg)
	return t
}

// SavePNG writes the current canvas to a PNG file.
func (t *Turtle) SavePNG(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, t.canvas)
}

// Image returns the underlying RGBA canvas (read/write).
func (t *Turtle) Image() *image.RGBA { return t.canvas }

// Starts Drawing Mode of Turtle
func (t *Turtle) PenUp() { t.penDown = false }

// Stops Drawing Mode of Turtle
func (t *Turtle) PenDown() { t.penDown = true }

// Set pen Color to color.Color type from "image/color" package
func (t *Turtle) SetColor(c color.Color) {
	if c != nil {
		t.penColor = c
	}
}

// Sets the Thickness or Width of the Pen
func (t *Turtle) SetWidth(w float64) {
	if w > 0 {
		t.penWidth = w
	}
}

// Set Turtle's Rotation towards (deg) Degrees
func (t *Turtle) SetHeading(deg float64) { t.headingDeg = deg }

// Turn Left (deg) Degrees
func (t *Turtle) Left(deg float64) { t.headingDeg += deg }

// Turn Right (deg) Degrees
func (t *Turtle) Right(deg float64) { t.headingDeg -= deg }

// Go To (0, 0) and Reset Direction to 0 Degrees, Keeps the Canvas state.
func (t *Turtle) Home() { t.GoTo(0, 0); t.headingDeg = 0 }

// Clear repaints the canvas with the background color but keeps turtle state.
func (t *Turtle) Clear() { t.fillCanvas(t.bg) }

// Reset clears the canvas and resets position/orientation/pen to defaults.
func (t *Turtle) Reset() {
	t.fillCanvas(t.bg)
	t.x, t.y = 0, 0
	t.headingDeg = 0
	t.penDown = true
	t.penColor = color.Black
	t.penWidth = 2
}

// Move Forward by (d) Steps
func (t *Turtle) Forward(d float64) {
	rad := t.headingDeg * math.Pi / 180
	nx := t.x + d*math.Cos(rad)
	ny := t.y + d*math.Sin(rad)
	if t.penDown {
		t.drawSegment(t.x, t.y, nx, ny, t.penWidth, t.penColor)
	}
	t.recordFillVertex(nx, ny)
	t.x, t.y = nx, ny
}

// Move Backwards by (d) Steps
func (t *Turtle) Backward(d float64) { t.Forward(-d) }

// GoTo moves to logical coords (x,y). If pen is down, draws a segment.
func (t *Turtle) GoTo(x, y float64) {
	if t.penDown {
		t.drawSegment(t.x, t.y, x, y, t.penWidth, t.penColor)
	}
	t.recordFillVertex(x, y)
	t.x, t.y = x, y
}

// Shapes (drawn at current position/orientation)
func (t *Turtle) Rect(w, h float64) {
	// Outline rectangle centered on the *path* starting corner (current pos)
	// and aligned to current heading.
	// We trace the perimeter and return to the start.
	orig := t.stateSnapshot()
	t.Forward(w)
	t.Left(90)
	t.Forward(h)
	t.Left(90)
	t.Forward(w)
	t.Left(90)
	t.Forward(h)
	t.restoreSnapshot(orig)
}

// Polygon draws an n-sided regular polygon with side length s.
func (t *Turtle) Polygon(n int, side float64) {
	if n < 3 {
		return
	}
	orig := t.stateSnapshot()
	turn := 360.0 / float64(n)
	for i := 0; i < n; i++ {
		t.Forward(side)
		t.Left(turn)
	}
	t.restoreSnapshot(orig)
}

// Circle draws an approximate circle with radius r using small segments.
func (t *Turtle) Circle(r float64) {
	circ := 2 * math.Pi * math.Abs(r)
	// segment length ~ 3 px (minimum 12 segments)
	segments := int(math.Max(12, circ/3))
	angle := 360.0 / float64(segments)
	// Shift center to the left of heading by r (turtle circle convention)
	orig := t.stateSnapshot()
	// Walk the polyline approximation
	stepLen := circ / float64(segments)
	turn := angle
	if r < 0 {
		turn = -angle
	}
	for i := 0; i < segments; i++ {
		t.Forward(stepLen)
		t.Left(turn)
	}
	t.restoreSnapshot(orig)
}

// BeginFill starts recording a polygon fill path
func (t *Turtle) BeginFill() {
	t.filling = true
	t.fillPath = nil
}

// FillColor sets the fill color
func (t *Turtle) FillColor(c color.Color) {
	if c != nil {
		t.fillColor = c
	}
}

// EndFill fills the collected polygon
func (t *Turtle) EndFill() {
	if !t.filling || len(t.fillPath) < 3 {
		t.filling = false
		t.fillPath = nil
		return
	}

	// Close polygon if needed
	first := t.fillPath[0]
	last := t.fillPath[len(t.fillPath)-1]
	if first != last {
		t.fillPath = append(t.fillPath, first)
	}

	// Fill polygon
	drawPolygon(t.canvas, t.fillPath, t.fillColor)

	// Reset fill state
	t.filling = false
	t.fillPath = nil
}
