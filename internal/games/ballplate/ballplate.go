// Package ballplate implements the Breakout-style "Ball and Plate" arcade game.
package ballplate

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/history"
)

type gameState int

const (
	stateTitle gameState = iota
	statePlaying
	stateWin
	stateOver
)

// Brick represents a destructible target block.
type Brick struct {
	X, Y   int // Top-left cell
	W, H   int // Width and height in cells
	Health int // Hits remaining (1 normal, 2 tough)
	Color  tcell.Color
	Score  int
}

// Profile defines concrete parameters per difficulty tier.
type Profile struct {
	Level             engine.Difficulty
	PlateWidth        int
	BaseBallSpeed     float64
	SpeedScaleOnClear float64
	Lives             int
	GappedBricks      bool
	ToughBricks       bool
	Interval          time.Duration
}

// DifficultyLevel implements engine.DifficultyProfile.
func (p Profile) DifficultyLevel() engine.Difficulty {
	return p.Level
}

// Profiles returns the documented tuning for Ball and Plate.
func Profiles(d engine.Difficulty) Profile {
	switch d {
	case engine.Easy:
		return Profile{
			Level:             engine.Easy,
			PlateWidth:        7,
			BaseBallSpeed:     0.32,
			SpeedScaleOnClear: 0.0,
			Lives:             4,
			GappedBricks:      true,
			ToughBricks:       false,
			Interval:          40 * time.Millisecond,
		}
	case engine.Hard:
		return Profile{
			Level:             engine.Hard,
			PlateWidth:        3,
			BaseBallSpeed:     0.48,
			SpeedScaleOnClear: 0.004,
			Lives:             2,
			GappedBricks:      false,
			ToughBricks:       true,
			Interval:          35 * time.Millisecond,
		}
	default: // Medium
		return Profile{
			Level:             engine.Medium,
			PlateWidth:        5,
			BaseBallSpeed:     0.40,
			SpeedScaleOnClear: 0.002,
			Lives:             3,
			GappedBricks:      false,
			ToughBricks:       false,
			Interval:          38 * time.Millisecond,
		}
	}
}

// BrickHitFlash represents in-cell visual feedback when a brick is destroyed.
type BrickHitFlash struct {
	X, Y, W, H int
	Text       string
	Ticks      int
}

// Game implements engine.Game for Ball and Plate.
type Game struct {
	cfg          engine.GameConfig
	profile      Profile
	state        gameState
	arenaW       int
	arenaH       int
	plateX       int
	plateY       int
	plateW       int
	ballX        float64
	ballY        float64
	ballVx       float64
	ballVy       float64
	attached     bool
	bricks       []Brick
	remBricks    int
	lives        int
	score        int
	bestScore    int
	bestInfo     history.PersonalBestInfo
	hasBest      bool
	reason       string
	bricksBroken int
	maxSpeed     int
	plateHits    int
	livesLost    int
	currentRally int
	longestRally int

	// Feedback visual timers
	plateHitFlashTicks int
	brickFlashes       []BrickHitFlash
}

// New creates an uninitialized Ball and Plate game.
func New() *Game {
	return &Game{
		arenaW: 42,
		arenaH: 20,
	}
}

// Name returns the canonical game identifier.
func (g *Game) Name() string { return "ballplate" }

// Score returns the current player points.
func (g *Game) Score() int { return g.score }

// IsOver reports whether play has concluded.
func (g *Game) IsOver() bool { return g.state == stateOver || g.state == stateWin }

// TickInterval returns the simulation tick duration.
func (g *Game) TickInterval() time.Duration { return g.profile.Interval }

// SetBestScore overrides high score for testing.
func (g *Game) SetBestScore(b int) { g.bestScore = b }

// Metrics returns the detailed session metrics satisfying engine.MetricsProvider.
func (g *Game) Metrics() map[string]int {
	return map[string]int{
		"bricks_broken":          g.bricksBroken,
		"max_ball_speed_reached": g.maxSpeed,
		"plate_hits":             g.plateHits,
		"lives_lost":             g.livesLost,
		"longest_rally":          g.longestRally,
	}
}

