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

// GetInfoAtDay parses the contents of a given log at a given day
func GetInfoAtDay(infoMap InfoMap, str string, day string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	layout := "02/01/2006"
	timeDay, _ := time.Parse(layout, day)
	d := timeDay.Day()
	m := int(timeDay.Month())
	y := timeDay.Year()
	maybeDate := DateRegex.FindString(str)

	info, err := CreateInfo(ip, aux[1], maybeDate)

	if err != nil {
		return none, err
	}

	if info.Date.Day() < d && m <= int(info.Date.Month()) && y <= info.Date.Year() {
		_, ok := infoMap.All[ip]

		if ok {
			return none, nil
		}

		if info.IP == "0" {
			return none, nil
		}

		infoMap.All[ip] = info

	} else if info.Date.Day() == d && m == int(info.Date.Month()) && y == info.Date.Year() {
		_, ok := infoMap.Day[ip]

		if ok {
			return none, nil
		}

		if info.IP == "0" {
			return none, nil
		}

		infoMap.Day[ip] = info
	}

	return info, nil
}

// GetAllInfo parses all contents of a log
func GetAllInfo(allMap map[string]Info, str string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	_, ok := allMap[ip]

	if ok {
		return none, nil
	}

	maybeDate := DateRegex.FindString(str)
	info, err := CreateInfo(ip, aux[1], maybeDate)

	if err != nil {
		return none, err
	}

	if info.IP == "0" {
		return none, nil
	}

	allMap[ip] = info

	return info, nil
}

// CreateInfo returns a Info type or error
func CreateInfo(ip string, str string, maybeDate string) (Info, error) {
	info := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	layout := "[02/Jan/2006:15:04:05 -0700]"

	date, err := time.Parse(layout, maybeDate)

	if err != nil {
		return info, err
	}

	botFlag := strings.Contains(str, "wp-admin")
	clientError := strings.Contains(str, "HTTP/1.1\" 4")
	userFlag := strings.Contains(str, "assets")

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

	return info, nil
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
