package engine

import "github.com/gdamore/tcell/v2"

// Position is an integer grid coordinate in terminal space.
type Position struct {
	X int
	Y int
}

// Add returns the vector sum of two positions.
func (p Position) Add(q Position) Position {
	return Position{X: p.X + q.X, Y: p.Y + q.Y}
}

// Equal reports whether two grid positions are identical.
func (p Position) Equal(q Position) bool {
	return p.X == q.X && p.Y == q.Y
}

// Sprite describes one terminal cell with a glyph and foreground color.
type Sprite struct {
	Glyph rune
	Color tcell.Color
}

// Rect is an axis-aligned bounding box on the grid.
type Rect struct {
	X int
	Y int
	W int
	H int
}

// Contains reports whether position p lies strictly within rectangle r.
func (r Rect) Contains(p Position) bool {
	return p.X >= r.X && p.X < r.X+r.W && p.Y >= r.Y && p.Y < r.Y+r.H
}

// Overlaps reports whether two grid rectangles overlap.
func (r Rect) Overlaps(o Rect) bool {
	return r.X < o.X+o.W && o.X < r.X+r.W && r.Y < o.Y+o.H && o.Y < r.Y+o.H
}
