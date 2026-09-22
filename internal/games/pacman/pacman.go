package pacman

import (
	"fmt"
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
	stateWin
	stateOver
)

// Profile defines Pacman's concrete difficulty parameters.
type Profile struct {
	Level          engine.Difficulty
	GhostCount     int
	PelletDuration int // ticks
	Lives          int
	Interval       time.Duration
}

// DifficultyLevel implements engine.DifficultyProfile.
func (p Profile) DifficultyLevel() engine.Difficulty {
	return p.Level
}

// Profiles returns the difficulty tuning for Pacman.
func Profiles(d engine.Difficulty) Profile {
	switch d {
	case engine.Easy:
		return Profile{
			Level:          engine.Easy,
			GhostCount:     2,
			PelletDuration: 50,
			Lives:          3,
			Interval:       140 * time.Millisecond,
		}
	case engine.Hard:
		return Profile{
			Level:          engine.Hard,
			GhostCount:     4,
			PelletDuration: 20,
			Lives:          2,
			Interval:       110 * time.Millisecond,
		}
	default: // Medium
		return Profile{
			Level:          engine.Medium,
			GhostCount:     3,
			PelletDuration: 35,
			Lives:          3,
			Interval:       125 * time.Millisecond,
		}
	}
}

// Game implements engine.Game for Pacman.
type Game struct {
	cfg             engine.GameConfig
	profile         Profile
	state           gameState
	maze            *Maze
	playerPos       engine.Position
	playerDir       engine.Position
	queuedDir       engine.Position
	ghosts          []*Ghost
	lives           int
	score           int
	bestScore       int
	bestInfo        history.PersonalBestInfo
	hasBest         bool
	vulnerableTicks int
	reason          string
	dotsEaten       int
	pelletsUsed     int
	ghostsEaten     int
	levelsCleared   int
	livesLost       int
	rng             *rand.Rand
}

