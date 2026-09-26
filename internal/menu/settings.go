package menu

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// SettingsScreen provides configuration for global theme defaults, audio, and keybindings.
type SettingsScreen struct {
	screen       *engine.Screen
	selectedItem int
}

// NewSettingsScreen constructs a settings screen.
func NewSettingsScreen(s *engine.Screen) *SettingsScreen {
	return &SettingsScreen{
		screen:       s,
		selectedItem: 0,
	}
}

// Show renders the settings screen and handles user interactions.
func (ss *SettingsScreen) Show() {
	for {
		ss.render()
		ss.screen.Flush()

		ev := ss.screen.Raw.PollEvent()
		if ev == nil {
			return
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			ss.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyUp:
				ss.selectedItem--
				if ss.selectedItem < 0 {
					ss.selectedItem = 2
				}
			case tcell.KeyDown:
				ss.selectedItem++
				if ss.selectedItem > 2 {
					ss.selectedItem = 0
				}
			case tcell.KeyEnter:
				if ss.selectedItem == 0 {
					engine.AudioBellEnabled = !engine.AudioBellEnabled
					if engine.AudioBellEnabled {
						engine.Beep()
					}
				} else if ss.selectedItem == 2 {
					return
				}
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			}

			switch e.Rune() {
			case 'w', 'W':
				ss.selectedItem--
				if ss.selectedItem < 0 {
					ss.selectedItem = 2
				}
			case 's', 'S':
				ss.selectedItem++
				if ss.selectedItem > 2 {
					ss.selectedItem = 0
				}
			case ' ':
				if ss.selectedItem == 0 {
					engine.AudioBellEnabled = !engine.AudioBellEnabled
					if engine.AudioBellEnabled {
						engine.Beep()
					}
				}
			case 'q', 'Q':
				return
			}
		case *tcell.EventMouse:
			if e.Buttons()&tcell.ButtonPrimary != 0 {
				_, my := e.Position()
				_, h := ss.screen.Size()
				boxH := 18
				boxY := (h - boxH) / 2
				if my == boxY+3 {
					ss.selectedItem = 0
					engine.AudioBellEnabled = !engine.AudioBellEnabled
				} else if my == boxY+5 {
					ss.selectedItem = 1
				} else if my == boxY+15 {
					return
				}
			}
		}
	}
}

func (ss *SettingsScreen) render() {
	ss.screen.Clear(tcell.ColorBlack)
	w, h := ss.screen.Size()

	boxW := 62
	boxH := 18
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2
	if boxX < 0 {
		boxX = 0
	}
	if boxY < 0 {
		boxY = 0
	}

	ss.screen.DoubleBoxWithTitle(boxX, boxY, boxW, boxH, "SYSTEM SETTINGS", tcell.ColorAqua, tcell.ColorBlack, tcell.ColorAqua)

	// 1. Audio Bell
	bellStatus := "ENABLED  [Bell sounds on food & events]"
	if !engine.AudioBellEnabled {
		bellStatus = "DISABLED [Silent terminal operation]"
	}
	row0 := fmt.Sprintf("1. Terminal Bell Audio: %s", bellStatus)
	fg0 := tcell.ColorWhite
	bg0 := tcell.ColorBlack
	if ss.selectedItem == 0 {
		fg0 = tcell.ColorYellow
		bg0 = tcell.ColorDarkBlue
		ss.screen.Fill(boxX+1, boxY+3, boxW-2, 1, ' ', fg0, bg0)
	}
	ss.screen.DrawText(boxX+3, boxY+3, row0, fg0, bg0)

	// 2. Refresh Rate
	row1 := "2. Render Engine:       60 FPS Equivalent (16ms decoupled)"
	fg1 := tcell.ColorWhite
	bg1 := tcell.ColorBlack
	if ss.selectedItem == 1 {
		fg1 = tcell.ColorYellow
		bg1 = tcell.ColorDarkBlue
		ss.screen.Fill(boxX+1, boxY+5, boxW-2, 1, ' ', fg1, bg1)
	}
	ss.screen.DrawText(boxX+3, boxY+5, row1, fg1, bg1)

	// 3. Keybindings Reference
	ss.screen.DrawText(boxX+3, boxY+7, "── KEYBINDINGS REFERENCE ──", tcell.ColorGold, tcell.ColorBlack)
	ss.screen.DrawText(boxX+3, boxY+9, "Movement:      ↑↓←→  or  W A S D", tcell.ColorLightCyan, tcell.ColorBlack)
	ss.screen.DrawText(boxX+3, boxY+10, "Confirm:       Enter  or  Spacebar", tcell.ColorLightCyan, tcell.ColorBlack)
	ss.screen.DrawText(boxX+3, boxY+11, "Pause Game:    P", tcell.ColorLightCyan, tcell.ColorBlack)
	ss.screen.DrawText(boxX+3, boxY+12, "Quit to Menu:  Q  or  Escape", tcell.ColorLightCyan, tcell.ColorBlack)
	ss.screen.DrawText(boxX+3, boxY+13, "Mouse:         Click menu items & buttons directly", tcell.ColorLightCyan, tcell.ColorBlack)

	// 4. Return item
	row2 := "3. Return to Main Menu"
	fg2 := tcell.ColorWhite
	bg2 := tcell.ColorBlack
	if ss.selectedItem == 2 {
		fg2 = tcell.ColorYellow
		bg2 = tcell.ColorDarkBlue
		ss.screen.Fill(boxX+1, boxY+15, boxW-2, 1, ' ', fg2, bg2)
	}
	ss.screen.DrawText(boxX+3, boxY+15, row2, fg2, bg2)

	ss.screen.CenterText(h-2, "↑/↓: Navigate  |  Enter/Space: Toggle/Select  |  Esc/Q: Back", tcell.ColorGray, tcell.ColorBlack)
}