// Init initializes the playfield, paddle, ball physics, and brick layout.
func (g *Game) Init(cfg engine.GameConfig) {
	g.cfg = cfg
	g.profile = Profiles(cfg.Difficulty)
	g.state = stateTitle
	g.score = 0
	g.reason = ""
	g.lives = g.profile.Lives
	g.plateW = g.profile.PlateWidth
	g.plateX = (g.arenaW - g.plateW) / 2
	g.plateY = g.arenaH - 2
	g.bricksBroken = 0
	g.maxSpeed = int(g.profile.BaseBallSpeed * 100)
	g.plateHits = 0
	g.livesLost = 0
	g.currentRally = 0
	g.longestRally = 0

	if store, err := history.Open(""); err == nil && store != nil {
		g.bestScore = history.HighScore(store.Records, "ballplate", string(cfg.Difficulty))
		if pb, ok := history.PersonalBests(store.Records, "ballplate"); ok {
			g.bestInfo = pb
			g.hasBest = true
		}
	}

	g.resetBallAttached()
	g.initBricks()
}

func (g *Game) resetBallAttached() {
	g.attached = true
	g.ballX = float64(g.plateX) + float64(g.plateW)/2.0
	g.ballY = float64(g.plateY - 1)
	g.ballVx = 0.2
	g.ballVy = -g.profile.BaseBallSpeed
}

func (g *Game) launchBall() {
	if !g.attached {
		return
	}
	g.attached = false
	g.ballVy = -g.profile.BaseBallSpeed
	if g.ballVx == 0 {
		g.ballVx = 0.2
	}
}

func (g *Game) initBricks() {
	g.bricks = make([]Brick, 0)
	rows := 4
	brickW := 4
	brickH := 1
	startY := 2
	startX := 1
	cols := (g.arenaW - 2) / brickW

	th := g.cfg.Theme
	palette := th.BrickColors
	if len(palette) == 0 {
		palette = []tcell.Color{tcell.ColorRed, tcell.ColorOrange, tcell.ColorYellow, tcell.ColorGreen}
	}

	for r := 0; r < rows; r++ {
		rowColor := palette[r%len(palette)]
		for c := 0; c < cols; c++ {
			// Easy mode leaves gaps
			if g.profile.GappedBricks && (r+c)%2 == 1 {
				continue
			}

			health := 1
			// Hard mode has a tough brick row
			if g.profile.ToughBricks && r == 0 {
				health = 2
			}

			bx := startX + c*brickW
			by := startY + r*brickH
			g.bricks = append(g.bricks, Brick{
				X:      bx,
				Y:      by,
				W:      brickW,
				H:      brickH,
				Health: health,
				Color:  rowColor,
				Score:  (rows - r) * 20,
			})
		}
	}
	g.remBricks = len(g.bricks)
}

// HandleInput processes paddle movement, ball launch, and confirmation actions.
func (g *Game) HandleInput(a engine.Action) {
	switch g.state {
	case stateTitle:
		if a == engine.ActionConfirm || a == engine.ActionUp || a == engine.ActionLeft || a == engine.ActionRight {
			g.state = statePlaying
		}
	case statePlaying:
		switch a {
		case engine.ActionLeft:
			g.movePlate(-2)
		case engine.ActionRight:
			g.movePlate(2)
		case engine.ActionConfirm, engine.ActionUp:
			if g.attached {
				g.launchBall()
			}
		}
	case stateWin, stateOver:
		if a == engine.ActionConfirm || a == engine.ActionQuit {
			// Ready to exit
		}
	}
}

func (g *Game) movePlate(dx int) {
	g.plateX += dx
	if g.plateX < 1 {
		g.plateX = 1
	}
	if g.plateX+g.plateW > g.arenaW-1 {
		g.plateX = g.arenaW - 1 - g.plateW
	}
	if g.attached {
		g.ballX = float64(g.plateX) + float64(g.plateW)/2.0
	}
}

