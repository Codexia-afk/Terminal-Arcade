package modes

import (
	"fmt"
	"math"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// SurvivalEnemy represents an active pursuit enemy in Survival mode.
type SurvivalEnemy struct {
	Pos       engine.Position
	PrevPos   engine.Position
	MoveDelay int // ticks per move
	TickCount int
}

// SurvivalMode implements 5-wave survival with chasing enemy AI, wave clear bonuses, and patrol resets.
type SurvivalMode struct {
	BaseMode
	Wave          int
	WaveFoodEaten int
	Enemies       []SurvivalEnemy
	PatrolTicks   int // ticks remaining in post-wave patrol phase
	EnemySpeed    int // ticks per enemy move based on tier
}

// NewSurvival constructs an uninitialized SurvivalMode.
func NewSurvival() *SurvivalMode {
	return &SurvivalMode{}
}

// Name returns the game name identifier.
func (su *SurvivalMode) Name() string { return "snake" }

// Mode returns the mode identifier string.
func (su *SurvivalMode) Mode() string { return "survival" }

// Score returns the current player score.
func (su *SurvivalMode) Score() int { return su.ScoreVal }

// IsOver reports whether the session has ended.
func (su *SurvivalMode) IsOver() bool { return su.State == StateOver || su.State == StateWin }

// TickInterval returns the snake tick duration (100ms baseline across tiers).
func (su *SurvivalMode) TickInterval() time.Duration { return su.Interval }

// Metrics returns session metrics for persistence.
func (su *SurvivalMode) Metrics() map[string]int {
	return map[string]int{
		"food_eaten":         su.FoodEaten,
		"max_length_reached": su.MaxLength,
		"ticks_survived":     su.Ticks,
		"near_misses":        su.NearMisses,
		"wave_reached":       su.Wave,
		"lives_remaining":    su.Lives,
	}
}

// Init configures Survival mode according to tier parameters.
func (su *SurvivalMode) Init(cfg engine.GameConfig) {
	var (
		w, h       int
		lives      int
		enemyDelay int
	)

	switch cfg.Difficulty {
	case engine.Easy:
		w, h = 45, 20
		lives = 3
		// Enemy tick 180ms relative to 100ms snake tick (~1.8 ticks)
		enemyDelay = 2
	case engine.Hard:
		w, h = 35, 15
		lives = 1
		// Enemy tick 100ms (1 tick per snake move)
		enemyDelay = 1
	default: // Medium
		w, h = 40, 18
		lives = 2
		// Enemy tick 140ms (~1.4 ticks)
		enemyDelay = 2
	}

	su.InitBase(cfg, w, h, 100*time.Millisecond, false, lives)
	su.Wave = 1
	su.WaveFoodEaten = 0
	su.EnemySpeed = enemyDelay
	su.PatrolTicks = 0
	su.initWaveEnemies(1)
	su.SpawnFood(0)
}

func (su *SurvivalMode) initWaveEnemies(wave int) {
	su.Enemies = nil
	count := 1
	delay := su.EnemySpeed

	switch wave {
	case 1:
		count = 1
	case 2:
		count = 1
		if delay > 1 {
			delay = su.EnemySpeed
		}
	case 3:
		count = 2
	case 4:
		count = 2
		delay = 1
	case 5:
		count = 3
		delay = 1
	}

	corners := []engine.Position{
		{X: 1, Y: 1},
		{X: su.ArenaW - 2, Y: 1},
		{X: 1, Y: su.ArenaH - 2},
		{X: su.ArenaW - 2, Y: su.ArenaH - 2},
	}

	for i := 0; i < count && i < len(corners); i++ {
		su.Enemies = append(su.Enemies, SurvivalEnemy{
			Pos:       corners[i],
			PrevPos:   corners[i],
			MoveDelay: delay,
			TickCount: 0,
		})
	}
}

// HandleInput processes player directional input.
func (su *SurvivalMode) HandleInput(a engine.Action) {
	if su.State != StatePlaying {
		return
	}
	if d, ok := a.Direction(); ok {
		su.PushDirection(d)
	}
}

// Tick executes one simulation step for snake and enemy AI.
func (su *SurvivalMode) Tick() engine.TickResult {
	if su.State == StateOver {
		return engine.TickResult{Continue: false, Reason: su.Reason}
	}
	if su.State == StateWin {
		return engine.TickResult{Continue: false, Reason: "won"}
	}

	if su.EatFlashTicks > 0 {
		su.EatFlashTicks--
	}
	if su.WrapFlashTicks > 0 {
		su.WrapFlashTicks--
	}
	if su.HitFlashTicks > 0 {
		su.HitFlashTicks--
	}

	if su.PatrolTicks > 0 {
		su.PatrolTicks--
	}

	// 1. Process player buffered direction
	prevDir, turned := su.PopNextDirection()
	if turned {
		su.CheckNearMiss(prevDir)
	}

	// 2. Next head position
	next := su.Snake[0].Add(su.Dir)

	// Wall collision (hard walls in Survival)
	if next.X < 0 || next.X >= su.ArenaW || next.Y < 0 || next.Y >= su.ArenaH {
		return su.handleFatalCollision("wall")
	}

	// Self collision
	for _, seg := range su.Snake {
		if seg.Equal(next) {
			return su.handleFatalCollision("self")
		}
	}

	// 3. Move snake head
	su.Snake = append([]engine.Position{next}, su.Snake...)

	// 4. Food consumption & Wave progression
	ate := false
	for i, f := range su.Foods {
		if next.Equal(f) {
			su.FoodEaten++
			su.WaveFoodEaten++
			su.ScoreVal += 10
			su.EatFlashTicks = 3
			engine.Beep()
			if len(su.Snake) > su.MaxLength {
				su.MaxLength = len(su.Snake)
			}
			su.Foods = append(su.Foods[:i], su.Foods[i+1:]...)
			ate = true

			// Wave progression: 5 food per wave
			if su.WaveFoodEaten >= 5 {
				su.ScoreVal += 50 // Wave clear bonus
				su.Wave++
				su.WaveFoodEaten = 0
				if su.Wave > 5 {
					su.State = StateWin
					su.Reason = "won"
					return engine.TickResult{Continue: false, Reason: "won"}
				}
				su.initWaveEnemies(su.Wave)
				// 3 seconds of patrol on wave clear (at 100ms = 30 ticks)
				su.PatrolTicks = 30
			}

			su.SpawnFood(0)
			break
		}
	}

	if !ate {
		su.Snake = su.Snake[:len(su.Snake)-1]
	}

	// 5. Update Enemies & Check Collision
	for i := range su.Enemies {
		su.Enemies[i].TickCount++
		if su.Enemies[i].TickCount >= su.Enemies[i].MoveDelay {
			su.Enemies[i].TickCount = 0
			su.Enemies[i].PrevPos = su.Enemies[i].Pos

			if su.PatrolTicks > 0 {
				// Patrol phase: random wandering
				su.Enemies[i].Pos = su.stepPatrolAI(su.Enemies[i].Pos)
			} else {
				// Pursuit phase: greedy Manhattan chase towards snake head
				su.Enemies[i].Pos = su.StepEnemyAI(su.Enemies[i].Pos, su.Snake[0])
			}
		}

		// Collision check with snake segments
		for _, seg := range su.Snake {
			if su.Enemies[i].Pos.Equal(seg) {
				return su.handleFatalCollision("enemy")
			}
		}
	}

	su.Ticks++
	return engine.TickResult{Continue: true, Reason: ""}
}

func (su *SurvivalMode) handleFatalCollision(reason string) engine.TickResult {
	engine.Beep()
	if su.Lives > 1 {
		su.Lives--
		su.RespawnSnake()
		// Reset enemies to corners to prevent instant re-collision
		su.initWaveEnemies(su.Wave)
		su.PatrolTicks = 15 // 1.5s grace patrol
		return engine.TickResult{Continue: true, Reason: ""}
	}
	su.Lives = 0
	su.State = StateOver
	su.Reason = reason
	return engine.TickResult{Continue: false, Reason: reason}
}

// StepEnemyAI advances one enemy position towards target using greedy Manhattan distance.
func (su *SurvivalMode) StepEnemyAI(from, to engine.Position) engine.Position {
	dx := to.X - from.X
	dy := to.Y - from.Y

	dirs := []engine.Position{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}

	best := from
	minDist := math.Abs(float64(dx)) + math.Abs(float64(dy))

	for _, d := range dirs {
		next := from.Add(d)
		if next.X < 0 || next.X >= su.ArenaW || next.Y < 0 || next.Y >= su.ArenaH {
			continue
		}
		if su.Obstacles[next] {
			continue
		}
		dist := math.Abs(float64(to.X-next.X)) + math.Abs(float64(to.Y-next.Y))
		if dist < minDist {
			minDist = dist
			best = next
		}
	}

	return best
}

func (su *SurvivalMode) stepPatrolAI(from engine.Position) engine.Position {
	dirs := []engine.Position{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}

	perm := su.Rng.Perm(len(dirs))
	for _, idx := range perm {
		next := from.Add(dirs[idx])
		if next.X >= 0 && next.X < su.ArenaW && next.Y >= 0 && next.Y < su.ArenaH && !su.Obstacles[next] {
			return next
		}
	}
	return from
}

// Render draws the Survival gameplay screen.
func (su *SurvivalMode) Render(s *engine.Screen) {
	th := su.Cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	arenaW := su.ArenaW + 2
	arenaH := su.ArenaH + 2
	totalH := arenaH + 4

	originX := (w - arenaW) / 2
	originY := (h - totalH) / 2
	if originX < 0 {
		originX = 0
	}
	if originY < 0 {
		originY = 0
	}

	statusExtra := fmt.Sprintf("Wave: %d/5 (%d/5)", su.Wave, su.WaveFoodEaten)
	if su.PatrolTicks > 0 {
		statusExtra += fmt.Sprintf(" [PATROL: %.1fs]", float64(su.PatrolTicks)*0.1)
	}
	su.DrawCommonHUD(s, originX, originY, arenaW, "SURVIVAL", statusExtra)

	borderY := originY + 3
	su.DrawArenaBorder(s, originX, borderY, arenaW, arenaH)
	su.DrawSnakeElements(s, originX+1, borderY+1)

	// Render Chasing Enemies
	enemyGlyph := th.EnemyGlyph
	if enemyGlyph == 0 {
		enemyGlyph = '◆'
	}
	for _, e := range su.Enemies {
		eColor := th.Enemy
		if su.PatrolTicks > 0 {
			eColor = th.Accent
		}
		s.DrawCell(originX+1+e.Pos.X, borderY+1+e.Pos.Y, enemyGlyph, eColor, th.Background)
	}

	s.CenterText(borderY+arenaH+1, "Survive 5 Waves of Enemy Pursuers  |  P: Pause  |  Q: Quit", th.Text, th.Background)

	if su.State == StateOver || su.State == StateWin {
		su.DrawGameOverOverlay(s, w, h, "SURVIVAL", fmt.Sprintf("Waves Completed: %d/5", su.Wave-1))
	}
}
