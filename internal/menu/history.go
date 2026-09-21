package menu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
	"goarcade/internal/history"
)

// HistoryScreen presents the browsable history log and statistics summary.
type HistoryScreen struct {
	screen       *engine.Screen
	store        *history.Store
	activeTab    int
	scrollOffset int
	tabs         []struct {
		Label    string
		GameName string
	}
}

// NewHistoryScreen constructs the history view.
func NewHistoryScreen(s *engine.Screen, store *history.Store) *HistoryScreen {
	return &HistoryScreen{
		screen:       s,
		store:        store,
		activeTab:    0,
		scrollOffset: 0,
		tabs: []struct {
			Label    string
			GameName string
		}{
			{Label: "All", GameName: ""},
			{Label: "Snake", GameName: "snake"},
			{Label: "Pacman", GameName: "pacman"},
			{Label: "Ball & Plate", GameName: "ballplate"},
		},
	}
}

// Show renders the history screen and handles keyboard interaction until exit.
func (hs *HistoryScreen) Show() {
	for {
		records := hs.getFilteredRecords()
		hs.render(records)
		hs.screen.Flush()

		ev := hs.screen.Raw.PollEvent()
		if ev == nil {
			return
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			hs.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyLeft:
				hs.prevTab()
			case tcell.KeyRight, tcell.KeyTab:
				hs.nextTab()
			case tcell.KeyUp:
				if hs.scrollOffset > 0 {
					hs.scrollOffset--
				}
			case tcell.KeyDown:
				if hs.scrollOffset < len(records)-1 {
					hs.scrollOffset++
				}
			case tcell.KeyEnter:
				return
			}

			switch e.Rune() {
			case 'a', 'A':
				hs.prevTab()
			case 'd', 'D':
				hs.nextTab()
			case 'w', 'W':
				if hs.scrollOffset > 0 {
					hs.scrollOffset--
				}
			case 's', 'S':
				if hs.scrollOffset < len(records)-1 {
					hs.scrollOffset++
				}
			case '1':
				hs.setTab(0)
			case '2':
				hs.setTab(1)
			case '3':
				hs.setTab(2)
			case '4':
				hs.setTab(3)
			case 'q', 'Q':
				return
			}
		}
	}
}

func (hs *HistoryScreen) prevTab() {
	hs.activeTab--
	if hs.activeTab < 0 {
		hs.activeTab = len(hs.tabs) - 1
	}
	hs.scrollOffset = 0
}

func (hs *HistoryScreen) nextTab() {
	hs.activeTab++
	if hs.activeTab >= len(hs.tabs) {
		hs.activeTab = 0
	}
	hs.scrollOffset = 0
}

func (hs *HistoryScreen) setTab(idx int) {
	if idx >= 0 && idx < len(hs.tabs) {
		hs.activeTab = idx
		hs.scrollOffset = 0
	}
}

func (hs *HistoryScreen) getFilteredRecords() []history.Record {
	if hs.store == nil {
		return nil
	}
	gameFilter := hs.tabs[hs.activeTab].GameName
	return history.Recent(hs.store.Records, gameFilter, 0)
}

