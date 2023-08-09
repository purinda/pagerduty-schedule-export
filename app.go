package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CSVHeader = "DTSTART,DTEND,ATTENDEE"
)

type Event struct {
	DTStart  string
	DTEnd    string
	Attendee string
}

func ParseICal(data string) ([]Event, error) {
	lines := strings.Split(data, "\n")
	var events []Event
	var currentEvent Event

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "BEGIN:VEVENT"):
			currentEvent = Event{}
		case strings.HasPrefix(line, "DTSTART"):
			currentEvent.DTStart = strings.Split(line, ":")[1]
		case strings.HasPrefix(line, "DTEND"):
			currentEvent.DTEnd = strings.Split(line, ":")[1]
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

func main() {
	url := flag.String("url", "https://example-domain.com/pd-sched.ical", "iCal URL")
	flag.Parse()

	data, err := FetchData(*url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	events, err := ParseICal(data)
	if err != nil {
		fmt.Printf("Error parsing iCal data: %v\n", err)
		return
	}

	csvData := EventsToCSV(events)
	fmt.Println(csvData)
}
