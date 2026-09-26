package menu

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// PreGameScreen presents a pre-launch confirmation summary detailing the selected
// game mode, difficulty parameters, visual theme, controls, and personal best.
type PreGameScreen struct {
	screen   *engine.Screen
	gameName string
	cfg      engine.GameConfig
	store    *history.Store
}

// NewPreGameScreen creates a new pre-game launch summary view.
func NewPreGameScreen(s *engine.Screen, gameName string, cfg engine.GameConfig, store *history.Store) *PreGameScreen {
	return &PreGameScreen{
		screen:   s,
		gameName: gameName,
		cfg:      cfg,
		store:    store,
	}
}

// Show renders the pre-game launch screen and awaits player confirmation.
// Returns true if player presses Enter/Space to launch, or false if cancelled.
func (p *PreGameScreen) Show() bool {
	for {
		p.render()
		p.screen.Flush()

		ev := p.screen.Raw.PollEvent()
		if ev == nil {
			return false
		}

		switch e := ev.(type) {
		case *tcell.EventResize:
			p.screen.Raw.Sync()
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyEnter:
				return true
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return false
			}

			switch e.Rune() {
			case ' ':
				return true
			case 'q', 'Q':
				return false
			}
		}
	}
}

func (p *PreGameScreen) render() {
	th := p.cfg.Theme
	p.screen.Clear(tcell.ColorBlack)
	w, h := p.screen.Size()

	banner := []string{
		`  ██████╗  ██████╗      █████╗ ██████╗  ██████╗ █████╗ ██████╗ ███████╗`,
		` ██╔════╝ ██╔═══██╗    ██╔══██╗██╔══██╗██╔════╝██╔══██╗██╔══██╗██╔════╝`,
		` ██║  ███╗██║   ██║    ███████║██████╔╝██║     ███████║██║  ██║█████╗  `,
		` ██║   ██║██║   ██║    ██╔══██║██╔══██╗██║     ██╔══██║██║  ██║██╔══╝  `,
		` ╚██████╔╝╚██████╔╝    ╚═╝  ╚═╝╚═╝  ╚═╝╚██████╗╚═╝  ╚═╝██████╔╝███████╗`,
	}

	startY := (h - 24) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Draw Title Banner
	for i, line := range banner {
		p.screen.CenterText(startY+i, line, th.Player, tcell.ColorBlack)
	}

	p.screen.CenterText(startY+6, "MISSION LAUNCH BRIEFING", th.Accent, tcell.ColorBlack)

	// 2. Main Summary Container Box
	boxW := 68
	boxH := 14
	boxX := (w - boxW) / 2
	boxY := startY + 8

	title := fmt.Sprintf("GAME: %s", strings.ToUpper(p.gameDisplayName()))
	p.screen.BoxWithTitle(boxX, boxY, boxW, boxH, title, th.Accent, tcell.ColorBlack, tcell.ColorWhite)


	sy := boxY + 1

	// Mode line (for Snake or others)
	modeDesc := p.getModeDescription()
	if p.gameName == "snake" {
		modeName := p.cfg.Mode
		if modeName == "" {
			modeName = "classic"
		}
		p.screen.DrawText(boxX+3, sy, fmt.Sprintf("Mode:       %s", strings.ToUpper(modeName)), th.Accent, tcell.ColorBlack)
		sy++
		p.screen.DrawText(boxX+5, sy, modeDesc, tcell.ColorWhite, tcell.ColorBlack)
		sy++
	}

	// Difficulty & Parameters
	diffText := fmt.Sprintf("Difficulty: %s", p.cfg.Difficulty.String())
	p.screen.DrawText(boxX+3, sy, diffText, th.Accent, tcell.ColorBlack)
	sy++
	paramText := p.getDifficultyParams()
	p.screen.DrawText(boxX+5, sy, paramText, tcell.ColorGray, tcell.ColorBlack)
	sy++

	// Theme
	p.screen.DrawText(boxX+3, sy, fmt.Sprintf("Theme:      %s", th.Name), th.Accent, tcell.ColorBlack)
	sy++

	// Personal Best
	pbText := "Personal Best: None yet — set the first record!"
	pbColor := tcell.ColorGray
	if p.store != nil {
		pbScore, found := history.PersonalBestForConfig(p.store.Records, p.gameName, p.cfg.Mode, string(p.cfg.Difficulty))
		if found && pbScore > 0 {
			pbText = fmt.Sprintf("★ Personal Best: %d points", pbScore)
			pbColor = tcell.ColorYellow
		}
	}
	p.screen.DrawText(boxX+3, sy, pbText, pbColor, tcell.ColorBlack)
	sy += 2

	// Controls Cheatsheet
	p.screen.DrawText(boxX+3, sy, "Controls:   ↑↓←→ / WASD: Steer/Move   |   P: Pause   |   Esc: Quit", th.HUD, tcell.ColorBlack)
	sy++

	// 3. Action Prompt
	launchPrompt := "[ SPACE / ENTER ] Start Game       [ ESC / Q ] Back to Menu"
	p.screen.CenterText(boxY+boxH+1, launchPrompt, tcell.ColorLime, tcell.ColorBlack)
}

