package history

import (
	"strings"
	"time"
)

// DayActivity models the session volume for one calendar day.
type DayActivity struct {
	Date  time.Time
	Count int
}

// GenerateHeatmap calculates daily session counts for the last n days leading up to now.
func GenerateHeatmap(records []Record, days int, now time.Time) []DayActivity {
	if days <= 0 {
		days = 30
	}

	loc := now.Location()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	counts := make(map[int64]int)
	for _, r := range records {
		localTime := r.PlayedAt.In(loc)
		midnight := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
		counts[midnight.Unix()]++
	}

	result := make([]DayActivity, days)
	for i := 0; i < days; i++ {
		offset := days - 1 - i
		d := today.AddDate(0, 0, -offset)
		result[i] = DayActivity{
			Date:  d,
			Count: counts[d.Unix()],
		}
	}

	return result
}

// HeatmapGlyph returns the appropriate ASCII/Unicode glyph for the given play count.
func HeatmapGlyph(count int, isMonochrome bool) rune {
	if count <= 0 {
		return '·'
	}
	if isMonochrome {
		if count <= 2 {
			return 'o'
		}
		return '#'
	}
	if count <= 2 {
		return '▪'
	}
	return '█'
}

// RenderHeatmapLine formats a row of daily activity blocks with day markers.
func RenderHeatmapLine(activity []DayActivity, isMonochrome bool) string {
	var b strings.Builder
	for _, a := range activity {
		b.WriteRune(HeatmapGlyph(a.Count, isMonochrome))
		b.WriteRune(' ')
	}
	return strings.TrimSpace(b.String())
}
