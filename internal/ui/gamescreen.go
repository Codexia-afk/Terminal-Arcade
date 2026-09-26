// Package ui provides reusable premium terminal UI layouts, centered HUDs, and screens.
package ui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// RenderHUD renders the standardized premium HUD strip above the gameplay area.
func RenderHUD(s *engine.Screen, x, y, width int, title, difficulty, themeName, extraInfo string, score, lives int, streak int) {
	th := engine.GetTheme(themeName)

	if width < 50 {
		width = 50
	}

	// 1. Top border with title tag
	topBar := fmt.Sprintf("┌── %s • %s • %s ", title, difficulty, themeName)
	for len(topBar) < width-1 {
		topBar += "─"
	}
	topBar += "┐"
	s.DrawText(x, y, topBar, th.HUD, th.Background)

	// 2. Middle content row: Score | Lives | Extra | Streak
	livesText := "—"
	if lives > 0 {
		livesText = fmt.Sprintf("%d", lives)
	}

	streakBadge := "[🔥 STREAK: 0 DAYS]"
	if themeName == "Monochrome" {
		streakBadge = fmt.Sprintf("[STREAK: %d DAYS]", streak)
	} else if streak > 0 {
		streakBadge = fmt.Sprintf("[🔥 STREAK: %d days]", streak)
	} else {
		streakBadge = "[STREAK: 0 days]"
	}

	midLeft := fmt.Sprintf("│ Score: %-5d | Lives: %-2s", score, livesText)
	if extraInfo != "" {
		midLeft += " | " + extraInfo
	}

	gap := width - len(midLeft) - len(streakBadge) - 2
	if gap < 1 {
		gap = 1
	}

	midRow := midLeft
	for i := 0; i < gap; i++ {
		midRow += " "
	}
	midRow += streakBadge + " │"
	if len(midRow) > width {
		midRow = midRow[:width-1] + "│"
	}

	s.DrawText(x, y+1, midRow, th.Text, th.Background)

	// 3. Bottom border
	bottomBar := "└"
	for i := 1; i < width-1; i++ {
		bottomBar += "─"
	}
	bottomBar += "┘"
	s.DrawText(x, y+2, bottomBar, th.HUD, th.Background)
}

// RenderBottomInstructions renders the single-line footer instruction bar below the play area.
func RenderBottomInstructions(s *engine.Screen, y int, customText string, th engine.Theme) {
	text := "↑↓←→ or WASD Move  |  P Pause  |  Q Quit to Menu"
	if customText != "" {
		text = customText
	}
	s.CenterText(y, text, th.Text, th.Background)
}

// RenderPauseOverlay draws an unmistakable centered pause badge over the game area.
func RenderPauseOverlay(s *engine.Screen, w, h int, th engine.Theme) {
	boxW := 36
	boxH := 5
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	s.Fill(boxX, boxY, boxW, boxH, ' ', th.Text, tcell.ColorBlack)
	s.DoubleBoxWithTitle(boxX, boxY, boxW, boxH, "PAUSED", tcell.ColorYellow, tcell.ColorBlack, tcell.ColorYellow)
	s.CenterText(boxY+2, "Press P to Resume Playing", tcell.ColorWhite, tcell.ColorBlack)
}

// RenderStreakBadge draws the persistent streak badge at a given row.
func RenderStreakBadge(s *engine.Screen, y int, store *history.Store, isMonochrome bool) {
	if store == nil || len(store.Records) == 0 {
		s.CenterText(y, "[STREAK: 0 DAYS]", tcell.ColorGray, tcell.ColorBlack)
		return
	}
	streaks := history.CalculateStreaks(store.Records, time.Now())
	if streaks.CurrentStreak > 0 {
		badge := fmt.Sprintf("[🔥 STREAK: %d days]", streaks.CurrentStreak)
		fg := tcell.ColorYellow
		if isMonochrome {
			badge = fmt.Sprintf("[STREAK: %d DAYS]", streaks.CurrentStreak)
			fg = tcell.ColorWhite
		}
		s.CenterText(y, badge, fg, tcell.ColorBlack)
	} else {
		s.CenterText(y, "[STREAK: 0 days]", tcell.ColorGray, tcell.ColorBlack)
	}
}

// CenterBox coordinates helper.
func CenterBox(screenW, screenH, boxW, boxH int) (x, y int) {
	x = (screenW - boxW) / 2
	y = (screenH - boxH) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return x, y
}

// DrawCenteredDoubleBox draws a centered double-line frame with title.
func DrawCenteredDoubleBox(s *engine.Screen, w, h, boxW, boxH int, title string, borderFg, bg, titleFg tcell.Color) (int, int) {
	bx, by := CenterBox(w, h, boxW, boxH)
	s.DoubleBoxWithTitle(bx, by, boxW, boxH, title, borderFg, bg, titleFg)
	return bx, by
}
