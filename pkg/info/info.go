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

// LogDate represents the date of a log entry
type LogDate struct {
	Day      [2]byte
	Month    [2]byte
	Year     [4]byte
	DateTime string
}

// Info represents an entire log entry
type Info struct {
	IP            string
	Date          LogDate
	IsBot         bool
	IsUser        bool
	IsClientError bool
}

// InfoMap represents a map with IP as key and Info as value
type InfoMap struct {
	All map[string]Info
	Day map[string]Info
}

// GetInfoAtDay parses the contents of a given log at a given day
func GetInfoAtDay(infoMap InfoMap, str string, day string) (Info, error) {
	none := Info{IP: "0", IsBot: false, IsUser: false, IsClientError: false}
	aux := strings.SplitN(str, " ", 2)
	ip := aux[0]
	d := [2]byte{day[0], day[1]}
	m := [2]byte{day[3], day[4]}
	y := [4]byte{day[6], day[7], day[8], day[9]}
	maybeDate := DateRegex.FindString(aux[1])
	info, err := CreateInfo(ip, aux[1], maybeDate)

	if err != nil {
		return none, err
	}

	if CompareDayOrMonth(info.Date.Day, d) < 0 && CompareDayOrMonth(info.Date.Month, m) <= 0 && CompareYear(info.Date.Year, y) <= 0 {
		_, ok := infoMap.All[ip]

		if ok {
			return none, nil
		}

		if info.IP == "0" {
			return none, nil
		}

		infoMap.All[ip] = info

	} else if CompareDayOrMonth(info.Date.Day, d) == 0 && CompareDayOrMonth(info.Date.Month, m) == 0 && CompareYear(info.Date.Year, y) == 0 {
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
	const dateSize = 23

	if len(maybeDate) < dateSize {
		return info, fmt.Errorf("Invalid date")
	}

	date := LogDate{
		Day:      [2]byte{maybeDate[1], maybeDate[2]},
		Month:    StringToMonth(maybeDate[4:7]),
		Year:     [4]byte{maybeDate[8], maybeDate[9], maybeDate[10], maybeDate[11]},
		DateTime: maybeDate[1:21],
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

// StringToMonth transforms a valid string to the correspondent month
func StringToMonth(str string) [2]byte {
	switch str {
	case "Jan":
		return [2]byte{'0', '1'}
	case "Feb":
		return [2]byte{'0', '2'}
	case "Mar":
		return [2]byte{'0', '3'}
	case "Apr":
		return [2]byte{'0', '4'}
	case "May":
		return [2]byte{'0', '5'}
	case "Jun":
		return [2]byte{'0', '6'}
	case "Jul":
		return [2]byte{'0', '7'}
	case "Aug":
		return [2]byte{'0', '8'}
	case "Sep":
		return [2]byte{'0', '9'}
	case "Oct":
		return [2]byte{'1', '0'}
	case "Nov":
		return [2]byte{'1', '1'}
	case "Dec":
		return [2]byte{'1', '2'}
	}

	return [2]byte{'0', '0'}
}

//CompareDayOrMonth does the comparison of 2 days/months represented by 2 bytes
func CompareDayOrMonth(fst [2]byte, snd [2]byte) int {
	ret := 0

	if fst[0] > snd[0] {
		ret = 1
	} else if fst[0] < snd[0] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] < snd[1] {
		ret = -1
	} else if fst[0] == snd[0] && fst[1] > snd[1] {
		ret = 1
	}

	return ret
}

//CompareYear does the comparison of two years represented by 4 bytes
func CompareYear(fst [4]byte, snd [4]byte) int {
	ret := 0

	for i := 0; i < 4; i++ {
		if fst[i] > snd[i] {
			ret = 1
			break
		} else if fst[i] < snd[i] {
			ret = -1
			break
		}
	}

	return ret
}

// String method of Info
func (i Info) String() string {
	str := ""

	if i.IP == "0" {
		return str
	}

	if i.IsBot {
		str += i.Date.DateTime + ": Found a bot -> " + i.IP
	}
	if i.IsUser {
		str += i.Date.DateTime + ": Found a user -> " + i.IP
	}
	if i.IsClientError {
		str += i.Date.DateTime + ": Found a client error request -> " + i.IP
	}
	/**if !i.IsBot && !i.IsUser && !i.IsClientError {
		str += i.Date.String() + ": Nothing -> " + i.IP
	}*/
	return str
}
