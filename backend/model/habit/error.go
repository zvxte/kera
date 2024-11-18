package habit

import "errors"

var (
	ErrTitleTooShort = errors.New("title is too short")
	ErrTitleTooLong  = errors.New("title is too long")
	ErrTitleInvalid  = errors.New("title is invalid")

	ErrDescriptionTooLong = errors.New("description is too long")
	ErrDescriptionInvalid = errors.New("description is invalid")

	ErrTrackedWeekDaysInvalid = errors.New(
		"tracked days of the week are invalid: unrecognized day specified",
	)
	ErrTrackedWeekDaysEmpty = errors.New(
		"at least one day of the week must be specified for tracking",
	)
)
