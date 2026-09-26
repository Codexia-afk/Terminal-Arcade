// Package menu provides terminal-native UI screens for game selection, configuration, and stats.
package menu

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// Choice represents a main menu selection.
type Choice int

const (
	ChoiceSnake Choice = iota
	ChoicePacman
	ChoiceBallPlate
	ChoiceHistory
	ChoiceExport
	ChoiceResetData
	ChoiceSettings
	ChoiceQuit
)

// MainMenuItem captures display details for a menu row.
type MainMenuItem struct {
	Choice Choice
	Label  string
	Desc   string
}

// MainMenu manages navigation and rendering for the title selection screen.
type MainMenu struct {
	screen   *engine.Screen
	store    *history.Store
	selected int
	items    []MainMenuItem
	lastBoxX int
	lastBoxY int
	lastBoxW int
	lastBoxH int
}

// NewMainMenu constructs the main menu with history awareness for streak display.
func NewMainMenu(s *engine.Screen, store *history.Store) *MainMenu {
	return &MainMenu{
		screen:   s,
		store:    store,
		selected: 0,
		items: []MainMenuItem{
			{ChoiceSnake, "1. Snake", "5 gameplay modes (Classic, Zen, Survival, Time Attack, Obstacle) × 3 tiers"},
			{ChoicePacman, "2. Pacman", "Authentic maze chase with 4 distinct ghost AI personalities"},
			{ChoiceBallPlate, "3. Ball & Plate", "Breakout with 3-zone paddle physics, speed gauge & tough bricks"},
			{ChoiceHistory, "4. History, Stats & Achievements", "View past sessions, 15 achievements, streaks & 30-day heatmap"},
			{ChoiceExport, "5. Export Data", "Save sessions and badges to local JSON or CSV file"},
			{ChoiceResetData, "6. Reset Saved Data", "Clear all recorded session history and achievements"},
			{ChoiceSettings, "7. Settings", "Customize theme, audio bell, 60fps refresh & view keybindings"},
			{ChoiceQuit, "8. Quit", "Exit back to terminal"},
		},
	}
}

// Show renders the main menu and blocks until the player makes a selection or quits.
func (m *MainMenu) Show() Choice {
	for {
		m.render()
		m.screen.Flush()

		ev := m.screen.Raw.PollEvent()
		if ev == nil {
			return ChoiceQuit
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			m.screen.Raw.Sync()
		case *tcell.EventMouse:
			if e.Buttons()&tcell.ButtonPrimary != 0 {
				mx, my := e.Position()
				// Check if click was inside menu container
				if mx >= m.lastBoxX && mx < m.lastBoxX+m.lastBoxW {
					for i := range m.items {
						itemY := m.lastBoxY + 1 + i*1
						if my == itemY {
							m.selected = i
							return m.items[i].Choice
						}
					}
				}
			}
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyUp:
				m.selected--
				if m.selected < 0 {
					m.selected = len(m.items) - 1
				}
			case tcell.KeyDown:
				m.selected++
				if m.selected >= len(m.items) {
					m.selected = 0
				}
			case tcell.KeyEnter:
				return m.items[m.selected].Choice
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return ChoiceQuit
			}

			switch e.Rune() {
			case 'w', 'W':
				m.selected--
				if m.selected < 0 {
					m.selected = len(m.items) - 1
				}
			case 's', 'S':
				m.selected++
				if m.selected >= len(m.items) {
					m.selected = 0
				}
			case ' ':
				return m.items[m.selected].Choice
			case '1':
				return ChoiceSnake
			case '2':
				return ChoicePacman
			case '3':
				return ChoiceBallPlate
			case '4':
				return ChoiceHistory
			case '5':
				return ChoiceExport
			case '6':
				return ChoiceResetData
			case '7':
				return ChoiceSettings
			case '8', 'q', 'Q':
				return ChoiceQuit
			}
		}
	}
}

func (m *MainMenu) render() {
	m.screen.Clear(tcell.ColorBlack)
	w, h := m.screen.Size()

	banner := []string{
		`  ██████╗  ██████╗      █████╗ ██████╗  ██████╗ █████╗ ██████╗ ███████╗`,
		` ██╔════╝ ██╔═══██╗    ██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝`,
		` ██║  ███╗██║   ██║    ███████║██████╔╝██║     ███████║██║  ██║█████╗  `,
		` ██║   ██║██║   ██║    ██╔══██║██╔══██╗██║     ██╔══██║██║  ██║██╔══╝  `,
		` ╚██████╔╝╚██████╔╝    ╚═╝  ╚═╝╚═╝  ╚═╝╚██████╗╚═╝  ╚═╝██████╔╝███████╗`,
	}

	totalH := 24
	startY := (h - totalH) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Draw Title Banner
	for i, line := range banner {
		m.screen.CenterText(startY+i, line, tcell.ColorAqua, tcell.ColorBlack)
	}

	m.screen.CenterText(startY+6, "Terminal Arcade Suite  •  Offline Edition", tcell.ColorWhite, tcell.ColorBlack)

	// Always visible streak badge (habit reminder)
	streak := 0
	if m.store != nil && len(m.store.Records) > 0 {
		st := history.CalculateStreaks(m.store.Records, time.Now())
		streak = st.CurrentStreak
	}
	streakBadge := fmt.Sprintf("[🔥 STREAK: %d days]", streak)
	m.screen.CenterText(startY+7, streakBadge, tcell.ColorYellow, tcell.ColorBlack)

	// 2. Centered Menu Container Box (~60 chars wide × 12 rows tall)
	boxW := 60
	boxH := len(m.items) + 3
	boxX := (w - boxW) / 2
	boxY := startY + 9
	if boxX < 0 {
		boxX = 0
	}

	m.lastBoxX = boxX
	m.lastBoxY = boxY
	m.lastBoxW = boxW
	m.lastBoxH = boxH

	m.screen.BoxWithTitle(boxX, boxY, boxW, boxH, "MAIN MENU", tcell.ColorBlueViolet, tcell.ColorBlack, tcell.ColorWhite)

	for i, item := range m.items {
		rowY := boxY + 1 + i
		prefix := "  "
		fg := tcell.ColorWhite
		bg := tcell.ColorBlack

		if i == m.selected {
			prefix = "► "
			fg = tcell.ColorYellow
			bg = tcell.ColorDarkBlue
			m.screen.Fill(boxX+1, rowY, boxW-2, 1, ' ', fg, bg)
		}

		m.screen.DrawText(boxX+3, rowY, prefix+item.Label, fg, bg)
	}

	// Active description below menu
	m.screen.CenterText(boxY+boxH+1, m.items[m.selected].Desc, tcell.ColorAqua, tcell.ColorBlack)

	// 3. Instruction footer
	m.screen.CenterText(h-2, "↑/↓ Navigate | Enter Select | Q Quit", tcell.ColorGray, tcell.ColorBlack)
}
