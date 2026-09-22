package history

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ExportPayload wraps the complete local archive for JSON export.
type ExportPayload struct {
	ExportedAt   time.Time     `json:"exported_at"`
	TotalRecords int           `json:"total_records"`
	Records      []Record      `json:"records"`
	Achievements []Achievement `json:"achievements"`
}

// ExportJSON writes all records and achievements to w as formatted JSON.
func ExportJSON(records []Record, achievements []Achievement, w io.Writer) error {
	payload := ExportPayload{
		ExportedAt:   time.Now(),
		TotalRecords: len(records),
		Records:      records,
		Achievements: achievements,
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize export data: %w", err)
	}

	_, err = w.Write(data)
	return err
}

// ExportCSV writes all session records to w as a CSV document with headers.
func ExportCSV(records []Record, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{
		"PlayedAt",
		"Game",
		"Difficulty",
		"Theme",
		"Score",
		"Outcome",
		"DurationSeconds",
		"Metrics",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for _, r := range records {
		var metricPairs []string
		if r.Metrics != nil {
			keys := make([]string, 0, len(r.Metrics))
			for k := range r.Metrics {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				metricPairs = append(metricPairs, fmt.Sprintf("%s=%d", k, r.Metrics[k]))
			}
		}

		row := []string{
			r.PlayedAt.Format(time.RFC3339),
			r.Game,
			r.Difficulty,
			r.Theme,
			strconv.Itoa(r.Score),
			r.Outcome,
			strconv.FormatFloat(r.Duration.Seconds(), 'f', 2, 64),
			strings.Join(metricPairs, ";"),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}
