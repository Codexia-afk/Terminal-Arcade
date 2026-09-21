// Package menu provides terminal-native UI screens for game selection, configuration, and stats.
package menu

import (
	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
)

// Choice represents a main menu selection.
type Choice int

const (
	ChoiceSnake Choice = iota
	ChoicePacman
	ChoiceBallPlate
	ChoiceHistory
	ChoiceQuit
)

// MainMenu manages navigation and rendering for the title selection screen.
type MainMenu struct {
	screen   *engine.Screen
	selected int
	items    []struct {
		Choice Choice
		Label  string
		Desc   string
	}
}

// NewMainMenu constructs the main menu.
func NewMainMenu(s *engine.Screen) *MainMenu {
	return &MainMenu{
		screen:   s,
		selected: 0,
		items: []struct {
			Choice Choice
			Label  string
			Desc   string
		}{
			{ChoiceSnake, "1. Snake", "Classic Nokia arcade snake with wrap and obstacle modes"},
			{ChoicePacman, "2. Pacman", "Arcade maze chase with distinct ghost personalities"},
			{ChoiceBallPlate, "3. Ball & Plate", "Breakout with angular deflection physics & tough bricks"},
			{ChoiceHistory, "4. History & Stats", "View past session logs, personal bests, and total playtime"},
			{ChoiceQuit, "5. Quit", "Exit back to terminal"},
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
			case '5', 'q', 'Q':
				return ChoiceQuit
			}
		}
	}
}

func (m *MainMenu) render() {
	m.screen.Clear(tcell.ColorBlack)
	w, h := m.screen.Size()

	banner := []string{
		`   ____  ___        _    ____   ____    _    ____  _____ `,
		`  / ___|/ _ \      / \  |  _ \ / ___|  / \  |  _ \| ____|`,
		` | |  _| | | |    / _ \ | |_) | |     / _ \ | | | |  _|  `,
		` | |_| | |_| |   / ___ \|  _ <| |___ / ___ \| |_| | |___ `,
		`  \____|\___/   /_/   \_\_| \_\\____/_/   \_\____/|_____|`,
	}

	startY := (h - 22) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Draw Title Banner
	for i, line := range banner {
		m.screen.CenterText(startY+i, line, tcell.ColorAqua, tcell.ColorBlack)
	}

	m.screen.CenterText(startY+6, "Terminal Arcade Suite  •  Offline Edition", tcell.ColorWhite, tcell.ColorBlack)

	// 2. Menu Items Container
	boxW := 60
	boxH := len(m.items)*2 + 3
	boxX := (w - boxW) / 2
	boxY := startY + 8

	m.screen.BoxWithTitle(boxX, boxY, boxW, boxH, "SELECT OPTION", tcell.ColorBlueViolet, tcell.ColorBlack, tcell.ColorWhite)

	for i, item := range m.items {
		rowY := boxY + 1 + i*2
		prefix := "  "
		fg := tcell.ColorWhite
		bg := tcell.ColorBlack

		if i == m.selected {
			prefix = "► "
			fg = tcell.ColorYellow
			bg = tcell.ColorDarkBlue
			// Fill highlight bar
			m.screen.Fill(boxX+1, rowY, boxW-2, 1, ' ', fg, bg)
		}

		m.screen.DrawText(boxX+3, rowY, prefix+item.Label, fg, bg)
		// Right-aligned or sub-description
		if i == m.selected {
			descY := boxY + boxH + 1
			m.screen.CenterText(descY, item.Desc, tcell.ColorAqua, tcell.ColorBlack)
		}
	}

	// 3. Footer instructions
	m.screen.CenterText(h-2, "Navigation: ↑/↓, W/S, 1-5  |  Select: Enter/Space  |  Quit: Q/Esc", tcell.ColorGray, tcell.ColorBlack)
}
