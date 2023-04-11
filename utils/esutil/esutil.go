package esutil

import (
	"fmt"
	"strings"
	"time"
)

// EMLA uses the layout `yyyy.MM.dd`.
const layoutISO = "2006.01.02"

func ResolveIndexNames(prefix string, start, end time.Time) string {
	if end.IsZero() {
		end = time.Now()
	}

	// If start is after end, ResolveIndexNames returns "".
	// However, query parameter validation will prevent it from happening at the very beginning (Bad Request).
	if start.After(end) {
		return ""
	}

	// In case of no start time or a broad query range over 30 days, search all indices.
	if start.IsZero() || end.Sub(start).Hours() > 24*30 {
		return fmt.Sprintf("%s*", prefix)
	}

	var indices []string
	// Elasticsearch creates indices based on UTC time every day.
	// Truncate(24 * time.Hour) returns the result of rounding 'end' down to 1d based on UTC time.
	end = end.Truncate(24 * time.Hour)
	suffix := end.Format(layoutISO)
	indices = append(indices, fmt.Sprintf("%s-%s", prefix, suffix))
	for start.Before(end) {
		end = end.Add(-24 * time.Hour)
		suffix = end.Format(layoutISO)
		indices = append(indices, fmt.Sprintf("%s-%s", prefix, suffix))
	}

	return strings.Join(indices, ",")
}
