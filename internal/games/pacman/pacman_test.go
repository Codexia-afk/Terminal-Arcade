package pacman

import (
	"math/rand"
	"testing"

	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
)

func newTestPacman(diff engine.Difficulty) *Game {
	g := New()
	cfg := engine.GameConfig{
		Difficulty: diff,
		Theme:      engine.GetTheme("Neon"),
	}
	g.Init(cfg)
	g.state = statePlaying
	return g
}

func TestDifficultyParameters(t *testing.T) {
	easy := Profiles(engine.Easy)
	med := Profiles(engine.Medium)
	hard := Profiles(engine.Hard)

	if easy.GhostCount != 2 || med.GhostCount != 3 || hard.GhostCount != 4 {
		t.Errorf("expected ghost counts 2, 3, 4; got %d, %d, %d", easy.GhostCount, med.GhostCount, hard.GhostCount)
	}
	if !(hard.PelletDuration < med.PelletDuration && med.PelletDuration < easy.PelletDuration) {
		t.Errorf("expected pellet duration Easy > Medium > Hard")
	}
	if easy.Lives != 3 || med.Lives != 3 || hard.Lives != 2 {
		t.Errorf("expected lives: Easy 3, Medium 3, Hard 2; got %d, %d, %d", easy.Lives, med.Lives, hard.Lives)
	}
}

func TestDotConsumptionAndWin(t *testing.T) {
	g := newTestPacman(engine.Medium)
	// Place a single dot directly in front of the player
	startPos := g.playerPos
	dotPos := startPos.Add(engine.Position{X: -1, Y: 0})
	g.maze.Tiles[dotPos.Y][dotPos.X] = TileDot
	g.maze.Remaining = 1 // Only 1 dot left in the entire maze

	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	res := g.Tick()
	if res.Continue {
		t.Fatalf("expected game to end with win when last dot is eaten")
	}
	if res.Reason != "won" {
		t.Fatalf("expected reason 'won', got '%s'", res.Reason)
	}
	if g.Score() != 10 {
		t.Errorf("expected score 10, got %d", g.Score())
	}
	if !g.IsOver() {
		t.Errorf("expected IsOver() to be true on win")
	}
}

func TestPowerPelletAndVulnerableTransitions(t *testing.T) {
	g := newTestPacman(engine.Medium)
	// Place power pellet directly in front of the player
	pelletPos := g.playerPos.Add(engine.Position{X: -1, Y: 0})
	g.maze.Tiles[pelletPos.Y][pelletPos.X] = TilePellet
	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	// Position ghosts safely away so they don't immediately collide
	for _, gh := range g.ghosts {
		gh.Pos = engine.Position{X: 1, Y: 1}
	}

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}
	if g.Score() != 50 {
		t.Errorf("expected 50 score for power pellet, got %d", g.Score())
	}
	if g.vulnerableTicks <= 0 {
		t.Fatalf("expected positive vulnerableTicks")
	}
	for _, gh := range g.ghosts {
		if !gh.Vulnerable {
			t.Errorf("expected ghost %s to be vulnerable", gh.Name)
		}
	}

	// Test eating a vulnerable ghost
	targetGhost := g.ghosts[0]
	nextPlayerPos := g.playerPos.Add(engine.Position{X: -1, Y: 0})
	// Ensure next cell has no dot so score increase is purely from eating the ghost
	g.maze.Tiles[nextPlayerPos.Y][nextPlayerPos.X] = TileEmpty
	targetGhost.Pos = nextPlayerPos
	targetGhost.Dir = engine.Position{X: 0, Y: 0}

	scoreBefore := g.Score()
	res = g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination when eating vulnerable ghost: %s", res.Reason)
	}
	if g.Score() != scoreBefore+200 {
		t.Errorf("expected +200 for eating vulnerable ghost, got score %d (before %d)", g.Score(), scoreBefore)
	}
	if !targetGhost.Eaten {
		t.Errorf("expected eaten ghost to have Eaten=true")
	}
}

