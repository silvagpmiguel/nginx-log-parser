package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/silvagpmiguel/nginx-log-parser/pkg/info"
)

// ReadLog is a function that reads a specific log (File) given in the argument
func ReadLog(file *os.File) (info.InfoMap, error) {
	infoMap := make(info.InfoMap)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		err := info.GetLogInfo(infoMap, scanner.Text())
		if err != nil {
			return infoMap, err
		}
	}

	if len(infoMap) == 0 {
		return nil, fmt.Errorf("Invalid log fields")
	}

	return infoMap, nil
}

func printCommandsInfo() {
	fmt.Println(
		"Nginx Log Parser Usage\n\n",
		"Usage: ./nginx-log-parser [OPTION]... $LOGPATH\n\n",
		"Reads information from a nginx access log located at $LOGPATH\n\n",
		"Provide LOGPATH as last argument\n\n",
		"OPTIONS\n",
		"\t-day,\t\tDisplay the number of users that accessed the website in a day <dd/mm/yyyy>\n",
		"\t-bots,\t\tDisplay the number of bots that accessed the website\n",
		"\t-detailed,\tDisplay more detailed information\n",
		"\t-h,\t\tDisplay this help and exit\n",
		"EXAMPLE\n",
		"\t./nginx-log-parser $LOGPATH,\tRead access log at $LOGPATH and display the total number of users that accessed a specific website",
	)
}

func main() {
	filepath := ""
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-h":
			printCommandsInfo()
			return
		case "-day":
			return
		case "-bots":
			return
		case "-detailed":
			return
		}
		if strings.Contains(os.Args[i], ".log") {
			filepath = os.Args[i]
		}
	}

	if filepath == "" {
		fmt.Println("Error: You must give a valid .log file")
		return
	}

	file, err := os.Open(filepath)

	if err != nil {
		fmt.Printf("Error: Couldn't open file: %v\n", err)
		return
	}

	infoMap, err := ReadLog(file)

	if err != nil {
		fmt.Printf("Error: Couldn't read log %s: %v\n", filepath, err)
		return
	}

	fmt.Printf("Found %d users that visited your site!\n", len(infoMap))

}
