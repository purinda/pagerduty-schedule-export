package ical

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var CSVHeader = []string{"DT Start", "DT End", "Owner"}

// Format events to CSV format
func EventsToCSV(events []Event) string {
	builder := strings.Builder{}
	for _, event := range events {
		builder.WriteString(fmt.Sprintf("%s,%s,%s\n", event.DTStart, event.DTEnd, event.Attendee))
	}
	return builder.String()
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
