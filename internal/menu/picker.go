package menu

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// ConfigPicker prompts the player to select Difficulty and Theme before game launch.
type ConfigPicker struct {
	screen   *engine.Screen
	gameName string
	store    *history.Store
}

// NewConfigPicker constructs a configuration picker for the given game.
func NewConfigPicker(s *engine.Screen, gameName string, store *history.Store) *ConfigPicker {
	return &ConfigPicker{
		screen:   s,
		gameName: gameName,
		store:    store,
	}
}

// Pick executes the difficulty and theme selection workflow.
// Returns the built GameConfig and true if confirmed, or false if cancelled.
func (cp *ConfigPicker) Pick() (engine.GameConfig, bool) {
	// 1. Pick Difficulty
	diff, ok := cp.pickDifficulty()
	if !ok {
		return engine.GameConfig{}, false
	}

	// 2. Pick Theme
	th, ok := cp.pickTheme()
	if !ok {
		return engine.GameConfig{}, false
	}

	return engine.GameConfig{
		Difficulty: diff,
		Theme:      th,
	}, true
}

func (cp *ConfigPicker) renderPersonalBest(boxY int) {
	if cp.store == nil || len(cp.store.Records) == 0 {
		return
	}
	pb, found := history.PersonalBests(cp.store.Records, cp.gameName)
	if !found || pb.HighScore <= 0 {
		return
	}

	timeAgo := history.FormatRelativeTime(pb.HighScoreDate, time.Now())
	bestText := fmt.Sprintf("★ Personal Best: %d pts (%s, %s)", pb.HighScore, pb.HighScoreDiff, timeAgo)
	if pb.KeyMetricName != "" && pb.KeyMetricValue > 0 {
		bestText += fmt.Sprintf("  •  %s: %d", pb.KeyMetricName, pb.KeyMetricValue)
	}
	cp.screen.CenterText(boxY-2, bestText, tcell.ColorYellow, tcell.ColorBlack)
}

func (cp *ConfigPicker) pickDifficulty() (engine.Difficulty, bool) {
	options := []struct {
		Level engine.Difficulty
		Desc  string
	}{
		{engine.Easy, cp.difficultyDescription(engine.Easy)},
		{engine.Medium, cp.difficultyDescription(engine.Medium)},
		{engine.Hard, cp.difficultyDescription(engine.Hard)},
	}

	selected := 1 // Default to Medium

	for {
		cp.screen.Clear(tcell.ColorBlack)
		w, h := cp.screen.Size()

		boxW := 68
		boxH := len(options)*3 + 4
		boxX := (w - boxW) / 2
		boxY := (h - boxH) / 2

		cp.renderPersonalBest(boxY)

		title := "SELECT DIFFICULTY — " + cp.gameDisplayName()
		cp.screen.BoxWithTitle(boxX, boxY, boxW, boxH, title, tcell.ColorAqua, tcell.ColorBlack, tcell.ColorWhite)

		for i, opt := range options {
			rowY := boxY + 2 + i*3
			prefix := "  "
			fg := tcell.ColorWhite
			bg := tcell.ColorBlack

			if i == selected {
				prefix = "► "
				fg = tcell.ColorYellow
				bg = tcell.ColorDarkBlue
				cp.screen.Fill(boxX+1, rowY, boxW-2, 2, ' ', fg, bg)
			}

			cp.screen.DrawText(boxX+3, rowY, prefix+opt.Level.String(), fg, bg)
			cp.screen.DrawText(boxX+6, rowY+1, opt.Desc, tcell.ColorGray, bg)
		}

		cp.screen.CenterText(h-3, "↑/↓ or W/S: Select  |  Enter: Confirm  |  Esc/Q: Back to Menu", tcell.ColorGray, tcell.ColorBlack)
		cp.screen.Flush()

		ev := cp.screen.Raw.PollEvent()
		if ev == nil {
			return engine.Medium, false
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			cp.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyUp:
				selected--
				if selected < 0 {
					selected = len(options) - 1
				}
			case tcell.KeyDown:
				selected++
				if selected >= len(options) {
					selected = 0
				}
			case tcell.KeyEnter:
				return options[selected].Level, true
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return engine.Medium, false
			}

			switch e.Rune() {
			case 'w', 'W':
				selected--
				if selected < 0 {
					selected = len(options) - 1
				}
			case 's', 'S':
				selected++
				if selected >= len(options) {
					selected = 0
				}
			case ' ':
				return options[selected].Level, true
			case 'q', 'Q':
				return engine.Medium, false
			}
		}
	}
}

