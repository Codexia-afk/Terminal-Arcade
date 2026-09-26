package modes

import (
	"testing"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

func TestClassicModeLivesAndScoring(t *testing.T) {
	c := NewClassic()
	c.Init(engine.GameConfig{Difficulty: engine.Medium, Theme: engine.GetTheme("Neon")})

	if c.Lives != 2 {
		t.Fatalf("expected 2 lives for Classic Medium, got %d", c.Lives)
	}

	// Eat food
	nextHead := c.Snake[0].Add(c.Dir)
	c.Foods = []engine.Position{nextHead}
	res := c.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}
	if c.Score() != 10 {
		t.Errorf("expected score 10, got %d", c.Score())
	}
	if len(c.Snake) != 4 {
		t.Errorf("expected snake length 4, got %d", len(c.Snake))
	}

	// Force 1st wall hit -> decrements life from 2 to 1
	c.Snake = []engine.Position{{X: 0, Y: 5}}
	c.Dir = engine.Position{X: -1, Y: 0}
	res = c.Tick()
	if !res.Continue {
		t.Fatalf("expected game to continue with 1 life remaining, got %s", res.Reason)
	}
	if c.Lives != 1 {
		t.Errorf("expected 1 life remaining, got %d", c.Lives)
	}

	// Force 2nd wall hit -> terminates
	c.Snake = []engine.Position{{X: 0, Y: 5}}
	c.Dir = engine.Position{X: -1, Y: 0}
	res = c.Tick()
	if res.Continue || res.Reason != "wall" {
		t.Fatalf("expected termination with 'wall', got %+v", res)
	}
	if !c.IsOver() {
		t.Errorf("expected game to be over")
	}
}

func TestZenModeFlowAndNoDeath(t *testing.T) {
	z := NewZen()
	z.Init(engine.GameConfig{Difficulty: engine.Easy, Theme: engine.GetTheme("Neon")})

	// Wrap across border
	z.Snake = []engine.Position{{X: z.ArenaW - 1, Y: 5}}
	z.Dir = engine.Position{X: 1, Y: 0}
	z.Foods = []engine.Position{{X: 10, Y: 10}}
	res := z.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination in Zen: %s", res.Reason)
	}
	if z.Snake[0].X != 0 {
		t.Errorf("expected head to wrap to X=0, got %d", z.Snake[0].X)
	}

	// Self-collision should be harmless
	z.Snake = []engine.Position{
		{X: 5, Y: 5},
		{X: 5, Y: 6},
		{X: 6, Y: 6},
		{X: 6, Y: 5},
	}
	z.Dir = engine.Position{X: 1, Y: 0}
	res = z.Tick()
	if !res.Continue {
		t.Fatalf("Zen mode must not die on self-collision: %s", res.Reason)
	}

	// Turning resets flow timer
	z.FlowDuration = 5 * time.Second
	z.HandleInput(engine.ActionUp)
	if z.FlowDuration != 0 {
		t.Errorf("expected flow duration to reset on turn, got %v", z.FlowDuration)
	}

	// Flow bonus +5 after 10s uninterrupted
	z.FlowDuration = 9900 * time.Millisecond
	z.Interval = 200 * time.Millisecond
	_ = z.Tick()
	if z.FlowBonusPoints < 5 {
		t.Errorf("expected flow bonus points, got %d", z.FlowBonusPoints)
	}
}

func TestSurvivalModeEnemyAIAndWaves(t *testing.T) {
	su := NewSurvival()
	su.Init(engine.GameConfig{Difficulty: engine.Hard, Theme: engine.GetTheme("Neon")})

	if su.Lives != 1 {
		t.Errorf("expected 1 life on Hard Survival, got %d", su.Lives)
	}
	if len(su.Enemies) == 0 {
		t.Fatalf("expected enemies initialized in wave 1")
	}

	// Test Greedy Manhattan step
	enemy := su.Enemies[0].Pos
	target := su.Snake[0]
	stepped := su.StepEnemyAI(enemy, target)
	oldDist := abs(target.X-enemy.X) + abs(target.Y-enemy.Y)
	newDist := abs(target.X-stepped.X) + abs(target.Y-stepped.Y)
	if newDist >= oldDist && oldDist > 0 {
		t.Errorf("expected enemy step to reduce Manhattan distance")
	}

	// Progress through wave 1: eat 5 foods
	for i := 0; i < 5; i++ {
		nextHead := su.Snake[0].Add(su.Dir)
		su.Foods = []engine.Position{nextHead}
		res := su.Tick()
		if !res.Continue {
			t.Fatalf("unexpected failure advancing wave: %s", res.Reason)
		}
	}

	if su.Wave != 2 {
		t.Errorf("expected wave 2 after 5 foods, got %d", su.Wave)
	}
	if su.PatrolTicks <= 0 {
		t.Errorf("expected patrol phase to be active after wave clear")
	}
}

func TestTimeAttackTimerAndSpeedBonus(t *testing.T) {
	ta := NewTimeAttack()
	ta.Init(engine.GameConfig{Difficulty: engine.Medium, Theme: engine.GetTheme("Neon")})

	if ta.TimeRemaining != 60*time.Second {
		t.Fatalf("expected 60s timer on Medium, got %v", ta.TimeRemaining)
	}

	// Eat food and check speed bonus (10 + remaining * 2)
	nextHead := ta.Snake[0].Add(ta.Dir)
	ta.Foods = []engine.Position{nextHead}
	res := ta.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	expectedBonus := int(ta.TimeRemaining.Seconds() * 2)
	expectedScore := 10 + expectedBonus
	if ta.Score() != expectedScore {
		t.Errorf("expected score %d, got %d", expectedScore, ta.Score())
	}

	// Test timer expiration
	ta.TimeRemaining = 10 * time.Millisecond
	res = ta.Tick()
	if res.Continue || res.Reason != "time_up" {
		t.Fatalf("expected 'time_up' on expiration, got %+v", res)
	}
}

func TestObstacleModeMazeValidation(t *testing.T) {
	ob := NewObstacle()
	ob.Init(engine.GameConfig{Difficulty: engine.Medium, Theme: engine.GetTheme("Neon")})

	if ob.MazeLevel != 1 {
		t.Errorf("expected starting maze level 1, got %d", ob.MazeLevel)
	}
	if len(ob.Obstacles) == 0 {
		t.Fatalf("expected maze obstacles to be generated")
	}

	// Check BFS connectivity of generated maze
	connected := ValidateMazeConnectivity(ob.Obstacles, ob.ArenaW, ob.ArenaH, engine.Position{X: ob.ArenaW / 2, Y: ob.ArenaH / 2})
	if !connected {
		t.Errorf("expected generated maze to satisfy connectivity validation")
	}

	// Skip preview for testing progression
	ob.PreviewTicks = 0

	// Eat 5 foods to clear level 1
	for i := 0; i < 5; i++ {
		ob.PreviewTicks = 0
		nextHead := ob.Snake[0].Add(ob.Dir)
		ob.Foods = []engine.Position{nextHead}
		res := ob.Tick()
		if !res.Continue {
			t.Fatalf("unexpected termination in maze: %s", res.Reason)
		}
	}

	if ob.MazeLevel != 2 {
		t.Errorf("expected maze level 2 after 5 foods, got %d", ob.MazeLevel)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
