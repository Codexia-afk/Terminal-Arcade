// Package snake implements the classic Nokia-style grid Snake game.
package snake

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
	"goarcade/internal/history"
)

type gameState int

const (
	stateTitle gameState = iota
	statePlaying
	stateOver
)

// Profile contains concrete per-difficulty mechanical parameters.
type Profile struct {
	Level     engine.Difficulty
	Interval  time.Duration
	Width     int
	Height    int
	Wrap      bool
	Obstacles bool
}

// DifficultyLevel implements engine.DifficultyProfile.
func (p Profile) DifficultyLevel() engine.Difficulty {
	return p.Level
}

// Profiles returns the documented tuning parameters for the requested difficulty.
func Profiles(d engine.Difficulty) Profile {
	switch d {
	case engine.Easy:
		return Profile{
			Level:     engine.Easy,
			Interval:  180 * time.Millisecond,
			Width:     40,
			Height:    20,
			Wrap:      true,
			Obstacles: false,
		}
	case engine.Hard:
		return Profile{
			Level:     engine.Hard,
			Interval:  90 * time.Millisecond,
			Width:     28,
			Height:    16,
			Wrap:      false,
			Obstacles: true,
		}
	default: // Medium
		return Profile{
			Level:     engine.Medium,
			Interval:  130 * time.Millisecond,
			Width:     34,
			Height:    18,
			Wrap:      false,
			Obstacles: false,
		}
	}
}

// Game implements engine.Game for Snake.
type Game struct {
	cfg       engine.GameConfig
	profile   Profile
	state     gameState
	snake     []engine.Position
	dir       engine.Position
	queuedDir engine.Position
	food      engine.Position
	obstacles map[engine.Position]bool
	score     int
	bestScore int
	reason    string
	ticks     int
	rng       *rand.Rand
}

// New creates an uninitialized Snake game.
func New() *Game {
	return &Game{
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		obstacles: make(map[engine.Position]bool),
	}
}

// Name returns the canonical identifier for this game.
func (g *Game) Name() string { return "snake" }

// Score returns the current player score.
func (g *Game) Score() int { return g.score }

// IsOver reports whether the game session has concluded.
func (g *Game) IsOver() bool { return g.state == stateOver }

// TickInterval returns the ticker duration specified by the current difficulty profile.
func (g *Game) TickInterval() time.Duration { return g.profile.Interval }

// SetBestScore manually overrides the personal best score (primarily for headless testing).
func (g *Game) SetBestScore(best int) { g.bestScore = best }

// Init initializes the game using the supplied configuration.
func (g *Game) Init(cfg engine.GameConfig) {
	g.cfg = cfg
	g.profile = Profiles(cfg.Difficulty)
	g.state = stateTitle
	g.score = 0
	g.reason = ""
	g.ticks = 0
	g.obstacles = make(map[engine.Position]bool)

	// Fetch personal best score from history if available
	if store, err := history.Open(""); err == nil && store != nil {
		g.bestScore = history.HighScore(store.Records, "snake", string(cfg.Difficulty))
	}

	cx := g.profile.Width / 2
	cy := g.profile.Height / 2
	g.snake = []engine.Position{
		{X: cx, Y: cy},
		{X: cx - 1, Y: cy},
		{X: cx - 2, Y: cy},
	}
	g.dir = engine.Position{X: 1, Y: 0}
	g.queuedDir = g.dir

	g.spawnFood()
}

func (g *Game) spawnFood() {
	for attempts := 0; attempts < 500; attempts++ {
		p := engine.Position{
			X: g.rng.Intn(g.profile.Width),
			Y: g.rng.Intn(g.profile.Height),
		}
		if g.obstacles[p] {
			continue
		}
		occupied := false
		for _, seg := range g.snake {
			if seg.Equal(p) {
				occupied = true
				break
			}
		}
		if !occupied {
			g.food = p
			return
		}
	}
}

