package utils

import "time"

func CurrentDate() time.Time {
	return time.Now().Truncate(24 * time.Hour).UTC()
}

func CurrentDateString() string {
	return CurrentDate().Format("02-01-2006")
}

func StringToTime(date_str string) (time.Time, error) {
	return time.Parse("02-01-2006", date_str)
}
