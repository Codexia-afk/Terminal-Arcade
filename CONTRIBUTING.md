# Contributing to Go Arcade

Thank you for your interest in extending Go Arcade! This document describes coding conventions, testing guidelines, and the exact pattern for contributing a new game to the suite.

---

## Code Style & Guidelines

1. **Idiomatic Go**: All code must adhere to standard Go formatting. Run `gofmt -s -w .` before submitting.
2. **Static Analysis**: Code must pass `go vet ./...` with zero warnings.
3. **Documentation**:
   - Every package must have a package-level doc comment (`// Package ...`).
   - Every exported type, constant, variable, and function must have a clear doc comment.
4. **Zero Unapproved Dependencies**:
   - The suite only uses Go standard library and `github.com/gdamore/tcell/v2`.
   - Do not introduce additional dependencies without clear justification.
5. **Completely Offline**:
   - No network calls, telemetry, or external web service integrations are permitted.
6. **Strict Dependency Directions**:
   - `internal/engine` and `internal/history` depend on nothing else in the project.
   - `internal/games/<game>` may only depend on `internal/engine` and `internal/history`.
   - Games must never import other games or `internal/menu`.
   - `internal/history` must never import `tcell` or any rendering logic.

---

## How to Add a New Game

Adding a new game to the suite is modular and straightforward. Follow these steps:

### 1. Create the Game Package
Create a new directory: `internal/games/<mygame>/`.

### 2. Implement `engine.DifficultyProfile`
Define a concrete configuration profile specifying the parameters for each difficulty level:

```go
package mygame

import (
	"time"
	"goarcade/internal/engine"
)

type Profile struct {
	Level    engine.Difficulty
	Speed    time.Duration
	Obstacles int
}

func (p Profile) DifficultyLevel() engine.Difficulty {
	return p.Level
}

func Profiles(d engine.Difficulty) Profile {
	switch d {
	case engine.Easy:
		return Profile{Level: engine.Easy, Speed: 150 * time.Millisecond, Obstacles: 2}
	case engine.Hard:
		return Profile{Level: engine.Hard, Speed: 80 * time.Millisecond, Obstacles: 8}
	default:
		return Profile{Level: engine.Medium, Speed: 110 * time.Millisecond, Obstacles: 5}
	}
}
```

### 3. Implement the `engine.Game` Interface
Implement the interface contract in `internal/games/<mygame>/<mygame>.go`:

```go
type Game struct {
	cfg     engine.GameConfig
	profile Profile
	score   int
	isOver  bool
	reason  string
}

func New() *Game { return &Game{} }

func (g *Game) Name() string { return "mygame" }
func (g *Game) Score() int   { return g.score }
func (g *Game) IsOver() bool { return g.isOver }

// Optional: specify tick rate if implementing engine.TickIntervalProvider
func (g *Game) TickInterval() time.Duration { return g.profile.Speed }

func (g *Game) Init(cfg engine.GameConfig) {
	g.cfg = cfg
	g.profile = Profiles(cfg.Difficulty)
	g.score = 0
	g.isOver = false
}

func (g *Game) HandleInput(a engine.Action) {
	// Translate engine actions (ActionUp, ActionDown, ActionConfirm, etc.)
}

func (g *Game) Tick() engine.TickResult {
	// Advance simulation step
	if g.isOver {
		return engine.TickResult{Continue: false, Reason: g.reason}
	}
	return engine.TickResult{Continue: true}
}

func (g *Game) Render(s *engine.Screen) {
	// Draw UI, borders, HUD, and game state using engine.Screen primitives
	s.Clear(g.cfg.Theme.Background)
	s.Box(0, 0, 40, 20, g.cfg.Theme.Wall, g.cfg.Theme.Background)
}
```

### 4. Register in Menu & Dispatch
1. In `internal/menu/mainmenu.go`:
   - Add a new `Choice` enum variant (e.g. `ChoiceMyGame`).
   - Add the item to `items` list in `NewMainMenu`.
   - Add number key shortcut mapping.
2. In `internal/menu/picker.go`:
   - Add game display name in `gameDisplayName()`.
   - Add difficulty descriptions in `difficultyDescription()`.
3. In `internal/menu/history.go`:
   - Add the game to the tabs list and high scores summary list.
4. In `cmd/arcade/main.go`:
   - Add the case in the game selection dispatch switch.

---

## Testing Expectations

All game simulation rules, collisions, scoring, state transitions, and history persistence must be 100% unit-testable headlessly without an interactive TTY:

- **No TTY requirement**: Unit tests must never open a real terminal screen or require user input.
- **Run all tests**:
  ```bash
  go test -v ./...
  ```
- **Coverage**: Ensure tests verify:
  1. Difficulty parameter application.
  2. Entity collisions and boundary conditions.
  3. Score increments and win/loss conditions.
  4. State transitions on key actions.