func (g *Game) spawnObstacle() {
	for attempts := 0; attempts < 100; attempts++ {
		p := engine.Position{
			X: g.rng.Intn(g.profile.Width),
			Y: g.rng.Intn(g.profile.Height),
		}
		if p.Equal(g.food) || g.obstacles[p] {
			continue
		}
		occupied := false
		for _, seg := range g.snake {
			if seg.Equal(p) {
				occupied = true
				break
			}
		}
		if !occupied {
			g.obstacles[p] = true
			return
		}
	}
}

// HandleInput handles user input buffering and screen transitions.
func (g *Game) HandleInput(a engine.Action) {
	switch g.state {
	case stateTitle:
		if a == engine.ActionConfirm || a == engine.ActionUp || a == engine.ActionDown || a == engine.ActionLeft || a == engine.ActionRight {
			g.state = statePlaying
		}
	case statePlaying:
		d, ok := a.Direction()
		if ok {
			// Buffer direction change; prevent 180° instant reversal
			if !(d.X == -g.dir.X && d.Y == -g.dir.Y) {
				g.queuedDir = d
			}
		}
	case stateOver:
		if a == engine.ActionConfirm || a == engine.ActionQuit {
			// Ready to exit back to menu
		}
	}
}

// Tick executes one discrete simulation step.
func (g *Game) Tick() engine.TickResult {
	switch g.state {
	case stateTitle:
		return engine.TickResult{Continue: true, Reason: ""}
	case stateOver:
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	g.dir = g.queuedDir
	next := g.snake[0].Add(g.dir)

	// Wall collision / Wrapping
	if g.profile.Wrap {
		next.X = (next.X + g.profile.Width) % g.profile.Width
		next.Y = (next.Y + g.profile.Height) % g.profile.Height
	} else if next.X < 0 || next.X >= g.profile.Width || next.Y < 0 || next.Y >= g.profile.Height {
		return g.terminate("wall")
	}

	// Obstacle collision
	if g.obstacles[next] {
		return g.terminate("obstacle")
	}

	// Self collision
	for _, seg := range g.snake {
		if seg.Equal(next) {
			return g.terminate("self")
		}
	}

	// Move head forward
	g.snake = append([]engine.Position{next}, g.snake...)

	// Food consumption
	if next.Equal(g.food) {
		g.score += 10
		g.spawnFood()
	} else {
		// Remove tail segment if food was not eaten
		g.snake = g.snake[:len(g.snake)-1]
	}

	g.ticks++
	if g.profile.Obstacles && g.ticks%30 == 0 && len(g.obstacles) < 8 {
		g.spawnObstacle()
	}

	return engine.TickResult{Continue: true, Reason: ""}
}

func (g *Game) terminate(reason string) engine.TickResult {
	g.state = stateOver
	g.reason = reason
	return engine.TickResult{Continue: false, Reason: reason}
}

// Render renders the game, HUD, title, or game over overlay using engine.Screen.
func (g *Game) Render(s *engine.Screen) {
	th := g.cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	switch g.state {
	case stateTitle:
		g.renderTitle(s, w, h)
		return
	case stateOver:
		g.renderPlaying(s, w, h)
		g.renderGameOver(s, w, h)
		return
	default:
		g.renderPlaying(s, w, h)
	}
}

func (g *Game) renderTitle(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	banner := []string{
		`  ____  _   _    _    _  _______ `,
		` / ___|| \ | |  / \  | |/ / ____|`,
		` \___ \|  \| | / _ \ | ' /|  _|  `,
		`  ___) | |\  |/ ___ \| . \| |___ `,
		` |____/|_| \_/_/   \_\_|\_\_____|`,
	}

	startY := (h - 14) / 2
	if startY < 2 {
		startY = 2
	}

	for i, line := range banner {
		s.CenterText(startY+i, line, th.Accent, th.Background)
	}

	s.CenterText(startY+7, "Classic Nokia Arcade Snake", th.HUD, th.Background)
	s.CenterText(startY+9, "Mode: "+g.profile.Level.String()+"  |  Theme: "+th.Name, th.Text, th.Background)
	s.CenterText(startY+11, "Controls: Arrow Keys / WASD  |  P: Pause  |  Q: Quit", th.Text, th.Background)
	s.CenterText(startY+13, "Press SPACE or ENTER to Start", th.Accent, th.Background)
}

func (g *Game) renderPlaying(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	arenaW := g.profile.Width + 2
	arenaH := g.profile.Height + 2
	hudH := 3
	totalH := arenaH + hudH + 1

	if w < arenaW+2 || h < totalH+2 {
		s.CenterText(h/2, "Terminal size too small! Please resize.", tcell.ColorYellow, th.Background)
		return
	}

	startX := (w - arenaW) / 2
	startY := (h - totalH) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Render HUD box
	s.Box(startX, startY, arenaW, hudH, th.HUD, th.Background)
	hudText := "Score: " + strconv.Itoa(g.score) + "   Len: " + strconv.Itoa(len(g.snake)) + "   [" + g.profile.Level.String() + "]"
	s.DrawText(startX+2, startY+1, hudText, th.HUD, th.Background)

	// 2. Render Arena Border
	arenaY := startY + hudH
	s.Box(startX, arenaY, arenaW, arenaH, th.Wall, th.Background)

	// 3. Render Food
	foodGlyph := th.ItemGlyph
	foodColor := th.Item
	if th.Name == "Neon" {
		foodGlyph = '●'
		foodColor = tcell.ColorYellow
	}
	s.DrawCell(startX+1+g.food.X, arenaY+1+g.food.Y, foodGlyph, foodColor, th.Background)

	// 4. Render Obstacles
	for obs := range g.obstacles {
		s.DrawCell(startX+1+obs.X, arenaY+1+obs.Y, '×', th.Enemy, th.Background)
	}

	// 5. Render Snake
	for i, seg := range g.snake {
		var glyph rune
		var color tcell.Color

		switch th.Name {
		case "Monochrome":
			if i == 0 {
				glyph = 'o'
			} else {
				glyph = '#'
			}
			color = tcell.ColorWhite
		case "Neon":
			glyph = '█'
			if i == 0 {
				color = tcell.ColorAqua
			} else {
				color = tcell.ColorFuchsia
			}
		default: // Retro Green
			glyph = '█'
			if i == 0 {
				color = tcell.ColorLime
			} else {
				color = tcell.ColorGreen
			}
		}

		s.DrawCell(startX+1+seg.X, arenaY+1+seg.Y, glyph, color, th.Background)
	}
}

func (g *Game) renderGameOver(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	boxW := 40
	boxH := 9
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	s.Fill(boxX, boxY, boxW, boxH, ' ', th.Text, tcell.ColorBlack)
	s.BoxWithTitle(boxX, boxY, boxW, boxH, "GAME OVER", th.Accent, tcell.ColorBlack, th.Accent)

	reasonText := "Outcome: Collision (" + g.reason + ")"
	s.CenterText(boxY+2, reasonText, th.Text, tcell.ColorBlack)
	s.CenterText(boxY+3, "Final Score: "+strconv.Itoa(g.score), th.Item, tcell.ColorBlack)

	if g.score > 0 && g.score > g.bestScore {
		s.CenterText(boxY+5, "★ NEW PERSONAL BEST! ★", tcell.ColorYellow, tcell.ColorBlack)
	} else if g.bestScore > 0 {
		s.CenterText(boxY+5, "Best Score: "+strconv.Itoa(g.bestScore), th.HUD, tcell.ColorBlack)
	}

	s.CenterText(boxY+7, "Press Q or ENTER for Menu", th.Text, tcell.ColorBlack)
}
