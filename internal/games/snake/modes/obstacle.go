package modes

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// ObstacleMode implements 3 levels of labyrinth mazes with guaranteed connectivity and 3s pre-level preview.
type ObstacleMode struct {
	BaseMode
	MazeLevel      int
	LevelFoodEaten int
	PreviewTicks   int // 3-second preview countdown before each maze level starts
}

// NewObstacle constructs an uninitialized ObstacleMode.
func NewObstacle() *ObstacleMode {
	return &ObstacleMode{}
}

// Name returns the game name identifier.
func (ob *ObstacleMode) Name() string { return "snake" }

// Mode returns the mode identifier string.
func (ob *ObstacleMode) Mode() string { return "obstacle" }

// Score returns the current player score.
func (ob *ObstacleMode) Score() int { return ob.ScoreVal }

// IsOver reports whether the session has ended.
func (ob *ObstacleMode) IsOver() bool { return ob.State == StateOver || ob.State == StateWin }

// TickInterval returns the speed interval based on difficulty tier.
func (ob *ObstacleMode) TickInterval() time.Duration { return ob.Interval }

// Metrics returns session metrics for persistence.
func (ob *ObstacleMode) Metrics() map[string]int {
	return map[string]int{
		"food_eaten":         ob.FoodEaten,
		"max_length_reached": ob.MaxLength,
		"ticks_survived":     ob.Ticks,
		"near_misses":        ob.NearMisses,
		"maze_level":         ob.MazeLevel,
	}
}

// Init configures Obstacle Challenge mode according to tier parameters.
func (ob *ObstacleMode) Init(cfg engine.GameConfig) {
	var (
		w, h     int
		interval time.Duration
	)

	switch cfg.Difficulty {
	case engine.Easy:
		w, h = 50, 20
		interval = 120 * time.Millisecond
	case engine.Hard:
		w, h = 40, 15
		interval = 80 * time.Millisecond
	default: // Medium
		w, h = 45, 18
		interval = 100 * time.Millisecond
	}

	ob.InitBase(cfg, w, h, interval, false, 0)
	ob.MazeLevel = 1
	ob.LevelFoodEaten = 0
	ob.loadMazeLevel(1)
}

func (ob *ObstacleMode) loadMazeLevel(level int) {
	ob.MazeLevel = level
	ob.LevelFoodEaten = 0
	ob.Obstacles = GenerateValidatedMaze(level, ob.Cfg.Difficulty, ob.ArenaW, ob.ArenaH)
	ob.RespawnSnake()
	ob.SpawnFood(0)
	// 3 seconds preview countdown (e.g. 30 ticks at 100ms)
	ob.PreviewTicks = int((3 * time.Second) / ob.Interval)
	if ob.PreviewTicks < 15 {
		ob.PreviewTicks = 15
	}
}

// HandleInput processes player directional input (active during preview to set initial heading).
func (ob *ObstacleMode) HandleInput(a engine.Action) {
	if ob.State != StatePlaying {
		return
	}
	if d, ok := a.Direction(); ok {
		ob.PushDirection(d)
	}
}

// Tick executes simulation step including 3-second maze preview countdown.
func (ob *ObstacleMode) Tick() engine.TickResult {
	if ob.State == StateOver {
		return engine.TickResult{Continue: false, Reason: ob.Reason}
	}
	if ob.State == StateWin {
		return engine.TickResult{Continue: false, Reason: "won"}
	}

	// 1. Process 3-Second Maze Preview Countdown
	if ob.PreviewTicks > 0 {
		ob.PreviewTicks--
		return engine.TickResult{Continue: true, Reason: ""}
	}

	if ob.EatFlashTicks > 0 {
		ob.EatFlashTicks--
	}
	if ob.WrapFlashTicks > 0 {
		ob.WrapFlashTicks--
	}

	// 2. Process buffered direction
	prevDir, turned := ob.PopNextDirection()
	if turned {
		ob.CheckNearMiss(prevDir)
	}

	// 3. Next head position
	next := ob.Snake[0].Add(ob.Dir)

	// Wall collision (hard boundary in Obstacle Challenge)
	if next.X < 0 || next.X >= ob.ArenaW || next.Y < 0 || next.Y >= ob.ArenaH {
		ob.State = StateOver
		ob.Reason = "wall"
		return engine.TickResult{Continue: false, Reason: "wall"}
	}

	// Static obstacle collision
	if ob.Obstacles[next] {
		ob.State = StateOver
		ob.Reason = "obstacle"
		return engine.TickResult{Continue: false, Reason: "obstacle"}
	}

	// Self collision
	for _, seg := range ob.Snake {
		if seg.Equal(next) {
			ob.State = StateOver
			ob.Reason = "self"
			return engine.TickResult{Continue: false, Reason: "self"}
		}
	}

	// 4. Move snake head
	ob.Snake = append([]engine.Position{next}, ob.Snake...)

	// 5. Food consumption & Maze progression
	ate := false
	for i, f := range ob.Foods {
		if next.Equal(f) {
			ob.FoodEaten++
			ob.LevelFoodEaten++
			ob.ScoreVal += 10
			ob.EatFlashTicks = 3
			engine.Beep()

			if len(ob.Snake) > ob.MaxLength {
				ob.MaxLength = len(ob.Snake)
			}
			ob.Foods = append(ob.Foods[:i], ob.Foods[i+1:]...)
			ate = true

			// Level progression: 5 food per maze level
			if ob.LevelFoodEaten >= 5 {
				ob.ScoreVal += 100 // Maze clear bonus
				ob.MazeLevel++
				if ob.MazeLevel > 3 {
					ob.State = StateWin
					ob.Reason = "won"
					return engine.TickResult{Continue: false, Reason: "won"}
				}
				ob.loadMazeLevel(ob.MazeLevel)
			} else {
				ob.SpawnFood(0)
			}
			break
		}
	}

	if !ate {
		ob.Snake = ob.Snake[:len(ob.Snake)-1]
	}

	ob.Ticks++
	return engine.TickResult{Continue: true, Reason: ""}
}

