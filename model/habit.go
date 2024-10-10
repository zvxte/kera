package model

import "time"

type Habit struct {
	ID              UUID
	Status          HabitStatus
	Title           string
	Description     string
	TrackedWeekDays TrackedWeekDays
	StartDate       time.Time
	EndDate         time.Time
	History         HabitHistory
}

type HabitStatus uint8

const (
	HabitActive HabitStatus = iota
	HabitEnded
)

type TrackedWeekDays uint8

const (
	Monday TrackedWeekDays = 1 << iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type HabitHistory []HabitDay

type HabitDay struct {
	Status HabitDayStatus
	Date   time.Time
}

type HabitDayStatus uint8

const (
	DayUntracked HabitDayStatus = iota
	DayDone
	DayMissed
	DayPending
)
