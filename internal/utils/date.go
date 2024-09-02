package utils

import "time"

func CurrentDate() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}

func CurrentDateString() string {
	return CurrentDate().Format("02-01-2006")
}
