// Package history persists and aggregates completed arcade game sessions.
// It is fully decoupled from rendering and tcell.
package history

import (
	"fmt"
	"sort"
	"time"
)

// Record represents one completed or quit game session.
type Record struct {
	Game       string         `json:"game"`                 // "snake" | "pacman" | "ballplate"
	Mode       string         `json:"mode,omitempty"`       // snake mode: "classic", "zen", "survival", "time_attack", "obstacle"
	Difficulty string         `json:"difficulty"`           // "easy" | "medium" | "hard"
	Theme      string         `json:"theme"`
	Score      int            `json:"score"`
	Outcome    string         `json:"outcome"` // "collision", "quit", "won", etc.
	PlayedAt   time.Time      `json:"played_at"`
	Duration   time.Duration  `json:"duration"`
	Metrics    map[string]int `json:"metrics,omitempty"`
}

// EffectiveMode returns the normalized mode (defaults to "classic" for snake if empty).
func (r Record) EffectiveMode() string {
	if r.Mode != "" {
		return r.Mode
	}
	if r.Game == "snake" {
		return "classic"
	}
	return ""
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

// PersonalBestInfo contains highlights of personal records for a given game.
type PersonalBestInfo struct {
	HighScore       int
	HighScoreDiff   string
	HighScoreDate   time.Time
	LongestDuration time.Duration
	KeyMetricName   string
	KeyMetricValue  int
}

// PersonalBests queries history for a game's standout achievements.
func PersonalBests(records []Record, game string) (PersonalBestInfo, bool) {
	var info PersonalBestInfo
	found := false

	for _, r := range records {
		if r.Game != game {
			continue
		}
		found = true

		if r.Score > info.HighScore || (r.Score == info.HighScore && info.HighScore == 0) {
			info.HighScore = r.Score
			info.HighScoreDiff = r.Difficulty
			info.HighScoreDate = r.PlayedAt
		}
		if r.Duration > info.LongestDuration {
			info.LongestDuration = r.Duration
		}

		if r.Metrics != nil {
			switch game {
			case "snake":
				info.KeyMetricName = "Max Length"
				if val := r.Metrics["max_length_reached"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			case "pacman":
				info.KeyMetricName = "Ghosts Eaten"
				if val := r.Metrics["ghosts_eaten"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			case "ballplate":
				info.KeyMetricName = "Longest Rally"
				if val := r.Metrics["longest_rally"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			}
		}
	}

	return info, found
}

// PersonalBestForConfig finds the highest score for a specific game, mode, and difficulty.
func PersonalBestForConfig(records []Record, game, mode, difficulty string) (int, bool) {
	best := -1
	for _, r := range records {
		if r.Game != game {
			continue
		}
		if mode != "" && r.EffectiveMode() != mode {
			continue
		}
		if difficulty != "" && r.Difficulty != difficulty {
			continue
		}
		if r.Score > best {
			best = r.Score
		}
	}
	if best >= 0 {
		return best, true
	}
	return 0, false
}

// PersonalBestsForMode queries history for a game and mode's standout achievements.
func PersonalBestsForMode(records []Record, game, mode string) (PersonalBestInfo, bool) {
	var info PersonalBestInfo
	found := false

	for _, r := range records {
		if r.Game != game {
			continue
		}
		if mode != "" && r.EffectiveMode() != mode {
			continue
		}
		found = true

		if r.Score > info.HighScore || (r.Score == info.HighScore && info.HighScore == 0) {
			info.HighScore = r.Score
			info.HighScoreDiff = r.Difficulty
			info.HighScoreDate = r.PlayedAt
		}
		if r.Duration > info.LongestDuration {
			info.LongestDuration = r.Duration
		}

		if r.Metrics != nil {
			switch game {
			case "snake":
				info.KeyMetricName = "Max Length"
				if val := r.Metrics["max_length_reached"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			case "pacman":
				info.KeyMetricName = "Ghosts Eaten"
				if val := r.Metrics["ghosts_eaten"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			case "ballplate":
				info.KeyMetricName = "Longest Rally"
				if val := r.Metrics["longest_rally"]; val > info.KeyMetricValue {
					info.KeyMetricValue = val
				}
			}
		}
	}

	return info, found
}


// FormatRelativeTime converts a timestamp into an intuitive human-readable offset.
func FormatRelativeTime(t time.Time, now time.Time) string {
	if t.IsZero() {
		return "never"
	}
	diff := now.Sub(t)
	if diff < 0 {
		diff = 0
	}

	if diff < time.Minute {
		return "just now"
	}
	if diff < time.Hour {
		m := int(diff.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", m)
	}
	if diff < 24*time.Hour {
		h := int(diff.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	}
	days := int(diff.Hours() / 24)
	if days == 1 {
		return "yesterday"
	}
	return fmt.Sprintf("%d days ago", days)
}

// FormatDuration formats a time.Duration into compact readable format like "2m14s" or "45s".
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	if m > 0 {
		return fmt.Sprintf("%dm%02ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

