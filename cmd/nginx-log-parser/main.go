package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/silvagpmiguel/nginx-log-parser/pkg/info"
)

func readLog(file *os.File, newFlag bool, oldFlag bool, botFlag bool, detailedFlag bool, verboseFlag bool, dayFlag bool, day string) (string, error) {
	infoMap := make(info.InfoMap)
	scanner := bufio.NewScanner(file)
	bots := 0
	totalViews := 0
	newcomers := 0
	existingUsers := 0
	clientErrors := 0
	str := ""
	noFlags := !botFlag && !dayFlag && !newFlag && !oldFlag && !detailedFlag
	onDay := ""

	if dayFlag {
		onDay = "on " + day
	}

	for scanner.Scan() {
		info, err := info.GetLogInfo(infoMap, scanner.Text())
		if err != nil {
			return "", err
		}

		if info.IP == "0" {
			continue
		}

		if dayFlag {
			layout := "02/01/2006:15:04:05 -0700"
			startDay, _ := time.Parse(layout, day+":00:00:00 +0000")
			endDay, _ := time.Parse(layout, day+":23:59:59 +0000")

			if !(info.Date.After(startDay) && info.Date.Before(endDay)) {
				continue
			}
		}

		if info.IsBot {
			bots++
			if verboseFlag && (botFlag || detailedFlag) {
				fmt.Println(info.String())
			}
		} else {
			if info.IsNewUser {
				newcomers++
				if verboseFlag && (newFlag || detailedFlag || dayFlag) {
					fmt.Println(info.String())
				}
			} else {
				if info.IsClientError {
					clientErrors++
				} else {
					existingUsers++
					if verboseFlag && (oldFlag || detailedFlag || dayFlag) {
						fmt.Println(info.String())
					}
				}
			}
			totalViews++
			if verboseFlag && (noFlags) {
				fmt.Println(info.String())
			}
		}
	}

	if verboseFlag {
		str += "\n"
	}

	if totalViews == 0 {
		return str, fmt.Errorf("Invalid log fields")
	}

	if detailedFlag {
		str += fmt.Sprintf("Detailed Information %s\n", onDay)
		str += fmt.Sprintf("Number of bots: %d\n", bots)
		str += fmt.Sprintf("Number of new users who accessed the site: %d\n", newcomers)
		str += fmt.Sprintf("Number of users who already had accessed the site: %d\n", existingUsers)
		str += fmt.Sprintf("Number of user error requests: %d\n", clientErrors)
		str += fmt.Sprintf("Total number of user views: %d\n", totalViews)
		return str, nil
	}
	if botFlag {
		str += fmt.Sprintf("Found %d bots who accessed the site %s\n", bots, onDay)
	}
	if newFlag {
		str += fmt.Sprintf("Found %d new users who accessed the site %s\n", totalViews, onDay)
	}
	if oldFlag {
		str += fmt.Sprintf("Found %d users who already had accessed the site %s\n", existingUsers, onDay)
	}
	if noFlags {
		str += fmt.Sprintf("Found %d users who accessed the site %s\n", totalViews, onDay)
	}
	if dayFlag {
		str += fmt.Sprintf("Found %d users who accessed the site on %s\n", totalViews, day)
	}

	return str, nil
}

func printCommandsInfo() {
	fmt.Println(
		"Nginx Log Parser Usage\n\n",
		"Usage: ./nginx-log-parser [OPTION]... $LOGPATH\n\n",
		"Reads information from a nginx access log located at $LOGPATH\n\n",
		"Provide LOGPATH as last argument\n\n",
		"OPTIONS\n",
		"\t-bots,\t\tDisplay the number of bots who accessed the website\n",
		"\t-day,\t\tDisplay the number of users who accessed the website in a day <dd/mm/yyyy>\n",
		"\t-detailed,\tDisplay more detailed information\n",
		"\t-new,\t\tDisplay the number of new users who accessed the website\n",
		"\t-old,\t\tDisplay the number of users who had already accessed the website\n",
		"\t-verbose,\tDisplay information about each line of the log\n",
		"\t-h,\t\tDisplay this help and exit\n",
		"EXAMPLE\n",
		"\t./nginx-log-parser $LOGPATH,\tRead access log at $LOGPATH and display the total number of users that accessed the website",
		"\t./nginx-log-parser -day 09/07/2020 $LOGPATH,\tRead access log at $LOGPATH and display the number of users that accessed the website at that day",
		"\t./nginx-log-parser -day 09/07/2020 $LOGPATH,\tRead access log at $LOGPATH and display the number of users that accessed the website at that day",
	)
}

func main() {
	filepath := ""
	argsLen := len(os.Args)
	dateRegex := regexp.MustCompile(`[0-9]{2}/[0-9]{2}/[0-9]{4}`)
	botFlag := false
	detailedFlag := false
	newFlag := false
	oldFlag := false
	verboseFlag := false
	dayFlag := false
	day := ""

	if argsLen == 1 {
		fmt.Println("Error. No arguments found\nWrite \"./nginx-log-parser -h\" to display the help")
		return
	}

	for i := 1; i < argsLen; i++ {
		switch os.Args[i] {
		case "-h":
			printCommandsInfo()
			return
		case "-day":
			if i+1 < argsLen {
				if !dateRegex.MatchString(os.Args[i+1]) {
					fmt.Println("Error. Insert a valid date <dd/mm/yyyy>")
					return
				}
				dayFlag = true
				day = os.Args[i+1]
				i++
			}
		case "-bots":
			botFlag = true
		case "-detailed":
			detailedFlag = true
		case "-old":
			oldFlag = true
		case "-new":
			newFlag = true
		case "-verbose":
			verboseFlag = true
		default:
			if strings.Contains(os.Args[i], ".log") {
				filepath = os.Args[i]
			} else {
				fmt.Println("Error. Wrong arguments")
				fmt.Println("Write \"./nginx-log-parser -h\" to display the help")
				return
			}
		}
	}

	file, err := os.Open(filepath)

	if err != nil {
		fmt.Printf("Error: Couldn't open file: %v\n", err)
		return
	}

	str, err := readLog(file, newFlag, oldFlag, botFlag, detailedFlag, verboseFlag, dayFlag, day)

	if err != nil {
		fmt.Printf("Error: Couldn't read log %s: %v\n", filepath, err)
		return
	}

	fmt.Printf(str)

}
