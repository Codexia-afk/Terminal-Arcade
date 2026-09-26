package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// PreGameScreen presents the pre-launch confirmation screen matching the exact design spec.
type PreGameScreen struct {
	screen   *engine.Screen
	gameName string
	cfg      engine.GameConfig
	store    *history.Store
}

// NewPreGameScreen constructs a new pre-game briefing view.
func NewPreGameScreen(s *engine.Screen, gameName string, cfg engine.GameConfig, store *history.Store) *PreGameScreen {
	return &PreGameScreen{
		screen:   s,
		gameName: gameName,
		cfg:      cfg,
		store:    store,
	}
}

// Show renders the pre-game summary screen and waits for Enter (start) or Esc/Q (back).
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

	boxW := 46
	boxH := 16
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2
	if boxX < 0 {
		boxX = 0
	}
	if boxY < 0 {
		boxY = 0
	}

	// 1. Draw outer double-line box with header
	p.screen.DoubleBox(boxX, boxY, boxW, boxH, th.Accent, tcell.ColorBlack)
	p.screen.CenterText(boxY+1, "READY TO PLAY?", th.Accent, tcell.ColorBlack)

	// Horizontal divider line
	p.screen.DrawCell(boxX, boxY+2, '╠', th.Accent, tcell.ColorBlack)
	p.screen.DrawCell(boxX+boxW-1, boxY+2, '╣', th.Accent, tcell.ColorBlack)
	for x := boxX + 1; x < boxX+boxW-1; x++ {
		p.screen.DrawCell(x, boxY+2, '═', th.Accent, tcell.ColorBlack)
	}

	sy := boxY + 4

	// Game and Mode line
	modeTitle := p.cfg.Mode
	if modeTitle == "" {
		modeTitle = "classic"
	}
	gameTitle := strings.ToUpper(p.gameName)
	if p.gameName == "snake" {
		gameTitle = fmt.Sprintf("SNAKE %s", strings.ToUpper(modeTitle))
	} else if p.gameName == "ballplate" {
		gameTitle = "BALL & PLATE"
	} else if p.gameName == "pacman" {
		gameTitle = "PACMAN"
	}

	p.screen.DrawText(boxX+4, sy, fmt.Sprintf("Mode:       %-22s", gameTitle), th.Text, tcell.ColorBlack)
	sy += 2
	p.screen.DrawText(boxX+4, sy, fmt.Sprintf("Difficulty: %-22s", strings.ToUpper(p.cfg.Difficulty.String())), th.Text, tcell.ColorBlack)
	sy += 2
	p.screen.DrawText(boxX+4, sy, fmt.Sprintf("Theme:      %-22s", strings.ToUpper(th.Name)), th.Text, tcell.ColorBlack)
	sy += 2

	// Personal Best summary
	p.screen.DrawText(boxX+4, sy, "Your Personal Best (this mode):", th.HUD, tcell.ColorBlack)
	sy++

	pbText := "None yet — set the first record!"
	pbColor := tcell.ColorGray
	if p.store != nil {
		pbScore, found := history.PersonalBestForConfig(p.store.Records, p.gameName, p.cfg.Mode, string(p.cfg.Difficulty))
		if found && pbScore > 0 {
			timeAgo := "recently"
			if pb, ok := history.PersonalBestsForMode(p.store.Records, p.gameName, p.cfg.Mode); ok && !pb.HighScoreDate.IsZero() {
				timeAgo = history.FormatRelativeTime(pb.HighScoreDate, time.Now())
			}
			pbText = fmt.Sprintf("%d points • %s", pbScore, timeAgo)
			pbColor = tcell.ColorYellow
		}
	}
	p.screen.DrawText(boxX+4, sy, pbText, pbColor, tcell.ColorBlack)
	sy += 2

	// Prompt actions
	p.screen.CenterText(sy, "[Enter] Start  [Esc] Back", th.Accent, tcell.ColorBlack)
}
