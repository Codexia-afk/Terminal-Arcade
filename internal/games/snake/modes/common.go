package modes

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

// BaseMode contains the common state and mechanics shared across Snake gameplay modes.
type BaseMode struct {
	Cfg        engine.GameConfig
	ArenaW     int
	ArenaH     int
	Interval   time.Duration
	Wrap       bool
	Lives      int
	MaxLives   int
	ScoreVal   int
	State      int // 0=title/preview, 1=playing, 2=over, 3=win
	Snake      []engine.Position
	Dir        engine.Position
	DirQueue   []engine.Position
	Foods      []engine.Position
	Obstacles  map[engine.Position]bool
	Ticks      int
	FoodEaten  int
	MaxLength  int
	NearMisses int
	Rng        *rand.Rand
	Reason     string

	// Audio-Visual feedback
	EatFlashTicks  int
	WrapFlashTicks int
	HitFlashTicks  int

	// Personal Best
	BestScore int
	BestInfo  history.PersonalBestInfo
	HasBest   bool
}

const (
	StatePreview = 0
	StatePlaying = 1
	StateOver    = 2
	StateWin     = 3
)

// InitBase sets up common fields for a Snake mode session.
func (b *BaseMode) InitBase(cfg engine.GameConfig, w, h int, interval time.Duration, wrap bool, lives int) {
	b.Cfg = cfg
	b.ArenaW = w
	b.ArenaH = h
	b.Interval = interval
	b.Wrap = wrap
	b.Lives = lives
	b.MaxLives = lives
	b.ScoreVal = 0
	b.State = StatePlaying
	b.Ticks = 0
	b.FoodEaten = 0
	b.MaxLength = 3
	b.NearMisses = 0
	b.Rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	b.Obstacles = make(map[engine.Position]bool)
	b.DirQueue = make([]engine.Position, 0, 2)
	b.EatFlashTicks = 0
	b.WrapFlashTicks = 0
	b.HitFlashTicks = 0
	b.Reason = ""

	// Initial 3-segment snake in arena center
	cx := w / 2
	cy := h / 2
	b.Snake = []engine.Position{
		{X: cx, Y: cy},
		{X: cx - 1, Y: cy},
		{X: cx - 2, Y: cy},
	}
	b.Dir = engine.Position{X: 1, Y: 0}

	// Load personal best
	if store, err := history.Open(""); err == nil && store != nil {
		b.BestScore = history.HighScore(store.Records, "snake", string(cfg.Difficulty))
		if pb, ok := history.PersonalBestsForMode(store.Records, "snake", cfg.Mode); ok {
			b.BestInfo = pb
			b.HasBest = true
		}
	}
}

// PushDirection buffers up to 2 direction changes preventing 180° neck reversals.
func (b *BaseMode) PushDirection(d engine.Position) {
	lastDir := b.Dir
	if len(b.DirQueue) > 0 {
		lastDir = b.DirQueue[len(b.DirQueue)-1]
	}
	// Prevent 180° reversal into neck
	if (d.X == -lastDir.X && d.Y == -lastDir.Y) || (d.X == lastDir.X && d.Y == lastDir.Y) {
		return
	}
	if len(b.DirQueue) < 2 {
		b.DirQueue = append(b.DirQueue, d)
	} else {
		b.DirQueue[1] = d
	}
}

// PopNextDirection retrieves and applies the next buffered directional step.
// Returns the previous direction and whether a turn occurred.
func (b *BaseMode) PopNextDirection() (prevDir engine.Position, turned bool) {
	prevDir = b.Dir
	if len(b.DirQueue) > 0 {
		b.Dir = b.DirQueue[0]
		b.DirQueue = b.DirQueue[1:]
		turned = prevDir != b.Dir
	}
	return prevDir, turned
}

