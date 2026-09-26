package engine

import (
	"testing"

	"github.com/gdamore/tcell/v2"
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

func TestFiveThemes(t *testing.T) {
	expected := []string{"Retro Green", "Neon", "Monochrome", "Cyberpunk", "Ocean"}
	names := ThemeNames()
	if len(names) != 5 {
		t.Fatalf("expected 5 themes, got %d", len(names))
	}
	for i, exp := range expected {
		if names[i] != exp {
			t.Errorf("expected theme %d to be %s, got %s", i, exp, names[i])
		}
		th := GetTheme(exp)
		if th.SnakeHeadGlyph == 0 || th.SnakeBodyGlyph == 0 || th.SnakeFoodGlyph == 0 {
			t.Errorf("theme %s missing snake element glyphs", exp)
		}
		// Verify directional head runes
		if th.SnakeHeadUp == 0 || th.SnakeHeadDown == 0 || th.SnakeHeadLeft == 0 || th.SnakeHeadRight == 0 {
			t.Errorf("theme %s missing directional head runes", exp)
		}

		headR := th.SnakeHeadForDir(Position{X: 1, Y: 0})
		if headR != th.SnakeHeadRight {
			t.Errorf("expected right head glyph %c, got %c", th.SnakeHeadRight, headR)
		}
		headL := th.SnakeHeadForDir(Position{X: -1, Y: 0})
		if headL != th.SnakeHeadLeft {
			t.Errorf("expected left head glyph %c, got %c", th.SnakeHeadLeft, headL)
		}
		headU := th.SnakeHeadForDir(Position{X: 0, Y: -1})
		if headU != th.SnakeHeadUp {
			t.Errorf("expected up head glyph %c, got %c", th.SnakeHeadUp, headU)
		}
		headD := th.SnakeHeadForDir(Position{X: 0, Y: 1})
		if headD != th.SnakeHeadDown {
			t.Errorf("expected down head glyph %c, got %c", th.SnakeHeadDown, headD)
		}
	}
}

func TestMockScreenRenderingAndBoxes(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simulation screen: %v", err)
	}
	defer simScreen.Fini()

	simScreen.SetSize(80, 24)
	screen := NewScreen(simScreen)

	// Test DoubleBox
	screen.DoubleBoxWithTitle(5, 5, 20, 10, "TEST", tcell.ColorAqua, tcell.ColorBlack, tcell.ColorWhite)
	screen.Flush()

	c, _, _, _ := simScreen.GetContent(5, 5)
	if c != '╔' {
		t.Errorf("expected top-left double border '╔', got '%c'", c)
	}
	c, _, _, _ = simScreen.GetContent(24, 5)
	if c != '╗' {
		t.Errorf("expected top-right double border '╗', got '%c'", c)
	}

	// Test Theme Rendering
	for _, themeName := range ThemeNames() {
		th := GetTheme(themeName)
		screen.Clear(th.Background)
		screen.DrawCell(10, 10, th.SnakeHeadGlyph, th.Player, th.Background)
		screen.DrawCell(11, 10, th.SnakeBodyGlyph, th.Player, th.Background)
		screen.DrawCell(12, 10, th.SnakeFoodGlyph, th.Item, th.Background)
		screen.Flush()

		rHead, _, _, _ := simScreen.GetContent(10, 10)
		if rHead != th.SnakeHeadGlyph {
			t.Errorf("theme %s: expected head '%c', got '%c'", themeName, th.SnakeHeadGlyph, rHead)
		}
		rBody, _, _, _ := simScreen.GetContent(11, 10)
		if rBody != th.SnakeBodyGlyph {
			t.Errorf("theme %s: expected body '%c', got '%c'", themeName, th.SnakeBodyGlyph, rBody)
		}
		rFood, _, _, _ := simScreen.GetContent(12, 10)
		if rFood != th.SnakeFoodGlyph {
			t.Errorf("theme %s: expected food '%c', got '%c'", themeName, th.SnakeFoodGlyph, rFood)
		}
	}
}

