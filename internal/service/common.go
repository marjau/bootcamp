package service

import "time"

// dateTimeNow returns the current UTC time in the form of date time with zero nanoseconds.
// NOTE: this is a work around to make the mock testing work, otherwise mocks dates will differ on nanoseconds values.
func dateTimeNow() time.Time {
	t := time.Now().UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
}