// SpawnFood places food on an open tile, avoiding snake body, cell ahead of head, and obstacles.
func (b *BaseMode) SpawnFood(extraCount int) {
	b.Foods = nil
	count := 1 + extraCount
	var nextHead engine.Position
	if len(b.Snake) > 0 {
		nextHead = b.Snake[0].Add(b.Dir)
	}

	for i := 0; i < count; i++ {
		for attempts := 0; attempts < 500; attempts++ {
			p := engine.Position{
				X: b.Rng.Intn(b.ArenaW),
				Y: b.Rng.Intn(b.ArenaH),
			}
			if b.Obstacles[p] {
				continue
			}
			if p.Equal(nextHead) {
				continue
			}
			occupied := false
			for _, seg := range b.Snake {
				if seg.Equal(p) {
					occupied = true
					break
				}
			}
			if occupied {
				continue
			}
			for _, existing := range b.Foods {
				if existing.Equal(p) {
					occupied = true
					break
				}
			}
			if !occupied {
				b.Foods = append(b.Foods, p)
				break
			}
		}
	}
}

// CheckNearMiss increments near-miss count if player turned away from a lethal cell directly ahead.
func (b *BaseMode) CheckNearMiss(prevDir engine.Position) {
	straight := b.Snake[0].Add(prevDir)
	isFatal := false
	if !b.Wrap && (straight.X < 0 || straight.X >= b.ArenaW || straight.Y < 0 || straight.Y >= b.ArenaH) {
		isFatal = true
	} else if b.Obstacles[straight] {
		isFatal = true
	} else {
		for _, seg := range b.Snake {
			if seg.Equal(straight) {
				isFatal = true
				break
			}
		}
	}
	if isFatal {
		b.NearMisses++
	}
}

// RespawnSnake resets snake back to center with initial length when a life is lost.
func (b *BaseMode) RespawnSnake() {
	cx := b.ArenaW / 2
	cy := b.ArenaH / 2
	b.Snake = []engine.Position{
		{X: cx, Y: cy},
		{X: cx - 1, Y: cy},
		{X: cx - 2, Y: cy},
	}
	b.Dir = engine.Position{X: 1, Y: 0}
	b.DirQueue = nil
	b.HitFlashTicks = 5
	b.EatFlashTicks = 0
}

// DrawCommonHUD draws the standardized HUD strip matching premium UI guidelines.
func (b *BaseMode) DrawCommonHUD(s *engine.Screen, originX, originY, arenaW int, modeTitle, extraHUD string) {
	th := b.Cfg.Theme
	streak := 0
	if store, err := history.Open(""); err == nil && store != nil && len(store.Records) > 0 {
		st := history.CalculateStreaks(store.Records, time.Now())
		streak = st.CurrentStreak
	}

	hudW := arenaW
	if hudW < 44 {
		hudW = 44
	}

	topBar := fmt.Sprintf("┌── SNAKE %s • %s • %s ", modeTitle, b.Cfg.Difficulty.String(), th.Name)
	for len(topBar) < hudW-1 {
		topBar += "─"
	}
	topBar += "┐"
	s.DrawText(originX, originY, topBar, th.HUD, th.Background)

	livesStr := "—"
	if b.Lives > 0 {
		livesStr = fmt.Sprintf("%d", b.Lives)
	}

	hudMid := fmt.Sprintf("│ Score: %-5d | Length: %-3d | Lives: %-2s | %s",
		b.ScoreVal, len(b.Snake), livesStr, extraHUD)

	streakStr := fmt.Sprintf("Streak: %d │", streak)
	gap := hudW - len(hudMid) - len(streakStr)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			hudMid += " "
		}
	}
	hudMid += streakStr
	if len(hudMid) > hudW {
		hudMid = hudMid[:hudW-1] + "│"
	}
	s.DrawText(originX, originY+1, hudMid, th.Text, th.Background)

	bottomBar := "└"
	for i := 1; i < hudW-1; i++ {
		bottomBar += "─"
	}
	bottomBar += "┘"
	s.DrawText(originX, originY+2, bottomBar, th.HUD, th.Background)
}

