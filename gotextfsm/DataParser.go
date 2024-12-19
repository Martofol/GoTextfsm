package gotextfsm

import (
	"bufio"
	"log"
	"os"
)

type Record struct {
	Name  string
	Value string
}

func ParseCLIOutput(filePath string, fields []TextFSMValue, startPatterns []StartPattern) []Record {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CLI output file: %s", err)
	}
	defer file.Close()

	var records []Record

	// Parse CLI output
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, pattern := range startPatterns {
			if pattern.Regex.MatchString(line) {
				match := pattern.Regex.FindStringSubmatch(line)
				for i, name := range pattern.Regex.SubexpNames() {
					if i != 0 && name != "" {
						records = append(records, Record{Name: name, Value: match[i]})
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading CLI output file: %s", err)
	}
	return records
}
