# PagerDuty Schedule Exporter
This solves one problem and one problem only, which is to take an export of PD schedule start, end times along with the owner.

# Usage

Download [latest release here](https://github.com/purinda/pagerduty-schedule-export/releases/latest) as per your system and run `--help` to view usage information.

```
./pd-export --help

Usage of ./pd-export:
  -timezone string
    	Target timezone for date conversions (e.g., 'Australia/Sydney' for AEST) (default "UTC")
  -url string
    	URL to fetch iCal data (default "https://example.com/path-to-ical-file.ics")
```

### How to Retrieve the PagerDuty iCal URL

Navigate to the PD schedule page and copy the URL as per screenshot.

![Schedule Page](docs/pd-schedule.png)

## Building

Prerequisites - Go (at least version 1.16) installed on your computer.

Compile the application:
```sh
go build -o app .
```

Run the compiled binary:
```sh
./app [arguments]
```

## Testing
To run the unit tests for the application

```sh
go test
```