// Command arcade launches the Go Terminal Game Suite with Snake, Pacman, and Ball & Plate.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/games/ballplate"
	"github.com/Codexia-afk/Terminal-Arcade/internal/games/pacman"
	"github.com/Codexia-afk/Terminal-Arcade/internal/games/snake"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
	"github.com/Codexia-afk/Terminal-Arcade/internal/menu"
)

func main() {
	// 1. CLI Flags definition
	exportFlag := flag.Bool("export", false, "Export all session records and achievements to a local file")
	formatFlag := flag.String("format", "json", "Export format: json or csv")
	outputFlag := flag.String("output", "", "Output filename for export (defaults to arcade_export.<format>)")
	resetFlag := flag.Bool("reset-data", false, "Reset all persistent session history and achievements")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Go Arcade — Production-Quality Terminal Arcade Suite\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  arcade [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  arcade                           # Launch interactive game suite\n")
		fmt.Fprintf(os.Stderr, "  arcade --export --format=json    # Export data to arcade_export.json\n")
		fmt.Fprintf(os.Stderr, "  arcade --export --format=csv     # Export sessions to arcade_export.csv\n")
		fmt.Fprintf(os.Stderr, "  arcade --reset-data              # Reset all saved history and badges\n")
	}

	flag.Parse()

	// 2. Handle CLI Reset
	if *resetFlag {
		fmt.Print("Warning: This will permanently delete all session history and achievements.\nAre you sure? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			if err := history.ResetAllData(""); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to reset data: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("All Go Arcade history and achievements have been reset.")
		} else {
			fmt.Println("Reset cancelled.")
		}
		os.Exit(0)
	}

	// 3. Initialize history & achievements storage
	store, err := history.Open("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not initialize history store: %v\n", err)
		store = &history.Store{Records: make([]history.Record, 0)}
	}

	achStore, err := history.OpenAchievements("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not initialize achievements store: %v\n", err)
		achStore = &history.AchievementsStore{Items: history.DefaultAchievements()}
	}

	// 4. Handle CLI Export
	if *exportFlag {
		outFile := *outputFlag
		format := strings.ToLower(strings.TrimSpace(*formatFlag))
		if outFile == "" {
			outFile = "arcade_export." + format
		}

		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to create output file %s: %v\n", outFile, err)
			os.Exit(1)
		}
		defer f.Close()

		switch format {
		case "csv":
			if err := history.ExportCSV(store.Records, f); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to export CSV: %v\n", err)
				os.Exit(1)
			}
		case "json":
			fallthrough
		default:
			if err := history.ExportJSON(store.Records, achStore.Items, f); err != nil {
				fmt.Fprintf(os.Stderr, "error: failed to export JSON: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Printf("Successfully exported %d session(s) and %d badge(s) to %s\n",
			len(store.Records), len(achStore.Items), outFile)
		os.Exit(0)
	}

	// Evaluate achievements on startup
	_ = achStore.Evaluate(store.Records, time.Now())

	// 5. Initialize terminal raw mode via tcell
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

	// Portable Signal handling: os.Interrupt ensures raw mode is cleaned on Ctrl-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		if screen != nil && screen.Raw != nil {
			screen.Raw.Fini()
		}
		os.Exit(0)
	}()

	mainMenu := menu.NewMainMenu(screen, store)

	for {
		choice := mainMenu.Show()
		if choice == menu.ChoiceQuit {
			break
		}

		if choice == menu.ChoiceHistory {
			historyScreen := menu.NewHistoryScreen(screen, store, achStore)
			historyScreen.Show()
			continue
		}

		if choice == menu.ChoiceExport {
			handleInteractiveExport(screen, store, achStore)
			continue
		}

		if choice == menu.ChoiceResetData {
			handleInteractiveReset(screen, store, achStore)
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
		picker := menu.NewConfigPicker(screen, game.Name(), store)
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

		// Extract metrics if game satisfies MetricsProvider
		var metrics map[string]int
		if mp, ok := game.(engine.MetricsProvider); ok {
			metrics = mp.Metrics()
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
			Metrics:    metrics,
		}
		_ = store.Append(record)

		// Evaluate achievements after session
		_ = achStore.Evaluate(store.Records, time.Now())
	}
}

