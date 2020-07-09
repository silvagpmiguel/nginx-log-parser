package main

import (
	"regexp"
	"strings"
	"time"
)

// ignore line contains:
// wp-admin

// DateRegex is a regex that catches the date from the log
var DateRegex *regexp.Regexp

func init() {
	DateRegex = regexp.MustCompile(`\[.* [+-][0-9]*\]`)
}

// Info represents a valid log entry
type Info struct {
	IP   string
	Date time.Time
}

// InfoMap represents an map with IP as key and Info as value
type InfoMap map[string]Info

// GetLogInfo parses the content of the log
func GetLogInfo(infoMap InfoMap, str string) error {
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	_, ok := infoMap[ip]

	if !ok && (strings.Contains(aux[0], ".") || strings.Contains(aux[0], ":")) {
		maybeDate := DateRegex.FindString(aux[1])
		if maybeDate != "" {
			date, err := CreateDate(maybeDate)

			if err != nil {
				return err
			}

			infoMap[ip] = Info{
				IP:   ip,
				Date: date,
			}
		}
	}

	return nil
}

// CreateDate returns a time.Time value from "02/Jan/2006:15:04:05 -0700" as layout
func CreateDate(str string) (time.Time, error) {
	layout := "[02/Jan/2006:15:04:05 -0700]"
	return time.Parse(layout, str)
}
