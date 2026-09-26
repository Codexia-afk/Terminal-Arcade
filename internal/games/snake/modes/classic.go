package modes

import (
	"fmt"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// ClassicMode implements Nokia-style Snake with lives and difficulty parameters.
type ClassicMode struct {
	BaseMode
}

// NewClassic constructs an uninitialized ClassicMode.
func NewClassic() *ClassicMode {
	return &ClassicMode{}
}

// Name returns the game name identifier.
func (c *ClassicMode) Name() string { return "snake" }

// Mode returns the mode identifier string.
func (c *ClassicMode) Mode() string { return "classic" }

// Score returns the current player score.
func (c *ClassicMode) Score() int { return c.ScoreVal }

// IsOver reports whether the session has ended.
func (c *ClassicMode) IsOver() bool { return c.State == StateOver }

// TickInterval returns the speed duration for the active tier.
func (c *ClassicMode) TickInterval() time.Duration { return c.Interval }

// Metrics returns the session metrics for history persistence.
func (c *ClassicMode) Metrics() map[string]int {
	return map[string]int{
		"food_eaten":         c.FoodEaten,
		"max_length_reached": c.MaxLength,
		"ticks_survived":     c.Ticks,
		"near_misses":        c.NearMisses,
		"lives_remaining":    c.Lives,
	}
}

// Init configures the mode according to difficulty parameters.
func (c *ClassicMode) Init(cfg engine.GameConfig) {
	var (
		w, h     int
		interval time.Duration
		wrap     bool
		lives    int
	)

	switch cfg.Difficulty {
	case engine.Easy:
		w, h = 50, 20
		interval = 120 * time.Millisecond
		wrap = true
		lives = 3
	case engine.Hard:
		w, h = 30, 15
		interval = 80 * time.Millisecond
		wrap = false
		lives = 1
	default: // Medium
		w, h = 40, 18
		interval = 100 * time.Millisecond
		wrap = false
		lives = 2
	}

	c.InitBase(cfg, w, h, interval, wrap, lives)
	c.SpawnFood(0)
}

// HandleInput handles directional input.
func (c *ClassicMode) HandleInput(a engine.Action) {
	if c.State != StatePlaying {
		return
	}
	if d, ok := a.Direction(); ok {
		c.PushDirection(d)
	}
}

// Tick steps simulation forward.
func (c *ClassicMode) Tick() engine.TickResult {
	if c.State == StateOver {
		return engine.TickResult{Continue: false, Reason: c.Reason}
	}

	if c.EatFlashTicks > 0 {
		c.EatFlashTicks--
	}
	if c.WrapFlashTicks > 0 {
		c.WrapFlashTicks--
	}
	if c.HitFlashTicks > 0 {
		c.HitFlashTicks--
	}

	// 1. Process 2-step buffered direction
	prevDir, turned := c.PopNextDirection()
	if turned {
		c.CheckNearMiss(prevDir)
	}

	// 2. Next head position
	next := c.Snake[0].Add(c.Dir)

	// 3. Wall collision / Wrapping
	if c.Wrap {
		wrapped := false
		if next.X < 0 || next.X >= c.ArenaW || next.Y < 0 || next.Y >= c.ArenaH {
			wrapped = true
		}
		next.X = (next.X + c.ArenaW) % c.ArenaW
		next.Y = (next.Y + c.ArenaH) % c.ArenaH
		if wrapped {
			c.WrapFlashTicks = 3
		}
	} else if next.X < 0 || next.X >= c.ArenaW || next.Y < 0 || next.Y >= c.ArenaH {
		return c.handleFatalCollision("wall")
	}

	// 4. Self collision
	for _, seg := range c.Snake {
		if seg.Equal(next) {
			return c.handleFatalCollision("self")
		}
	}

	// 5. Move head
	c.Snake = append([]engine.Position{next}, c.Snake...)

	// 6. Food consumption
	ate := false
	for i, f := range c.Foods {
		if next.Equal(f) {
			c.FoodEaten++
			c.ScoreVal += 10
			c.EatFlashTicks = 3
			engine.Beep()
			if len(c.Snake) > c.MaxLength {
				c.MaxLength = len(c.Snake)
			}
			c.Foods = append(c.Foods[:i], c.Foods[i+1:]...)
			ate = true
			c.SpawnFood(0)
			break
		}
	}

	if !ate {
		c.Snake = c.Snake[:len(c.Snake)-1]
	}

	c.Ticks++
	return engine.TickResult{Continue: true, Reason: ""}
}

func (c *ClassicMode) handleFatalCollision(reason string) engine.TickResult {
	engine.Beep()
	if c.Lives > 1 {
		c.Lives--
		c.RespawnSnake()
		return engine.TickResult{Continue: true, Reason: ""}
	}
	c.Lives = 0
	c.State = StateOver
	c.Reason = reason
	return engine.TickResult{Continue: false, Reason: reason}
}

// Render draws the HUD, playfield, and game elements.
func (c *ClassicMode) Render(s *engine.Screen) {
	th := c.Cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	arenaW := c.ArenaW + 2
	arenaH := c.ArenaH + 2
	totalH := arenaH + 4

	originX := (w - arenaW) / 2
	originY := (h - totalH) / 2
	if originX < 0 {
		originX = 0
	}
	if originY < 0 {
		originY = 0
	}

	extraHUD := fmt.Sprintf("Mode: Classic | Eaten: %d", c.FoodEaten)
	c.DrawCommonHUD(s, originX, originY, arenaW, "CLASSIC", extraHUD)

	borderY := originY + 3
	c.DrawArenaBorder(s, originX, borderY, arenaW, arenaH)
	c.DrawSnakeElements(s, originX+1, borderY+1)

	// Instruction line below arena
	s.CenterText(borderY+arenaH+1, "↑↓←→ or WASD Move  |  P Pause  |  Q Quit to Menu", th.Text, th.Background)

	if c.State == StateOver {
		c.DrawGameOverOverlay(s, w, h, "CLASSIC", "")
	}
}