// Tick advances the ball physics and resolves collisions.
func (g *Game) Tick() engine.TickResult {
	switch g.state {
	case stateTitle:
		return engine.TickResult{Continue: true, Reason: ""}
	case stateWin:
		return engine.TickResult{Continue: false, Reason: "won"}
	case stateOver:
		return engine.TickResult{Continue: false, Reason: g.reason}
	}

	if g.attached {
		return engine.TickResult{Continue: true, Reason: ""}
	}

	// 1. Advance ball position with floating-point math
	nextX := g.ballX + g.ballVx
	nextY := g.ballY + g.ballVy

	// 2. Arena wall collisions
	if nextX <= 1.0 {
		nextX = 1.0
		g.ballVx = math.Abs(g.ballVx)
	} else if nextX >= float64(g.arenaW-2) {
		nextX = float64(g.arenaW - 2)
		g.ballVx = -math.Abs(g.ballVx)
	}

	if nextY <= 1.0 {
		nextY = 1.0
		g.ballVy = math.Abs(g.ballVy)
	}

	// 3. Plate collision and angle reflection
	plateTop := float64(g.plateY)
	if nextY >= plateTop-0.5 && nextY <= plateTop+0.8 && g.ballVy > 0 {
		plateLeft := float64(g.plateX)
		plateRight := float64(g.plateX + g.plateW)

		if nextX >= plateLeft-0.5 && nextX <= plateRight+0.5 {
			g.resolvePlateBounce(nextX)
			nextY = plateTop - 0.5
			g.plateHits++
			g.currentRally++
			if g.currentRally > g.longestRally {
				g.longestRally = g.currentRally
			}
		}
	}

	// 4. Brick collision detection
	g.ballX = nextX
	g.ballY = nextY
	g.checkBrickCollisions()

	// Update transient feedback timers
	if g.plateHitFlashTicks > 0 {
		g.plateHitFlashTicks--
	}
	if len(g.brickFlashes) > 0 {
		active := g.brickFlashes[:0]
		for _, bf := range g.brickFlashes {
			bf.Ticks--
			if bf.Ticks > 0 {
				active = append(active, bf)
			}
		}
		g.brickFlashes = active
	}

	// 5. Level win check
	if g.remBricks <= 0 {
		g.state = stateWin
		g.reason = "won"
		return engine.TickResult{Continue: false, Reason: "won"}
	}

	// 6. Bottom boundary: life lost
	if g.ballY >= float64(g.arenaH-1) {
		g.lives--
		g.livesLost++
		g.currentRally = 0
		if g.lives <= 0 {
			g.state = stateOver
			g.reason = "lost"
			return engine.TickResult{Continue: false, Reason: "lost"}
		}
		g.resetBallAttached()
	}

	return engine.TickResult{Continue: true, Reason: ""}
}

// ResolvePlateBounce computes variable reflection angle based on paddle hit point.
// Exported for headless testing.
func (g *Game) ResolvePlateBounce(hitX float64) {
	g.resolvePlateBounce(hitX)
}

func (g *Game) resolvePlateBounce(hitX float64) {
	g.plateHitFlashTicks = 3 // Briefly invert plate colors for 50ms
	plateCenter := float64(g.plateX) + float64(g.plateW)/2.0
	halfWidth := float64(g.plateW) / 2.0
	offset := (hitX - plateCenter) / halfWidth

	// Clamp offset between -1.0 and 1.0
	if offset < -1.0 {
		offset = -1.0
	} else if offset > 1.0 {
		offset = 1.0
	}

	// 3-zone deflection physics:
	// Left zone (< -0.33): -60° to -20°
	// Center zone (-0.33 to +0.33): -20° to +20°
	// Right zone (> 0.33): +20° to +60°
	var angle float64
	degToRad := math.Pi / 180.0
	if offset < -0.33 {
		norm := (offset - (-0.33)) / (-1.0 - (-0.33))
		deg := -20.0 - norm*40.0
		angle = deg * degToRad
	} else if offset > 0.33 {
		norm := (offset - 0.33) / (1.0 - 0.33)
		deg := 20.0 + norm*40.0
		angle = deg * degToRad
	} else {
		norm := offset / 0.33
		deg := norm * 20.0
		angle = deg * degToRad
	}

	// Calculate current speed with difficulty-scaling
	clearedCount := len(g.bricks) - g.remBricks
	speed := g.profile.BaseBallSpeed + float64(clearedCount)*g.profile.SpeedScaleOnClear
	if speed < 0.25 {
		speed = 0.25
	} else if speed > 0.85 {
		speed = 0.85
	}
	speedInt := int(speed * 100)
	if speedInt > g.maxSpeed {
		g.maxSpeed = speedInt
	}

	g.ballVx = speed * math.Sin(angle)
	g.ballVy = -math.Abs(speed * math.Cos(angle))
	engine.Beep()
}

