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

func TestStreaksCalculation(t *testing.T) {
	// Fixed clock: 2026-06-10 14:00:00 UTC
	now := time.Date(2026, 6, 10, 14, 0, 0, 0, time.UTC)

	// Scenario 1: Empty records
	s0 := CalculateStreaks(nil, now)
	if s0.CurrentStreak != 0 || s0.LongestStreak != 0 {
		t.Errorf("expected 0/0 for empty records, got %+v", s0)
	}

	// Scenario 2: Played today, yesterday, and day before yesterday (3-day streak)
	records3 := []Record{
		{PlayedAt: time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 6, 9, 20, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 6, 10, 9, 0, 0, 0, time.UTC)},
	}
	s3 := CalculateStreaks(records3, now)
	if s3.CurrentStreak != 3 || s3.LongestStreak != 3 {
		t.Errorf("expected 3/3 streak, got %+v", s3)
	}

	// Scenario 3: Played yesterday, but hasn't played today yet -> streak still active (2 days)
	recordsActive := []Record{
		{PlayedAt: time.Date(2026, 6, 8, 10, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 6, 9, 20, 0, 0, 0, time.UTC)},
	}
	sActive := CalculateStreaks(recordsActive, now)
	if sActive.CurrentStreak != 2 || sActive.LongestStreak != 2 {
		t.Errorf("expected active 2/2 streak when last played yesterday, got %+v", sActive)
	}

	// Scenario 4: Gap in play -> longest preserved, current reset to 0
	recordsBroken := []Record{
		// 4 consecutive days in May
		{PlayedAt: time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 5, 3, 12, 0, 0, 0, time.UTC)},
		{PlayedAt: time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)},
		// Gap until June 7 (3 days before now -> streak broken)
		{PlayedAt: time.Date(2026, 6, 7, 12, 0, 0, 0, time.UTC)},
	}
	sBroken := CalculateStreaks(recordsBroken, now)
	if sBroken.CurrentStreak != 0 || sBroken.LongestStreak != 4 {
		t.Errorf("expected 0 current streak and 4 longest streak, got %+v", sBroken)
	}
}

func TestAchievementsUnlockLogic(t *testing.T) {
	tempDir := t.TempDir()
	achFile := filepath.Join(tempDir, "achievements.json")
	now := time.Date(2026, 6, 10, 14, 0, 0, 0, time.UTC)

	store, err := OpenAchievements(achFile)
	if err != nil {
		t.Fatalf("OpenAchievements failed: %v", err)
	}

	// Verify all 10 achievements exist initially locked
	if len(store.Items) < 10 {
		t.Fatalf("expected at least 10 achievements, got %d", len(store.Items))
	}
	for _, a := range store.Items {
		if a.IsUnlocked() {
			t.Errorf("expected achievement %s to be initially locked", a.ID)
		}
	}

	// Craft synthetic records to test unlocking each achievement
	records := []Record{
		// 1. "first_blood": session exists
		// 7. "completionist": snake, pacman, ballplate
		// 2. "century_club": snake score >= 100
		// 9. "speed_demon": snake hard duration >= 60s
		{
			Game:       "snake",
			Difficulty: "hard",
			Score:      120,
			Duration:   65 * time.Second,
			PlayedAt:   time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
			Metrics:    map[string]int{"max_length_reached": 15},
		},
		// 3. "ghost_hunter": 20+ ghosts
		// 8. "perfectionist": won pacman with 0 lives lost
		// 5. "marathon": duration >= 5 minutes (300s)
		{
			Game:       "pacman",
			Difficulty: "medium",
			Score:      3000,
			Outcome:    "won",
			Duration:   310 * time.Second,
			PlayedAt:   time.Date(2026, 6, 10, 13, 0, 0, 0, time.UTC),
			Metrics:    map[string]int{"ghosts_eaten": 22, "lives_lost": 0},
		},
		// 4. "brick_breaker": won ballplate on hard
		{
			Game:       "ballplate",
			Difficulty: "hard",
			Score:      800,
			Outcome:    "won",
			Duration:   120 * time.Second,
			PlayedAt:   time.Date(2026, 6, 10, 14, 0, 0, 0, time.UTC),
			Metrics:    map[string]int{"longest_rally": 18},
		},
		// 10. "night_owl": played between 00:00 and 04:00 local time
		{
			Game:       "snake",
			Difficulty: "easy",
			Score:      40,
			Duration:   20 * time.Second,
			PlayedAt:   time.Date(2026, 6, 10, 2, 30, 0, 0, time.UTC),
		},
	}

	// Add 7 consecutive days of sessions for 6. "dedicated"
	for day := 3; day <= 9; day++ {
		records = append(records, Record{
			Game:     "snake",
			PlayedAt: time.Date(2026, 6, day, 12, 0, 0, 0, time.UTC),
			Duration: 10 * time.Second,
		})
	}

	unlocked := store.Evaluate(records, now)
	if len(unlocked) != 10 {
		t.Fatalf("expected all 10 achievements to unlock, got %d", len(unlocked))
	}

	for _, a := range store.Items {
		if !a.IsUnlocked() {
			t.Errorf("achievement %s failed to unlock", a.ID)
		}
	}

	// Reopen store from disk to verify atomic persistence of achievements
	reloaded, err := OpenAchievements(achFile)
	if err != nil {
		t.Fatalf("reloading achievements failed: %v", err)
	}
	for _, a := range reloaded.Items {
		if !a.IsUnlocked() {
			t.Errorf("reloaded achievement %s was not persisted as unlocked", a.ID)
		}
	}
}

