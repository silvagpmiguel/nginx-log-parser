package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/silvagpmiguel/nginx-log-parser/pkg/info"
)

func getResultsFromDate(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool, day string, month string, year string) (string, error) {
	infoMap := info.Data{
		All:      make(map[string]info.Info),
		FromDate: make(map[string]info.Info),
	}
	scanner := bufio.NewScanner(file)
	bots := 0
	totalAccesses := 0
	users := 0
	clientErrors := 0
	str := ""
	dayFlag := day != ""
	monthFlag := month != ""
	yearFlag := year != ""
	existingUsers := 0
	noFlags := !botFlag && !detailedFlag
	line := ""
	onDate := ""
	for scanner.Scan() {
		if dayFlag {
			_, err := info.GetInfoAtDay(infoMap, scanner.Text(), day)
			if err != nil {
				return "", err
			}
			onDate = "on " + day
		} else if monthFlag {
			_, err := info.GetInfoAtMonth(infoMap, scanner.Text(), month)
			if err != nil {
				return "", err
			}
			onDate = "on " + month
		} else if yearFlag {
			_, err := info.GetInfoAtYear(infoMap, scanner.Text(), year)
			if err != nil {
				return "", err
			}
			onDate = "on " + year
		}
	}

	for _, v := range infoMap.FromDate {

		if v.IsBot {
			bots++
			line = v.String()
			if verboseFlag {
				fmt.Println(line)
			}
		}
		if v.IsUser {
			line = v.String()
			if verboseFlag {
				fmt.Println(line)
			}
			_, ok := infoMap.All[v.IP]
			if ok {
				existingUsers++
			} else {
				users++
			}
		}
		if v.IsClientError {
			clientErrors++
			line = v.String()
			if verboseFlag {
				fmt.Println(line)
			}
		}
		totalAccesses++
	}

	if verboseFlag {
		str += "\n"
	}

	if detailedFlag {
		str += fmt.Sprintf("Detailed Information %s\n", onDate)
		str += fmt.Sprintf("Number of unique bots: %d\n", bots)
		str += fmt.Sprintf("Number of new unique users: %d\n", users)
		str += fmt.Sprintf("Number of unique users who already had accessed the site: %d\n", existingUsers)
		str += fmt.Sprintf("Number of unique user error requests: %d\n", clientErrors)
		str += fmt.Sprintf("Total number of unique accesses: %d\n", totalAccesses)
		return str, nil
	}
	if botFlag {
		str += fmt.Sprintf("Found %d unique bots which accessed the site %s\n", bots, onDate)
	}
	if noFlags {
		str += fmt.Sprintf("Found %d unique users who accessed the site %s\n", users, onDate)
	}

	return str, nil
}

func getAllResults(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool) (string, error) {
	infoMap := make(map[string]info.Info)
	scanner := bufio.NewScanner(file)
	bots := 0
	totalAccesses := 0
	users := 0
	clientErrors := 0
	str := ""
	noFlags := !botFlag && !detailedFlag
	line := ""

	for scanner.Scan() {
		v, err := info.GetAllInfo(infoMap, scanner.Text())
		if err != nil {
			return "", err
		}

		if v.IP == "0" {
			continue
		}

		if v.IsBot {
			bots++
			line = v.String()
			if verboseFlag && line != "" {
				fmt.Println(line)
			}
		}
		if v.IsUser {
			users++
			line = v.String()
			if verboseFlag && line != "" {
				fmt.Println(line)
			}
		}
		if v.IsClientError {
			clientErrors++
			line = v.String()
			if verboseFlag && line != "" {
				fmt.Println(line)
			}
		}
		totalAccesses++
	}

	if verboseFlag {
		str += "\n"
	}

	if detailedFlag {
		str += fmt.Sprintf("Detailed Information\n")
		str += fmt.Sprintf("Number of unique bots: %d\n", bots)
		str += fmt.Sprintf("Number of unique users: %d\n", users)
		str += fmt.Sprintf("Number of unique user error requests: %d\n", clientErrors)
		str += fmt.Sprintf("Total number of unique accesses: %d\n", totalAccesses)
		return str, nil
	}
	if botFlag {
		str += fmt.Sprintf("Found %d unique bots which accessed the site\n", bots)
	}
	if noFlags {
		str += fmt.Sprintf("Found %d unique users who accessed the site\n", users)
	}

	return str, nil
}

func readLog(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool, day string, month string, year string) (string, error) {
	if day != "" || month != "" || year != "" {
		str, err := getResultsFromDate(file, botFlag, detailedFlag, verboseFlag, day, month, year)
		if err != nil {
			return str, err
		}
		return str, nil
	}

	str, err := getAllResults(file, botFlag, detailedFlag, verboseFlag)
	if err != nil {
		return str, err
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
		"\t-bots,\t\tDisplay the number of bots which accessed the website\n",
		"\t-day,\t\tDisplay the number of users who accessed the website on <dd/mm/yyyy>\n",
		"\t-month,\t\tDisplay the number of users who accessed the website on <mm/yyyy>\n",
		"\t-year,\t\tDisplay the number of users who accessed the website on <yyyy>\n",
		"\t-detailed,\tDisplay more detailed information\n",
		"\t-verbose,\tDisplay information about each line of the log\n",
		"\t-h,\t\tDisplay this help and exit\n",
		"EXAMPLE\n",
		"\t./nginx-log-parser $LOGPATH,\t\t\t\tRead access log at $LOGPATH and display the total number of users that accessed the website\n",
		"\t./nginx-log-parser -day 09/07/2020 $LOGPATH,\t\tRead access log at $LOGPATH and display the users that accessed the website on that day\n",
		"\t./nginx-log-parser -detailed -day 09/07/2020 $LOGPATH,\tRead access log at $LOGPATH and display more detailed information on that day\n",
	)
}

func main() {
	filepath := ""
	argsLen := len(os.Args)
	dayRegex := regexp.MustCompile(`(3[0-1]|[1-2][0-9]|0[1-9])/(1[0-2]|0[1-9])/[0-9]{4}`)
	monthRegex := regexp.MustCompile(`1[0-2]|0[1-9]/[0-9]{4}`)
	yearRegex := regexp.MustCompile(`[0-9]{4}`)
	botFlag := false
	detailedFlag := false
	verboseFlag := false
	day := ""
	month := ""
	year := ""

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
				if !dayRegex.MatchString(os.Args[i+1]) {
					fmt.Println("Error. Insert a valid date <dd/mm/yyyy>")
					return
				}
				day = os.Args[i+1]
				i++
			}
		case "-month":
			if i+1 < argsLen {
				if !monthRegex.MatchString(os.Args[i+1]) {
					fmt.Println("Error. Insert a valid date <mm/yyyy>")
					return
				}
				month = os.Args[i+1]
				i++
			}
		case "-year":
			if i+1 < argsLen {
				if !yearRegex.MatchString(os.Args[i+1]) {
					fmt.Println("Error. Insert a valid date <yyyy>")
					return
				}
				year = os.Args[i+1]
				i++
			}
		case "-bots":
			botFlag = true
		case "-detailed":
			detailedFlag = true
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

	str, err := readLog(file, botFlag, detailedFlag, verboseFlag, day, month, year)

	if err != nil {
		fmt.Printf("Error: Couldn't read log %s: %v\n", filepath, err)
		return
	}

	fmt.Printf(str)

}