func (g *Game) checkBrickCollisions() {
	cellX := int(math.Round(g.ballX))
	cellY := int(math.Round(g.ballY))

	for i := range g.bricks {
		b := &g.bricks[i]
		if b.Health <= 0 {
			continue
		}

		if cellX >= b.X && cellX < b.X+b.W && cellY >= b.Y && cellY < b.Y+b.H {
			b.Health--
			if b.Health <= 0 {
				g.remBricks--
				g.score += b.Score
				g.bricksBroken++
				g.brickFlashes = append(g.brickFlashes, BrickHitFlash{
					X:     b.X,
					Y:     b.Y,
					W:     b.W,
					H:     b.H,
					Text:  fmt.Sprintf("+%d", b.Score),
					Ticks: 4,
				})
			}
			engine.Beep()
			// Reverse vertical direction
			g.ballVy = -g.ballVy
			break
		}
	}
}


// Render draws playfield, HUD, bricks, plate, ball, and overlays.
func (g *Game) Render(s *engine.Screen) {
	th := g.cfg.Theme
	s.Clear(th.Background)
	w, h := s.Size()

	switch g.state {
	case stateTitle:
		g.renderTitle(s, w, h)
	case stateWin:
		g.renderField(s, w, h)
		g.renderOutcomeOverlay(s, w, h, "LEVEL CLEARED! VICTORY", tcell.ColorLime)
	case stateOver:
		g.renderField(s, w, h)
		g.renderOutcomeOverlay(s, w, h, "GAME OVER - OUT OF LIVES", tcell.ColorRed)
	default:
		g.renderField(s, w, h)
	}
}

func (g *Game) renderTitle(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	banner := []string{
		`  ____    _    _     _        _   ____  _        _  _____ _____ `,
		` | __ )  / \  | |   | |      / \ |  _ \| |      / \|_   _| ____|`,
		` |  _ \ / _ \ | |   | |     / _ \| |_) | |     / _ \ | | |  _|  `,
		` | |_) / ___ \| |___| |___ / ___ \  __/| |___ / ___ \| | | |___ `,
		` |____/_/   \_\_____|_____/_/   \_\_|   |_____/_/   \_\_| |_____|`,
	}

	startY := (h - 16) / 2
	if startY < 1 {
		startY = 1
	}

	for i, line := range banner {
		s.CenterText(startY+i, line, th.Item, th.Background)
	}

	s.CenterText(startY+6, "Arcade Breakout with Angular Physics", th.HUD, th.Background)
	s.CenterText(startY+8, "Difficulty: "+g.profile.Level.String()+"  |  Theme: "+th.Name, th.Text, th.Background)

	if g.hasBest {
		bestText := fmt.Sprintf("Your best: %d pts (%s)", g.bestInfo.HighScore, history.FormatRelativeTime(g.bestInfo.HighScoreDate, time.Now()))
		s.CenterText(startY+10, bestText, th.Item, th.Background)
	}

	s.CenterText(startY+12, "Controls: Arrow Keys / A/D  |  Space: Launch  |  P: Pause  |  Q: Quit", th.Text, th.Background)
	s.CenterText(startY+14, "Press SPACE or ENTER to Start", th.Accent, th.Background)
}

