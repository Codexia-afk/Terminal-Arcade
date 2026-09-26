// Package snake provides the Nokia-inspired 5-mode snake game with refined mechanics.
package snake

import (
	"math"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
)

// DirectionQueue buffers up to 2 player direction inputs so rapid key taps are never lost.
type DirectionQueue struct {
	items []engine.Position
}

// NewDirectionQueue creates an empty 2-step direction queue.
func NewDirectionQueue() *DirectionQueue {
	return &DirectionQueue{items: make([]engine.Position, 0, 2)}
}

// Push adds a new direction if it is not an immediate 180° neck reversal.
func (q *DirectionQueue) Push(curDir, newDir engine.Position) {
	last := curDir
	if len(q.items) > 0 {
		last = q.items[len(q.items)-1]
	}

	// Reject exact reverse (180° reversal into snake body)
	if (newDir.X == -last.X && newDir.Y == -last.Y) || (newDir.X == last.X && newDir.Y == last.Y) {
		return
	}

	if len(q.items) < 2 {
		q.items = append(q.items, newDir)
	} else {
		// Overwrite the pending second input with the newest player intent
		q.items[1] = newDir
	}
}

// Pop extracts the next queued direction.
func (q *DirectionQueue) Pop(curDir engine.Position) (engine.Position, bool) {
	if len(q.items) == 0 {
		return curDir, false
	}
	next := q.items[0]
	q.items = q.items[1:]
	return next, true
}

// Clear empties the buffered directions.
func (q *DirectionQueue) Clear() {
	q.items = q.items[:0]
}

// SubCellCollision tests if two entities collide at cell centers with sub-cell tolerance.
func SubCellCollision(p1, p2 engine.Position, tolerance float64) bool {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist < tolerance
}

// ManhattenDistance computes grid distance between two positions.
func ManhattanDistance(a, b engine.Position) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// PredictNextCell returns where an entity will be on the following tick.
func PredictNextCell(pos, dir engine.Position) engine.Position {
	return pos.Add(dir)
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