// DrawArenaBorder draws the playfield border with wrap/hit flash.
func (b *BaseMode) DrawArenaBorder(s *engine.Screen, x, y, w, h int) {
	th := b.Cfg.Theme
	wallFg := th.Wall
	if b.WrapFlashTicks > 0 {
		wallFg = th.Accent
	} else if b.HitFlashTicks > 0 {
		wallFg = tcell.ColorRed
	}

	s.Box(x, y, w, h, wallFg, th.Background)
}

// DrawSnakeElements renders food, obstacles, and snake segments with directional head.
func (b *BaseMode) DrawSnakeElements(s *engine.Screen, startX, startY int) {
	th := b.Cfg.Theme

	// 1. Obstacles
	obsGlyph := th.SnakeObstacleGlyph
	if obsGlyph == 0 {
		obsGlyph = '◆'
	}
	for obs := range b.Obstacles {
		s.DrawCell(startX+obs.X, startY+obs.Y, obsGlyph, th.Enemy, th.Background)
	}

	// 2. Food
	foodGlyph := th.SnakeFoodGlyph
	if foodGlyph == 0 {
		foodGlyph = '☆'
	}
	foodFg := th.Item
	if b.Ticks%2 == 0 {
		foodFg = th.Accent
	}
	for _, food := range b.Foods {
		s.DrawCell(startX+food.X, startY+food.Y, foodGlyph, foodFg, th.Background)
	}

	// 3. Snake Segments
	headGlyph := th.SnakeHeadForDir(b.Dir)
	bodyGlyph := th.SnakeBodyGlyph
	if bodyGlyph == 0 {
		bodyGlyph = '▪'
	}

	for i, seg := range b.Snake {
		cx := startX + seg.X
		cy := startY + seg.Y

		if i == 0 {
			headFg := th.Player
			if b.EatFlashTicks > 0 {
				headFg = tcell.ColorYellow
			} else if b.HitFlashTicks > 0 {
				headFg = tcell.ColorRed
			}
			s.DrawCell(cx, cy, headGlyph, headFg, th.Background)
		} else {
			bodyFg := th.Player
			if i%2 == 0 {
				bodyFg = th.Accent
			}
			s.DrawCell(cx, cy, bodyGlyph, bodyFg, th.Background)
		}
	}
}

// DrawGameOverOverlay renders the centered victory or game over overlay box.
func (b *BaseMode) DrawGameOverOverlay(s *engine.Screen, w, h int, modeTitle string, extraDetails string) {
	th := b.Cfg.Theme
	boxW := 52
	boxH := 11
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	title := " GAME OVER "
	titleFg := tcell.ColorRed
	if b.State == StateWin || b.Reason == "won" {
		title = " VICTORY! "
		titleFg = tcell.ColorGreen
	}

	s.DoubleBoxWithTitle(boxX, boxY, boxW, boxH, title, titleFg, th.Background, titleFg)

	var reasonText string
	switch b.Reason {
	case "wall":
		reasonText = "You collided with the boundary wall!"
	case "self":
		reasonText = "You collided with your own tail!"
	case "obstacle":
		reasonText = "You crashed into an obstacle!"
	case "enemy":
		reasonText = "An enemy caught you!"
	case "time_up":
		reasonText = "Time expired!"
	case "won":
		reasonText = "Congratulations! You completed the challenge!"
	default:
		reasonText = "Session concluded."
	}

	s.CenterText(boxY+2, reasonText, th.Text, th.Background)
	s.CenterText(boxY+4, fmt.Sprintf("Final Score: %d   Length: %d", b.ScoreVal, len(b.Snake)), th.HUD, th.Background)
	s.CenterText(boxY+5, fmt.Sprintf("Food Eaten: %d   Near Misses: %d", b.FoodEaten, b.NearMisses), th.Text, th.Background)
	if extraDetails != "" {
		s.CenterText(boxY+6, extraDetails, th.Item, th.Background)
	}

	if b.ScoreVal > b.BestScore && b.ScoreVal > 0 {
		s.CenterText(boxY+7, "★ NEW PERSONAL BEST! ★", tcell.ColorYellow, th.Background)
	}

	s.CenterText(boxY+9, "Press SPACE, ENTER or ESC to return to Menu", th.Accent, th.Background)
}
