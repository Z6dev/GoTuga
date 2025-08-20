package gotuga

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
)

type Turtle struct {
	canvas     *image.RGBA
	W, H       int
	bg         color.Color
	x, y       float64 // logical coords, origin at center, +y up
	headingDeg float64 // 0° = +x (east), CCW positive (like standard math)
	penDown    bool
	penColor   color.Color
	penWidth   float64
}

// New creates a new turtle with a W×H canvas and a background color.
// The turtle starts at (0,0) facing 0° (east), pen down, black ink, width 2px.
func New(W, H int, bg color.Color) *Turtle {
	if bg == nil {
		bg = color.White
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

// State setters
func (t *Turtle) PenUp()   { t.penDown = false }
func (t *Turtle) PenDown() { t.penDown = true }
func (t *Turtle) SetColor(c color.Color) {
	if c != nil {
		t.penColor = c
	}
}
func (t *Turtle) SetWidth(w float64) {
	if w > 0 {
		t.penWidth = w
	}
}
func (t *Turtle) SetHeading(deg float64) { t.headingDeg = deg }
func (t *Turtle) Left(deg float64)       { t.headingDeg += deg }
func (t *Turtle) Right(deg float64)      { t.headingDeg -= deg }
func (t *Turtle) Home()                  { t.GoTo(0, 0); t.headingDeg = 0 }

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

// Movement
func (t *Turtle) Forward(d float64) {
	rad := t.headingDeg * math.Pi / 180
	nx := t.x + d*math.Cos(rad)
	ny := t.y + d*math.Sin(rad)
	if t.penDown {
		t.drawSegment(t.x, t.y, nx, ny, t.penWidth, t.penColor)
	}
	t.x, t.y = nx, ny
}
func (t *Turtle) Backward(d float64) { t.Forward(-d) }

// GoTo moves to logical coords (x,y). If pen is down, draws a segment.
func (t *Turtle) GoTo(x, y float64) {
	if t.penDown {
		t.drawSegment(t.x, t.y, x, y, t.penWidth, t.penColor)
	}
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
// Positive r draws CCW, negative r draws CW (like classic turtle).
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

// --- Internals ---

type snapshot struct {
	x, y       float64
	headingDeg float64
}

func (t *Turtle) stateSnapshot() snapshot {
	return snapshot{t.x, t.y, t.headingDeg}
}
func (t *Turtle) restoreSnapshot(s snapshot) {
	t.x, t.y = s.x, s.y
	t.headingDeg = s.headingDeg
}

func (t *Turtle) fillCanvas(c color.Color) {
	draw.Draw(t.canvas, t.canvas.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)
}

// Map logical (x,y) where origin is center and +y up, to image pixel coords.
func (t *Turtle) mapToPixel(x, y float64) (int, int) {
	ix := int(math.Round(x + float64(t.W)/2))
	iy := int(math.Round(float64(t.H)/2 - y))
	return ix, iy
}

// Draw a thick anti-aliased-ish segment by stamping discs along the path.
// This is simple and dependency-free; good enough for turtle graphics.
func (t *Turtle) drawSegment(x0, y0, x1, y1 float64, width float64, col color.Color) {
	dx := x1 - x0
	dy := y1 - y0
	dist := math.Hypot(dx, dy)
	if dist == 0 {
		t.stampDisc(x0, y0, width/2, col)
		return
	}
	// One sample per pixel of distance, plus endpoints
	steps := int(math.Ceil(dist)) + 1
	for i := 0; i <= steps; i++ {
		tt := float64(i) / float64(steps)
		x := x0 + tt*dx
		y := y0 + tt*dy
		t.stampDisc(x, y, width/2, col)
	}
}

// stampDisc draws a filled circle of radius r in logical coordinates.
func (t *Turtle) stampDisc(cx, cy, r float64, col color.Color) {
	if r <= 0 {
		return
	}
	px, py := t.mapToPixel(cx, cy)
	rr := int(math.Ceil(r))
	minX := clamp(px-rr, 0, t.W-1)
	maxX := clamp(px+rr, 0, t.W-1)
	minY := clamp(py-rr, 0, t.H-1)
	maxY := clamp(py+rr, 0, t.H-1)

	r2 := r * r
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x-px) + 0.5
			dy := float64(y-py) + 0.5
			if dx*dx+dy*dy <= r2 {
				t.canvas.Set(x, y, col)
			}
		}
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
