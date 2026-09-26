package snake

import (
	"time"

	"github.com/Codexia-afk/Terminal-Arcade/internal/engine"
	"github.com/Codexia-afk/Terminal-Arcade/internal/games/snake/modes"
)

// Mode identifies one of the 5 distinct Snake gameplay modes.
type Mode string

const (
	ModeClassic    Mode = "classic"
	ModeZen        Mode = "zen"
	ModeSurvival   Mode = "survival"
	ModeTimeAttack Mode = "time_attack"
	ModeObstacle   Mode = "obstacle"
)

// ModeInfo describes a Snake mode for selection and UI presentation.
type ModeInfo struct {
	Mode        Mode
	Title       string
	Description string
	Icon        string
}

// Modes returns all 5 available Snake gameplay modes in canonical order.
func Modes() []ModeInfo {
	return []ModeInfo{
		{
			Mode:        ModeClassic,
			Title:       "Classic",
			Description: "Pure Nokia snake — survive as long as possible (the original).",
			Icon:        "🐍",
		},
		{
			Mode:        ModeZen,
			Title:       "Zen",
			Description: "No deaths, no walls, just peaceful growth. Score by length & flow.",
			Icon:        "🧘",
		},
		{
			Mode:        ModeSurvival,
			Title:       "Survival",
			Description: "Enemies chase you with pursuit AI. Clear 5 waves to win.",
			Icon:        "⚔️",
		},
		{
			Mode:        ModeTimeAttack,
			Title:       "Time Attack",
			Description: "High-intensity sprint against the clock. Time bonus for speed.",
			Icon:        "⏱️",
		},
		{
			Mode:        ModeObstacle,
			Title:       "Obstacle Challenge",
			Description: "Navigate escalating labyrinth mazes. 3 levels to victory.",
			Icon:        "🧩",
		},
	}
}

// ModeFromString safely parses a mode string with default fallback to Classic.
func ModeFromString(s string) Mode {
	switch Mode(s) {
	case ModeZen:
		return ModeZen
	case ModeSurvival:
		return ModeSurvival
	case ModeTimeAttack:
		return ModeTimeAttack
	case ModeObstacle:
		return ModeObstacle
	default:
		return ModeClassic
	}
}

// Enemy represents an AI chaser in Survival mode.
type Enemy = modes.SurvivalEnemy

// GenerateMaze generates fixed labyrinth obstacle layouts for Obstacle Challenge mode.
func GenerateMaze(level int, width, height int) map[engine.Position]bool {
	return modes.GenerateValidatedMaze(level, engine.Medium, width, height)
}

// TimeAttackDuration is the initial sprint time limit.
const TimeAttackDuration = 60 * time.Second
