package menu

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// ConfigPicker prompts the player to select Mode (if applicable), Difficulty, and Theme before launch.
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

// Pick executes the mode, difficulty and theme selection workflow.
// Returns the built GameConfig and true if confirmed, or false if cancelled.
func (cp *ConfigPicker) Pick() (engine.GameConfig, bool) {
	mode := "classic"

	// 1. Pick Mode if Snake
	if cp.gameName == "snake" {
		m, ok := cp.pickSnakeMode()
		if !ok {
			return engine.GameConfig{}, false
		}
		mode = m
	}

	// 2. Pick Difficulty
	diff, ok := cp.pickDifficulty(mode)
	if !ok {
		return engine.GameConfig{}, false
	}

	// 3. Pick Theme with Live Preview
	th, ok := cp.pickTheme()
	if !ok {
		return engine.GameConfig{}, false
	}

	return engine.GameConfig{
		Mode:       mode,
		Difficulty: diff,
		Theme:      th,
	}, true
}

func (cp *ConfigPicker) renderPersonalBest(boxY int, mode string) {
	if cp.store == nil || len(cp.store.Records) == 0 {
		return
	}
	pb, found := history.PersonalBestsForMode(cp.store.Records, cp.gameName, mode)
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

func (cp *ConfigPicker) pickSnakeMode() (string, bool) {
	modes := []struct {
		ID    string
		Title string
		Desc  string
	}{
		{"classic", "1. CLASSIC", "Survive waves of self-collision (pure Nokia baseline)"},
		{"zen", "2. ZEN", "Peaceful growth — no threats, continuous flow scoring"},
		{"survival", "3. SURVIVAL", "Dodge enemies, clear 5 waves of pursuit AI"},
		{"time_attack", "4. TIME ATTACK", "Race against the clock — score big with speed bonuses"},
		{"obstacle", "5. OBSTACLE CHALLENGE", "Navigate 3 puzzle labyrinth mazes to victory"},
	}

	selected := 0

	for {
		cp.screen.Clear(tcell.ColorBlack)
		w, h := cp.screen.Size()

		boxW := 68
		boxH := len(modes)*3 + 4
		boxX := (w - boxW) / 2
		boxY := (h - boxH) / 2

		cp.renderPersonalBest(boxY, "")

		title := "SELECT GAMEPLAY MODE — SNAKE"
		cp.screen.BoxWithTitle(boxX, boxY, boxW, boxH, title, tcell.ColorLime, tcell.ColorBlack, tcell.ColorWhite)

		for i, m := range modes {
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

			cp.screen.DrawText(boxX+3, rowY, prefix+m.Title, fg, bg)
			cp.screen.DrawText(boxX+6, rowY+1, m.Desc, tcell.ColorGray, bg)
		}

		cp.screen.CenterText(h-3, "↑/↓ or W/S: Select  |  Enter: Confirm Mode  |  Esc/Q: Back", tcell.ColorGray, tcell.ColorBlack)
		cp.screen.Flush()

		ev := cp.screen.Raw.PollEvent()
		if ev == nil {
			return "classic", false
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			cp.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyUp:
				selected--
				if selected < 0 {
					selected = len(modes) - 1
				}
			case tcell.KeyDown:
				selected++
				if selected >= len(modes) {
					selected = 0
				}
			case tcell.KeyEnter:
				return modes[selected].ID, true
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return "classic", false
			}

			switch e.Rune() {
			case 'w', 'W':
				selected--
				if selected < 0 {
					selected = len(modes) - 1
				}
			case 's', 'S':
				selected++
				if selected >= len(modes) {
					selected = 0
				}
			case '1':
				return "classic", true
			case '2':
				return "zen", true
			case '3':
				return "survival", true
			case '4':
				return "time_attack", true
			case '5':
				return "obstacle", true
			case ' ':
				return modes[selected].ID, true
			case 'q', 'Q':
				return "classic", false
			}
		}
	}
}

