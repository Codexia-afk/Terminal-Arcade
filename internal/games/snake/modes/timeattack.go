package modes

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// TimeAttackMode implements the high-speed sprint against the clock with visual timer & speed bonus.
type TimeAttackMode struct {
	BaseMode
	TimeRemaining time.Duration
	TotalTime     time.Duration
}

// NewTimeAttack constructs an uninitialized TimeAttackMode.
func NewTimeAttack() *TimeAttackMode {
	return &TimeAttackMode{}
}

// Name returns the game name identifier.
func (ta *TimeAttackMode) Name() string { return "snake" }

// Mode returns the mode identifier string.
func (ta *TimeAttackMode) Mode() string { return "time_attack" }

// Score returns the current player score.
func (ta *TimeAttackMode) Score() int { return ta.ScoreVal }

// IsOver reports whether the time attack session has ended.
func (ta *TimeAttackMode) IsOver() bool { return ta.State == StateOver }

// TickInterval returns the speed interval based on difficulty tier.
func (ta *TimeAttackMode) TickInterval() time.Duration { return ta.Interval }

// Metrics returns session metrics for persistence.
func (ta *TimeAttackMode) Metrics() map[string]int {
	return map[string]int{
		"food_eaten":         ta.FoodEaten,
		"max_length_reached": ta.MaxLength,
		"ticks_survived":     ta.Ticks,
		"near_misses":        ta.NearMisses,
		"time_remaining_s":   int(ta.TimeRemaining.Seconds()),
	}
}

// Init configures Time Attack mode with timer duration and arena parameters per tier.
func (ta *TimeAttackMode) Init(cfg engine.GameConfig) {
	var (
		w, h     int
		interval time.Duration
		timer    time.Duration
		wrap     bool
	)

	switch cfg.Difficulty {
	case engine.Easy:
		w, h = 50, 20
		interval = 120 * time.Millisecond
		timer = 90 * time.Second
		wrap = true
	case engine.Hard:
		w, h = 35, 15
		interval = 80 * time.Millisecond
		timer = 45 * time.Second
		wrap = false
	default: // Medium
		w, h = 40, 18
		interval = 100 * time.Millisecond
		timer = 60 * time.Second
		wrap = false
	}

	ta.InitBase(cfg, w, h, interval, wrap, 0)
	ta.TimeRemaining = timer
	ta.TotalTime = timer
	ta.SpawnFood(0)
}

// HandleInput handles directional input.
func (ta *TimeAttackMode) HandleInput(a engine.Action) {
	if ta.State != StatePlaying {
		return
	}
	if d, ok := a.Direction(); ok {
		ta.PushDirection(d)
	}
}

