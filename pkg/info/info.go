package info

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
	IP            string
	Date          time.Time
	IsBot         bool
	IsNewUser     bool
	IsClientError bool
}

// InfoMap represents an map with IP as key and Info as value
type InfoMap map[string]Info

// GetLogInfo parses the content of the log
func GetLogInfo(infoMap InfoMap, str string) (Info, error) {
	info := Info{IP: "0", IsBot: false, IsNewUser: false, IsClientError: false}

	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	_, ok := infoMap[ip]

	if !ok && (strings.Contains(aux[0], ".") || strings.Contains(aux[0], ":")) {
		maybeDate := DateRegex.FindString(aux[1])

		if maybeDate != "" {
			date, err := CreateDate(maybeDate)

			if err != nil {
				return info, err
			}

			botFlag := strings.Contains(aux[1], "wp-admin")
			newFlag := strings.Contains(aux[1], "HTTP/1.1\" 200")
			clientError := strings.Contains(aux[1], "HTTP/1.1\" 4")
			info = Info{
				IP:            ip,
				Date:          date,
				IsBot:         botFlag,
				IsNewUser:     newFlag,
				IsClientError: clientError,
			}

			infoMap[ip] = info
		}
	}
	return info, nil
}

// CreateDate returns a time.Time value from "02/Jan/2006:15:04:05 -0700" as layout
func CreateDate(str string) (time.Time, error) {
	layout := "[02/Jan/2006:15:04:05 -0700]"
	return time.Parse(layout, str)
}

func (i Info) String() string {
	str := ""

	if i.IP == "0" {
		return str
	}

	if i.IsBot {
		str += i.Date.String() + ": Found a bot -> " + i.IP
	} else {
		if i.IsNewUser {
			str += i.Date.String() + ": Found a new user -> " + i.IP
		} else {
			if i.IsClientError {
				str += i.Date.String() + ": Found a client error request -> " + i.IP
			} else {
				str += i.Date.String() + ": Found a user who had already accessed this site -> " + i.IP
			}
		}
	}
	return str
}
