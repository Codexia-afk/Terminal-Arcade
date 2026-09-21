// Package history persists and aggregates completed arcade game sessions.
// It is fully decoupled from rendering and tcell.
package history

import (
	"sort"
	"time"
)

// Record represents one completed or quit game session.
type Record struct {
	Game       string        `json:"game"`       // "snake" | "pacman" | "ballplate"
	Difficulty string        `json:"difficulty"` // "easy" | "medium" | "hard"
	Theme      string        `json:"theme"`
	Score      int           `json:"score"`
	Outcome    string        `json:"outcome"` // "collision", "quit", "won", etc.
	PlayedAt   time.Time     `json:"played_at"`
	Duration   time.Duration `json:"duration"`
}

// HighScore returns the highest score for a specific game and difficulty.
func HighScore(records []Record, game, difficulty string) int {
	best := 0
	for _, r := range records {
		if r.Game == game && r.Difficulty == difficulty && r.Score > best {
			best = r.Score
		}
	}
	return best
}

// GameHighScore returns a game's all-time highest score across all difficulties.
func GameHighScore(records []Record, game string) int {
	best := 0
	for _, r := range records {
		if r.Game == game && r.Score > best {
			best = r.Score
		}
	}
	return best
}

// TotalPlaytime returns the cumulative playtime, optionally filtered to a specific game.
// Pass game="" for all games.
func TotalPlaytime(records []Record, game string) time.Duration {
	var total time.Duration
	for _, r := range records {
		if game == "" || r.Game == game {
			total += r.Duration
		}
	}
	return total
}

// GamesPlayed returns the total number of sessions, optionally filtered to a specific game.
// Pass game="" for all games.
func GamesPlayed(records []Record, game string) int {
	count := 0
	for _, r := range records {
		if game == "" || r.Game == game {
			count++
		}
	}
	return count
}

// MostRecent returns the most recent session record, optionally filtered to a specific game.
// Pass game="" for all games.
func MostRecent(records []Record, game string) (Record, bool) {
	var latest Record
	found := false
	for _, r := range records {
		if game == "" || r.Game == game {
			if !found || r.PlayedAt.After(latest.PlayedAt) {
				latest = r
				found = true
			}
		}
	}
	return latest, found
}

// Recent returns up to n records in reverse chronological order (newest first),
// optionally filtered by game. Pass game="" for all games.
func Recent(records []Record, game string, n int) []Record {
	filtered := make([]Record, 0, len(records))
	for _, r := range records {
		if game == "" || r.Game == game {
			filtered = append(filtered, r)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PlayedAt.After(filtered[j].PlayedAt)
	})

	if n > 0 && len(filtered) > n {
		return filtered[:n]
	}
	return filtered
}
