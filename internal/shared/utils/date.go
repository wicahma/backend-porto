package utils

import (
	"fmt"
	"time"
)

// FormatDate formats a time as "Jan 2006".
func FormatDate(t time.Time) string {
	return t.Format("Jan 2006")
}

// FormatDateRange returns a human-readable date range like "Jan 2020 - Present".
func FormatDateRange(start, end time.Time, isCurrent bool) string {
	startStr := FormatDate(start)

	if isCurrent {
		return startStr + " - Present"
	}

	endStr := FormatDate(end)
	return startStr + " - " + endStr
}

// CalculateDuration computes the duration between two times as a human-readable string
// like "2 yrs 6 mos".
func CalculateDuration(start, end time.Time) string {
	if end.Before(start) {
		end = start
	}

	years := 0
	months := 0

	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	if monthDiff < 0 {
		yearDiff--
		monthDiff += 12
	}

	years = yearDiff
	months = monthDiff

	// Also account for days
	if end.Day() < start.Day() {
		if months == 0 {
			years--
			months = 11
		} else {
			months--
		}
	}

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d yr%s", years, plural(years)))
	}
	if months > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d mo%s", months, plural(months)))
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}

	return result
}

// CalculateTotalDuration sums the durations of multiple job periods and returns
// a human-readable string like "2 yrs 6 mos".
func CalculateTotalDuration(startEnds ...[2]time.Time) string {
	if len(startEnds) == 0 {
		return "0 mos"
	}

	totalMonths := 0
	for _, se := range startEnds {
		start := se[0]
		end := se[1]
		if end.Before(start) {
			end = start
		}

		months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month())
		if end.Day() < start.Day() {
			months--
		}
		if months < 0 {
			months = 0
		}
		totalMonths += months
	}

	years := totalMonths / 12
	months := totalMonths % 12

	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d yr%s", years, plural(years)))
	}
	if months > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%d mo%s", months, plural(months)))
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}

	return result
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
