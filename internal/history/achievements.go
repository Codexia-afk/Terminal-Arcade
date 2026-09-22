package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Achievement represents a badge or milestone players can unlock.
type Achievement struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Hint        string     `json:"hint"`
	UnlockedAt  *time.Time `json:"unlocked_at,omitempty"`
}

// IsUnlocked reports whether this achievement has been unlocked.
func (a Achievement) IsUnlocked() bool {
	return a.UnlockedAt != nil && !a.UnlockedAt.IsZero()
}

// DefaultAchievements returns the complete set of supported achievements in initial locked state.
func DefaultAchievements() []Achievement {
	return []Achievement{
		{
			ID:          "first_blood",
			Title:       "First Blood",
			Description: "Complete your first session of any game.",
			Hint:        "Play and finish or exit a session in any game.",
		},
		{
			ID:          "century_club",
			Title:       "Century Club",
			Description: "Score 100+ points in any single Snake session.",
			Hint:        "Eat at least 10 food pellets in one Snake game.",
		},
		{
			ID:          "ghost_hunter",
			Title:       "Ghost Hunter",
			Description: "Eat 20 ghosts total across all Pacman sessions.",
			Hint:        "Consume frightened ghosts during power pellet periods.",
		},
		{
			ID:          "brick_breaker",
			Title:       "Brick Breaker",
			Description: "Clear a full Ball & Plate level on Hard.",
			Hint:        "Destroy every brick including tough bricks on Hard difficulty.",
		},
		{
			ID:          "marathon",
			Title:       "Marathon",
			Description: "Survive in a single session for 5+ minutes.",
			Hint:        "Keep any game session actively playing for at least 300 seconds.",
		},
		{
			ID:          "dedicated",
			Title:       "Dedicated",
			Description: "Achieve a 7-day play streak.",
			Hint:        "Play at least one game session each day for 7 consecutive days.",
		},
		{
			ID:          "completionist",
			Title:       "Completionist",
			Description: "Play all three arcade games at least once.",
			Hint:        "Play at least one round of Snake, Pacman, and Ball & Plate.",
		},
		{
			ID:          "perfectionist",
			Title:       "Perfectionist",
			Description: "Clear a Pacman level without losing a single life.",
			Hint:        "Eat all dots and win a Pacman game with 0 lives lost.",
		},
		{
			ID:          "speed_demon",
			Title:       "Speed Demon",
			Description: "Survive 60+ seconds on Snake Hard difficulty.",
			Hint:        "Dodge walls and obstacles at top speed for over a minute.",
		},
		{
			ID:          "night_owl",
			Title:       "Night Owl",
			Description: "Play a session between midnight and 4:00 AM.",
			Hint:        "Launch and play a round in the quiet late hours (00:00 - 04:00 local time).",
		},
	}
}

// AchievementsStore manages persistence for player achievements.
type AchievementsStore struct {
	mu    sync.Mutex
	Path  string
	Items []Achievement
}

// DefaultAchievementsPath returns os.UserConfigDir()/goarcade/achievements.json.
func DefaultAchievementsPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(configDir, "goarcade", "achievements.json"), nil
}

// OpenAchievements loads achievements from disk, repairing corruption if detected.
func OpenAchievements(path string) (*AchievementsStore, error) {
	if path == "" {
		var err error
		path, err = DefaultAchievementsPath()
		if err != nil {
			return nil, err
		}
	}

	store := &AchievementsStore{
		Path:  path,
		Items: DefaultAchievements(),
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read achievements file: %w", err)
	}
	if len(data) == 0 {
		return store, nil
	}

	var saved []Achievement
	if err := json.Unmarshal(data, &saved); err != nil {
		// File corrupted: backup and restore defaults
		timestamp := time.Now().Format("20060102-150405")
		backupPath := fmt.Sprintf("%s.corrupt-%s", path, timestamp)
		_ = os.Rename(path, backupPath)
		return store, nil
	}

	// Merge saved unlocked timestamps into defaults
	savedMap := make(map[string]*time.Time)
	for _, a := range saved {
		if a.UnlockedAt != nil {
			savedMap[a.ID] = a.UnlockedAt
		}
	}

	for i := range store.Items {
		if ut, ok := savedMap[store.Items[i].ID]; ok {
			store.Items[i].UnlockedAt = ut
		}
	}

	return store, nil
}

