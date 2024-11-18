package habit

import (
	"reflect"
	"testing"

	"github.com/zvxte/kera/model/date"
)

func TestLoadHistoryFromBitmap(t *testing.T) {
	tests := []struct {
		name        string
		historyDate date.Date
		days        uint
		tracked     TrackedWeekDays
		startDate   date.Date
		endDate     date.Date
		expected    History
	}{
		{
			"Valid: first week done, entire week tracked, month starts from monday",
			date.New(2024, 7, 1),
			0b_01111111,
			0b_01111111,
			date.New(2024, 7, 1),
			date.New(2024, 7, 7),
			[]Day{
				{DayDone, date.New(2024, 7, 1)},
				{DayDone, date.New(2024, 7, 2)},
				{DayDone, date.New(2024, 7, 3)},
				{DayDone, date.New(2024, 7, 4)},
				{DayDone, date.New(2024, 7, 5)},
				{DayDone, date.New(2024, 7, 6)},
				{DayDone, date.New(2024, 7, 7)},
				{DayUntracked, date.New(2024, 7, 8)},
				{DayUntracked, date.New(2024, 7, 9)},
				{DayUntracked, date.New(2024, 7, 10)},
				{DayUntracked, date.New(2024, 7, 11)},
				{DayUntracked, date.New(2024, 7, 12)},
				{DayUntracked, date.New(2024, 7, 13)},
				{DayUntracked, date.New(2024, 7, 14)},
				{DayUntracked, date.New(2024, 7, 15)},
				{DayUntracked, date.New(2024, 7, 16)},
				{DayUntracked, date.New(2024, 7, 17)},
				{DayUntracked, date.New(2024, 7, 18)},
				{DayUntracked, date.New(2024, 7, 19)},
				{DayUntracked, date.New(2024, 7, 20)},
				{DayUntracked, date.New(2024, 7, 21)},
				{DayUntracked, date.New(2024, 7, 22)},
				{DayUntracked, date.New(2024, 7, 23)},
				{DayUntracked, date.New(2024, 7, 24)},
				{DayUntracked, date.New(2024, 7, 25)},
				{DayUntracked, date.New(2024, 7, 26)},
				{DayUntracked, date.New(2024, 7, 27)},
				{DayUntracked, date.New(2024, 7, 28)},
				{DayUntracked, date.New(2024, 7, 29)},
				{DayUntracked, date.New(2024, 7, 30)},
				{DayUntracked, date.New(2024, 7, 31)},
			},
		},
		{
			`Valid: first week missed, mon-fri tracked,
			month starts from monday, habit ends at day 10`,
			date.New(2024, 4, 1),
			0b_11_10000000,
			0b_00011111,
			date.New(2024, 4, 1),
			date.New(2024, 4, 10),
			[]Day{
				{DayMissed, date.New(2024, 4, 1)},
				{DayMissed, date.New(2024, 4, 2)},
				{DayMissed, date.New(2024, 4, 3)},
				{DayMissed, date.New(2024, 4, 4)},
				{DayMissed, date.New(2024, 4, 5)},
				{DayUntracked, date.New(2024, 4, 6)},
				{DayUntracked, date.New(2024, 4, 7)},
				{DayDone, date.New(2024, 4, 8)},
				{DayDone, date.New(2024, 4, 9)},
				{DayDone, date.New(2024, 4, 10)},
				{DayUntracked, date.New(2024, 4, 11)},
				{DayUntracked, date.New(2024, 4, 12)},
				{DayUntracked, date.New(2024, 4, 13)},
				{DayUntracked, date.New(2024, 4, 14)},
				{DayUntracked, date.New(2024, 4, 15)},
				{DayUntracked, date.New(2024, 4, 16)},
				{DayUntracked, date.New(2024, 4, 17)},
				{DayUntracked, date.New(2024, 4, 18)},
				{DayUntracked, date.New(2024, 4, 19)},
				{DayUntracked, date.New(2024, 4, 20)},
				{DayUntracked, date.New(2024, 4, 21)},
				{DayUntracked, date.New(2024, 4, 22)},
				{DayUntracked, date.New(2024, 4, 23)},
				{DayUntracked, date.New(2024, 4, 24)},
				{DayUntracked, date.New(2024, 4, 25)},
				{DayUntracked, date.New(2024, 4, 26)},
				{DayUntracked, date.New(2024, 4, 27)},
				{DayUntracked, date.New(2024, 4, 28)},
				{DayUntracked, date.New(2024, 4, 29)},
				{DayUntracked, date.New(2024, 4, 30)},
			},
		},
		{
			`Valid: mon-wed + sunday tracked,
			month starts from wednesday, habit starts at day 5,
			habit ends at day 25
			`,
			date.New(2024, 5, 1),
			0b_0000010_1111100_1000011_10000,
			0b_1000111,
			date.New(2024, 5, 5),
			date.New(2024, 5, 25),
			[]Day{
				{DayUntracked, date.New(2024, 5, 1)},
				{DayUntracked, date.New(2024, 5, 2)},
				{DayUntracked, date.New(2024, 5, 3)},
				{DayUntracked, date.New(2024, 5, 4)},
				{DayDone, date.New(2024, 5, 5)},

				{DayDone, date.New(2024, 5, 6)},
				{DayDone, date.New(2024, 5, 7)},
				{DayMissed, date.New(2024, 5, 8)},
				{DayUntracked, date.New(2024, 5, 9)},
				{DayUntracked, date.New(2024, 5, 10)},
				{DayUntracked, date.New(2024, 5, 11)},
				{DayDone, date.New(2024, 5, 12)},

				{DayMissed, date.New(2024, 5, 13)},
				{DayMissed, date.New(2024, 5, 14)},
				{DayDone, date.New(2024, 5, 15)},
				{DayUntracked, date.New(2024, 5, 16)},
				{DayUntracked, date.New(2024, 5, 17)},
				{DayUntracked, date.New(2024, 5, 18)},
				{DayDone, date.New(2024, 5, 19)},

				{DayMissed, date.New(2024, 5, 20)},
				{DayDone, date.New(2024, 5, 21)},
				{DayMissed, date.New(2024, 5, 22)},
				{DayUntracked, date.New(2024, 5, 23)},
				{DayUntracked, date.New(2024, 5, 24)},
				{DayUntracked, date.New(2024, 5, 25)},

				{DayUntracked, date.New(2024, 5, 26)},
				{DayUntracked, date.New(2024, 5, 27)},
				{DayUntracked, date.New(2024, 5, 28)},
				{DayUntracked, date.New(2024, 5, 29)},
				{DayUntracked, date.New(2024, 5, 30)},
				{DayUntracked, date.New(2024, 5, 31)},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := LoadHistoryFromBitmap(
				test.historyDate, test.days, test.tracked,
				test.startDate, test.endDate,
			)
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf(
					"LoadHistoryFromBitmap(%q, %v, %v, %q, %q), \ngot=%v, \nexpected=%v",
					test.historyDate, test.days, test.tracked, test.startDate, test.endDate,
					got, test.expected,
				)
			}
		})
	}
}
