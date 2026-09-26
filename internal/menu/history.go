package menu

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// HistoryScreen presents the browsable history log, achievements, and statistics summary.
type HistoryScreen struct {
	screen       *engine.Screen
	store        *history.Store
	achStore     *history.AchievementsStore
	activeTab    int
	scrollOffset int
	tabs         []struct {
		Label    string
		GameName string
	}
}

// NewHistoryScreen constructs the history view with achievements and activity heatmap support.
func NewHistoryScreen(s *engine.Screen, store *history.Store, achStore *history.AchievementsStore) *HistoryScreen {
	return &HistoryScreen{
		screen:       s,
		store:        store,
		achStore:     achStore,
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
			{Label: "Achievements", GameName: "achievements"},
		},
	}
}

// Show renders the history screen and handles keyboard interaction until exit.
func (hs *HistoryScreen) Show() {
	for {
		if hs.isAchievementsTab() {
			hs.renderAchievements()
		} else {
			records := hs.getFilteredRecords()
			hs.render(records)
		}
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
				hs.scrollOffset++
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
				hs.scrollOffset++
			case '1':
				hs.setTab(0)
			case '2':
				hs.setTab(1)
			case '3':
				hs.setTab(2)
			case '4':
				hs.setTab(3)
			case '5':
				hs.setTab(4)
			case 'q', 'Q':
				return
			}
		}
	}
}

