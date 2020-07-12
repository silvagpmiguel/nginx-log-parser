package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/silvagpmiguel/nginx-log-parser/pkg/info"
)

func getResultsFromDay(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool, dayFlag bool, day string) (string, error) {
	infoMap := info.InfoMap{
		All: make(map[string]info.Info),
		Day: make(map[string]info.Info),
	}
	scanner := bufio.NewScanner(file)
	bots := 0
	totalAccesses := 0
	users := 0
	clientErrors := 0
	str := ""
	existingUsers := 0
	noFlags := !botFlag && !dayFlag && !detailedFlag
	line := ""
	onDay := "on " + day

	for scanner.Scan() {
		_, err := info.GetInfoAtDay(infoMap, scanner.Text(), day)
		if err != nil {
			return "", err
		}
	}

	for _, v := range infoMap.Day {
		if v.IP == "0" {
			continue
		}

		if v.IsBot {
			bots++
			line = v.String()
			if verboseFlag && (botFlag || detailedFlag) && line != "" {
				fmt.Println(line)
			}
		}
		if v.IsUser {
			line = v.String()
			if verboseFlag && (noFlags || detailedFlag || dayFlag) && line != "" {
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
			if verboseFlag && (detailedFlag || dayFlag) && line != "" {
				fmt.Println(line)
			}
		}
		totalAccesses++
	}

	if verboseFlag {
		str += "\n"
	}

	if detailedFlag {
		str += fmt.Sprintf("Detailed Information %s\n", onDay)
		str += fmt.Sprintf("Number of unique bots: %d\n", bots)
		str += fmt.Sprintf("Number of new unique users: %d\n", users)
		str += fmt.Sprintf("Number of unique users who already had accessed the site: %d\n", existingUsers)
		str += fmt.Sprintf("Number of unique user error requests: %d\n", clientErrors)
		str += fmt.Sprintf("Total number of unique accesses: %d\n", totalAccesses)
		return str, nil
	}
	if botFlag {
		str += fmt.Sprintf("Found %d unique bots which accessed the site %s\n", bots, onDay)
	}
	if noFlags || dayFlag {
		str += fmt.Sprintf("Found %d unique users who accessed the site %s\n", users, onDay)
	}

	return str, nil
}

func getAllResults(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool, dayFlag bool, day string) (string, error) {
	infoMap := make(map[string]info.Info)
	scanner := bufio.NewScanner(file)
	bots := 0
	totalAccesses := 0
	users := 0
	clientErrors := 0
	str := ""
	noFlags := !botFlag && !dayFlag && !detailedFlag
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
			if verboseFlag && (botFlag || detailedFlag) && line != "" {
				fmt.Println(line)
			}
		}
		if v.IsUser {
			users++
			line = v.String()
			if verboseFlag && (noFlags || detailedFlag || dayFlag) && line != "" {
				fmt.Println(line)
			}
		}
		if v.IsClientError {
			clientErrors++
			line = v.String()
			if verboseFlag && (detailedFlag || dayFlag) && line != "" {
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
	if noFlags || dayFlag {
		str += fmt.Sprintf("Found %d unique users who accessed the site\n", users)
	}

	return str, nil
}

func readLog(file *os.File, botFlag bool, detailedFlag bool, verboseFlag bool, dayFlag bool, day string) (string, error) {
	if dayFlag {
		str, err := getResultsFromDay(file, botFlag, detailedFlag, verboseFlag, dayFlag, day)
		if err != nil {
			return str, err
		}
		return str, nil
	}

	str, err := getAllResults(file, botFlag, detailedFlag, verboseFlag, dayFlag, day)
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
	dateRegex := regexp.MustCompile(`[0-9]{2}/[0-9]{2}/[0-9]{4}`)
	botFlag := false
	detailedFlag := false
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

	str, err := readLog(file, botFlag, detailedFlag, verboseFlag, dayFlag, day)

	if err != nil {
		fmt.Printf("Error: Couldn't read log %s: %v\n", filepath, err)
		return
	}

	fmt.Printf(str)

}
