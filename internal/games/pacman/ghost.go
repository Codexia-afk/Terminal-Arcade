package pacman

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"goarcade/internal/engine"
)

// Ghost represents an enemy ghost in the maze.
type Ghost struct {
	ID          int
	Name        string
	Home        engine.Position
	Pos         engine.Position
	Dir         engine.Position
	Color       tcell.Color
	Personality string // "random", "chase", "ambush"
	Vulnerable  bool
	Eaten       bool
}

// NewGhost initializes a ghost with a distinct color and AI archetype.
func NewGhost(id int, name string, home engine.Position, color tcell.Color, personality string) *Ghost {
	return &Ghost{
		ID:          id,
		Name:        name,
		Home:        home,
		Pos:         home,
		Dir:         engine.Position{X: 0, Y: -1},
		Color:       color,
		Personality: personality,
	}
}

// DecideNextMove calculates the next grid step based on the ghost's AI personality.
func (g *Ghost) DecideNextMove(m *Maze, pacmanPos, pacmanDir engine.Position, rng *rand.Rand) engine.Position {
	dirs := []engine.Position{
		{X: 0, Y: -1}, // Up
		{X: 0, Y: 1},  // Down
		{X: -1, Y: 0}, // Left
		{X: 1, Y: 0},  // Right
	}

	// 1. Gather all passable directions
	var legal []engine.Position
	var nonReverse []engine.Position

	for _, d := range dirs {
		next := m.WrapPosition(g.Pos.Add(d))
		if m.IsGhostPassable(next) {
			legal = append(legal, d)
			// Avoid instant 180° reversal if other moves exist
			if !(d.X == -g.Dir.X && d.Y == -g.Dir.Y) {
				nonReverse = append(nonReverse, d)
			}
		}
	}

	if len(legal) == 0 {
		return engine.Position{X: 0, Y: 0}
	}

	choices := nonReverse
	if len(choices) == 0 {
		choices = legal
	}

	// 2. If eaten, ghost moves directly back to ghost house
	if g.Eaten {
		return g.bestMoveTowards(choices, g.Home, m)
	}

	// 3. If vulnerable, move randomly or away from Pacman
	if g.Vulnerable {
		if rng != nil && len(choices) > 0 {
			return choices[rng.Intn(len(choices))]
		}
		return choices[0]
	}

	// 4. Apply personality targeting
	switch g.Personality {
	case "chase":
		return g.bestMoveTowards(choices, pacmanPos, m)

	case "ambush":
		// Ambush targets 4 tiles ahead of Pacman's current direction
		target := engine.Position{
			X: pacmanPos.X + pacmanDir.X*4,
			Y: pacmanPos.Y + pacmanDir.Y*4,
		}
		return g.bestMoveTowards(choices, target, m)

	default: // "random"
		if rng != nil && len(choices) > 0 {
			return choices[rng.Intn(len(choices))]
		}
		return choices[0]
	}
}

func (g *Ghost) bestMoveTowards(choices []engine.Position, target engine.Position, m *Maze) engine.Position {
	bestDir := choices[0]
	minDist := math.MaxFloat64

	for _, d := range choices {
		next := m.WrapPosition(g.Pos.Add(d))
		dist := manhattanDistance(next, target)
		if dist < minDist {
			minDist = dist
			bestDir = d
		}
	}
	return bestDir
}

func manhattanDistance(p1, p2 engine.Position) float64 {
	return math.Abs(float64(p1.X-p2.X)) + math.Abs(float64(p1.Y-p2.Y))
}
