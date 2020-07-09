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

func main() {
	filepath := ""
	for i := 1; i < len(os.Args); i++ {
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
