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
	IsUser        bool
	IsClientError bool
}

// InfoMap represents an map with IP as key and Info as value
type InfoMap struct {
	All map[string]Info
	Day map[string]Info
}

// GetLogInfo parses the content of the log
func GetLogInfo(infoMap InfoMap, str string, day string) (Info, error) {
	info := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	botFlag := false
	clientError := false
	layout := ""
	timeDay := time.Now()
	d := 0
	m := 0
	y := 0

	maybeDate := DateRegex.FindString(aux[1])

	if maybeDate != "" {
		date, err := CreateDate(maybeDate)

		if err != nil {
			return info, err
		}

		botFlag = strings.Contains(aux[1], "wp-admin")
		clientError = strings.Contains(aux[1], "HTTP/1.1\" 4")
		userFlag := strings.Contains(aux[1], "assets")

		if !botFlag && !userFlag && !clientError {
			return info, nil
		}

		info = Info{
			IP:            ip,
			Date:          date,
			IsBot:         botFlag,
			IsUser:        userFlag,
			IsClientError: clientError,
		}

		if day != "" {
			layout = "02/01/2006"
			timeDay, _ = time.Parse(layout, day)
			d = timeDay.Day()
			m = int(timeDay.Month())
			y = timeDay.Year()

			if date.Day() < d && m <= int(date.Month()) && y <= date.Year() {
				_, ok := infoMap.All[ip]

				if ok {
					return info, nil
				}

				infoMap.All[ip] = info
			} else if date.Day() == d && m == int(date.Month()) && y == date.Year() {
				_, ok := infoMap.Day[ip]

				if ok {
					return info, nil
				}

				infoMap.Day[ip] = info
			} else {
				return info, nil
			}
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
	}
	if i.IsUser {
		str += i.Date.String() + ": Found a user -> " + i.IP
	}
	if i.IsClientError {
		str += i.Date.String() + ": Found a client error request -> " + i.IP
	}
	/**if !i.IsBot && !i.IsUser && !i.IsClientError {
		str += i.Date.String() + ": Nothing -> " + i.IP
	}*/
	return str
}