// Save writes achievements to disk atomically using temporary file + rename.
func (s *AchievementsStore) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked()
}

func (s *AchievementsStore) saveLocked() error {
	dir := filepath.Dir(s.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(s.Items, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode achievements: %w", err)
	}

	tmpPath := s.Path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write temporary achievements file: %w", err)
	}

	if err := os.Rename(tmpPath, s.Path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically rename achievements file: %w", err)
	}

	return nil
}

// Evaluate evaluates all achievements against game history and unlocks newly qualified badges.
// Returns a slice of newly unlocked achievements.
func (s *AchievementsStore) Evaluate(records []Record, now time.Time) []Achievement {
	s.mu.Lock()
	defer s.mu.Unlock()

	newlyUnlocked := make([]Achievement, 0)
	unlockedMap := make(map[string]bool)
	for _, a := range s.Items {
		if a.IsUnlocked() {
			unlockedMap[a.ID] = true
		}
	}

	// Precompute totals and checks
	totalGhostsEaten := 0
	playedGames := make(map[string]bool)
	streaks := CalculateStreaks(records, now)

	hasFirstBlood := len(records) > 0
	hasCenturyClub := false
	hasBrickBreaker := false
	hasMarathon := false
	hasPerfectionist := false
	hasSpeedDemon := false
	hasNightOwl := false

	for _, r := range records {
		playedGames[r.Game] = true

		if r.Game == "snake" && r.Score >= 100 {
			hasCenturyClub = true
		}
		if r.Game == "pacman" && r.Metrics != nil {
			totalGhostsEaten += r.Metrics["ghosts_eaten"]
		}
		if r.Game == "ballplate" && r.Difficulty == "hard" && r.Outcome == "won" {
			hasBrickBreaker = true
		}
		if r.Duration >= 5*time.Minute {
			hasMarathon = true
		}
		if r.Game == "pacman" && r.Outcome == "won" && r.Metrics != nil && r.Metrics["lives_lost"] == 0 {
			hasPerfectionist = true
		}
		if r.Game == "snake" && r.Difficulty == "hard" && r.Duration >= 60*time.Second {
			hasSpeedDemon = true
		}

		// Night owl check (PlayedAt in local time between 00:00 and 04:00)
		localHour := r.PlayedAt.In(now.Location()).Hour()
		if localHour >= 0 && localHour < 4 {
			hasNightOwl = true
		}
	}

	checkUnlock := func(id string, condition bool, unlockTime time.Time) {
		if condition && !unlockedMap[id] {
			tCopy := unlockTime
			for i := range s.Items {
				if s.Items[i].ID == id {
					s.Items[i].UnlockedAt = &tCopy
					newlyUnlocked = append(newlyUnlocked, s.Items[i])
					unlockedMap[id] = true
					break
				}
			}
		}
	}

	checkUnlock("first_blood", hasFirstBlood, now)
	checkUnlock("century_club", hasCenturyClub, now)
	checkUnlock("ghost_hunter", totalGhostsEaten >= 20, now)
	checkUnlock("brick_breaker", hasBrickBreaker, now)
	checkUnlock("marathon", hasMarathon, now)
	checkUnlock("dedicated", streaks.LongestStreak >= 7 || streaks.CurrentStreak >= 7, now)
	checkUnlock("completionist", playedGames["snake"] && playedGames["pacman"] && playedGames["ballplate"], now)
	checkUnlock("perfectionist", hasPerfectionist, now)
	checkUnlock("speed_demon", hasSpeedDemon, now)
	checkUnlock("night_owl", hasNightOwl, now)

	if len(newlyUnlocked) > 0 {
		_ = s.saveLocked()
	}

	return newlyUnlocked
}
