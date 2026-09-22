package engine

import "time"

// GameConfig captures selections injected into every game.
type GameConfig struct {
	Difficulty Difficulty
	Theme      Theme
}

// TickResult represents the outcome of one game step.
type TickResult struct {
	Continue bool   // false when the game has ended this tick
	Reason   string // "collision", "quit", "won", etc. — for game-over screen + history record
}

// Game is the universal interface implemented by every playable title in the suite.
type Game interface {
	Init(cfg GameConfig)
	Tick() TickResult
	HandleInput(a Action)
	Render(s *Screen)
	IsOver() bool
	Score() int
	Name() string
}

// TickIntervalProvider is an optional interface games can implement to specify their tick rate.
type TickIntervalProvider interface {
	TickInterval() time.Duration
}

// MetricsProvider is an optional interface games can implement to supply detailed per-game session metrics.
type MetricsProvider interface {
	Metrics() map[string]int
}
