package convertors

import (
	"errors"
	"time"
)

// ParseTimeWithDate parses a string like "11:30 13-05-2025" into a time.Time value
func ParseTimeWithDate(input string) (time.Time, error) {
	layout := "15:04 02-01-2006" // format: hour:minute day-month-year
	return time.Parse(layout, input)
}

// ValidateDateTime checks if the input time is not in the past.
// It allows the current date, but not a time earlier than now.
func ValidateDateTime(t time.Time) error {
	now := time.Now()

	// If the date is different, ensure it's not before today
	if t.Year() != now.Year() || t.Month() != now.Month() || t.Day() != now.Day() {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if t.Before(today) {
			return errors.New("time must not be earlier than today")
		}

		return nil
	}

	// If the date is today, ensure it's not before the current time
	if t.Before(now) {
		return errors.New("time must not be earlier than the current time")
	}

	return nil
}

// ParseAndValidateTime parses a string like "11:30 13-05-2025" into a time.Time value
// and validates that the time is not in the past.
func ParseAndValidateTime(input string) (time.Time, error) {
	// Define the layout for parsing
	layout := "15:04 02-01-2006" // format: hour:minute day-month-year

	// Parse the input string
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		return time.Time{}, err
	}

	// Validate the parsed time
	now := time.Now()

	// If the date is different, ensure it's not before today
	if parsedTime.Year() != now.Year() || parsedTime.Month() != now.Month() || parsedTime.Day() != now.Day() {
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if parsedTime.Before(today) {
			return time.Time{}, errors.New("time must not be earlier than today")
		}
	} else {
		// If the date is today, ensure it's not before the current time
		if parsedTime.Before(now) {
			return time.Time{}, errors.New("time must not be earlier than the current time")
		}
	}

	return parsedTime, nil
}
