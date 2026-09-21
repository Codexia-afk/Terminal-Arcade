// Command arcade launches the Go Terminal Game Suite with Snake, Pacman, and Ball & Plate.
package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
	"goarcade/internal/games/ballplate"
	"goarcade/internal/games/pacman"
	"goarcade/internal/games/snake"
	"goarcade/internal/history"
	"goarcade/internal/menu"
)

func main() {
	// Initialize history storage
	store, err := history.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not initialize history store: %v\n", err)
		store = &history.Store{Records: make([]history.Record, 0)}
	}

	// Initialize terminal raw mode via tcell
	rawScreen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to create tcell screen: %v\n", err)
		os.Exit(1)
	}

	if err := rawScreen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}

	screen := engine.NewScreen(rawScreen)

	// Ensure terminal raw mode is ALWAYS restored on exit, error, or unexpected panic
	defer func() {
		if screen != nil && screen.Raw != nil {
			screen.Raw.Fini()
		}
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "\n[FATAL] Recovered from panic in Go Arcade: %v\n", r)
			debug.PrintStack()
			os.Exit(2)
		}
	}()

	mainMenu := menu.NewMainMenu(screen)

	for {
		choice := mainMenu.Show()
		if choice == menu.ChoiceQuit {
			break
		}

		if choice == menu.ChoiceHistory {
			historyScreen := menu.NewHistoryScreen(screen, store)
			historyScreen.Show()
			continue
		}

		var game engine.Game
		switch choice {
		case menu.ChoiceSnake:
			game = snake.New()
		case menu.ChoicePacman:
			game = pacman.New()
		case menu.ChoiceBallPlate:
			game = ballplate.New()
		default:
			continue
		}

		// Prompt user for Difficulty and Theme configuration
		picker := menu.NewConfigPicker(screen, game.Name())
		cfg, confirmed := picker.Pick()
		if !confirmed {
			continue // User canceled out of picker; return to main menu
		}

		// Initialize selected game with chosen configuration
		game.Init(cfg)

		startTime := time.Now()
		loop := engine.NewLoop(screen, 0)
		result := loop.Run(game)
		duration := time.Since(startTime)

		outcome := result.Reason
		if outcome == "" {
			outcome = "quit"
		}

		// Save completed or quit session to history
		record := history.Record{
			Game:       game.Name(),
			Difficulty: string(cfg.Difficulty),
			Theme:      cfg.Theme.Name,
			Score:      game.Score(),
			Outcome:    outcome,
			PlayedAt:   startTime,
			Duration:   duration,
		}
		_ = store.Append(record)
	}
}