func (cp *ConfigPicker) pickDifficulty(mode string) (engine.Difficulty, bool) {
	options := []struct {
		Level engine.Difficulty
		Desc  string
	}{
		{engine.Easy, cp.difficultyDescription(engine.Easy, mode)},
		{engine.Medium, cp.difficultyDescription(engine.Medium, mode)},
		{engine.Hard, cp.difficultyDescription(engine.Hard, mode)},
	}

	selected := 1 // Default to Medium

	for {
		cp.screen.Clear(tcell.ColorBlack)
		w, h := cp.screen.Size()

		boxW := 68
		boxH := len(options)*3 + 4
		boxX := (w - boxW) / 2
		boxY := (h - boxH) / 2

		cp.renderPersonalBest(boxY, mode)

		title := "SELECT DIFFICULTY — " + cp.gameDisplayName()
		if mode != "" && mode != "classic" {
			title += " [" + mode + "]"
		}
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

		cp.screen.CenterText(h-3, "↑/↓ or W/S: Select  |  Enter: Confirm  |  Esc/Q: Back", tcell.ColorGray, tcell.ColorBlack)
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
		"Retro Green": "Monochrome Nokia LCD green aesthetic",
		"Neon":        "Arcade neon (cyan, yellow, magenta)",
		"Monochrome":  "High-contrast pure ASCII b&w",
		"Cyberpunk":   "Electric cyan, dark grid & magenta",
		"Ocean":       "Deep navy, aquamarine & seafoam",
	}

	selected := 1 // Default to Neon

	for {
		cp.screen.Clear(tcell.ColorBlack)
		w, h := cp.screen.Size()

		totalW := 74
		listW := 40
		previewW := 32
		boxH := 18

		startX := (w - totalW) / 2
		startY := (h - boxH) / 2
		if startY < 1 {
			startY = 1
		}

		// Left: Theme list box
		title := "SELECT THEME"
		cp.screen.BoxWithTitle(startX, startY, listW, boxH, title, tcell.ColorFuchsia, tcell.ColorBlack, tcell.ColorWhite)

		for i, name := range themeNames {
			rowY := startY + 2 + i*3
			prefix := "  "
			fg := tcell.ColorWhite
			bg := tcell.ColorBlack

			if i == selected {
				prefix = "► "
				fg = tcell.ColorYellow
				bg = tcell.ColorDarkBlue
				cp.screen.Fill(startX+1, rowY, listW-2, 2, ' ', fg, bg)
			}

			cp.screen.DrawText(startX+3, rowY, prefix+name, fg, bg)
			cp.screen.DrawText(startX+6, rowY+1, descriptions[name], tcell.ColorGray, bg)
		}

		// Right: Live Theme Previewer!
		currentTheme := engine.GetTheme(themeNames[selected])
		previewX := startX + listW + 2
		RenderThemePreview(cp.screen, currentTheme, previewX, startY, previewW, boxH, cp.gameName)

		cp.screen.CenterText(h-2, "↑/↓ or W/S: Select  |  Enter: Confirm Theme  |  Esc/Q: Back", tcell.ColorGray, tcell.ColorBlack)
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

func (cp *ConfigPicker) difficultyDescription(d engine.Difficulty, mode string) string {
	switch cp.gameName {
	case "snake":
		switch mode {
		case "zen":
			switch d {
			case engine.Easy:
				return "Arena: 60×25 | Speed: 150ms | Walls wrap | Meditative pace"
			case engine.Hard:
				return "Arena: 40×15 | Speed: 100ms | Walls wrap | Tighter flow"
			default: // Medium
				return "Arena: 50×20 | Speed: 120ms | Walls wrap | Balanced flow"
			}
		case "survival":
			switch d {
			case engine.Easy:
				return "Arena: 45×20 | Speed: 100ms | 3 Lives | Enemies: 180ms"
			case engine.Hard:
				return "Arena: 35×15 | Speed: 100ms | 1 Life  | Enemies: 100ms (Fast)"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 2 Lives | Enemies: 140ms"
			}
		case "time_attack":
			switch d {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 90s Timer | Wrapping walls"
			case engine.Hard:
				return "Arena: 35×15 | Speed: 80ms  | 45s Timer | Hard walls"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 60s Timer | Hard walls"
			}
		case "obstacle":
			switch d {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 3 Mazes (~20% obstacle coverage)"
			case engine.Hard:
				return "Arena: 40×15 | Speed: 80ms  | 3 Mazes (~50% obstacle coverage)"
			default: // Medium
				return "Arena: 45×18 | Speed: 100ms | 3 Mazes (~35% obstacle coverage)"
			}
		default: // Classic
			switch d {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 3 Lives | Wrapping walls"
			case engine.Hard:
				return "Arena: 30×15 | Speed: 80ms  | 1 Life  | Hard walls"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 2 Lives | Hard walls"
			}
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