func (p *PreGameScreen) gameDisplayName() string {
	switch p.gameName {
	case "snake":
		return "Snake Enhanced"
	case "pacman":
		return "Pacman Maze"
	case "ballplate":
		return "Ball & Plate Breakout"
	default:
		return p.gameName
	}
}

func (p *PreGameScreen) getModeDescription() string {
	switch p.cfg.Mode {
	case "zen":
		return "Peaceful flow growth without death — flow time bonuses every 10s."
	case "survival":
		return "Chased by enemy pursuers — clear 5 waves to claim victory!"
	case "time_attack":
		return "60-second sprint — eat fast to maximize time bonus!"
	case "obstacle":
		return "Navigate 3 escalating procedural labyrinth maze levels!"
	case "classic":
		fallthrough
	default:
		return "Classic survival — grow as long as possible without crashing."
	}
}

func (p *PreGameScreen) getDifficultyParams() string {
	switch p.gameName {
	case "snake":
		switch p.cfg.Mode {
		case "zen":
			switch p.cfg.Difficulty {
			case engine.Easy:
				return "Arena: 60×25 | Speed: 150ms | Borders: Endless Wrap"
			case engine.Hard:
				return "Arena: 40×15 | Speed: 100ms | Borders: Endless Wrap"
			default: // Medium
				return "Arena: 50×20 | Speed: 120ms | Borders: Endless Wrap"
			}
		case "survival":
			switch p.cfg.Difficulty {
			case engine.Easy:
				return "Arena: 45×20 | Speed: 100ms | 3 Lives | Enemies: 180ms"
			case engine.Hard:
				return "Arena: 35×15 | Speed: 100ms | 1 Life  | Enemies: 100ms"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 2 Lives | Enemies: 140ms"
			}
		case "time_attack":
			switch p.cfg.Difficulty {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 90s Timer | Borders: Wrap"
			case engine.Hard:
				return "Arena: 35×15 | Speed: 80ms  | 45s Timer | Borders: Hard"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 60s Timer | Borders: Hard"
			}
		case "obstacle":
			switch p.cfg.Difficulty {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 3 Mazes (~20% coverage)"
			case engine.Hard:
				return "Arena: 40×15 | Speed: 80ms  | 3 Mazes (~50% coverage)"
			default: // Medium
				return "Arena: 45×18 | Speed: 100ms | 3 Mazes (~35% coverage)"
			}
		default: // Classic
			switch p.cfg.Difficulty {
			case engine.Easy:
				return "Arena: 50×20 | Speed: 120ms | 3 Lives | Borders: Wrap"
			case engine.Hard:
				return "Arena: 30×15 | Speed: 80ms  | 1 Life  | Borders: Hard"
			default: // Medium
				return "Arena: 40×18 | Speed: 100ms | 2 Lives | Borders: Hard"
			}
		}
	case "pacman":
		switch p.cfg.Difficulty {
		case engine.Easy:
			return "Ghosts: 2 (Random Walk) | Power Pellet: 50 ticks | Lives: 3"
		case engine.Hard:
			return "Ghosts: 4 (Ambush + Chase) | Power Pellet: 20 ticks | Lives: 2"
		default:
			return "Ghosts: 3 (Manhattan Chase) | Power Pellet: 35 ticks | Lives: 3"
		}
	case "ballplate":
		switch p.cfg.Difficulty {
		case engine.Easy:
			return "Plate: 7 cells | Speed: 0.32x | Lives: 4 | Gapped Bricks"
		case engine.Hard:
			return "Plate: 3 cells | Speed: 0.48x | Lives: 2 | Tough Bricks"
		default:
			return "Plate: 5 cells | Speed: 0.40x | Lives: 3 | Standard Grid"
		}
	default:
		return ""
	}
}
