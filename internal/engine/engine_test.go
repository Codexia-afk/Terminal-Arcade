package engine

import (
	"testing"
)

func TestDifficulty(t *testing.T) {
	diffs := Difficulties()
	if len(diffs) != 3 {
		t.Fatalf("expected 3 difficulties, got %d", len(diffs))
	}
	for _, d := range diffs {
		if !d.Valid() {
			t.Errorf("expected difficulty %s to be valid", d)
		}
		if d.String() == "Unknown" {
			t.Errorf("expected valid display string for %s", d)
		}
	}

	invalid := Difficulty("insane")
	if invalid.Valid() {
		t.Errorf("expected invalid difficulty to return false on Valid()")
	}
	if invalid.String() != "Unknown" {
		t.Errorf("expected Unknown display string for invalid difficulty")
	}
}

func TestPositionAndRect(t *testing.T) {
	p1 := Position{X: 3, Y: 4}
	p2 := Position{X: 1, Y: -2}
	sum := p1.Add(p2)
	if sum.X != 4 || sum.Y != 2 {
		t.Fatalf("unexpected sum: %+v", sum)
	}
	if !p1.Equal(Position{X: 3, Y: 4}) {
		t.Fatalf("expected equality for identical positions")
	}
	if p1.Equal(p2) {
		t.Fatalf("expected inequality for different positions")
	}

	rect := Rect{X: 5, Y: 5, W: 10, H: 10}
	if !rect.Contains(Position{X: 5, Y: 5}) {
		t.Errorf("expected top-left corner to be contained")
	}
	if !rect.Contains(Position{X: 14, Y: 14}) {
		t.Errorf("expected bottom-right inside cell to be contained")
	}
	if rect.Contains(Position{X: 15, Y: 5}) {
		t.Errorf("expected right boundary to be outside")
	}
	if rect.Contains(Position{X: 4, Y: 5}) {
		t.Errorf("expected outside point to not be contained")
	}

	// Overlaps
	r2 := Rect{X: 10, Y: 10, W: 10, H: 10}
	if !rect.Overlaps(r2) {
		t.Errorf("expected overlapping rectangles to return true")
	}
	r3 := Rect{X: 20, Y: 20, W: 5, H: 5}
	if rect.Overlaps(r3) {
		t.Errorf("expected disjoint rectangles to not overlap")
	}
}

func TestActionDirection(t *testing.T) {
	tests := []struct {
		action Action
		dx, dy int
		isDir  bool
	}{
		{ActionUp, 0, -1, true},
		{ActionDown, 0, 1, true},
		{ActionLeft, -1, 0, true},
		{ActionRight, 1, 0, true},
		{ActionPause, 0, 0, false},
		{ActionQuit, 0, 0, false},
		{ActionConfirm, 0, 0, false},
	}

	for _, tt := range tests {
		pos, ok := tt.action.Direction()
		if ok != tt.isDir {
			t.Errorf("action %v isDir expected %v, got %v", tt.action, tt.isDir, ok)
		}
		if ok && (pos.X != tt.dx || pos.Y != tt.dy) {
			t.Errorf("action %v delta expected (%d,%d), got (%d,%d)", tt.action, tt.dx, tt.dy, pos.X, pos.Y)
		}
	}
}

func TestThemes(t *testing.T) {
	names := ThemeNames()
	if len(names) < 3 {
		t.Fatalf("expected at least 3 themes, got %d", len(names))
	}
	for _, name := range names {
		th := GetTheme(name)
		if th.Name != name {
			t.Errorf("expected theme name %s, got %s", name, th.Name)
		}
	}
	fallback := GetTheme("NonExistent")
	if fallback.Name == "" {
		t.Errorf("expected valid fallback theme")
	}
}
