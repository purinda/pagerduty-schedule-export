package main

import (
	"testing"
	"time"
)

func TestParseICal(t *testing.T) {
	data := `BEGIN:VEVENT
DTEND;VALUE=DATE-TIME:20230709T220000Z
DTSTART;VALUE=DATE-TIME:20230708T220000Z
ATTENDEE:hello.kitty@testdomain.com
UID:Q1S46N8B8ADP9D
URL:https://yourorg.pagerduty.com/schedules#KXQAOAO
SUMMARY:On Call - Hello Kitty - Alpha
END:VEVENT`

	location, _ := time.LoadLocation("UTC")
	events, err := ParseICal(data, location)
	if err != nil {
		t.Errorf("Failed to parse iCal: %v", err)
	}

	expectedStartTime := time.Date(2023, 7, 8, 22, 0, 0, 0, location)
	expectedEndTime := time.Date(2023, 7, 9, 22, 0, 0, 0, location)
	expectedAttendee := "hello.kitty@testdomain.com"

	if !events[0].DTStart.Equal(expectedStartTime) {
		t.Errorf("Expected %v, got %v", expectedStartTime, events[0].DTStart)
	}
	if !events[0].DTEnd.Equal(expectedEndTime) {
		t.Errorf("Expected %v, got %v", expectedEndTime, events[0].DTEnd)
	}
	if events[0].Attendee != expectedAttendee {
		t.Errorf("Expected %s, got %s", expectedAttendee, events[0].Attendee)
	}
}
