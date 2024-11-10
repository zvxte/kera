package habit

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zvxte/kera/model"
)

const (
	titleMinChars = 2
	titleMaxChars = 64

	descriptionMaxChars = 256

	trackedWeekDaysMin = 1
	trackedWeekDaysMax = (1 << 7) - 1

	statusMin = 0
	statusMax = 1
)

// ValidateTitle fails if the provided title
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateTitle(title string) error {
	// Prevents from counting runes on a large string
	if len(title) > titleMaxChars*4 {
		return ErrTitleTooLong
	}

	length := utf8.RuneCountInString(title)
	if length < titleMinChars {
		return ErrTitleTooShort
	}
	if length > titleMaxChars {
		return ErrTitleTooLong
	}

	for _, c := range title {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return ErrTitleInvalid
		}
	}

	spaceCount := 0
	for _, r := range title {
		if r == ' ' {
			spaceCount++
		}
	}
	if (length - spaceCount) < titleMinChars {
		return ErrTitleTooShort
	}

	if strings.HasPrefix(title, " ") ||
		strings.HasSuffix(title, " ") {
		return ErrTitleInvalid
	}

	return nil
}

// ValidateDescription fails if the provided description
// does not meet the application requirements.
// The returned error is safe for client-side message.
func ValidateDescription(description string) error {
	// Prevents from counting runes on a large string
	if len(description) > descriptionMaxChars*4 {
		return ErrDescriptionTooLong
	}

	length := utf8.RuneCountInString(description)
	if length > descriptionMaxChars {
		return ErrDescriptionTooLong
	}

	for _, c := range description {
		if unicode.IsControl(c) || (unicode.IsSpace(c) && c != ' ') {
			return ErrDescriptionInvalid
		}
	}

	return nil
}

// validateTrackedWeekDays fails if the provided tracked week days
// do not meet the application requirements.
// The returned error is safe for client-side message.
func validateTrackedWeekDays(trackedWeekDays TrackedWeekDays) error {
	if trackedWeekDays < trackedWeekDaysMin {
		return ErrTrackedWeekDaysEmpty
	}
	if trackedWeekDays > trackedWeekDaysMax {
		return ErrTrackedWeekDaysInvalid
	}

	return nil
}

// validateHabitStatus fails if the provided status
// does not meet the application requirements.
// The returned error is safe for client-side message.
func validateStatus(s Status) error {
	if s < statusMin || s > statusMax {
		return model.ErrUnexpected
	}

	return nil
}
