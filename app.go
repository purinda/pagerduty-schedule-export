package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Event struct {
	DTStart  time.Time
	DTEnd    time.Time
	Attendee string
}

var CSVHeader = []string{"DT Start", "DT End", "Owner"}

func main() {
	// Command-line argument processing
	icalUrl := flag.String("url", "", "URL to fetch the iCal data from.")
	timezone := flag.String("timezone", "UTC", "Timezone for converting the iCal dates.")
	flag.Parse()

	if *icalUrl == "" {
		fmt.Println("Please provide a URL using the -url flag.")
		return
	}

	// Fetching iCal data from the URL
	data, err := FetchData(*icalUrl)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	// Parsing the iCal data
	location, _ := time.LoadLocation(*timezone)
	events, err := ParseICal(data, location)
	if err != nil {
		fmt.Printf("Error parsing iCal data: %v\n", err)
		return
	}

	// Determine the filename for the output CSV
	parsedURL, err := url.Parse(*icalUrl)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return
	}
	filenameSegment := path.Base(parsedURL.Path)
	csvFilename := filenameSegment + ".csv"

	// Writing the parsed data to the CSV
	err = WriteToCSV(events, csvFilename)
	if err != nil {
		fmt.Printf("Error writing to CSV: %v\n", err)
		return
	}
	fmt.Println("Schedule exported to:", csvFilename)
}

// Parse the iCal content provided via 'data' variable and transform
// datetime entries to required TZ format.
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

// Format events to CSV format
func EventsToCSV(events []Event) string {
	builder := strings.Builder{}
	for _, event := range events {
		builder.WriteString(fmt.Sprintf("%s,%s,%s\n", event.DTStart, event.DTEnd, event.Attendee))
	}
	return builder.String()
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

// Export the events extracted from the iCal as a CSV file
func WriteToCSV(events []Event, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// Ensure every field is comma separated and enclosed with quotes
	writer.Comma = ','    // Field delimiter
	writer.UseCRLF = true // Use \n as line terminator

	if err := writer.Write(CSVHeader); err != nil {
		return err
	}

	for _, event := range events {
		if err := writer.Write([]string{event.DTStart.String(), event.DTEnd.String(), event.Attendee}); err != nil {
			return err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return err
	}

	return nil
}