func (hs *HistoryScreen) isAchievementsTab() bool {
	return hs.activeTab == 4
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

func (hs *HistoryScreen) renderTabBar(tabY int) {
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
}

func (hs *HistoryScreen) renderSummaryPanel(rightX, boxY, summaryW, boxH int) {
	hs.screen.BoxWithTitle(rightX, boxY, summaryW, boxH, "STATISTICS", tcell.ColorGold, tcell.ColorBlack, tcell.ColorWhite)

	sy := boxY + 1
	var allRecords []history.Record
	if hs.store != nil {
		allRecords = hs.store.Records
	}

	// Playtime & Games
	hs.screen.DrawText(rightX+2, sy, "Games:   "+strconv.Itoa(history.GamesPlayed(allRecords, "")), tcell.ColorWhite, tcell.ColorBlack)
	sy++
	hs.screen.DrawText(rightX+2, sy, "Playtime: "+formatDuration(history.TotalPlaytime(allRecords, "")), tcell.ColorWhite, tcell.ColorBlack)
	sy += 2

	// Streaks
	streaks := history.CalculateStreaks(allRecords, time.Now())
	hs.screen.DrawText(rightX+2, sy, "── STREAKS ──", tcell.ColorAqua, tcell.ColorBlack)
	sy++
	hs.screen.DrawText(rightX+2, sy, fmt.Sprintf("Current: %d days", streaks.CurrentStreak), tcell.ColorYellow, tcell.ColorBlack)
	sy++
	hs.screen.DrawText(rightX+2, sy, fmt.Sprintf("Longest: %d days", streaks.LongestStreak), tcell.ColorWhite, tcell.ColorBlack)
	sy += 2

	// 30-Day Activity Heatmap
	hs.screen.DrawText(rightX+2, sy, "── 30-DAY HEATMAP ──", tcell.ColorAqua, tcell.ColorBlack)
	sy++
	heatmap := history.GenerateHeatmap(allRecords, 20, time.Now())
	var glyphRow string
	for _, day := range heatmap {
		glyphRow += string(history.HeatmapGlyph(day.Count, false))
	}
	hs.screen.DrawText(rightX+2, sy, glyphRow, tcell.ColorLime, tcell.ColorBlack)
	sy++
	hs.screen.DrawText(rightX+2, sy, "· none  ▪ 1-2  █ 3+", tcell.ColorGray, tcell.ColorBlack)
	sy += 2

	// High Scores
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

		hs.screen.DrawText(rightX+4, sy, fmt.Sprintf("E: %d  M: %d  H: %d", easy, med, hard), tcell.ColorWhite, tcell.ColorBlack)
		sy++
	}
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
	hs.renderTabBar(3)

	// 3. Layout: Left Table (width w - 30), Right Summary Box (width 26)
	summaryW := 26
	leftW := w - summaryW - 5
	boxY := 5
	boxH := h - 8
	if boxH < 12 {
		boxH = 12
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
		maxOffset := len(records) - visibleRows
		if maxOffset < 0 {
			maxOffset = 0
		}
		if hs.scrollOffset > maxOffset {
			hs.scrollOffset = maxOffset
		}

		for i := 0; i < visibleRows; i++ {
			recIdx := hs.scrollOffset + i
			if recIdx >= len(records) {
				break
			}
			r := records[recIdx]
			rowY := boxY + 2 + i

			dateStr := r.PlayedAt.Format("01/02 15:04")
			durStr := formatDuration(r.Duration)
			gameStr := r.Game
			if r.Game == "snake" && r.Mode != "" && r.Mode != "classic" {
				gameStr = "snk:" + r.Mode
				if len(gameStr) > 9 {
					gameStr = gameStr[:9]
				}
			}
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
				dateStr, gameStr, diffStr, themeStr, r.Score, outcomeStr, durStr)

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
	hs.renderSummaryPanel(rightX, boxY, summaryW, boxH)

	// 4. Footer
	hs.screen.CenterText(h-2, "←/→ or 1-5: Tabs  |  ↑/↓: Scroll  |  Esc/Q/Enter: Main Menu", tcell.ColorGray, tcell.ColorBlack)
}

func (hs *HistoryScreen) renderAchievements() {
	hs.screen.Clear(tcell.ColorBlack)
	w, h := hs.screen.Size()

	if w < 76 || h < 20 {
		hs.screen.CenterText(h/2, "Terminal too small to show achievements. Please resize.", tcell.ColorYellow, tcell.ColorBlack)
		return
	}

	hs.screen.CenterText(1, "ARCADE ACHIEVEMENTS & BADGES", tcell.ColorAqua, tcell.ColorBlack)
	hs.renderTabBar(3)

	summaryW := 26
	leftW := w - summaryW - 5
	boxY := 5
	boxH := h - 8
	if boxH < 12 {
		boxH = 12
	}

	hs.screen.BoxWithTitle(2, boxY, leftW, boxH, "ACHIEVEMENT BADGES (15)", tcell.ColorGold, tcell.ColorBlack, tcell.ColorWhite)

	var items []history.Achievement
	if hs.achStore != nil {
		items = hs.achStore.Items
	} else {
		items = history.DefaultAchievements()
	}

	var records []history.Record
	if hs.store != nil {
		records = hs.store.Records
	}
	now := time.Now()

	// Render each achievement (takes 3-4 lines each)
	achLines := make([]struct {
		Text string
		Fg   tcell.Color
	}, 0)

	for _, a := range items {
		progress := history.CalculateProgressHint(a, records, now)
		var statusText string
		var statusColor tcell.Color
		if a.IsUnlocked() {
			statusText = fmt.Sprintf("★ [UNLOCKED] %s", a.Title)
			statusColor = tcell.ColorLime
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{statusText, statusColor})
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{"   " + a.Description, tcell.ColorWhite})
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{"   " + progress, tcell.ColorLimeGreen})
		} else {
			statusText = fmt.Sprintf("○ [LOCKED] %s", a.Title)
			statusColor = tcell.ColorGray
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{statusText, statusColor})
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{"   " + a.Description, tcell.ColorWhite})
			achLines = append(achLines, struct {
				Text string
				Fg   tcell.Color
			}{"   Progress: " + progress, tcell.ColorAqua})
			if a.Hint != "" {
				achLines = append(achLines, struct {
					Text string
					Fg   tcell.Color
				}{"   Hint: " + a.Hint, tcell.ColorGray})
			}
		}

		achLines = append(achLines, struct {
			Text string
			Fg   tcell.Color
		}{"", tcell.ColorBlack}) // separator
	}


	visibleRows := boxH - 2
	maxOffset := len(achLines) - visibleRows
	if maxOffset < 0 {
		maxOffset = 0
	}
	if hs.scrollOffset > maxOffset {
		hs.scrollOffset = maxOffset
	}

	for i := 0; i < visibleRows; i++ {
		idx := hs.scrollOffset + i
		if idx >= len(achLines) {
			break
		}
		line := achLines[idx]
		if len(line.Text) > leftW-4 {
			line.Text = line.Text[:leftW-4]
		}
		hs.screen.DrawText(4, boxY+1+i, line.Text, line.Fg, tcell.ColorBlack)
	}

	rightX := leftW + 3
	hs.renderSummaryPanel(rightX, boxY, summaryW, boxH)

	hs.screen.CenterText(h-2, "←/→ or 1-5: Tabs  |  ↑/↓: Scroll Badges  |  Esc/Q/Enter: Main Menu", tcell.ColorGray, tcell.ColorBlack)
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
