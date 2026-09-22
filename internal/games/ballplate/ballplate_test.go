package ballplate

import (
	"math"
	"testing"

	"goarcade/internal/engine"
)

func newTestBallPlate(diff engine.Difficulty) *Game {
	g := New()
	cfg := engine.GameConfig{
		Difficulty: diff,
		Theme:      engine.GetTheme("Neon"),
	}
	g.Init(cfg)
	g.state = statePlaying
	g.attached = false // Unattach ball for testing physics
	return g
}

func TestDifficultyParameters(t *testing.T) {
	easy := Profiles(engine.Easy)
	med := Profiles(engine.Medium)
	hard := Profiles(engine.Hard)

	if easy.PlateWidth != 7 || med.PlateWidth != 5 || hard.PlateWidth != 3 {
		t.Errorf("expected plate widths 7, 5, 3; got %d, %d, %d", easy.PlateWidth, med.PlateWidth, hard.PlateWidth)
	}
	if !(easy.BaseBallSpeed < med.BaseBallSpeed && med.BaseBallSpeed < hard.BaseBallSpeed) {
		t.Errorf("expected speed Easy < Medium < Hard")
	}
	if easy.Lives != 4 || med.Lives != 3 || hard.Lives != 2 {
		t.Errorf("expected lives 4, 3, 2; got %d, %d, %d", easy.Lives, med.Lives, hard.Lives)
	}
	if !easy.GappedBricks || med.GappedBricks || hard.GappedBricks {
		t.Errorf("expected only Easy to have gapped bricks")
	}
	if !hard.ToughBricks || easy.ToughBricks || med.ToughBricks {
		t.Errorf("expected only Hard to have tough bricks")
	}
}

func TestPlateAngleReflection(t *testing.T) {
	g := newTestBallPlate(engine.Medium)
	center := float64(g.plateX) + float64(g.plateW)/2.0

	// 1. Hit directly at center of plate
	g.ResolvePlateBounce(center)
	if math.Abs(g.ballVx) > 0.05 {
		t.Errorf("expected center hit to yield ~0 horizontal velocity, got %f", g.ballVx)
	}
	if g.ballVy >= 0 {
		t.Errorf("expected upward bounce Vy < 0, got %f", g.ballVy)
	}

	// 2. Hit at right side of plate
	g.ResolvePlateBounce(center + 2.0)
	if g.ballVx <= 0 {
		t.Errorf("expected right side hit to yield Vx > 0, got %f", g.ballVx)
	}
	if g.ballVy >= 0 {
		t.Errorf("expected upward bounce Vy < 0, got %f", g.ballVy)
	}

	// 3. Hit at left side of plate
	g.ResolvePlateBounce(center - 2.0)
	if g.ballVx >= 0 {
		t.Errorf("expected left side hit to yield Vx < 0, got %f", g.ballVx)
	}
	if g.ballVy >= 0 {
		t.Errorf("expected upward bounce Vy < 0, got %f", g.ballVy)
	}
}

func TestBrickHitDetectionAndToughBrick(t *testing.T) {
	g := newTestBallPlate(engine.Hard)
	initRem := g.remBricks

	// Test tough brick with health 2
	toughBrick := &g.bricks[0]
	toughBrick.Health = 2
	g.ballX = float64(toughBrick.X + 1)
	g.ballY = float64(toughBrick.Y)
	g.ballVy = -0.4

	// First hit: health reduces to 1, brick not destroyed yet
	g.checkBrickCollisions()
	if toughBrick.Health != 1 {
		t.Fatalf("expected tough brick health 1 after first hit, got %d", toughBrick.Health)
	}
	if g.remBricks != initRem {
		t.Errorf("expected remaining count unchanged after first tough hit")
	}
	if g.ballVy <= 0 {
		t.Errorf("expected vertical velocity to reverse downward after hit, got %f", g.ballVy)
	}

	// Second hit on same brick
	g.ballVy = -0.4
	g.checkBrickCollisions()
	if toughBrick.Health != 0 {
		t.Fatalf("expected tough brick destroyed on second hit")
	}
	if g.remBricks != initRem-1 {
		t.Errorf("expected remaining count decremented by 1, got %d", g.remBricks)
	}
	if g.Score() <= 0 {
		t.Errorf("expected score to increase on brick destruction")
	}
}

func TestLifeLossAndGameOver(t *testing.T) {
	g := newTestBallPlate(engine.Hard) // 2 lives
	g.lives = 1

	// Ball drops below the bottom boundary
	g.ballX = 10.0
	g.ballY = float64(g.arenaH + 1)
	g.ballVx = 0.0
	g.ballVy = 0.4

	res := g.Tick()
	if res.Continue {
		t.Fatalf("expected game to end on 0 lives")
	}
	if res.Reason != "lost" {
		t.Errorf("expected reason 'lost', got '%s'", res.Reason)
	}
	if !g.IsOver() {
		t.Errorf("expected IsOver() to be true")
	}
}

func TestBallPlateMetrics(t *testing.T) {
	g := newTestBallPlate(engine.Medium)

	// 1. Test Plate Hits and Longest Rally
	plateCenter := float64(g.plateX) + float64(g.plateW)/2.0
	g.ballX = plateCenter
	g.ballY = float64(g.plateY) - 0.4
	g.ballVy = 0.4

	res := g.Tick()
	if !res.Continue {
		t.Fatalf("unexpected termination: %s", res.Reason)
	}

	metrics := g.Metrics()
	if metrics["plate_hits"] != 1 {
		t.Errorf("expected 1 plate hit, got %d", metrics["plate_hits"])
	}
	if metrics["longest_rally"] != 1 {
		t.Errorf("expected longest rally 1, got %d", metrics["longest_rally"])
	}

	// 2. Test Brick Broken metric
	b := &g.bricks[0]
	b.Health = 1
	g.ballX = float64(b.X)
	g.ballY = float64(b.Y)
	g.checkBrickCollisions()

	metrics = g.Metrics()
	if metrics["bricks_broken"] != 1 {
		t.Errorf("expected 1 brick broken, got %d", metrics["bricks_broken"])
	}

	// 3. Test Life Lost metric resets current rally
	g.ballY = float64(g.arenaH + 1)
	g.Tick()

	metrics = g.Metrics()
	if metrics["lives_lost"] != 1 {
		t.Errorf("expected 1 life lost, got %d", metrics["lives_lost"])
	}
	// Longest rally remains 1 even though life was lost
	if metrics["longest_rally"] != 1 {
		t.Errorf("expected longest rally 1 preserved, got %d", metrics["longest_rally"])
	}
}
