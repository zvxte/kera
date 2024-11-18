package habit

import (
	"time"

	"github.com/zvxte/kera/model/date"
)

const (
	// HistoryPatchWindow represents a timeframe,
	// in which the history can be patched.
	HistoryPatchWindow = 7 * 24 * time.Hour
)

// History represents a history of a habit.
type History []Day

// HabitDay represents a single day record in a history of a habit.
// It contains the status and the date of that record.
type Day struct {
	Status DayStatus
	Date   date.Date
}

func newDay(status DayStatus, date date.Date) Day {
	return Day{Status: status, Date: date}
}

// DayStatus represents a status of a single day record in a history of a habit.
type DayStatus uint8

const (
	DayUntracked DayStatus = iota
	DayDone
	DayMissed
	DayPending
)

func LoadHistoryFromBitmap(
	historyDate date.Date, days uint,
	trackedWeekDays TrackedWeekDays,
	startDate, endDate date.Date,
) History {
	firstOfMonth := historyDate.FirstOfMonth()

	maxDays := historyDate.MaxDays()
	history := make([]Day, maxDays)

	for i := 0; i < maxDays; i++ {
		bit := (days >> i) & 1
		day := firstOfMonth.Add(time.Duration(i) * 24 * time.Hour)
		isTracked := trackedWeekDays.Tracked(WeekDay(day.WeekDay()))
		now := date.Now()

		switch {
		case
			day.Before(startDate),
			day.After(now),
			!endDate.IsZero() && day.After(endDate),
			!isTracked:
			history[i] = newDay(DayUntracked, day)

		case
			bit == 1:
			history[i] = newDay(DayDone, day)

		case
			day.Equal(now):
			history[i] = newDay(DayPending, day)

		default:
			history[i] = newDay(DayMissed, day)
		}
	}

	return History(history)
}

func NewUntrackedHistory(historyDate date.Date) History {
	firstOfMonth := historyDate.FirstOfMonth()

	maxDays := historyDate.MaxDays()
	history := make([]Day, maxDays)

	for i := 0; i < maxDays; i++ {
		dayDate := firstOfMonth.Add(time.Duration(i) * 24 * time.Hour)
		history[i] = newDay(DayUntracked, dayDate)
	}

	return history
}
