package date

import "time"

// Date represents the date in UTC
// with the time portion truncated to the midnight (00:00:00)
type Date time.Time

func New(year int, month time.Month, day int) Date {
	return Date(
		time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	)
}

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

func (d Date) After(other Date) bool {
	return time.Time(d).After(time.Time(other))
}

func (d Date) Sub(other Date) time.Duration {
	return time.Time(d).Sub(time.Time(other))
}

func (d Date) FirstOfMonth() Date {
	return Date(
		time.Time(d).AddDate(0, 0, 1-time.Time(d).Day()),
	)
}

func (d Date) MaxDays() int {
	return time.Date(
		time.Time(d).Year(),
		time.Time(d).Month()+1,
		0, 0, 0, 0, 0, time.UTC,
	).Day()
}

func (d Date) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d Date) Equal(other Date) bool {
	return time.Time(d).Equal(time.Time(other))
}

func (d Date) WeekDay() uint8 {
	switch day := uint8(time.Time(d).Weekday()); day {
	case 1, 2, 3, 4, 5, 6:
		return day - 1
	default:
		return 6
	}
}

func (d Date) Year() int {
	return time.Time(d).Year()
}

func (d Date) String() string {
	return time.Time(d).String()
}
