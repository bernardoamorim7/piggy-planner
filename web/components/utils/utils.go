package utils

import "time"

// ParseDate returns a formatted date string
func ParseDate(date time.Time) string {
	if date.IsZero() || date.Equal(time.Time{}) {
		return ""
	}
	return date.Format("2006-01-02")
}
