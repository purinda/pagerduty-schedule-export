# PagerDuty Schedule Exporter
This solves one problem and one problem only, which is to take an export of PD schedule start, end times along with the owner.

# Usage

Download a release as per your system and run `--help` to view usage information.

```
./pd-export --help

Usage of ./pd-export:
  -timezone string
    	Target timezone for date conversions (e.g., 'Australia/Sydney' for AEST) (default "UTC")
  -url string
    	URL to fetch iCal data (default "https://example.com/path-to-ical-file.ics")
```

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