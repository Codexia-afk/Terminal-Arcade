package menu

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// TitleScreen displays the premium title splash screen shown on application launch.
type TitleScreen struct {
	screen  *engine.Screen
	store   *history.Store
	version string
}

// NewTitleScreen creates a TitleScreen.
func NewTitleScreen(s *engine.Screen, store *history.Store, version string) *TitleScreen {
	return &TitleScreen{
		screen:  s,
		store:   store,
		version: version,
	}
}

// Show renders the title screen and blocks until any key or mouse click is received.
func (ts *TitleScreen) Show() {
	for {
		ts.render()
		ts.screen.Flush()

		ev := ts.screen.Raw.PollEvent()
		if ev == nil {
			return
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			ts.screen.Raw.Sync()
		case *tcell.EventKey:
			return // Any key enters main menu
		case *tcell.EventMouse:
			if e.Buttons()&tcell.ButtonPrimary != 0 {
				return
			}
		}
	}
}

func (ts *TitleScreen) render() {
	ts.screen.Clear(tcell.ColorBlack)
	_, h := ts.screen.Size()

	banner := []string{
		`   ██████╗  ██████╗      █████╗ ██████╗  ██████╗ █████╗ ██████╗ ███████╗`,
		`  ██╔════╝ ██╔═══██╗    ██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝`,
		`  ██║  ███╗██║   ██║    ███████║██████╔╝██║     ███████║██║  ██║█████╗  `,
		`  ██║   ██║██║   ██║    ██╔══██║██╔══██╗██║     ██╔══██║██║  ██║██╔══╝  `,
		`  ╚██████╔╝╚██████╔╝    ╚═╝  ╚═╝╚═╝  ╚═╝╚██████╗╚═╝  ╚═╝██████╔╝███████╗`,
	}

	totalH := 20
	startY := (h - totalH) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. ASCII Block Letters
	for i, line := range banner {
		ts.screen.CenterText(startY+i, line, tcell.ColorAqua, tcell.ColorBlack)
	}

	// 2. Tagline
	ts.screen.CenterText(startY+6, "Terminal Arcade Suite  •  Offline Edition", tcell.ColorWhite, tcell.ColorBlack)

	// 3. Current streak badge (habit reminder)
	streak := 0
	if ts.store != nil && len(ts.store.Records) > 0 {
		st := history.CalculateStreaks(ts.store.Records, time.Now())
		streak = st.CurrentStreak
	}
	streakBadge := fmt.Sprintf("[🔥 STREAK: %d days]", streak)
	if streak == 0 {
		streakBadge = "[🔥 STREAK: 0 days]"
	}
	ts.screen.CenterText(startY+8, streakBadge, tcell.ColorYellow, tcell.ColorBlack)

	// 4. Personal Best summary
	pbSummary := "Your best: None yet — play your first session!"
	if ts.store != nil && len(ts.store.Records) > 0 {
		bestScore := 0
		bestGame := ""
		for _, r := range ts.store.Records {
			if r.Score > bestScore {
				bestScore = r.Score
				bestGame = r.Game
				if r.Game == "snake" {
					bestGame = fmt.Sprintf("Snake %s", stringsTitle(r.EffectiveMode()))
				} else if r.Game == "ballplate" {
					bestGame = "Ball & Plate"
				} else if r.Game == "pacman" {
					bestGame = "Pacman"
				}
			}
		}
		if bestScore > 0 {
			pbSummary = fmt.Sprintf("Your best: %d points (%s)", bestScore, bestGame)
		}
	}
	ts.screen.CenterText(startY+10, pbSummary, tcell.ColorLightCyan, tcell.ColorBlack)

	// 5. One-line prompt
	ts.screen.CenterText(startY+13, "Press any key to enter...", tcell.ColorLime, tcell.ColorBlack)

	// 6. Centered footer with version and copyright
	footer := fmt.Sprintf("v%s  •  © 2025-2026 Go Arcade Project", ts.version)
	ts.screen.CenterText(h-2, footer, tcell.ColorGray, tcell.ColorBlack)
}

func stringsTitle(s string) string {
	if len(s) == 0 {
		return ""
	}
	first := s[0]
	if first >= 'a' && first <= 'z' {
		first -= 32
	}
	return string(first) + s[1:]
}
