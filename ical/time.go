package ical

import "time"

func floatToDuration(f float64) time.Duration {
	hours := int(f)
	minutes := int((f - float64(hours)) * 60)
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute
}