// New constructs an uninitialized Pacman game.
func New() *Game {
	return &Game{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Name returns the canonical identifier.
func (g *Game) Name() string { return "pacman" }

// Score returns current points.
func (g *Game) Score() int { return g.score }

// IsOver reports whether the game session has concluded.
func (g *Game) IsOver() bool { return g.state == stateOver || g.state == stateWin }

// TickInterval returns the simulation tick duration.
func (g *Game) TickInterval() time.Duration { return g.profile.Interval }

// SetBestScore overrides high score for testing.
func (g *Game) SetBestScore(b int) { g.bestScore = b }

// Metrics returns the detailed session metrics satisfying engine.MetricsProvider.
func (g *Game) Metrics() map[string]int {
	return map[string]int{
		"dots_eaten":         g.dotsEaten,
		"power_pellets_used": g.pelletsUsed,
		"ghosts_eaten":       g.ghostsEaten,
		"levels_cleared":     g.levelsCleared,
		"lives_lost":         g.livesLost,
	}
}

// Init initializes the maze, ghost AI, and game state according to config.
func (g *Game) Init(cfg engine.GameConfig) {
	g.cfg = cfg
	g.profile = Profiles(cfg.Difficulty)
	g.state = stateTitle
	g.score = 0
	g.reason = ""
	g.lives = g.profile.Lives
	g.vulnerableTicks = 0
	g.dotsEaten = 0
	g.pelletsUsed = 0
	g.ghostsEaten = 0
	g.levelsCleared = 0
	g.livesLost = 0

	// Retrieve personal best
	if store, err := history.Open(""); err == nil && store != nil {
		g.bestScore = history.HighScore(store.Records, "pacman", string(cfg.Difficulty))
		if pb, ok := history.PersonalBests(store.Records, "pacman"); ok {
			g.bestInfo = pb
			g.hasBest = true
		}
	}

	g.maze = NewMaze()
	g.playerPos = g.maze.PlayerStart
	g.playerDir = engine.Position{X: -1, Y: 0}
	g.queuedDir = g.playerDir

	g.initGhosts()
}

func (g *Game) initGhosts() {
	ghostColors := []tcell.Color{
		tcell.ColorRed,        // Blinky
		tcell.ColorFuchsia,    // Pinky
		tcell.ColorAqua,       // Inky
		tcell.ColorDarkOrange, // Clyde
	}

	g.ghosts = make([]*Ghost, 0, g.profile.GhostCount)

	configs := []struct {
		name        string
		personality string
	}{
		{"Blinky", "chase"},
		{"Pinky", "random"},
		{"Inky", "chase"},
		{"Clyde", "ambush"},
	}

	if g.profile.Level == engine.Easy {
		configs[0].personality = "random"
		configs[1].personality = "random"
	} else if g.profile.Level == engine.Hard {
		configs[1].personality = "ambush"
	}

	for i := 0; i < g.profile.GhostCount; i++ {
		home := g.maze.GhostStarts[i%len(g.maze.GhostStarts)]
		gh := NewGhost(i+1, configs[i].name, home, ghostColors[i%len(ghostColors)], configs[i].personality)
		g.ghosts = append(g.ghosts, gh)
	}
}

// HandleInput processes keyboard direction changes and screen confirmations.
func (g *Game) HandleInput(a engine.Action) {
	switch g.state {
	case stateTitle:
		if a == engine.ActionConfirm || a == engine.ActionUp || a == engine.ActionDown || a == engine.ActionLeft || a == engine.ActionRight {
			g.state = statePlaying
		}
	case statePlaying:
		if d, ok := a.Direction(); ok {
			g.queuedDir = d
		}
	case stateWin, stateOver:
		if a == engine.ActionConfirm || a == engine.ActionQuit {
			// Ready to exit back to menu
		}
	}
}

// Tick steps the maze simulation forward.
func (g *Game) Tick() engine.TickResult {
	switch g.state {
	case stateTitle:
		return engine.TickResult{Continue: true, Reason: ""}
	case stateWin:
		return engine.TickResult{Continue: false, Reason: "won"}
	case stateOver:
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	prevPlayerPos := g.playerPos

	// 1. Move Player
	candNext := g.maze.WrapPosition(g.playerPos.Add(g.queuedDir))
	if !g.maze.IsWall(candNext) {
		g.playerDir = g.queuedDir
		g.playerPos = candNext
	} else {
		curNext := g.maze.WrapPosition(g.playerPos.Add(g.playerDir))
		if !g.maze.IsWall(curNext) {
			g.playerPos = curNext
		}
	}

	// Immediate collision check after player moves
	if g.resolveCollisions(prevPlayerPos, nil) {
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	// 2. Tile interactions (Dot / Power Pellet)
	tile := g.maze.Tiles[g.playerPos.Y][g.playerPos.X]
	if tile == TileDot {
		g.maze.Tiles[g.playerPos.Y][g.playerPos.X] = TileEmpty
		g.maze.Remaining--
		g.score += 10
		g.dotsEaten++
	} else if tile == TilePellet {
		g.maze.Tiles[g.playerPos.Y][g.playerPos.X] = TileEmpty
		g.maze.Remaining--
		g.score += 50
		g.pelletsUsed++
		g.vulnerableTicks = g.profile.PelletDuration
		for _, gh := range g.ghosts {
			if !gh.Eaten {
				gh.Vulnerable = true
			}
		}
	}

	// 3. Check for Level Win
	if g.maze.Remaining <= 0 {
		g.state = stateWin
		g.reason = "won"
		g.levelsCleared++
		return engine.TickResult{Continue: false, Reason: "won"}
	}

	// 4. Update Vulnerable Timer
	if g.vulnerableTicks > 0 {
		g.vulnerableTicks--
		if g.vulnerableTicks == 0 {
			for _, gh := range g.ghosts {
				gh.Vulnerable = false
			}
		}
	}

	// 5. Move Ghosts
	prevGhostPositions := make([]engine.Position, len(g.ghosts))
	for i, gh := range g.ghosts {
		prevGhostPositions[i] = gh.Pos
		dir := gh.DecideNextMove(g.maze, g.playerPos, g.playerDir, g.rng)
		gh.Dir = dir
		gh.Pos = g.maze.WrapPosition(gh.Pos.Add(dir))

		if gh.Eaten && gh.Pos.Equal(gh.Home) {
			gh.Eaten = false
			gh.Vulnerable = false
		}
	}

	// 6. Collision Check after ghosts move (including crossover check)
	if g.resolveCollisions(prevPlayerPos, prevGhostPositions) {
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	return engine.TickResult{Continue: true, Reason: ""}
}

func (g *Game) resolveCollisions(prevPlayerPos engine.Position, prevGhostPositions []engine.Position) bool {
	for i, gh := range g.ghosts {
		directHit := gh.Pos.Equal(g.playerPos)
		crossHit := false
		if prevGhostPositions != nil && i < len(prevGhostPositions) {
			// Crossed paths during this tick
			if prevGhostPositions[i].Equal(g.playerPos) && gh.Pos.Equal(prevPlayerPos) {
				crossHit = true
			}
		}

		if directHit || crossHit {
			if gh.Vulnerable && !gh.Eaten {
				g.score += 200
				gh.Eaten = true
				gh.Vulnerable = false
				g.ghostsEaten++
			} else if !gh.Eaten {
				g.lives--
				g.livesLost++
				if g.lives <= 0 {
					g.state = stateOver
					g.reason = "died"
					return true
				}
				// Reset to starting positions for next life
				g.playerPos = g.maze.PlayerStart
				g.playerDir = engine.Position{X: -1, Y: 0}
				g.queuedDir = g.playerDir
				for _, gReset := range g.ghosts {
					gReset.Pos = gReset.Home
					gReset.Vulnerable = false
					gReset.Eaten = false
				}
				break
			}
		}
	}
	return false
}

// Render presents the maze, HUD, and game state.
func (g *Game) Render(s *engine.Screen) {
	th := g.cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	switch g.state {
	case stateTitle:
		g.renderTitle(s, w, h)
	case stateWin:
		g.renderMaze(s, w, h)
		g.renderOutcomeOverlay(s, w, h, "VICTORY! ALL DOTS CLEARED", tcell.ColorLime)
	case stateOver:
		g.renderMaze(s, w, h)
		g.renderOutcomeOverlay(s, w, h, "GAME OVER - OUT OF LIVES", tcell.ColorRed)
	default:
		g.renderMaze(s, w, h)
	}
}

func (g *Game) renderTitle(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	banner := []string{
		`  ____   _    ____ __  __    _    _   _ `,
		` |  _ \ / \  / ___|  \/  |  / \  | \ | |`,
		` | |_) / _ \| |   | |\/| | / _ \ |  \| |`,
		` |  __/ ___ \ |___| |  | |/ ___ \| |\  |`,
		` |_| /_/   \_\____|_|  |_/_/   \_\_| \_|`,
	}

	startY := (h - 16) / 2
	if startY < 1 {
		startY = 1
	}

	for i, line := range banner {
		s.CenterText(startY+i, line, th.Player, th.Background)
	}

	s.CenterText(startY+6, "Arcade Maze Chase", th.HUD, th.Background)
	s.CenterText(startY+8, "Difficulty: "+g.profile.Level.String()+"  |  Theme: "+th.Name, th.Text, th.Background)

	if g.hasBest {
		bestText := fmt.Sprintf("Your best: %d pts (%s)", g.bestInfo.HighScore, history.FormatRelativeTime(g.bestInfo.HighScoreDate, time.Now()))
		s.CenterText(startY+10, bestText, th.Item, th.Background)
	}

	s.CenterText(startY+12, "Controls: Arrow Keys / WASD  |  P: Pause  |  Q: Quit", th.Text, th.Background)
	s.CenterText(startY+14, "Press SPACE or ENTER to Start", th.Accent, th.Background)
}

func (g *Game) renderMaze(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	mazeW := g.maze.Width
	mazeH := g.maze.Height
	hudH := 3
	totalH := mazeH + hudH

	if w < mazeW+4 || h < totalH+2 {
		s.CenterText(h/2, "Terminal size too small! Please expand window.", tcell.ColorYellow, th.Background)
		return
	}

	startX := (w - mazeW) / 2
	startY := (h - totalH) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Render HUD
	s.Box(startX, startY, mazeW, hudH, th.HUD, th.Background)
	livesIcon := ""
	for i := 0; i < g.lives; i++ {
		livesIcon += "C "
	}
	hudText := "Score: " + strconv.Itoa(g.score) + "   Lives: " + livesIcon + " [" + g.profile.Level.String() + "]"
	s.DrawText(startX+2, startY+1, hudText, th.HUD, th.Background)

	// 2. Render Maze Tiles
	boardY := startY + hudH
	for y := 0; y < mazeH; y++ {
		for x := 0; x < mazeW; x++ {
			t := g.maze.Tiles[y][x]
			posX := startX + x
			posY := boardY + y

			switch t {
			case TileWall:
				s.DrawCell(posX, posY, th.WallGlyph, th.Wall, th.Background)
			case TileGate:
				s.DrawCell(posX, posY, '-', th.HUD, th.Background)
			case TileDot:
				s.DrawCell(posX, posY, th.ItemGlyph, th.Item, th.Background)
			case TilePellet:
				s.DrawCell(posX, posY, th.PelletGlyph, th.Accent, th.Background)
			default:
				s.DrawCell(posX, posY, ' ', th.Text, th.Background)
			}
		}
	}

	// 3. Render Ghosts
	for _, gh := range g.ghosts {
		glyph := th.EnemyGlyph
		color := gh.Color

		if th.Name == "Retro Green" {
			color = th.Enemy
		} else if th.Name == "Monochrome" {
			color = tcell.ColorWhite
		}

		if gh.Eaten {
			glyph = '"'
			color = tcell.ColorGray
		} else if gh.Vulnerable {
			glyph = 'w'
			color = th.Vulnerable
			if g.vulnerableTicks < 10 && (g.vulnerableTicks%2 == 0) {
				color = tcell.ColorWhite
			}
		}

		s.DrawCell(startX+gh.Pos.X, boardY+gh.Pos.Y, glyph, color, th.Background)
	}

	// 4. Render Pacman
	playerGlyph := th.PlayerGlyph
	switch th.Name {
	case "Monochrome":
		playerGlyph = 'C'
	case "Retro Green":
		playerGlyph = '█'
	default:
		if g.playerDir.X == 1 {
			playerGlyph = '>'
		} else if g.playerDir.X == -1 {
			playerGlyph = '<'
		} else if g.playerDir.Y == 1 {
			playerGlyph = 'v'
		} else if g.playerDir.Y == -1 {
			playerGlyph = '^'
		}
	}
	s.DrawCell(startX+g.playerPos.X, boardY+g.playerPos.Y, playerGlyph, th.Player, th.Background)
}

func (g *Game) renderOutcomeOverlay(s *engine.Screen, w, h int, title string, titleColor tcell.Color) {
	th := g.cfg.Theme
	boxW := 42
	boxH := 9
	boxX := (w - boxW) / 2
	boxY := (h - boxH) / 2

	s.Fill(boxX, boxY, boxW, boxH, ' ', th.Text, tcell.ColorBlack)
	s.BoxWithTitle(boxX, boxY, boxW, boxH, title, titleColor, tcell.ColorBlack, titleColor)

	s.CenterText(boxY+2, "Final Score: "+strconv.Itoa(g.score), th.Item, tcell.ColorBlack)

	if g.score > 0 && g.score > g.bestScore {
		s.CenterText(boxY+4, "★ NEW PERSONAL BEST! ★", tcell.ColorYellow, tcell.ColorBlack)
	} else if g.bestScore > 0 {
		s.CenterText(boxY+4, "High Score: "+strconv.Itoa(g.bestScore), th.HUD, tcell.ColorBlack)
	}

	s.CenterText(boxY+6, "Press Q or ENTER for Menu", th.Text, tcell.ColorBlack)
}
