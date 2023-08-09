package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	CSVHeader = "DTSTART,DTEND,ATTENDEE"
)

type Event struct {
	DTStart  time.Time
	DTEnd    time.Time
	Attendee string
}

func ParseICal(data string, location *time.Location) ([]Event, error) {
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
			currentEvent.DTStart = t.In(location)
		case strings.HasPrefix(line, "DTEND"):
			t, err := time.Parse(iCalLayout, strings.Split(line, ":")[1])
			if err != nil {
				return nil, err
			}
			currentEvent.DTEnd = t.In(location)
		case strings.HasPrefix(line, "ATTENDEE"):
			currentEvent.Attendee = strings.Split(line, ":")[1]
		case strings.HasPrefix(line, "END:VEVENT"):
			events = append(events, currentEvent)
		}
	}

	return events, nil
}

func EventsToCSV(events []Event) string {
	builder := strings.Builder{}
	builder.WriteString(CSVHeader + "\n")
	for _, event := range events {
		builder.WriteString(fmt.Sprintf("%s,%s,%s\n", event.DTStart, event.DTEnd, event.Attendee))
	}
	return builder.String()
}

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

func main() {
	url := flag.String("url", "https://example.com/path-to-ical-file.ics", "URL to fetch iCal data")
	tz := flag.String("timezone", "UTC", "Target timezone for date conversions (e.g., 'Australia/Sydney' for AEST)")
	flag.Parse()

	location, err := time.LoadLocation(*tz)
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		return
	}

	data, err := FetchData(*url)
	if err != nil {
		fmt.Printf("Error fetching iCal: %v\n", err)
		return
	}

	events, err := ParseICal(data, location)
	if err != nil {
		fmt.Printf("Error parsing iCal data: %v\n", err)
		return
	}

	csvData := EventsToCSV(events)
	fmt.Println(csvData)
}
