package info

import (
	"fmt"
	"regexp"
	"strings"
)

// ignore line contains:
// wp-admin

// DateRegex is a regex that catches the date from the log
var DateRegex *regexp.Regexp

func init() {
	DateRegex = regexp.MustCompile(`\[.* [+-][0-9]*\]`)
}

// Info represents an entire log entry
type Info struct {
	IP            string
	Date          Date
	IsBot         bool
	IsUser        bool
	IsClientError bool
}

// Data represents all the access log data organized into two maps with IP as key and Info as value
type Data struct {
	All      map[string]Info
	FromDate map[string]Info
}

// GetInfoAtDay parses the contents of a given log at a given day
func GetInfoAtDay(infoMap Data, str string, day string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	d := [2]byte{day[0], day[1]}
	m := [2]byte{day[3], day[4]}
	y := [4]byte{day[6], day[7], day[8], day[9]}
	maybeDate := DateRegex.FindString(aux[1])

	if len(maybeDate) < 23 {
		return none, fmt.Errorf("Invalid date")
	}

	info := CreateInfo(ip, aux[1], maybeDate)

	if info.IP == "0" {
		return none, nil
	}

	if info.Date.CompareDay(d) < 0 && info.Date.CompareMonth(m) <= 0 && info.Date.CompareYear(y) <= 0 {
		_, ok := infoMap.All[ip]
		if ok {
			return none, nil
		}
		infoMap.All[ip] = info

	} else if info.Date.CompareDay(d) == 0 && info.Date.CompareMonth(m) == 0 && info.Date.CompareYear(y) == 0 {
		_, ok := infoMap.FromDate[ip]
		if ok {
			return none, nil
		}
		infoMap.FromDate[ip] = info
	}

	return info, nil
}

// GetInfoAtMonth parses the contents of a given log at a given month
func GetInfoAtMonth(infoMap Data, str string, month string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	m := [2]byte{month[0], month[1]}
	y := [4]byte{month[3], month[4], month[5], month[6]}
	maybeDate := DateRegex.FindString(aux[1])

	if len(maybeDate) < 23 {
		return none, fmt.Errorf("Invalid date")
	}

	info := CreateInfo(ip, aux[1], maybeDate)

	if info.IP == "0" {
		return none, nil
	}

	if info.Date.CompareMonth(m) < 0 && info.Date.CompareYear(y) <= 0 {
		_, ok := infoMap.All[ip]
		if ok {
			return none, nil
		}
		infoMap.All[ip] = info
	} else if info.Date.CompareMonth(m) == 0 && info.Date.CompareYear(y) == 0 {
		_, ok := infoMap.FromDate[ip]
		if ok {
			return none, nil
		}
		infoMap.FromDate[ip] = info
	}

	return info, nil
}

// GetInfoAtYear parses the contents of a given log at a given year
func GetInfoAtYear(infoMap Data, str string, year string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	y := [4]byte{year[0], year[1], year[2], year[3]}
	maybeDate := DateRegex.FindString(aux[1])

	if len(maybeDate) < 23 {
		return none, fmt.Errorf("Invalid date")
	}

	info := CreateInfo(ip, aux[1], maybeDate)

	if info.IP == "0" {
		return none, nil
	}

	if info.Date.CompareYear(y) < 0 {
		_, ok := infoMap.All[ip]

		if ok {
			return none, nil
		}

		infoMap.All[ip] = info

	} else if info.Date.CompareYear(y) == 0 {
		_, ok := infoMap.FromDate[ip]

		if ok {
			return none, nil
		}

		infoMap.FromDate[ip] = info
	}

	return info, nil
}

// GetAllInfo parses all contents of a given log
func GetAllInfo(allMap map[string]Info, str string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	_, ok := allMap[ip]

	if ok {
		return none, nil
	}

	maybeDate := DateRegex.FindString(aux[1])

	if len(maybeDate) < 23 {
		return none, fmt.Errorf("Invalid date")
	}

	info := CreateInfo(ip, aux[1], maybeDate)

	if info.IP == "0" {
		return none, nil
	}

	allMap[ip] = info

	return info, nil
}

// CreateInfo returns a Info type or error
func CreateInfo(ip string, str string, maybeDate string) Info {
	info := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}

	date := Date{
		Day:      [2]byte{maybeDate[1], maybeDate[2]},
		Month:    StringToMonth(maybeDate[4:7]),
		Year:     [4]byte{maybeDate[8], maybeDate[9], maybeDate[10], maybeDate[11]},
		DateTime: maybeDate[1:21],
	}

	botFlag := strings.Contains(str, "wp-admin")
	clientError := strings.Contains(str, "HTTP/1.1\" 4")
	userFlag := strings.Contains(str, "assets")

	if !botFlag && !userFlag && !clientError {
		return info
	}

	info = Info{
		IP:            ip,
		Date:          date,
		IsBot:         botFlag,
		IsUser:        userFlag,
		IsClientError: clientError,
	}

	return info
}

// String method of Info
func (i Info) String() string {
	str := ""

	if i.IP == "0" {
		return str
	}
	if i.IsBot {
		str += i.Date.DateTime + ": Found a bot -> " + i.IP
	} else if i.IsUser {
		str += i.Date.DateTime + ": Found a user -> " + i.IP
	} else if i.IsClientError {
		str += i.Date.DateTime + ": Found a client error request -> " + i.IP
	}

	return str
}
