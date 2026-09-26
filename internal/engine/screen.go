package engine

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

// Screen wraps a tcell.Screen with higher-level text and box drawing primitives.
type Screen struct {
	Raw tcell.Screen
}

// NewScreen creates a Screen wrapper around an existing tcell.Screen.
func NewScreen(raw tcell.Screen) *Screen {
	return &Screen{Raw: raw}
}

// Size returns the terminal dimensions as (width, height).
func (s *Screen) Size() (int, int) {
	if s.Raw == nil {
		return 80, 24
	}
	return s.Raw.Size()
}

// Clear fills the entire screen with the specified background color.
func (s *Screen) Clear(bg tcell.Color) {
	if s.Raw == nil {
		return
	}
	s.Raw.SetStyle(tcell.StyleDefault.Background(bg))
	s.Raw.Clear()
}

// DrawCell draws a single rune at grid coordinates (x, y) if inside screen bounds.
func (s *Screen) DrawCell(x, y int, r rune, fg, bg tcell.Color) {
	if s.Raw == nil {
		return
	}
	w, h := s.Size()
	if x >= 0 && y >= 0 && x < w && y < h {
		st := tcell.StyleDefault.Foreground(fg).Background(bg)
		s.Raw.SetContent(x, y, r, nil, st)
	}
}

// DrawText renders a horizontal string starting at (x, y).
func (s *Screen) DrawText(x, y int, text string, fg, bg tcell.Color) {
	curX := x
	for _, r := range text {
		s.DrawCell(curX, y, r, fg, bg)
		curX += runewidth.RuneWidth(r)
	}
}

// CenterText draws text horizontally centered at row y.
func (s *Screen) CenterText(y int, text string, fg, bg tcell.Color) {
	w, _ := s.Size()
	textWidth := runewidth.StringWidth(text)
	startX := (w - textWidth) / 2
	if startX < 0 {
		startX = 0
	}
	s.DrawText(startX, y, text, fg, bg)
}

// Box draws a rectangular border using Unicode box-drawing characters (┌─┐│└─┘).
func (s *Screen) Box(x, y, w, h int, fg, bg tcell.Color) {
	if w < 2 || h < 2 {
		return
	}
	s.DrawCell(x, y, '┌', fg, bg)
	s.DrawCell(x+w-1, y, '┐', fg, bg)
	s.DrawCell(x, y+h-1, '└', fg, bg)
	s.DrawCell(x+w-1, y+h-1, '┘', fg, bg)

	for i := 1; i < w-1; i++ {
		s.DrawCell(x+i, y, '─', fg, bg)
		s.DrawCell(x+i, y+h-1, '─', fg, bg)
	}
	for i := 1; i < h-1; i++ {
		s.DrawCell(x, y+i, '│', fg, bg)
		s.DrawCell(x+w-1, y+i, '│', fg, bg)
	}
}

// BoxWithTitle draws a border with a centered title embedded in the top border.
func (s *Screen) BoxWithTitle(x, y, w, h int, title string, fg, bg, titleFg tcell.Color) {
	s.Box(x, y, w, h, fg, bg)
	if title != "" && w > 4 {
		titleText := " " + title + " "
		tw := runewidth.StringWidth(titleText)
		if tw > w-4 {
			titleText = titleText[:w-4]
			tw = runewidth.StringWidth(titleText)
		}
		tx := x + (w-tw)/2
		s.DrawText(tx, y, titleText, titleFg, bg)
	}
}

// Fill fills an area of the screen with a rune.
func (s *Screen) Fill(x, y, w, h int, r rune, fg, bg tcell.Color) {
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			s.DrawCell(x+col, y+row, r, fg, bg)
		}
	}
}

// Flush flushes pending draws to the terminal display.
func (s *Screen) Flush() {
	if s.Raw != nil {
		s.Raw.Show()
	}
}

// DoubleBox draws a border using double-line Unicode characters (╔═╗║╚═╝).
func (s *Screen) DoubleBox(x, y, w, h int, fg, bg tcell.Color) {
	if w < 2 || h < 2 {
		return
	}
	s.DrawCell(x, y, '╔', fg, bg)
	s.DrawCell(x+w-1, y, '╗', fg, bg)
	s.DrawCell(x, y+h-1, '╚', fg, bg)
	s.DrawCell(x+w-1, y+h-1, '╝', fg, bg)

	for i := 1; i < w-1; i++ {
		s.DrawCell(x+i, y, '═', fg, bg)
		s.DrawCell(x+i, y+h-1, '═', fg, bg)
	}
	for i := 1; i < h-1; i++ {
		s.DrawCell(x, y+i, '║', fg, bg)
		s.DrawCell(x+w-1, y+i, '║', fg, bg)
	}
}

// DoubleBoxWithTitle draws a double-line border with a centered title.
func (s *Screen) DoubleBoxWithTitle(x, y, w, h int, title string, fg, bg, titleFg tcell.Color) {
	s.DoubleBox(x, y, w, h, fg, bg)
	if title != "" && w > 4 {
		titleText := " " + title + " "
		tw := runewidth.StringWidth(titleText)
		if tw > w-4 {
			titleText = titleText[:w-4]
			tw = runewidth.StringWidth(titleText)
		}
		tx := x + (w-tw)/2
		s.DrawText(tx, y, titleText, titleFg, bg)
	}
}

// AudioBellEnabled toggles audio bell feedback across the suite.
var AudioBellEnabled = true

// Beep sends an audible terminal bell alert if supported by the terminal emulator.
func (s *Screen) Beep() {
	if !AudioBellEnabled {
		return
	}
	if s.Raw != nil {
		_ = s.Raw.Beep()
	}
}

// Beep emits an audible terminal bell alert.
func Beep() {
	if !AudioBellEnabled {
		return
	}
	// Emit ASCII bell control character
	print("\a")
}


