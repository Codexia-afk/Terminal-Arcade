package snake

import (
	"testing"
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

func newTestSnake(diff engine.Difficulty, mode Mode) *Game {
	g := New()
	cfg := engine.GameConfig{
		Difficulty: diff,
		Theme:      engine.GetTheme("Neon"),
		Mode:       string(mode),
	}
	g.Init(cfg)
	g.state = statePlaying
	return g
}

func TestDifficultyParameters(t *testing.T) {
	easy := Profiles(engine.Easy)
	med := Profiles(engine.Medium)
	hard := Profiles(engine.Hard)

	if !easy.Wrap || med.Wrap || hard.Wrap {
		t.Errorf("expected only Easy to have wrap enabled")
	}
	if !hard.Obstacles || easy.Obstacles || med.Obstacles {
		t.Errorf("expected only Hard to have obstacles enabled")
	}
	if !(hard.Interval < med.Interval && med.Interval < easy.Interval) {
		t.Errorf("expected Hard to be fastest and Easy slowest")
	}
	if !(hard.Width < med.Width && med.Width < easy.Width) {
		t.Errorf("expected Hard arena to be smallest")
	}
}

func TestSnakeModesList(t *testing.T) {
	modes := Modes()
	if len(modes) != 5 {
		t.Fatalf("expected 5 Snake modes, got %d", len(modes))
	}
	expected := []Mode{ModeClassic, ModeZen, ModeSurvival, ModeTimeAttack, ModeObstacle}
	for i, exp := range expected {
		if modes[i].Mode != exp {
			t.Errorf("expected mode %d to be %s, got %s", i, exp, modes[i].Mode)
		}
	}
}

func TestWrapMode(t *testing.T) {
	g := newTestSnake(engine.Easy, ModeClassic)
	g.snake = []engine.Position{{X: g.profile.Width - 1, Y: 5}}
	g.dir = engine.Position{X: 1, Y: 0}
	g.food = engine.Position{X: 2, Y: 2}

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("expected snake to wrap in Easy mode, but ended: %s", res.Reason)
	}
	if g.snake[0].X != 0 || g.snake[0].Y != 5 {
		t.Fatalf("expected head to wrap to X=0, got (%d, %d)", g.snake[0].X, g.snake[0].Y)
	}
}

func TestWallCollision(t *testing.T) {
	// Hard tier has 1 life: immediate termination on collision
	g := newTestSnake(engine.Hard, ModeClassic)
	g.snake = []engine.Position{{X: 0, Y: 5}}
	g.dir = engine.Position{X: -1, Y: 0}

	res := g.Tick()
	if res.Continue || res.Reason != "wall" {
		t.Fatalf("expected wall collision on 1 life, got res: %+v", res)
	}
	if !g.IsOver() {
		t.Fatalf("expected game to be over on wall collision")
	}

	// Medium tier has 2 lives: 1st collision decrements life, 2nd terminates
	g2 := newTestSnake(engine.Medium, ModeClassic)
	g2.snake = []engine.Position{{X: 0, Y: 5}}
	g2.dir = engine.Position{X: -1, Y: 0}
	res2 := g2.Tick()
	if !res2.Continue {
		t.Fatalf("expected life decrement on first collision, got termination: %s", res2.Reason)
	}
	// 2nd collision terminates
	g2.snake = []engine.Position{{X: 0, Y: 5}}
	g2.dir = engine.Position{X: -1, Y: 0}
	res3 := g2.Tick()
	if res3.Continue || res3.Reason != "wall" {
		t.Fatalf("expected wall collision termination on 0 lives, got: %+v", res3)
	}
}

func TestSelfCollision(t *testing.T) {
	// Hard tier has 1 life: immediate self collision termination
	g := newTestSnake(engine.Hard, ModeClassic)
	g.snake = []engine.Position{
		{X: 5, Y: 5},
		{X: 5, Y: 6},
		{X: 6, Y: 6},
		{X: 6, Y: 5},
	}
	g.dir = engine.Position{X: 1, Y: 0}

	res := g.Tick()
	if res.Continue || res.Reason != "self" {
		t.Fatalf("expected self collision on 1 life, got res: %+v", res)
	}
}

func TestZenModeSafetyAndFlow(t *testing.T) {
	g := newTestSnake(engine.Hard, ModeZen)
	// Zen must force wrap to true
	if !g.profile.Wrap {
		t.Errorf("expected Zen mode to force wrap=true")
	}

	// Self-collision in Zen mode should NOT terminate the game
	g.snake = []engine.Position{
		{X: 5, Y: 5},
		{X: 5, Y: 6},
		{X: 6, Y: 6},
		{X: 6, Y: 5},
	}
	g.dir = engine.Position{X: 1, Y: 0}
	g.food = engine.Position{X: 10, Y: 10}

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("expected Zen mode to survive self collision, got: %s", res.Reason)
	}

	// Test Flow Bonus: Simulate 10 seconds of ticks
	tickCount := int((10 * time.Second) / g.profile.Interval)
	for i := 0; i < tickCount+1; i++ {
		_ = g.Tick()
	}
	if g.flowBonusPoints < 5 {
		t.Errorf("expected flow bonus points in Zen mode, got %d", g.flowBonusPoints)
	}
}