func handleInteractiveExport(screen *engine.Screen, store *history.Store, achStore *history.AchievementsStore) {
	screen.Clear(tcell.ColorBlack)
	w, h := screen.Size()

	boxW := 58
	boxH := 9
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	screen.BoxWithTitle(boxX, boxY, boxW, boxH, "EXPORT LOCAL ARCHIVE", tcell.ColorAqua, tcell.ColorBlack, tcell.ColorWhite)
	screen.CenterText(boxY+2, "Select format to export local archive:", tcell.ColorWhite, tcell.ColorBlack)
	screen.CenterText(boxY+4, "[J] Export as JSON (arcade_export.json)", tcell.ColorYellow, tcell.ColorBlack)
	screen.CenterText(boxY+5, "[C] Export as CSV  (arcade_export.csv)", tcell.ColorAqua, tcell.ColorBlack)
	screen.CenterText(boxY+7, "Press J, C, or Esc/Q to Cancel", tcell.ColorGray, tcell.ColorBlack)
	screen.Flush()

	for {
		ev := screen.Raw.PollEvent()
		if ev == nil {
			return
		}
		if key, ok := ev.(*tcell.EventKey); ok {
			if key.Key() == tcell.KeyEscape || key.Key() == tcell.KeyCtrlC || key.Rune() == 'q' || key.Rune() == 'Q' {
				return
			}
			if key.Rune() == 'j' || key.Rune() == 'J' {
				f, err := os.Create("arcade_export.json")
				if err == nil {
					_ = history.ExportJSON(store.Records, achStore.Items, f)
					_ = f.Close()
					showExportConfirmation(screen, "arcade_export.json")
				}
				return
			}
			if key.Rune() == 'c' || key.Rune() == 'C' {
				f, err := os.Create("arcade_export.csv")
				if err == nil {
					_ = history.ExportCSV(store.Records, f)
					_ = f.Close()
					showExportConfirmation(screen, "arcade_export.csv")
				}
				return
			}
		}
	}
}

func showExportConfirmation(screen *engine.Screen, filename string) {
	screen.Clear(tcell.ColorBlack)
	w, h := screen.Size()
	boxW := 52
	boxH := 7
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	screen.BoxWithTitle(boxX, boxY, boxW, boxH, "EXPORT COMPLETE", tcell.ColorLime, tcell.ColorBlack, tcell.ColorLime)
	screen.CenterText(boxY+2, "Saved to "+filename, tcell.ColorWhite, tcell.ColorBlack)
	screen.CenterText(boxY+4, "Press any key to return to menu", tcell.ColorAqua, tcell.ColorBlack)
	screen.Flush()

	for {
		ev := screen.Raw.PollEvent()
		if ev != nil {
			return
		}
	}
}

func handleInteractiveReset(screen *engine.Screen, store *history.Store, achStore *history.AchievementsStore) {
	screen.Clear(tcell.ColorBlack)
	w, h := screen.Size()

	boxW := 56
	boxH := 9
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	screen.BoxWithTitle(boxX, boxY, boxW, boxH, "RESET SAVED DATA", tcell.ColorRed, tcell.ColorBlack, tcell.ColorRed)
	screen.CenterText(boxY+2, "Are you sure you want to erase all data?", tcell.ColorWhite, tcell.ColorBlack)
	screen.CenterText(boxY+3, "This will delete all past sessions and achievements.", tcell.ColorOrangeRed, tcell.ColorBlack)
	screen.CenterText(boxY+5, "Press Y to confirm  |  Press any other key to cancel", tcell.ColorYellow, tcell.ColorBlack)
	screen.Flush()

	for {
		ev := screen.Raw.PollEvent()
		if ev == nil {
			return
		}
		if key, ok := ev.(*tcell.EventKey); ok {
			if key.Rune() == 'y' || key.Rune() == 'Y' {
				_ = history.ResetAllData("")
				store.Records = make([]history.Record, 0)
				achStore.Items = history.DefaultAchievements()
				return
			}
			return
		}
	}
}