func (g *Game) renderField(s *engine.Screen, w, h int) {
	th := g.cfg.Theme
	arenaW := g.arenaW
	arenaH := g.arenaH
	hudH := 3
	totalH := arenaH + hudH

	if w < arenaW+4 || h < totalH+2 {
		s.CenterText(h/2, "Terminal size too small! Please expand window.", tcell.ColorYellow, th.Background)
		return
	}

	startX := (w - arenaW) / 2
	startY := (h - totalH) / 2
	if startY < 1 {
		startY = 1
	}

	// 1. Render HUD with live speed indicator
	s.Box(startX, startY, arenaW, hudH, th.HUD, th.Background)
	livesIcon := ""
	for i := 0; i < g.lives; i++ {
		livesIcon += "● "
	}
	clearedCount := len(g.bricks) - g.remBricks
	currentSpeed := g.profile.BaseBallSpeed + float64(clearedCount)*g.profile.SpeedScaleOnClear
	speedMultiplier := currentSpeed / g.profile.BaseBallSpeed

	speedGauge := "▂"
	if speedMultiplier >= 1.5 {
		speedGauge = "▂▃▄▅"
	} else if speedMultiplier >= 1.3 {
		speedGauge = "▂▃▄"
	} else if speedMultiplier >= 1.15 {
		speedGauge = "▂▃"
	}

	hudText := fmt.Sprintf("Score: %-5d Lives: %s Bricks: %-2d Speed: %s (%.1fx) [%s]",
		g.score, livesIcon, g.remBricks, speedGauge, speedMultiplier, g.profile.Level.String())
	s.DrawText(startX+2, startY+1, hudText, th.HUD, th.Background)

	// 2. Render Playfield Border
	fieldY := startY + hudH
	s.Box(startX, fieldY, arenaW, arenaH, th.Wall, th.Background)

	// 3. Render Bricks
	for _, b := range g.bricks {
		if b.Health <= 0 {
			continue
		}
		glyph := th.BrickGlyph
		color := b.Color
		if b.Health > 1 {
			glyph = th.ToughBrick
		}
		for bx := 0; bx < b.W; bx++ {
			s.DrawCell(startX+b.X+bx, fieldY+b.Y, glyph, color, th.Background)
		}
	}

	// In-cell brick destruction flashes and score popups
	for _, bf := range g.brickFlashes {
		for bx := 0; bx < bf.W; bx++ {
			s.DrawCell(startX+bf.X+bx, fieldY+bf.Y, '█', tcell.ColorWhite, th.Background)
		}
		s.DrawText(startX+bf.X+1, fieldY+bf.Y, bf.Text, tcell.ColorYellow, tcell.ColorWhite)
	}

	// 4. Render Plate (3 zones visually distinguished with hit inversion)
	plateColor := th.Player
	plateBg := th.Background
	if g.plateHitFlashTicks > 0 {
		plateColor = tcell.ColorWhite
		plateBg = th.Accent
	} else if th.Name == "Monochrome" {
		plateColor = tcell.ColorWhite
	}

	for px := 0; px < g.plateW; px++ {
		glyph := th.PlateGlyph
		cellColor := plateColor
		if px == 0 {
			glyph = '◄'
			if g.plateHitFlashTicks == 0 {
				cellColor = th.Accent
			}
		} else if px == g.plateW-1 {
			glyph = '►'
			if g.plateHitFlashTicks == 0 {
				cellColor = th.Accent
			}
		}
		s.DrawCell(startX+g.plateX+px, fieldY+g.plateY, glyph, cellColor, plateBg)
	}


	// 5. Render Ball (Sub-cell rounded interpolation)
	renderBallX := int(math.Round(g.ballX))
	renderBallY := int(math.Round(g.ballY))

	ballGlyph := th.BallGlyph
	ballColor := th.Item
	if th.Name == "Retro Green" {
		ballColor = tcell.ColorLime
	} else if th.Name == "Monochrome" {
		ballColor = tcell.ColorWhite
	}

	s.DrawCell(startX+renderBallX, fieldY+renderBallY, ballGlyph, ballColor, th.Background)

	if g.attached {
		s.CenterText(fieldY+g.plateY-2, "Press SPACE to Launch", th.Accent, th.Background)
	}
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
