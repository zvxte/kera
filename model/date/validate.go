package date

const (
	yearMin = 2024

	monthMin = 1
	monthMax = 12
)

var (
	yearMax = Now().Year() + 1
)

func ValidateYear(year int) error {
	if year < yearMin || year > yearMax {
		return ErrInvalidYear
	}
	return nil
}

func ValidateMonth(month int) error {
	if month < monthMin || month > monthMax {
		return ErrInvalidMonth
	}
	return nil
}
