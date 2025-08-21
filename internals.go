package gotuga

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
)

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

// Very simple polygon fill using draw.DrawMask
func drawPolygon(img *image.RGBA, pts []image.Point, col color.Color) {
	// Make a mask the same size
	mask := image.NewAlpha(img.Bounds())

	// Rasterize polygon edges into the mask
	// (We’ll use a basic scanline fill here)
	for y := mask.Bounds().Min.Y; y < mask.Bounds().Max.Y; y++ {
		var intersections []int
		for i := 0; i < len(pts); i++ {
			j := (i + 1) % len(pts)
			x0, y0 := pts[i].X, pts[i].Y
			x1, y1 := pts[j].X, pts[j].Y
			if (y0 <= y && y1 > y) || (y1 <= y && y0 > y) {
				x := x0 + (y-y0)*(x1-x0)/(y1-y0)
				intersections = append(intersections, x)
			}
		}
		// sort intersections
		sort.Ints(intersections)
		for i := 0; i+1 < len(intersections); i += 2 {
			for x := intersections[i]; x <= intersections[i+1]; x++ {
				mask.SetAlpha(x, y, color.Alpha{A: 255})
			}
		}
	}

	// Apply fill
	draw.DrawMask(img, img.Bounds(), &image.Uniform{C: col}, image.Point{}, mask, image.Point{}, draw.Over)
}

// recordFillVertex adds a vertex if filling is active
func (t *Turtle) recordFillVertex(x, y float64) {
	if t.filling {
		px, py := t.mapToPixel(x, y)
		t.fillPath = append(t.fillPath, image.Pt(px, py))
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
