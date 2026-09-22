package history

import (
	"sort"
	"time"
)

// StreakInfo contains current and all-time consecutive day play streaks.
type StreakInfo struct {
	CurrentStreak int
	LongestStreak int
}

// CalculateStreaks determines the player's consecutive calendar day play streaks.
// now is injected to ensure deterministic streak calculation and reliable testing across day/timezone boundaries.
func CalculateStreaks(records []Record, now time.Time) StreakInfo {
	if len(records) == 0 {
		return StreakInfo{CurrentStreak: 0, LongestStreak: 0}
	}

	loc := now.Location()
	// Map unique calendar days normalized to midnight
	dayMap := make(map[int64]time.Time)
	for _, r := range records {
		localTime := r.PlayedAt.In(loc)
		midnight := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
		dayMap[midnight.Unix()] = midnight
	}

	if len(dayMap) == 0 {
		return StreakInfo{CurrentStreak: 0, LongestStreak: 0}
	}

	days := make([]time.Time, 0, len(dayMap))
	for _, d := range dayMap {
		days = append(days, d)
	}

	// Sort days ascending
	sort.Slice(days, func(i, j int) bool {
		return days[i].Before(days[j])
	})

	// Calculate longest streak
	longest := 1
	currentRun := 1
	for i := 1; i < len(days); i++ {
		prevExpectedNext := days[i-1].AddDate(0, 0, 1)
		if isSameCalendarDay(days[i], prevExpectedNext) {
			currentRun++
		} else {
			currentRun = 1
		}
		if currentRun > longest {
			longest = currentRun
		}
	}

	// Calculate current streak
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	yesterday := today.AddDate(0, 0, -1)

	// Check if played today or yesterday
	lastDay := days[len(days)-1]
	currentStreak := 0

	var checkDay time.Time
	if isSameCalendarDay(lastDay, today) {
		checkDay = today
	} else if isSameCalendarDay(lastDay, yesterday) {
		// Player hasn't played yet today, but played yesterday so streak remains active
		checkDay = yesterday
	} else {
		// Last played before yesterday: streak broken
		return StreakInfo{CurrentStreak: 0, LongestStreak: longest}
	}

	// Count backwards from checkDay
	for i := len(days) - 1; i >= 0; i-- {
		if isSameCalendarDay(days[i], checkDay) {
			currentStreak++
			checkDay = checkDay.AddDate(0, 0, -1)
		} else if days[i].Before(checkDay) {
			break
		}
	}

	return StreakInfo{
		CurrentStreak: currentStreak,
		LongestStreak: longest,
	}
}

func isSameCalendarDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
