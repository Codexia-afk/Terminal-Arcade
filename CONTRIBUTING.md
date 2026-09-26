# Contributing to Go Arcade

Thank you for your interest in extending Go Arcade! This document describes coding conventions, testing guidelines, and the exact patterns for adding a new Snake gameplay mode or an entirely new arcade game to the suite.

---

## Code Style & Architectural Guidelines

1. **Idiomatic Go**: All code must adhere to standard Go formatting. Run `gofmt -s -w .` before submitting.
2. **Static Analysis**: Code must pass `go vet ./...` with zero warnings.
3. **Documentation**:
   - Every package must have a package-level doc comment (`// Package ...`).
   - Every exported type, constant, variable, and function must have a clear doc comment.
4. **Zero Unapproved Dependencies**:
   - The suite uses only the Go standard library and `github.com/gdamore/tcell/v2`.
   - Do not introduce additional dependencies without clear justification.
5. **Completely Offline**:
   - Zero network calls, telemetry, update checks, or external web service integrations are permitted.
6. **Strict Dependency Directions**:
   - `internal/engine` and `internal/history` depend on nothing else in the project.
   - `internal/games/<game>` may only depend on `internal/engine` and `internal/history`.
   - Games must never import other games or `internal/menu`.
   - `internal/history` must never import `tcell` or any rendering logic.
7. **Deterministic Time / Clock Injection**:
   - Streak and achievement calculations must accept reference `now time.Time` parameters rather than calling `time.Now()` directly, enabling deterministic unit testing across timezone and day boundaries.

---

## How to Add a New Snake Gameplay Mode

Snake modes in Go Arcade are modular structs implementing `engine.Game` (and optionally `engine.TickIntervalProvider`), inheriting shared behavior from `BaseMode`.

Follow these steps to add a new mode (e.g. `PortalMode`):

### 1. Create the Mode File
Create `internal/games/snake/modes/portal.go`:

```go
package modes

import (
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// PortalMode represents a custom Snake gameplay mode.
type PortalMode struct {
	BaseMode
	// Custom mode fields (e.g., portals, timers, hazards)
	Portals [2]engine.Position
}

// NewPortalMode creates a new uninitialized PortalMode.
func NewPortalMode() *PortalMode {
	return &PortalMode{}
}

func (m *PortalMode) Name() string { return "portal" }
func (m *PortalMode) Score() int   { return m.score }
func (m *PortalMode) IsOver() bool { return m.isOver }

func (m *PortalMode) TickInterval() time.Duration {
	// Pacing based on difficulty
	switch m.cfg.Difficulty {
	case engine.Hard:
		return 80 * time.Millisecond
	case engine.Medium:
		return 100 * time.Millisecond
	default:
		return 120 * time.Millisecond
	}
}

func (m *PortalMode) Init(cfg engine.GameConfig) {
	// 1. Initialize arena dimensions and base state
	arenaW, arenaH := 36, 18
	m.InitBase(cfg, arenaW, arenaH)
	m.lives = 2

	// 2. Set mode-specific state
	m.Portals[0] = engine.Position{X: 5, Y: 5}
	m.Portals[1] = engine.Position{X: arenaW - 6, Y: arenaH - 6}

	// 3. Spawn initial food using BaseMode helper
	m.SpawnFood()
}

func (m *PortalMode) HandleInput(a engine.Action) {
	m.BaseMode.HandleInput(a)
}

func (m *PortalMode) Tick() engine.TickResult {
	if m.isOver {
		return engine.TickResult{Continue: false, Reason: m.reason}
	}

	// Advance direction from 2-step queue
	m.AdvanceDirection()

	// Calculate target position
	head := m.body[0]
	next := engine.Position{X: head.X + m.dir.X, Y: head.Y + m.dir.Y}

	// Mode-specific teleport mechanic:
	if next == m.Portals[0] {
		next = m.Portals[1]
	} else if next == m.Portals[1] {
		next = m.Portals[0]
	}

	// Check boundary collisions
	if next.X <= 0 || next.X >= m.arenaW-1 || next.Y <= 0 || next.Y >= m.arenaH-1 {
		m.lives--
		if m.lives <= 0 {
			m.isOver = true
			m.reason = "Collided with wall"
			return engine.TickResult{Continue: false, Reason: m.reason}
		}
		// Reset snake to center on remaining lives
		m.ResetSnake()
		return engine.TickResult{Continue: true}
	}

	// Check self-collision
	for _, segment := range m.body[:len(m.body)-1] {
		if next == segment {
			m.isOver = true
			m.reason = "Collided with self"
			return engine.TickResult{Continue: false, Reason: m.reason}
		}
	}

	// Move snake body
	m.body = append([]engine.Position{next}, m.body...)

	// Food consumption
	if next == m.food {
		m.score += 10
		m.foodEaten++
		m.SpawnFood()
	} else {
		m.body = m.body[:len(m.body)-1]
	}

	return engine.TickResult{Continue: true}
}

func (m *PortalMode) Render(s *engine.Screen) {
	// 1. Draw arena box and base layout
	m.RenderArena(s)

	// 2. Render portals
	s.Set(m.originX+m.Portals[0].X, m.originY+m.Portals[0].Y, 'O', m.cfg.Theme.Food, m.cfg.Theme.Background)
	s.Set(m.originX+m.Portals[1].X, m.originY+m.Portals[1].Y, 'O', m.cfg.Theme.Food, m.cfg.Theme.Background)

	// 3. Render snake and food using BaseMode helper
	m.RenderSnakeAndFood(s)

	// 4. Render centered HUD strip
	extraStatus := "Portals: Active"
	m.RenderHUD(s, "Portal", extraStatus)
}
```

