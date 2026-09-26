package modes

import (
	"fmt"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// ZenMode implements relaxed, peaceful snake play with flow bonuses and continuous wrapping.
type ZenMode struct {
	BaseMode
	FlowDuration     time.Duration
	FlowBonusPoints  int
	PersonalLongest  int
	BestZenLength    int
}

// NewZen constructs an uninitialized ZenMode.
func NewZen() *ZenMode {
	return &ZenMode{}
}

// Name returns the game name identifier.
func (z *ZenMode) Name() string { return "snake" }

// Mode returns the mode identifier string.
func (z *ZenMode) Mode() string { return "zen" }

// Score returns the current player score.
func (z *ZenMode) Score() int { return z.ScoreVal }

// IsOver reports whether the session has ended.
func (z *ZenMode) IsOver() bool { return z.State == StateOver }

// TickInterval returns the speed duration for the active tier.
func (z *ZenMode) TickInterval() time.Duration { return z.Interval }

// Metrics returns session metrics for persistence.
func (z *ZenMode) Metrics() map[string]int {
	return map[string]int{
		"food_eaten":         z.FoodEaten,
		"max_length_reached": z.MaxLength,
		"ticks_survived":     z.Ticks,
		"flow_bonus":         z.FlowBonusPoints,
	}
}

// Init configures Zen mode with wide arenas and peaceful parameters.
func (z *ZenMode) Init(cfg engine.GameConfig) {
	var (
		w, h     int
		interval time.Duration
	)

	switch cfg.Difficulty {
	case engine.Easy:
		w, h = 60, 25
		interval = 150 * time.Millisecond
	case engine.Hard:
		w, h = 40, 15
		interval = 100 * time.Millisecond
	default: // Medium
		w, h = 50, 20
		interval = 120 * time.Millisecond
	}

	z.InitBase(cfg, w, h, interval, true, 0)
	z.FlowDuration = 0
	z.FlowBonusPoints = 0
	z.PersonalLongest = 3

	if z.HasBest && z.BestInfo.KeyMetricName == "Max Length" {
		z.BestZenLength = z.BestInfo.KeyMetricValue
	}

	z.SpawnFood(0)
}

// HandleInput handles directional input and resets flow if direction changes.
func (z *ZenMode) HandleInput(a engine.Action) {
	if z.State != StatePlaying {
		return
	}
	if d, ok := a.Direction(); ok {
		lastDir := z.Dir
		if len(z.DirQueue) > 0 {
			lastDir = z.DirQueue[len(z.DirQueue)-1]
		}
		// If input actively turns away from current heading, reset 10s unbroken flow
		if d != lastDir && !(d.X == -lastDir.X && d.Y == -lastDir.Y) {
			z.FlowDuration = 0
		}
		z.PushDirection(d)
	}
}

// Tick executes a discrete simulation step without death hazards.
func (z *ZenMode) Tick() engine.TickResult {
	if z.State == StateOver {
		return engine.TickResult{Continue: false, Reason: z.Reason}
	}

	if z.EatFlashTicks > 0 {
		z.EatFlashTicks--
	}
	if z.WrapFlashTicks > 0 {
		z.WrapFlashTicks--
	}

	// 1. Process unbroken flow bonus (+5 points every 10s without direction change)
	z.FlowDuration += z.Interval
	if z.FlowDuration >= 10*time.Second {
		z.FlowDuration -= 10 * time.Second
		z.ScoreVal += 5
		z.FlowBonusPoints += 5
		z.EatFlashTicks = 3
		engine.Beep()
	}

	// 2. Process buffered direction
	_, turned := z.PopNextDirection()
	if turned {
		z.FlowDuration = 0
	}

	// 3. Next head position with endless safe wrapping
	next := z.Snake[0].Add(z.Dir)
	wrapped := false
	if next.X < 0 || next.X >= z.ArenaW || next.Y < 0 || next.Y >= z.ArenaH {
		wrapped = true
	}
	next.X = (next.X + z.ArenaW) % z.ArenaW
	next.Y = (next.Y + z.ArenaH) % z.ArenaH
	if wrapped {
		z.WrapFlashTicks = 3
	}

	// 4. Move head (no self collision death in Zen)
	z.Snake = append([]engine.Position{next}, z.Snake...)

	// 5. Food consumption (1 point per food in Zen)
	ate := false
	for i, f := range z.Foods {
		if next.Equal(f) {
			z.FoodEaten++
			z.ScoreVal += 1
			z.EatFlashTicks = 3
			engine.Beep()
			if len(z.Snake) > z.MaxLength {
				z.MaxLength = len(z.Snake)
			}
			if len(z.Snake) > z.PersonalLongest {
				z.PersonalLongest = len(z.Snake)
			}
			z.Foods = append(z.Foods[:i], z.Foods[i+1:]...)
			ate = true
			z.SpawnFood(0)
			break
		}
	}

	if !ate {
		z.Snake = z.Snake[:len(z.Snake)-1]
	}

	z.Ticks++
	return engine.TickResult{Continue: true, Reason: ""}
}

// Render renders the Zen mode interface.
func (z *ZenMode) Render(s *engine.Screen) {
	th := z.Cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	arenaW := z.ArenaW + 2
	arenaH := z.ArenaH + 2
	totalH := arenaH + 4

	originX := (w - arenaW) / 2
	originY := (h - totalH) / 2
	if originX < 0 {
		originX = 0
	}
	if originY < 0 {
		originY = 0
	}

	extraHUD := fmt.Sprintf("Flow: +%d pts | Longest: %d", z.FlowBonusPoints, z.PersonalLongest)
	z.DrawCommonHUD(s, originX, originY, arenaW, "ZEN", extraHUD)

	borderY := originY + 3
	z.DrawArenaBorder(s, originX, borderY, arenaW, arenaH)
	z.DrawSnakeElements(s, originX+1, borderY+1)

	// Meditative footer
	s.CenterText(borderY+arenaH+1, "Zen: No Deaths • Walls Wrap • Maintain Flow for Bonus Points  |  Q: Exit", th.Text, th.Background)

	if z.State == StateOver {
		longestMsg := fmt.Sprintf("Personal Longest: %d segments", z.PersonalLongest)
		if z.BestZenLength > 0 {
			longestMsg += fmt.Sprintf(" (Best: %d)", z.BestZenLength)
		}
		z.DrawGameOverOverlay(s, w, h, "ZEN", longestMsg)
	}
}
