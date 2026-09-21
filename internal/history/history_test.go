package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStoreRoundTripAndAtomicWrite(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.json")

	store, err := Open(historyFile)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	if len(store.Records) != 0 {
		t.Fatalf("expected 0 initial records, got %d", len(store.Records))
	}

	now := time.Now()
	rec := Record{
		Game:       "snake",
		Difficulty: "medium",
		Theme:      "Neon",
		Score:      150,
		Outcome:    "collision",
		PlayedAt:   now,
		Duration:   45 * time.Second,
	}

	if err := store.Append(rec); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Verify file exists on disk
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Fatalf("expected history file to exist at %s", historyFile)
	}

	// Reopen store from disk
	reloaded, err := Open(historyFile)
	if err != nil {
		t.Fatalf("reloading store failed: %v", err)
	}

	if len(reloaded.Records) != 1 {
		t.Fatalf("expected 1 record after reload, got %d", len(reloaded.Records))
	}
	r := reloaded.Records[0]
	if r.Game != "snake" || r.Score != 150 || r.Outcome != "collision" {
		t.Fatalf("reloaded record mismatch: %+v", r)
	}
}

func TestCorruptionRecovery(t *testing.T) {
	tempDir := t.TempDir()
	historyFile := filepath.Join(tempDir, "history.json")

	// Write invalid JSON
	corruptContent := []byte("{ this is not valid JSON content :::")
	if err := os.WriteFile(historyFile, corruptContent, 0644); err != nil {
		t.Fatalf("failed to write corrupted file: %v", err)
	}

	// Open should detect corruption, backup the file, and return empty store without error
	store, err := Open(historyFile)
	if err != nil {
		t.Fatalf("expected Open to recover from corruption, but got error: %v", err)
	}

	if len(store.Records) != 0 {
		t.Fatalf("expected empty records after corruption recovery, got %d", len(store.Records))
	}

	// Verify a backup file was created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	foundBackup := false
	for _, f := range files {
		if f.Name() != "history.json" {
			foundBackup = true
			break
		}
	}
	if !foundBackup {
		t.Errorf("expected backup file to be created in %s", tempDir)
	}
}

func TestAggregationQueries(t *testing.T) {
	now := time.Now()
	records := []Record{
		{Game: "snake", Difficulty: "easy", Score: 50, PlayedAt: now.Add(-5 * time.Hour), Duration: 30 * time.Second},
		{Game: "snake", Difficulty: "easy", Score: 80, PlayedAt: now.Add(-4 * time.Hour), Duration: 40 * time.Second},
		{Game: "snake", Difficulty: "medium", Score: 120, PlayedAt: now.Add(-3 * time.Hour), Duration: 60 * time.Second},
		{Game: "pacman", Difficulty: "easy", Score: 500, PlayedAt: now.Add(-2 * time.Hour), Duration: 120 * time.Second},
		{Game: "pacman", Difficulty: "hard", Score: 1200, PlayedAt: now.Add(-1 * time.Hour), Duration: 150 * time.Second},
		{Game: "ballplate", Difficulty: "medium", Score: 300, PlayedAt: now, Duration: 90 * time.Second},
	}

	// HighScore per (game, difficulty)
	if HighScore(records, "snake", "easy") != 80 {
		t.Errorf("expected snake easy highscore 80, got %d", HighScore(records, "snake", "easy"))
	}
	if HighScore(records, "snake", "medium") != 120 {
		t.Errorf("expected snake medium highscore 120, got %d", HighScore(records, "snake", "medium"))
	}
	if HighScore(records, "snake", "hard") != 0 {
		t.Errorf("expected snake hard highscore 0, got %d", HighScore(records, "snake", "hard"))
	}

	// GameHighScore across difficulties
	if GameHighScore(records, "snake") != 120 {
		t.Errorf("expected snake all-time highscore 120, got %d", GameHighScore(records, "snake"))
	}
	if GameHighScore(records, "pacman") != 1200 {
		t.Errorf("expected pacman all-time highscore 1200, got %d", GameHighScore(records, "pacman"))
	}

	// TotalPlaytime
	expectedSnakeTime := 130 * time.Second
	if TotalPlaytime(records, "snake") != expectedSnakeTime {
		t.Errorf("expected snake total playtime %v, got %v", expectedSnakeTime, TotalPlaytime(records, "snake"))
	}
	expectedTotalTime := 490 * time.Second
	if TotalPlaytime(records, "") != expectedTotalTime {
		t.Errorf("expected cumulative playtime %v, got %v", expectedTotalTime, TotalPlaytime(records, ""))
	}

	// GamesPlayed
	if GamesPlayed(records, "snake") != 3 {
		t.Errorf("expected snake count 3, got %d", GamesPlayed(records, "snake"))
	}
	if GamesPlayed(records, "") != 6 {
		t.Errorf("expected total count 6, got %d", GamesPlayed(records, ""))
	}

	// MostRecent
	latestSnake, ok := MostRecent(records, "snake")
	if !ok || latestSnake.Score != 120 {
		t.Errorf("expected latest snake score 120, got %+v", latestSnake)
	}
	latestOverall, ok := MostRecent(records, "")
	if !ok || latestOverall.Game != "ballplate" {
		t.Errorf("expected latest overall game ballplate, got %+v", latestOverall)
	}

	// Recent list ordering and limit
	recent3 := Recent(records, "", 3)
	if len(recent3) != 3 {
		t.Fatalf("expected 3 recent records, got %d", len(recent3))
	}
	if recent3[0].Game != "ballplate" || recent3[1].Game != "pacman" || recent3[2].Game != "pacman" {
		t.Errorf("unexpected ordering in Recent: %+v", recent3)
	}
}
