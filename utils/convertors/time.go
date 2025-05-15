package convertors

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func RoundDurationToNearestFive(d time.Duration) time.Duration {
	totalMinutes := int(d.Minutes())
	roundedMinutes := ((totalMinutes + 4) / 5) * 5

	return time.Duration(roundedMinutes) * time.Minute
}

func ConvertDateStringToTime(input string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format (%s): %v", input, err)
	}

	return parsedTime, nil
}

func ParseDurationToPGInterval(durationStr string) (pgtype.Interval, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return pgtype.Interval{}, fmt.Errorf("invalid duration format: %v", err)
	}

	return pgtype.Interval{
		Microseconds: duration.Microseconds(),
		Valid:        true,
	}, nil
}