### 2. Register the Mode Constant & Metadata
In `internal/games/snake/modes.go`:
```go
const (
	// ... existing modes ...
	ModePortal SnakeMode = "portal"
)

// Add to ModeList:
var ModeList = []SnakeMode{
	ModeClassic,
	ModeZen,
	ModeSurvival,
	ModeTimeAttack,
	ModeObstacle,
	ModePortal, // New mode
}

// Add title, description, and difficulty table in modes.go
```

### 3. Register in Mode Factory & Game Orchestrator
In `internal/games/snake/modes/factory.go` (or `internal/games/snake/game.go` `NewMode`):
```go
case ModePortal:
	g.activeMode = modes.NewPortalMode()
```

### 4. Register in Menu Mode Picker
In `internal/menu/picker.go`:
Add the difficulty configuration details in `snakeModeDetails()` to display descriptions in the UI carousel.

### 5. Write Headless Unit Tests
Add unit tests in `internal/games/snake/modes/portal_test.go`:
- Verify initialization under Easy, Medium, and Hard.
- Verify custom mechanics (teleporting, scoring, lives decrement).
- Ensure no TTY is required.

---

## How to Add an Entirely New Game

Adding a completely new arcade game to the suite (e.g. `Tetris` or `Space Invaders`):

### 1. Create the Game Package
Create `internal/games/<mygame>/<mygame>.go`.

### 2. Implement the `engine.Game` Interface
```go
type Game interface {
	Name() string
	Init(cfg GameConfig)
	HandleInput(action Action)
	Tick() TickResult
	Render(screen *Screen)
	Score() int
	IsOver() bool
}
```

Optionally implement:
- `engine.TickIntervalProvider` (`TickInterval() time.Duration`) for dynamic pacing.
- `engine.MetricsProvider` (`Metrics() map[string]int`) for deep telemetry tracking.

### 3. Register in Menu & CLI Dispatch
1. In `internal/menu/mainmenu.go` or `internal/menu/picker.go`: add the game to the game selection carousel.
2. In `internal/menu/history.go`: add high-score display entries.
3. In `cmd/arcade/main.go`: instantiate the game in the game execution loop.

---

## Testing & Quality Expectations

All game simulation rules, collisions, scoring, state transitions, physics, ghost AI, and history persistence must be 100% unit-testable headlessly without an interactive TTY or `sudo` permissions:

- **No TTY Requirement**: Unit tests must never open a real terminal screen or block for keyboard input.
- **Run all unit tests**:
  ```bash
  make test
  # or
  go test -buildvcs=false -v ./...
  ```
- **Run static analysis**:
  ```bash
  make vet
  ```
- **Cross-compile validation**:
  ```bash
  make cross-compile
  ```