func TestSurvivalModeWavesAndEnemyAI(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeSurvival)
	if g.currentWave != 1 {
		t.Fatalf("expected wave 1, got %d", g.currentWave)
	}
	if len(g.enemies) == 0 {
		t.Fatalf("expected at least 1 enemy in Survival mode")
	}

	// Test Enemy AI movement towards snake
	enemy := g.enemies[0]
	snakeHead := g.snake[0]
	stepped := g.stepEnemyAI(enemy.Pos, snakeHead)
	// Stepped position should be closer to snakeHead than before
	oldDist := abs(snakeHead.X-enemy.Pos.X) + abs(snakeHead.Y-enemy.Pos.Y)
	newDist := abs(snakeHead.X-stepped.X) + abs(snakeHead.Y-stepped.Y)
	if newDist >= oldDist && oldDist > 0 {
		t.Errorf("expected enemy to move closer to snake head")
	}

	// Test Wave progression: eating 5 food advances to wave 2
	for i := 0; i < 5; i++ {
		nextHead := g.snake[0].Add(g.dir)
		g.food = nextHead
		res := g.Tick()
		if !res.Continue {
			t.Fatalf("unexpected termination while advancing wave: %s", res.Reason)
		}
	}
	if g.currentWave != 2 {
		t.Errorf("expected wave 2 after 5 food, got %d", g.currentWave)
	}
}

func TestTimeAttackModeTimerAndBonus(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeTimeAttack)
	if g.timeRemaining != TimeAttackDuration {
		t.Fatalf("expected initial time %v, got %v", TimeAttackDuration, g.timeRemaining)
	}

	// Eat food and check time bonus
	nextHead := g.snake[0].Add(g.dir)
	g.food = nextHead
	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}
	// Score should be 10 + timeRemaining.Seconds() * 2
	expectedMinScore := 10 + int(50*2) // >= 110 pts
	if g.Score() < expectedMinScore {
		t.Errorf("expected Time Attack score >= %d, got %d", expectedMinScore, g.Score())
	}

	// Fast-forward timer to 0 and verify time_up termination
	g.timeRemaining = 10 * time.Millisecond
	res = g.Tick()
	if res.Continue || res.Reason != "time_up" {
		t.Fatalf("expected time_up termination when timer expires, got %+v", res)
	}
}

func TestObstacleChallengeModeProgression(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeObstacle)
	if g.mazeLevel != 1 {
		t.Fatalf("expected maze level 1, got %d", g.mazeLevel)
	}
	if len(g.obstacles) == 0 {
		t.Fatalf("expected maze obstacles to be generated")
	}

	// Eat 5 food to advance to Level 2
	for i := 0; i < 5; i++ {
		nextHead := g.snake[0].Add(g.dir)
		g.food = nextHead
		res := g.Tick()
		if !res.Continue {
			t.Fatalf("unexpected termination advancing maze: %s", res.Reason)
		}
	}
	if g.mazeLevel != 2 {
		t.Errorf("expected maze level 2 after 5 food, got %d", g.mazeLevel)
	}
}

func TestTwoStepInputQueue(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeClassic)
	// Snake is moving Right (1, 0)
	// Player quickly taps Up then Left
	g.HandleInput(engine.ActionUp)
	g.HandleInput(engine.ActionLeft)

	if len(g.dirQueue) != 2 {
		t.Fatalf("expected 2 directions queued, got %d", len(g.dirQueue))
	}
	if g.dirQueue[0] != (engine.Position{X: 0, Y: -1}) {
		t.Errorf("expected first queued dir to be Up, got %+v", g.dirQueue[0])
	}
	if g.dirQueue[1] != (engine.Position{X: -1, Y: 0}) {
		t.Errorf("expected second queued dir to be Left, got %+v", g.dirQueue[1])
	}

	// First tick should consume Up
	g.food = engine.Position{X: 15, Y: 15}
	_ = g.Tick()
	if g.dir != (engine.Position{X: 0, Y: -1}) {
		t.Errorf("expected snake to turn Up on first tick, got %+v", g.dir)
	}

	// Second tick should consume Left
	_ = g.Tick()
	if g.dir != (engine.Position{X: -1, Y: 0}) {
		t.Errorf("expected snake to turn Left on second tick, got %+v", g.dir)
	}
}

func TestSmartFoodSpawner(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeClassic)
	for i := 0; i < 50; i++ {
		g.spawnFood()
		// Must not spawn on snake head + next cell ahead
		nextCell := g.snake[0].Add(g.dir)
		if g.food.Equal(nextCell) {
			t.Errorf("smart food spawner spawned on cell directly ahead of snake head")
		}
		for _, seg := range g.snake {
			if g.food.Equal(seg) {
				t.Errorf("smart food spawner spawned on snake segment: %+v", seg)
			}
		}
	}
}

func TestPersonalBestTracking(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeClassic)
	g.SetBestScore(50)
	g.score = 60
	g.terminate("wall")

	if !g.IsOver() {
		t.Fatalf("expected game to be over")
	}
	if g.score <= g.bestScore {
		t.Errorf("expected current score %d to beat best %d", g.score, g.bestScore)
	}
}

func TestSnakeMetrics(t *testing.T) {
	g := newTestSnake(engine.Medium, ModeClassic)
	g.snake = []engine.Position{{X: 1, Y: 5}, {X: 0, Y: 5}}
	g.dir = engine.Position{X: 1, Y: 0}
	g.obstacles[engine.Position{X: 2, Y: 5}] = true

	// Turn Up to avoid obstacle
	g.dirQueue = []engine.Position{{X: 0, Y: -1}}
	g.food = engine.Position{X: 1, Y: 4}

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	metrics := g.Metrics()
	if metrics["near_misses"] != 1 {
		t.Errorf("expected 1 near miss, got %d", metrics["near_misses"])
	}
	if metrics["food_eaten"] != 1 {
		t.Errorf("expected 1 food eaten, got %d", metrics["food_eaten"])
	}
	if metrics["max_length_reached"] != 3 {
		t.Errorf("expected max length 3, got %d", metrics["max_length_reached"])
	}
	if metrics["ticks_survived"] != 1 {
		t.Errorf("expected 1 tick survived, got %d", metrics["ticks_survived"])
	}
}
