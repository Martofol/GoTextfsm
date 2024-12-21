package gotextfsm

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type ParserState int

const (
	VALUE_DECLERATION ParserState = iota
	START
	END
)

func (r ParserState) String() string {
	return [...]string{"DEFAULT", "VALUE_DECLERATİON", "START"}[r]
}

var parserCurrentState ParserState

func ParseTemplateFile(filePath string) ([]TextFSMValue, []StartPattern) {
	//we start with default state wich is equavelent to empty space
	ParserStateChanger(VALUE_DECLERATION)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open template file: %s", err)
	}
	defer file.Close()
	var fields []TextFSMValue
	var StartPatterns []StartPattern
	lineNum := 0
	scanner := bufio.NewScanner(file)
OuterLoop:
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			//We skip the comment and empty lines
			continue
		}
		switch ParserStateController() {
		case VALUE_DECLERATION:
			err := CheckTheLineForValueField(&fields, line, lineNum)
			if err != nil {
				log.Fatalf("\nError occured : %s ", err)
			}
		case START:
			err := CheckTheLineForPatterns(&StartPatterns, fields, line, lineNum)
			if err != nil {
				log.Fatalf("\nError occured : %s ", err)
			}
		case END:
			log.Println("States has been finished!!")
			break OuterLoop
		default:
			log.Fatalf("\nLine %d : State for parser Couldnt Found", lineNum)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading template file: %s", err)
	}

	return fields, StartPatterns
}

func ParserStateController() ParserState {
	return parserCurrentState
}

func ParserStateChanger(newState ParserState) {
	parserCurrentState = newState
}