// Tick steps the timer countdown and snake movement.
func (ta *TimeAttackMode) Tick() engine.TickResult {
	if ta.State == StateOver {
		return engine.TickResult{Continue: false, Reason: ta.Reason}
	}

	// 1. Decrement Countdown Timer
	ta.TimeRemaining -= ta.Interval
	if ta.TimeRemaining <= 0 {
		ta.TimeRemaining = 0
		ta.State = StateOver
		ta.Reason = "time_up"
		engine.Beep()
		return engine.TickResult{Continue: false, Reason: "time_up"}
	}

	// Beep alert when crossing into last 10 seconds
	remSec := ta.TimeRemaining.Seconds()
	if remSec < 10.0 && remSec+ta.Interval.Seconds() >= 10.0 {
		engine.Beep()
	}

	if ta.EatFlashTicks > 0 {
		ta.EatFlashTicks--
	}
	if ta.WrapFlashTicks > 0 {
		ta.WrapFlashTicks--
	}

	// 2. Process buffered direction
	prevDir, turned := ta.PopNextDirection()
	if turned {
		ta.CheckNearMiss(prevDir)
	}

	// 3. Next head position
	next := ta.Snake[0].Add(ta.Dir)

	// Wall collision / Wrapping
	if ta.Wrap {
		wrapped := false
		if next.X < 0 || next.X >= ta.ArenaW || next.Y < 0 || next.Y >= ta.ArenaH {
			wrapped = true
		}
		next.X = (next.X + ta.ArenaW) % ta.ArenaW
		next.Y = (next.Y + ta.ArenaH) % ta.ArenaH
		if wrapped {
			ta.WrapFlashTicks = 3
		}
	} else if next.X < 0 || next.X >= ta.ArenaW || next.Y < 0 || next.Y >= ta.ArenaH {
		// In Time Attack on hard walls: bumping wall rebounds/respawns with short penalty
		ta.RespawnSnake()
		return engine.TickResult{Continue: true, Reason: ""}
	}

	// Self collision penalty in Time Attack: resets snake length to 3 without terminating
	for _, seg := range ta.Snake {
		if seg.Equal(next) {
			ta.RespawnSnake()
			return engine.TickResult{Continue: true, Reason: ""}
		}
	}

	// 4. Move snake head
	ta.Snake = append([]engine.Position{next}, ta.Snake...)

	// 5. Food consumption with time remaining bonus (Time remaining * 2)
	ate := false
	for i, f := range ta.Foods {
		if next.Equal(f) {
			ta.FoodEaten++
			bonus := int(ta.TimeRemaining.Seconds() * 2)
			ta.ScoreVal += 10 + bonus
			ta.EatFlashTicks = 3
			engine.Beep()

			if len(ta.Snake) > ta.MaxLength {
				ta.MaxLength = len(ta.Snake)
			}
			ta.Foods = append(ta.Foods[:i], ta.Foods[i+1:]...)
			ate = true

			// Last 10 seconds frenzy: spawn 2 active foods
			extraFoods := 0
			if ta.TimeRemaining.Seconds() < 10.0 {
				extraFoods = 1
			}
			ta.SpawnFood(extraFoods)
			break
		}
	}

	if !ate {
		ta.Snake = ta.Snake[:len(ta.Snake)-1]
	}

	ta.Ticks++
	return engine.TickResult{Continue: true, Reason: ""}
}

// Render draws the Time Attack gameplay screen with prominent visual timer.
func (ta *TimeAttackMode) Render(s *engine.Screen) {
	th := ta.Cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	arenaW := ta.ArenaW + 2
	arenaH := ta.ArenaH + 2
	totalH := arenaH + 4

	originX := (w - arenaW) / 2
	originY := (h - totalH) / 2
	if originX < 0 {
		originX = 0
	}
	if originY < 0 {
		originY = 0
	}

	// Timer color changes: green (>30s) -> yellow (10s-30s) -> red (<10s)
	remSec := ta.TimeRemaining.Seconds()
	timerColor := tcell.ColorLime
	if remSec < 10.0 {
		timerColor = tcell.ColorRed
	} else if remSec <= 30.0 {
		timerColor = tcell.ColorYellow
	}

	timerStr := fmt.Sprintf("⏱ TIME: %04.1fs", remSec)
	statusExtra := fmt.Sprintf("%s | Bonus: +%d", timerStr, int(remSec*2))
	ta.DrawCommonHUD(s, originX, originY, arenaW, "TIME ATTACK", statusExtra)

	borderY := originY + 3
	ta.DrawArenaBorder(s, originX, borderY, arenaW, arenaH)
	ta.DrawSnakeElements(s, originX+1, borderY+1)

	// Bottom instructions with visual timer banner
	instruction := fmt.Sprintf("Race the Clock! [Remaining: %.1fs]  |  P: Pause  |  Q: Quit", remSec)
	s.CenterText(borderY+arenaH+1, instruction, timerColor, th.Background)

	if ta.State == StateOver {
		resultDetail := fmt.Sprintf("Final Sprint Score: %d points", ta.ScoreVal)
		if ta.BestScore > 0 {
			if ta.ScoreVal > ta.BestScore {
				resultDetail += " (★ Beat Personal Best!)"
			} else {
				resultDetail += fmt.Sprintf(" (Personal Best: %d)", ta.BestScore)
			}
		}
		ta.DrawGameOverOverlay(s, w, h, "TIME ATTACK", resultDetail)
	}
}
