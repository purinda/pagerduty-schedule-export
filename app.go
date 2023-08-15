package main

import (
	"flag"
	"fmt"
	"net/url"
	"path"

	ical "github.com/purinda/pagerduty-schedule-export/pkg"
)

func main() {
	// Command-line argument processing
	icalUrl := flag.String("url", "", "URL to fetch the iCal data from.")
	tzHours := flag.Float64("timezone", 0, "Timezone offset in hours from UTC.")
	flag.Parse()

	if *icalUrl == "" {
		fmt.Println("Please provide a URL using the -url flag.")
		return
	}

	// Fetching iCal data from the URL
	data, err := ical.FetchData(*icalUrl)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	// Parse the iCal data
	events, err := ical.ParseICal(data, ical.FloatToDuration(*tzHours))
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
	err = ical.WriteToCSV(events, csvFilename)
	if err != nil {
		fmt.Printf("Error writing to CSV: %v\n", err)
		return
	}
	fmt.Println("Schedule exported to:", csvFilename)
}
