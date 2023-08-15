package ical

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	DTStart  time.Time
	DTEnd    time.Time
	Attendee string
}

// Retrieve iCal file from the provided URL
func FetchData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// Parse the iCal content provided via 'data' variable and transform
// datetime entries to required TZ format as an offset.
func ParseICal(data string, offset time.Duration) ([]Event, error) {
	lines := strings.Split(data, "\n")
	var events []Event
	var currentEvent Event
	const iCalLayout = "20060102T150405Z"

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "BEGIN:VEVENT"):
			currentEvent = Event{}
		case strings.HasPrefix(line, "DTSTART"):
			t, err := time.Parse(iCalLayout, strings.Split(line, ":")[1])
			if err != nil {
				return nil, err
			}

			currentEvent.DTStart = t.Add(offset)
		case strings.HasPrefix(line, "DTEND"):
			t, err := time.Parse(iCalLayout, strings.Split(line, ":")[1])
			if err != nil {
				return nil, err
			}
			currentEvent.DTEnd = t.Add(offset)
		case strings.HasPrefix(line, "ATTENDEE"):
			currentEvent.Attendee = strings.Split(line, ":")[1]
		case strings.HasPrefix(line, "END:VEVENT"):
			events = append(events, currentEvent)
		}
	}

	return events, nil
}
