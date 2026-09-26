// Package snake implements the 5-mode Nokia-style grid Snake game.
package snake

import (
	"math/rand"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/games/snake/modes"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

type gameState int

const (
	stateTitle gameState = iota
	statePlaying
	stateOver
)

// Profile contains concrete per-difficulty mechanical parameters.
type Profile struct {
	Level     engine.Difficulty
	Interval  time.Duration
	Width     int
	Height    int
	Wrap      bool
	Obstacles bool
	Lives     int
}

// DifficultyLevel implements engine.DifficultyProfile.
func (p Profile) DifficultyLevel() engine.Difficulty {
	return p.Level
}

// Profiles returns tuning parameters for the requested difficulty (Classic defaults).
func Profiles(d engine.Difficulty) Profile {
	switch d {
	case engine.Easy:
		return Profile{
			Level:     engine.Easy,
			Interval:  120 * time.Millisecond,
			Width:     50,
			Height:    20,
			Wrap:      true,
			Obstacles: false,
			Lives:     3,
		}
	case engine.Hard:
		return Profile{
			Level:     engine.Hard,
			Interval:  80 * time.Millisecond,
			Width:     30,
			Height:    15,
			Wrap:      false,
			Obstacles: true,
			Lives:     1,
		}
	default: // Medium
		return Profile{
			Level:     engine.Medium,
			Interval:  100 * time.Millisecond,
			Width:     40,
			Height:    18,
			Wrap:      false,
			Obstacles: false,
			Lives:     2,
		}
	}
}

// ProfilesForMode returns exact tuning parameters for a mode and difficulty tier.
func ProfilesForMode(m Mode, d engine.Difficulty) Profile {
	switch m {
	case ModeZen:
		switch d {
		case engine.Easy:
			return Profile{Level: engine.Easy, Interval: 150 * time.Millisecond, Width: 60, Height: 25, Wrap: true, Lives: 0}
		case engine.Hard:
			return Profile{Level: engine.Hard, Interval: 100 * time.Millisecond, Width: 40, Height: 15, Wrap: true, Lives: 0}
		default:
			return Profile{Level: engine.Medium, Interval: 120 * time.Millisecond, Width: 50, Height: 20, Wrap: true, Lives: 0}
		}
	case ModeSurvival:
		switch d {
		case engine.Easy:
			return Profile{Level: engine.Easy, Interval: 100 * time.Millisecond, Width: 45, Height: 20, Wrap: false, Lives: 3}
		case engine.Hard:
			return Profile{Level: engine.Hard, Interval: 100 * time.Millisecond, Width: 35, Height: 15, Wrap: false, Lives: 1}
		default:
			return Profile{Level: engine.Medium, Interval: 100 * time.Millisecond, Width: 40, Height: 18, Wrap: false, Lives: 2}
		}
	case ModeTimeAttack:
		switch d {
		case engine.Easy:
			return Profile{Level: engine.Easy, Interval: 120 * time.Millisecond, Width: 50, Height: 20, Wrap: true, Lives: 0}
		case engine.Hard:
			return Profile{Level: engine.Hard, Interval: 80 * time.Millisecond, Width: 35, Height: 15, Wrap: false, Lives: 0}
		default:
			return Profile{Level: engine.Medium, Interval: 100 * time.Millisecond, Width: 40, Height: 18, Wrap: false, Lives: 0}
		}
	case ModeObstacle:
		switch d {
		case engine.Easy:
			return Profile{Level: engine.Easy, Interval: 120 * time.Millisecond, Width: 50, Height: 20, Wrap: false, Lives: 0}
		case engine.Hard:
			return Profile{Level: engine.Hard, Interval: 80 * time.Millisecond, Width: 40, Height: 15, Wrap: false, Lives: 0}
		default:
			return Profile{Level: engine.Medium, Interval: 100 * time.Millisecond, Width: 45, Height: 18, Wrap: false, Lives: 0}
		}
	default: // Classic
		return Profiles(d)
	}
}

// Game implements engine.Game and delegates to the active mode struct.
type Game struct {
	cfg        engine.GameConfig
	profile    Profile
	mode       Mode
	state      gameState
	snake      []engine.Position
	dir        engine.Position
	dirQueue   []engine.Position
	food       engine.Position
	obstacles  map[engine.Position]bool
	enemies    []Enemy
	score      int
	bestScore  int
	bestInfo   history.PersonalBestInfo
	hasBest    bool
	reason     string
	ticks      int
	foodEaten  int
	maxLength  int
	nearMisses int
	rng        *rand.Rand

	// Visual feedback counters
	eatFlashTicks  int
	wrapFlashTicks int

	// Mode-specific state
	currentWave     int
	waveFoodEaten   int
	mazeLevel       int
	levelFoodEaten  int
	timeRemaining   time.Duration
	flowDuration    time.Duration
	flowBonusPoints int

	// Active mode struct implementing engine.Game
	activeMode engine.Game
}

// New creates an uninitialized Snake game.
func New() *Game {
	return &Game{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		obstacles: make(map[engine.Position]bool),
		mode:      ModeClassic,
		dirQueue:  make([]engine.Position, 0, 2),
		maxLength: 3,
	}
}

// NewMode instantiates a specific mode implementing engine.Game.
func NewMode(m Mode) engine.Game {
	switch m {
	case ModeZen:
		return modes.NewZen()
	case ModeSurvival:
		return modes.NewSurvival()
	case ModeTimeAttack:
		return modes.NewTimeAttack()
	case ModeObstacle:
		return modes.NewObstacle()
	default:
		return modes.NewClassic()
	}
}

// Name returns the canonical game identifier.
func (g *Game) Name() string { return "snake" }

// Mode returns the active gameplay mode string.
func (g *Game) Mode() string { return string(g.mode) }

// Score returns the current player score.
func (g *Game) Score() int {
	if g.activeMode != nil {
		return g.activeMode.Score()
	}
	return g.score
}

// IsOver reports whether the game session has concluded.
func (g *Game) IsOver() bool {
	if g.state == stateOver {
		return true
	}
	if g.activeMode != nil {
		return g.activeMode.IsOver()
	}
	return false
}

// TickInterval returns ticker duration specified by active profile.
func (g *Game) TickInterval() time.Duration {
	if tip, ok := g.activeMode.(engine.TickIntervalProvider); ok {
		return tip.TickInterval()
	}
	return g.profile.Interval
}

// SetBestScore overrides high score for testing.
func (g *Game) SetBestScore(best int) { g.bestScore = best }

// Metrics returns session metrics satisfying engine.MetricsProvider.
func (g *Game) Metrics() map[string]int {
	if mp, ok := g.activeMode.(engine.MetricsProvider); ok {
		m := mp.Metrics()
		// Merge any custom metrics
		if g.nearMisses > 0 && m["near_misses"] == 0 {
			m["near_misses"] = g.nearMisses
		}
		if g.foodEaten > 0 && m["food_eaten"] == 0 {
			m["food_eaten"] = g.foodEaten
		}
		if g.maxLength > 0 && m["max_length_reached"] == 0 {
			m["max_length_reached"] = g.maxLength
		}
		if g.ticks > 0 && m["ticks_survived"] == 0 {
			m["ticks_survived"] = g.ticks
		}
		return m
	}
	return map[string]int{
		"food_eaten":         g.foodEaten,
		"max_length_reached": g.maxLength,
		"ticks_survived":     g.ticks,
		"near_misses":        g.nearMisses,
		"survival_wave":      g.currentWave,
		"maze_level":         g.mazeLevel,
		"time_remaining_s":   int(g.timeRemaining.Seconds()),
		"flow_bonus":         g.flowBonusPoints,
	}
}

// Init initializes the game and instantiates the underlying mode struct.
func (g *Game) Init(cfg engine.GameConfig) {
	g.cfg = cfg
	g.mode = ModeFromString(cfg.Mode)
	g.profile = ProfilesForMode(g.mode, cfg.Difficulty)
	g.state = statePlaying
	g.score = 0
	g.reason = ""
	g.ticks = 0
	g.foodEaten = 0
	g.maxLength = 3
	g.nearMisses = 0
	g.obstacles = make(map[engine.Position]bool)
	g.enemies = nil
	g.dirQueue = make([]engine.Position, 0, 2)
	g.currentWave = 1
	g.mazeLevel = 1
	g.timeRemaining = TimeAttackDuration

	cx := g.profile.Width / 2
	cy := g.profile.Height / 2
	g.snake = []engine.Position{
		{X: cx, Y: cy},
		{X: cx - 1, Y: cy},
		{X: cx - 2, Y: cy},
	}
	g.dir = engine.Position{X: 1, Y: 0}

	// Instantiate specialized mode struct
	g.activeMode = NewMode(g.mode)
	g.activeMode.Init(cfg)

	// Fetch personal best
	if store, err := history.Open(""); err == nil && store != nil {
		g.bestScore = history.HighScore(store.Records, "snake", string(cfg.Difficulty))
		if pb, ok := history.PersonalBests(store.Records, "snake"); ok {
			g.bestInfo = pb
			g.hasBest = true
		}
	}

	g.syncFromActiveMode()
}

func (g *Game) syncToActiveMode() {
	if g.activeMode == nil {
		return
	}
	switch m := g.activeMode.(type) {
	case *modes.ClassicMode:
		if len(g.snake) > 0 {
			m.Snake = g.snake
		}
		m.Dir = g.dir
		if len(g.dirQueue) > 0 {
			m.DirQueue = g.dirQueue
		}
		if g.food != (engine.Position{}) {
			m.Foods = []engine.Position{g.food}
		}
		m.Obstacles = g.obstacles
	case *modes.ZenMode:
		if len(g.snake) > 0 {
			m.Snake = g.snake
		}
		m.Dir = g.dir
		if len(g.dirQueue) > 0 {
			m.DirQueue = g.dirQueue
		}
		if g.food != (engine.Position{}) {
			m.Foods = []engine.Position{g.food}
		}
	case *modes.SurvivalMode:
		if len(g.snake) > 0 {
			m.Snake = g.snake
		}
		m.Dir = g.dir
		if len(g.dirQueue) > 0 {
			m.DirQueue = g.dirQueue
		}
		if g.food != (engine.Position{}) {
			m.Foods = []engine.Position{g.food}
		}
	case *modes.TimeAttackMode:
		if len(g.snake) > 0 {
			m.Snake = g.snake
		}
		m.Dir = g.dir
		if len(g.dirQueue) > 0 {
			m.DirQueue = g.dirQueue
		}
		if g.food != (engine.Position{}) {
			m.Foods = []engine.Position{g.food}
		}
		if g.timeRemaining > 0 {
			m.TimeRemaining = g.timeRemaining
		}
	case *modes.ObstacleMode:
		if len(g.snake) > 0 {
			m.Snake = g.snake
		}
		m.Dir = g.dir
		if len(g.dirQueue) > 0 {
			m.DirQueue = g.dirQueue
		}
		if g.food != (engine.Position{}) {
			m.Foods = []engine.Position{g.food}
			m.PreviewTicks = 0
		}
		if len(g.obstacles) > 0 {
			m.Obstacles = g.obstacles
		}
	}
}

func (g *Game) syncFromActiveMode() {
	if g.activeMode == nil {
		return
	}
	switch m := g.activeMode.(type) {
	case *modes.ClassicMode:
		g.snake = m.Snake
		g.dir = m.Dir
		g.dirQueue = m.DirQueue
		if len(m.Foods) > 0 {
			g.food = m.Foods[0]
		}
		g.obstacles = m.Obstacles
		g.score = m.ScoreVal
		g.ticks = m.Ticks
		g.foodEaten = m.FoodEaten
		g.maxLength = m.MaxLength
		g.nearMisses = m.NearMisses
		g.reason = m.Reason
		if m.State == modes.StateOver {
			g.state = stateOver
		}
	case *modes.ZenMode:
		g.snake = m.Snake
		g.dir = m.Dir
		g.dirQueue = m.DirQueue
		if len(m.Foods) > 0 {
			g.food = m.Foods[0]
		}
		g.score = m.ScoreVal
		g.ticks = m.Ticks
		g.foodEaten = m.FoodEaten
		g.maxLength = m.MaxLength
		g.flowBonusPoints = m.FlowBonusPoints
		g.reason = m.Reason
		if m.State == modes.StateOver {
			g.state = stateOver
		}
	case *modes.SurvivalMode:
		g.snake = m.Snake
		g.dir = m.Dir
		g.dirQueue = m.DirQueue
		if len(m.Foods) > 0 {
			g.food = m.Foods[0]
		}
		g.score = m.ScoreVal
		g.ticks = m.Ticks
		g.foodEaten = m.FoodEaten
		g.maxLength = m.MaxLength
		g.currentWave = m.Wave
		g.waveFoodEaten = m.WaveFoodEaten
		g.enemies = m.Enemies
		g.nearMisses = m.NearMisses
		g.reason = m.Reason
		if m.State == modes.StateOver || m.State == modes.StateWin {
			g.state = stateOver
		}
	case *modes.TimeAttackMode:
		g.snake = m.Snake
		g.dir = m.Dir
		g.dirQueue = m.DirQueue
		if len(m.Foods) > 0 {
			g.food = m.Foods[0]
		}
		g.score = m.ScoreVal
		g.ticks = m.Ticks
		g.foodEaten = m.FoodEaten
		g.maxLength = m.MaxLength
		g.timeRemaining = m.TimeRemaining
		g.nearMisses = m.NearMisses
		g.reason = m.Reason
		if m.State == modes.StateOver {
			g.state = stateOver
		}
	case *modes.ObstacleMode:
		g.snake = m.Snake
		g.dir = m.Dir
		g.dirQueue = m.DirQueue
		if len(m.Foods) > 0 {
			g.food = m.Foods[0]
		}
		g.obstacles = m.Obstacles
		g.score = m.ScoreVal
		g.ticks = m.Ticks
		g.foodEaten = m.FoodEaten
		g.maxLength = m.MaxLength
		g.mazeLevel = m.MazeLevel
		g.levelFoodEaten = m.LevelFoodEaten
		g.nearMisses = m.NearMisses
		g.reason = m.Reason
		if m.State == modes.StateOver || m.State == modes.StateWin {
			g.state = stateOver
		}
	}
}

// HandleInput processes keyboard direction actions.
func (g *Game) HandleInput(a engine.Action) {
	if g.state == stateTitle {
		if a == engine.ActionConfirm || a == engine.ActionUp || a == engine.ActionDown || a == engine.ActionLeft || a == engine.ActionRight {
			g.state = statePlaying
		}
		return
	}

	if d, ok := a.Direction(); ok {
		lastDir := g.dir
		if len(g.dirQueue) > 0 {
			lastDir = g.dirQueue[len(g.dirQueue)-1]
		}
		if !(d.X == -lastDir.X && d.Y == -lastDir.Y) && !(d.X == lastDir.X && d.Y == lastDir.Y) {
			if len(g.dirQueue) < 2 {
				g.dirQueue = append(g.dirQueue, d)
			} else {
				g.dirQueue[1] = d
			}
		}
	}

	if g.activeMode != nil {
		g.syncToActiveMode()
		g.activeMode.HandleInput(a)
		g.syncFromActiveMode()
	}
}

// Tick executes one discrete simulation step.
func (g *Game) Tick() engine.TickResult {
	if g.state == stateTitle {
		return engine.TickResult{Continue: true, Reason: ""}
	}
	if g.state == stateOver {
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	if g.activeMode != nil {
		g.syncToActiveMode()
		res := g.activeMode.Tick()
		g.syncFromActiveMode()
		if !res.Continue {
			g.state = stateOver
			g.reason = res.Reason
		}
		return res
	}

	return engine.TickResult{Continue: true, Reason: ""}
}

// Render presents the HUD and playfield.
func (g *Game) Render(s *engine.Screen) {
	if g.activeMode != nil {
		g.activeMode.Render(s)
	}
}

func (g *Game) stepEnemyAI(from, to engine.Position) engine.Position {
	if sm, ok := g.activeMode.(*modes.SurvivalMode); ok {
		return sm.StepEnemyAI(from, to)
	}
	return from
}

func (g *Game) spawnFood() {
	if g.activeMode != nil {
		switch m := g.activeMode.(type) {
		case *modes.ClassicMode:
			m.SpawnFood(0)
		case *modes.ZenMode:
			m.SpawnFood(0)
		case *modes.SurvivalMode:
			m.SpawnFood(0)
		case *modes.TimeAttackMode:
			m.SpawnFood(0)
		case *modes.ObstacleMode:
			m.SpawnFood(0)
		}
		g.syncFromActiveMode()
	}
}

func (g *Game) terminate(reason string) engine.TickResult {
	g.state = stateOver
	g.reason = reason
	if g.activeMode != nil {
		switch m := g.activeMode.(type) {
		case *modes.ClassicMode:
			m.State = modes.StateOver
			m.Reason = reason
		case *modes.ZenMode:
			m.State = modes.StateOver
			m.Reason = reason
		case *modes.SurvivalMode:
			m.State = modes.StateOver
			m.Reason = reason
		case *modes.TimeAttackMode:
			m.State = modes.StateOver
			m.Reason = reason
		case *modes.ObstacleMode:
			m.State = modes.StateOver
			m.Reason = reason
		}
	}
	return engine.TickResult{Continue: false, Reason: reason}
}
