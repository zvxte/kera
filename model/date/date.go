package date

import "time"

// Date represents the date in UTC
// with the time portion truncated to the midnight (00:00:00)
type Date time.Time

// Now returns the current Date value.
func Now() Date {
	return Date(
		time.Now().UTC().Truncate(24 * time.Hour),
	)
}

// Load returns Date value from provided time.Time.
func Load(t time.Time) Date {
	return Date(
		t.UTC().Truncate(24 * time.Hour),
	)
}

// Add returns the Date value + provided duration.
func (d Date) Add(duration time.Duration) Date {
	return Date(
		time.Time(d).Add(duration).Truncate(24 * time.Hour),
	)
}

// Before returns true if the other is before d.
func (d Date) Before(other Date) bool {
	return time.Time(d).Before(time.Time(other))
}

func (d Date) String() string {
	return time.Time(d).String()
}