func (hs *HistoryScreen) render(records []history.Record) {
	hs.screen.Clear(tcell.ColorBlack)
	w, h := hs.screen.Size()

	if w < 76 || h < 20 {
		hs.screen.CenterText(h/2, "Terminal too small to show history. Please resize.", tcell.ColorYellow, tcell.ColorBlack)
		return
	}

	// 1. Header
	hs.screen.CenterText(1, "ARCADE HISTORY & STATISTICS", tcell.ColorAqua, tcell.ColorBlack)

	// 2. Tab Bar
	tabY := 3
	tabX := 2
	for i, tab := range hs.tabs {
		tabText := fmt.Sprintf(" [%d] %s ", i+1, tab.Label)
		fg := tcell.ColorWhite
		bg := tcell.ColorBlack
		if i == hs.activeTab {
			fg = tcell.ColorBlack
			bg = tcell.ColorAqua
		}
		hs.screen.DrawText(tabX, tabY, tabText, fg, bg)
		tabX += len(tabText) + 2
	}

	// 3. Layout: Left Table (width w - 28), Right Summary Box (width 24)
	summaryW := 24
	leftW := w - summaryW - 5
	boxY := 5
	boxH := h - 8
	if boxH < 10 {
		boxH = 10
	}

	// Left: Session Table
	hs.screen.BoxWithTitle(2, boxY, leftW, boxH, "RECENT SESSIONS", tcell.ColorBlueViolet, tcell.ColorBlack, tcell.ColorWhite)

	// Table column headers
	headerY := boxY + 1
	header := fmt.Sprintf("  %-10s %-9s %-7s %-10s %-6s %-7s %s", "DATE", "GAME", "DIFF", "THEME", "SCORE", "RESULT", "TIME")
	hs.screen.DrawText(3, headerY, header, tcell.ColorYellow, tcell.ColorBlack)

	visibleRows := boxH - 3
	if len(records) == 0 {
		hs.screen.DrawText(4, boxY+3, "No sessions recorded yet.", tcell.ColorGray, tcell.ColorBlack)
	} else {
		for i := 0; i < visibleRows; i++ {
			recIdx := hs.scrollOffset + i
			if recIdx >= len(records) {
				break
			}
			r := records[recIdx]
			rowY := boxY + 2 + i

			dateStr := r.PlayedAt.Format("01/02 15:04")
			durStr := formatDuration(r.Duration)
			diffStr := r.Difficulty
			if len(diffStr) > 6 {
				diffStr = diffStr[:6]
			}
			themeStr := r.Theme
			if len(themeStr) > 9 {
				themeStr = themeStr[:9]
			}
			outcomeStr := r.Outcome
			if len(outcomeStr) > 6 {
				outcomeStr = outcomeStr[:6]
			}

			line := fmt.Sprintf("  %-10s %-9s %-7s %-10s %-6d %-7s %s",
				dateStr, r.Game, diffStr, themeStr, r.Score, outcomeStr, durStr)

			fg := tcell.ColorWhite
			if r.Outcome == "won" {
				fg = tcell.ColorLime
			} else if r.Outcome == "died" || r.Outcome == "wall" || r.Outcome == "self" || r.Outcome == "lost" {
				fg = tcell.ColorOrangeRed
			}

			hs.screen.DrawText(3, rowY, line, fg, tcell.ColorBlack)
		}
	}

	// Right: Summary & Stats Box
	rightX := leftW + 3
	hs.screen.BoxWithTitle(rightX, boxY, summaryW, boxH, "STATISTICS", tcell.ColorGold, tcell.ColorBlack, tcell.ColorWhite)

	sy := boxY + 1
	allRecords := hs.store.Records
	hs.screen.DrawText(rightX+2, sy, "Total Games: "+strconv.Itoa(history.GamesPlayed(allRecords, "")), tcell.ColorWhite, tcell.ColorBlack)
	sy++
	hs.screen.DrawText(rightX+2, sy, "Total Time:  "+formatDuration(history.TotalPlaytime(allRecords, "")), tcell.ColorWhite, tcell.ColorBlack)
	sy += 2

	hs.screen.DrawText(rightX+2, sy, "── HIGH SCORES ──", tcell.ColorAqua, tcell.ColorBlack)
	sy++

	games := []struct {
		name    string
		display string
	}{
		{"snake", "Snake"},
		{"pacman", "Pacman"},
		{"ballplate", "Ball & Plate"},
	}

	for _, g := range games {
		hs.screen.DrawText(rightX+2, sy, g.display+":", tcell.ColorYellow, tcell.ColorBlack)
		sy++
		easy := history.HighScore(allRecords, g.name, "easy")
		med := history.HighScore(allRecords, g.name, "medium")
		hard := history.HighScore(allRecords, g.name, "hard")

		hs.screen.DrawText(rightX+4, sy, fmt.Sprintf("Easy:   %d", easy), tcell.ColorWhite, tcell.ColorBlack)
		sy++
		hs.screen.DrawText(rightX+4, sy, fmt.Sprintf("Medium: %d", med), tcell.ColorWhite, tcell.ColorBlack)
		sy++
		hs.screen.DrawText(rightX+4, sy, fmt.Sprintf("Hard:   %d", hard), tcell.ColorWhite, tcell.ColorBlack)
		sy++
	}

	// 4. Footer
	hs.screen.CenterText(h-2, "←/→ or 1-4: Switch Tabs  |  ↑/↓: Scroll  |  Esc/Q/Enter: Back to Main Menu", tcell.ColorGray, tcell.ColorBlack)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}