func TestGhostAIDecisionLogic(t *testing.T) {
	m := NewMaze()
	rng := rand.New(rand.NewSource(42))

	// 1. Chase ghost AI: in open hallway row 13, should pick step that minimizes distance to Pacman
	chaseGhost := NewGhost(1, "Blinky", engine.Position{X: 8, Y: 13}, tcell.ColorRed, "chase")
	chaseGhost.Pos = engine.Position{X: 8, Y: 13}
	pacmanPos := engine.Position{X: 4, Y: 13} // Pacman is to the left
	pacmanDir := engine.Position{X: -1, Y: 0}

	dir := chaseGhost.DecideNextMove(m, pacmanPos, pacmanDir, rng)
	if dir.X != -1 || dir.Y != 0 {
		t.Errorf("chase ghost expected to move Left (-1, 0), got (%d, %d)", dir.X, dir.Y)
	}

	// 2. Ambush ghost AI: targets 4 tiles ahead of Pacman
	ambushGhost := NewGhost(2, "Pinky", engine.Position{X: 8, Y: 13}, tcell.ColorFuchsia, "ambush")
	ambushGhost.Pos = engine.Position{X: 8, Y: 13}
	pacmanPos = engine.Position{X: 2, Y: 13}
	pacmanDir = engine.Position{X: 1, Y: 0}

	dirAmbush := ambushGhost.DecideNextMove(m, pacmanPos, pacmanDir, rng)
	if dirAmbush.X != -1 || dirAmbush.Y != 0 {
		t.Errorf("ambush ghost expected to move towards ambush target Left (-1, 0), got (%d, %d)", dirAmbush.X, dirAmbush.Y)
	}

	// 3. Eaten ghost AI: must move towards Home
	eatenGhost := NewGhost(3, "Inky", engine.Position{X: 14, Y: 13}, tcell.ColorAqua, "chase")
	eatenGhost.Pos = engine.Position{X: 8, Y: 13}
	eatenGhost.Eaten = true

	dirHome := eatenGhost.DecideNextMove(m, pacmanPos, pacmanDir, rng)
	if dirHome.X != 1 || dirHome.Y != 0 {
		t.Errorf("eaten ghost expected to move Right (+1, 0) towards Home, got (%d, %d)", dirHome.X, dirHome.Y)
	}
}

func TestPacmanMetrics(t *testing.T) {
	g := newTestPacman(engine.Medium)
	// Safely position ghosts
	for _, gh := range g.ghosts {
		gh.Pos = engine.Position{X: 1, Y: 1}
	}

	// 1. Eat a dot
	startPos := g.playerPos
	dotPos := startPos.Add(engine.Position{X: -1, Y: 0})
	g.maze.Tiles[dotPos.Y][dotPos.X] = TileDot
	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	metrics := g.Metrics()
	if metrics["dots_eaten"] != 1 {
		t.Errorf("expected 1 dot eaten, got %d", metrics["dots_eaten"])
	}

	// 2. Eat a pellet
	pelletPos := g.playerPos.Add(engine.Position{X: -1, Y: 0})
	g.maze.Tiles[pelletPos.Y][pelletPos.X] = TilePellet
	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	res = g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	metrics = g.Metrics()
	if metrics["power_pellets_used"] != 1 {
		t.Errorf("expected 1 power pellet used, got %d", metrics["power_pellets_used"])
	}

	// 3. Eat a vulnerable ghost
	g.ghosts[0].Pos = g.playerPos.Add(engine.Position{X: -1, Y: 0})
	g.ghosts[0].Vulnerable = true
	g.maze.Tiles[g.ghosts[0].Pos.Y][g.ghosts[0].Pos.X] = TileEmpty
	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	res = g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	metrics = g.Metrics()
	if metrics["ghosts_eaten"] != 1 {
		t.Errorf("expected 1 ghost eaten, got %d", metrics["ghosts_eaten"])
	}

	// 4. Test losing a life
	g.ghosts[0].Vulnerable = false
	g.ghosts[0].Eaten = false
	g.ghosts[0].Pos = g.playerPos

	// Collision resolution on next tick or direct resolution
	g.resolveCollisions(g.playerPos, nil)
	metrics = g.Metrics()
	if metrics["lives_lost"] != 1 {
		t.Errorf("expected 1 life lost, got %d", metrics["lives_lost"])
	}
}
