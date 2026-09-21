package snake

import (
	"testing"

	"goarcade/internal/engine"
)

func newTestSnake(diff engine.Difficulty) *Game {
	g := New()
	cfg := engine.GameConfig{
		Difficulty: diff,
		Theme:      engine.GetTheme("Neon"),
	}
	g.Init(cfg)
	// Skip title screen directly to playing state
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

func TestWrapMode(t *testing.T) {
	g := newTestSnake(engine.Easy)
	// Place snake at the right edge, moving right
	g.snake = []engine.Position{{X: g.profile.Width - 1, Y: 5}}
	g.dir = engine.Position{X: 1, Y: 0}
	g.queuedDir = g.dir
	g.food = engine.Position{X: 2, Y: 2} // Far away from snake

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("expected snake to wrap in Easy mode, but ended: %s", res.Reason)
	}
	if g.snake[0].X != 0 || g.snake[0].Y != 5 {
		t.Fatalf("expected head to wrap to X=0, got (%d, %d)", g.snake[0].X, g.snake[0].Y)
	}
}

func TestWallCollision(t *testing.T) {
	g := newTestSnake(engine.Medium)
	// Place snake at left edge, moving left
	g.snake = []engine.Position{{X: 0, Y: 5}}
	g.dir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.dir

	res := g.Tick()
	if res.Continue || res.Reason != "wall" {
		t.Fatalf("expected wall collision, got res: %+v", res)
	}
	if !g.IsOver() {
		t.Fatalf("expected game to be over on wall collision")
	}
}

func TestSelfCollision(t *testing.T) {
	g := newTestSnake(engine.Medium)
	// Snake arranged in a loop where moving right collides with its own body segment
	g.snake = []engine.Position{
		{X: 5, Y: 5},
		{X: 5, Y: 6},
		{X: 6, Y: 6},
		{X: 6, Y: 5},
	}
	g.dir = engine.Position{X: 1, Y: 0}
	g.queuedDir = g.dir

	res := g.Tick()
	if res.Continue || res.Reason != "self" {
		t.Fatalf("expected self collision, got res: %+v", res)
	}
}

func TestFoodConsumptionAndGrowth(t *testing.T) {
	g := newTestSnake(engine.Medium)
	g.snake = []engine.Position{
		{X: 5, Y: 5},
		{X: 4, Y: 5},
		{X: 3, Y: 5},
	}
	initialLen := len(g.snake)
	// Place food directly ahead
	g.food = engine.Position{X: 6, Y: 5}
	g.dir = engine.Position{X: 1, Y: 0}
	g.queuedDir = g.dir

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected game termination: %s", res.Reason)
	}
	if g.Score() != 10 {
		t.Errorf("expected score 10, got %d", g.Score())
	}
	if len(g.snake) != initialLen+1 {
		t.Errorf("expected snake length to grow to %d, got %d", initialLen+1, len(g.snake))
	}
}

func TestObstacleSpawnAndCollision(t *testing.T) {
	g := newTestSnake(engine.Hard)
	g.snake = []engine.Position{{X: 5, Y: 5}}
	obs := engine.Position{X: 6, Y: 5}
	g.obstacles[obs] = true
	g.dir = engine.Position{X: 1, Y: 0}
	g.queuedDir = g.dir

	res := g.Tick()
	if res.Continue || res.Reason != "obstacle" {
		t.Fatalf("expected obstacle collision, got: %+v", res)
	}
}

func TestPersonalBestTracking(t *testing.T) {
	g := newTestSnake(engine.Medium)
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