func TestHeatmapGeneration(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	records := []Record{
		{PlayedAt: now},                                              // 1 session today
		{PlayedAt: now.AddDate(0, 0, -1)},                           // 1 session yesterday
		{PlayedAt: now.AddDate(0, 0, -1).Add(2 * time.Hour)},        // 2nd session yesterday
		{PlayedAt: now.AddDate(0, 0, -2)},                           // 1 session 2 days ago
		{PlayedAt: now.AddDate(0, 0, -2).Add(time.Hour)},            // 2
		{PlayedAt: now.AddDate(0, 0, -2).Add(2 * time.Hour)},        // 3rd session (high activity)
	}

	activity := GenerateHeatmap(records, 30, now)
	if len(activity) != 30 {
		t.Fatalf("expected 30 days of activity, got %d", len(activity))
	}

	todayAct := activity[29]
	if todayAct.Count != 1 {
		t.Errorf("expected today count 1, got %d", todayAct.Count)
	}
	yestAct := activity[28]
	if yestAct.Count != 2 {
		t.Errorf("expected yesterday count 2, got %d", yestAct.Count)
	}
	twoDaysAgo := activity[27]
	if twoDaysAgo.Count != 3 {
		t.Errorf("expected 2 days ago count 3, got %d", twoDaysAgo.Count)
	}

	// Verify glyphs
	if HeatmapGlyph(0, false) != '·' {
		t.Errorf("expected · for 0 count")
	}
	if HeatmapGlyph(2, false) != '▪' {
		t.Errorf("expected ▪ for 2 count")
	}
	if HeatmapGlyph(4, false) != '█' {
		t.Errorf("expected █ for 4 count")
	}
	if HeatmapGlyph(3, true) != '#' {
		t.Errorf("expected # for 3 count in monochrome")
	}
}

func TestExportAndReset(t *testing.T) {
	tempDir := t.TempDir()
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	records := []Record{
		{
			Game:       "snake",
			Difficulty: "medium",
			Theme:      "Neon",
			Score:      150,
			Outcome:    "wall",
			PlayedAt:   now,
			Duration:   45 * time.Second,
			Metrics:    map[string]int{"food_eaten": 15},
		},
	}
	achievements := DefaultAchievements()

	// 1. Test JSON Export
	jsonPath := filepath.Join(tempDir, "export.json")
	jf, err := os.Create(jsonPath)
	if err != nil {
		t.Fatalf("failed to create json file: %v", err)
	}
	if err := ExportJSON(records, achievements, jf); err != nil {
		_ = jf.Close()
		t.Fatalf("ExportJSON failed: %v", err)
	}
	_ = jf.Close()

	data, err := os.ReadFile(jsonPath)
	if err != nil || len(data) == 0 {
		t.Fatalf("expected non-empty JSON export")
	}

	// 2. Test CSV Export
	csvPath := filepath.Join(tempDir, "export.csv")
	cf, err := os.Create(csvPath)
	if err != nil {
		t.Fatalf("failed to create csv file: %v", err)
	}
	if err := ExportCSV(records, cf); err != nil {
		_ = cf.Close()
		t.Fatalf("ExportCSV failed: %v", err)
	}
	_ = cf.Close()

	csvData, err := os.ReadFile(csvPath)
	if err != nil || len(csvData) == 0 {
		t.Fatalf("expected non-empty CSV export")
	}

	// 3. Test ResetAllData
	// Write dummy history.json and achievements.json
	_ = os.WriteFile(filepath.Join(tempDir, "history.json"), []byte("[]"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "achievements.json"), []byte("[]"), 0644)

	if err := ResetAllData(tempDir); err != nil {
		t.Fatalf("ResetAllData failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "history.json")); !os.IsNotExist(err) {
		t.Errorf("expected history.json to be deleted after ResetAllData")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "achievements.json")); !os.IsNotExist(err) {
		t.Errorf("expected achievements.json to be deleted after ResetAllData")
	}
}

func TestPersonalBestsAndRelativeTime(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	records := []Record{
		{
			Game:       "snake",
			Difficulty: "hard",
			Score:      90,
			PlayedAt:   now.Add(-2 * time.Hour),
			Duration:   50 * time.Second,
			Metrics:    map[string]int{"max_length_reached": 12},
		},
		{
			Game:       "snake",
			Difficulty: "medium",
			Score:      150,
			PlayedAt:   now.AddDate(0, 0, -2),
			Duration:   70 * time.Second,
			Metrics:    map[string]int{"max_length_reached": 18},
		},
	}

	pb, ok := PersonalBests(records, "snake")
	if !ok {
		t.Fatalf("expected personal bests found for snake")
	}
	if pb.HighScore != 150 {
		t.Errorf("expected high score 150, got %d", pb.HighScore)
	}
	if pb.KeyMetricValue != 18 {
		t.Errorf("expected key metric 18, got %d", pb.KeyMetricValue)
	}

	// Test FormatRelativeTime
	if FormatRelativeTime(now.Add(-30*time.Second), now) != "just now" {
		t.Errorf("expected 'just now'")
	}
	if FormatRelativeTime(now.Add(-10*time.Minute), now) != "10 mins ago" {
		t.Errorf("expected '10 mins ago'")
	}
	if FormatRelativeTime(now.Add(-2*time.Hour), now) != "2 hours ago" {
		t.Errorf("expected '2 hours ago'")
	}
	if FormatRelativeTime(now.AddDate(0, 0, -1), now) != "yesterday" {
		t.Errorf("expected 'yesterday'")
	}
	if FormatRelativeTime(now.AddDate(0, 0, -5), now) != "5 days ago" {
		t.Errorf("expected '5 days ago'")
	}
}
