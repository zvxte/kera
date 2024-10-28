package model

import "time"

// dateNow returns the current date in UTC
// with the time portion truncated to the midnight (00:00:00).
func dateNow() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour)
}