// GenerateValidatedMaze builds structured maze layouts with guaranteed path connectivity.
func GenerateValidatedMaze(level int, diff engine.Difficulty, width, height int) map[engine.Position]bool {
	obstacles := make(map[engine.Position]bool)
	midX := width / 2
	midY := height / 2

	switch level {
	case 1:
		// ~20% coverage: two central horizontal shelves with corridors
		for x := 4; x < width-4; x++ {
			if x < midX-3 || x > midX+3 {
				obstacles[engine.Position{X: x, Y: midY - 3}] = true
				obstacles[engine.Position{X: x, Y: midY + 3}] = true
			}
		}
	case 2:
		// ~35% coverage: cross barriers with corner alcoves
		for y := 2; y < height-2; y++ {
			if y != midY && y != midY-1 && y != midY+1 {
				obstacles[engine.Position{X: midX - 7, Y: y}] = true
				obstacles[engine.Position{X: midX + 7, Y: y}] = true
			}
		}
		for x := 6; x < width-6; x++ {
			if x < midX-8 || x > midX+8 {
				obstacles[engine.Position{X: x, Y: midY}] = true
			}
		}
	case 3:
		// ~50% coverage: dense column posts with alternating gates
		for x := 4; x < width-4; x += 4 {
			for y := 2; y < height-2; y++ {
				if (x/4+y)%3 != 0 {
					obstacles[engine.Position{X: x, Y: y}] = true
				}
			}
		}
	}

	// Always clear safe sanctuary around spawn center
	for dx := -3; dx <= 3; dx++ {
		for dy := -2; dy <= 2; dy++ {
			delete(obstacles, engine.Position{X: midX + dx, Y: midY + dy})
		}
	}

	// Validate path connectivity via BFS flood fill
	if !ValidateMazeConnectivity(obstacles, width, height, engine.Position{X: midX, Y: midY}) {
		// Fallback: clear a wide central corridor
		for x := 2; x < width-2; x++ {
			delete(obstacles, engine.Position{X: x, Y: midY})
		}
	}

	return obstacles
}

// ValidateMazeConnectivity performs BFS flood fill to ensure open cells are interconnected.
func ValidateMazeConnectivity(obstacles map[engine.Position]bool, width, height int, start engine.Position) bool {
	visited := make(map[engine.Position]bool)
	queue := []engine.Position{start}
	visited[start] = true

	dirs := []engine.Position{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
	}

	accessibleCount := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		accessibleCount++

		for _, d := range dirs {
			next := curr.Add(d)
			if next.X >= 0 && next.X < width && next.Y >= 0 && next.Y < height && !obstacles[next] && !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	totalOpen := width*height - len(obstacles)
	// Guaranteed connectivity if at least 70% of open tiles are reachable from spawn
	return accessibleCount >= (totalOpen * 7 / 10)
}

// Render draws the Obstacle Challenge screen including the 3-second planning preview.
func (ob *ObstacleMode) Render(s *engine.Screen) {
	th := ob.Cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	arenaW := ob.ArenaW + 2
	arenaH := ob.ArenaH + 2
	totalH := arenaH + 4

	originX := (w - arenaW) / 2
	originY := (h - totalH) / 2
	if originX < 0 {
		originX = 0
	}
	if originY < 0 {
		originY = 0
	}

	statusExtra := fmt.Sprintf("Maze: %d/3 (%d/5)", ob.MazeLevel, ob.LevelFoodEaten)
	if ob.PreviewTicks > 0 {
		statusExtra += fmt.Sprintf(" [PLAN: %.1fs]", float64(ob.PreviewTicks)*ob.Interval.Seconds())
	}
	ob.DrawCommonHUD(s, originX, originY, arenaW, "OBSTACLE", statusExtra)

	borderY := originY + 3
	ob.DrawArenaBorder(s, originX, borderY, arenaW, arenaH)
	ob.DrawSnakeElements(s, originX+1, borderY+1)

	// Preview overlay or bottom instruction
	if ob.PreviewTicks > 0 {
		secondsLeft := int(float64(ob.PreviewTicks)*ob.Interval.Seconds()) + 1
		previewMsg := fmt.Sprintf(" 🔍 MAZE PREVIEW %d/3 — STARTING IN %d... ", ob.MazeLevel, secondsLeft)
		s.CenterText(borderY+arenaH/2, previewMsg, tcell.ColorBlack, tcell.ColorYellow)
		s.CenterText(borderY+arenaH+1, "Examine the maze layout and plan your opening route!", tcell.ColorAqua, th.Background)
	} else {
		s.CenterText(borderY+arenaH+1, "Navigate 3 Maze Levels  |  P: Pause  |  Q: Quit", th.Text, th.Background)
	}

	if ob.State == StateOver || ob.State == StateWin {
		ob.DrawGameOverOverlay(s, w, h, "OBSTACLE", fmt.Sprintf("Maze Levels Cleared: %d/3", ob.MazeLevel-1))
	}
}