func (cp *ConfigPicker) pickTheme() (engine.Theme, bool) {
	themeNames := engine.ThemeNames()
	descriptions := map[string]string{
		"Retro Green": "Monochrome Nokia LCD green-on-black aesthetic",
		"Neon":        "Cyberpunk arcade cabinet palette (cyan, yellow, magenta)",
		"Monochrome":  "High-contrast pure ASCII (white on black, no color reliance)",
	}

	selected := 1 // Default to Neon

	for {
		cp.screen.Clear(tcell.ColorBlack)
		w, h := cp.screen.Size()

		boxW := 68
		boxH := len(themeNames)*3 + 4
		boxX := (w - boxW) / 2
		boxY := (h - boxH) / 2

		cp.renderPersonalBest(boxY)

		title := "SELECT VISUAL THEME — " + cp.gameDisplayName()
		cp.screen.BoxWithTitle(boxX, boxY, boxW, boxH, title, tcell.ColorFuchsia, tcell.ColorBlack, tcell.ColorWhite)

		for i, name := range themeNames {
			rowY := boxY + 2 + i*3
			prefix := "  "
			fg := tcell.ColorWhite
			bg := tcell.ColorBlack

			if i == selected {
				prefix = "► "
				fg = tcell.ColorYellow
				bg = tcell.ColorDarkBlue
				cp.screen.Fill(boxX+1, rowY, boxW-2, 2, ' ', fg, bg)
			}

			cp.screen.DrawText(boxX+3, rowY, prefix+name, fg, bg)
			cp.screen.DrawText(boxX+6, rowY+1, descriptions[name], tcell.ColorGray, bg)
		}

		cp.screen.CenterText(h-3, "↑/↓ or W/S: Select  |  Enter: Launch Game  |  Esc/Q: Back", tcell.ColorGray, tcell.ColorBlack)
		cp.screen.Flush()

		ev := cp.screen.Raw.PollEvent()
		if ev == nil {
			return engine.GetTheme("Neon"), false
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			cp.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyUp:
				selected--
				if selected < 0 {
					selected = len(themeNames) - 1
				}
			case tcell.KeyDown:
				selected++
				if selected >= len(themeNames) {
					selected = 0
				}
			case tcell.KeyEnter:
				return engine.GetTheme(themeNames[selected]), true
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return engine.GetTheme("Neon"), false
			}

			switch e.Rune() {
			case 'w', 'W':
				selected--
				if selected < 0 {
					selected = len(themeNames) - 1
				}
			case 's', 'S':
				selected++
				if selected >= len(themeNames) {
					selected = 0
				}
			case ' ':
				return engine.GetTheme(themeNames[selected]), true
			case 'q', 'Q':
				return engine.GetTheme("Neon"), false
			}
		}
	}
}

func (cp *ConfigPicker) gameDisplayName() string {
	switch cp.gameName {
	case "snake":
		return "SNAKE"
	case "pacman":
		return "PACMAN"
	case "ballplate":
		return "BALL & PLATE"
	default:
		return "GAME"
	}
}

func (cp *ConfigPicker) difficultyDescription(d engine.Difficulty) string {
	switch cp.gameName {
	case "snake":
		switch d {
		case engine.Easy:
			return "180ms tick, 40x20 arena, walls wrap around (safe)"
		case engine.Hard:
			return "90ms tick, 28x16 arena, walls kill, deadly obstacles spawn"
		default:
			return "130ms tick, 34x18 arena, walls kill"
		}
	case "pacman":
		switch d {
		case engine.Easy:
			return "2 ghosts, random-walk AI, 50-tick power pellet, 3 lives"
		case engine.Hard:
			return "4 ghosts, chase + ambush AI, 20-tick power pellet, 2 lives"
		default:
			return "3 ghosts, Manhattan chase AI, 35-tick power pellet, 3 lives"
		}
	case "ballplate":
		switch d {
		case engine.Easy:
			return "Wide plate (7 cells), slow ball speed, 4 lives, gapped bricks"
		case engine.Hard:
			return "Narrow plate (3 cells), fast ball, 2 lives, 2-hit tough bricks"
		default:
			return "Standard plate (5 cells), moderate ball speed, 3 lives, full grid"
		}
	default:
		return ""
	}
}
