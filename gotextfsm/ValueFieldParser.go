package gotextfsm

import (
	"fmt"
	"log"
	"strings"
)

var emptyLineCount int

func CheckTheLineForValueField(fields *[]TextFSMValue, line string, lineNum int) error {
	if strings.HasPrefix(line, "Start") {
		if len(*fields) <= 0 {
			return fmt.Errorf("\n In Line : %d / Value lines needs declared before start", lineNum)
		}
		if emptyLineCount < 1 {
			return fmt.Errorf("\n In Line : %d / There sould be one empty line between Values and Start", lineNum)
		}
		ParserStateChanger(START)
		return nil
	}
	if len(*fields) <= 0 {
		emptyLineCount = 0
	}
	if emptyLineCount >= 1 {
		return fmt.Errorf("\n In Line : %d / There shouldnt be empty lines between values or more than one empty lines before Start key Word", lineNum)
	}
	if len(line) <= 0 {
		if len(*fields) <= 0 {
			return fmt.Errorf("\n In Line : %d / Value lines needs to be start with Value key word", lineNum)
		}
		emptyLineCount++
		return nil
	}
	if ParserStateController() == VALUE_DECLERATION && strings.HasPrefix(line, "Value") {
		parts := strings.Fields(line)
		if len(parts) < 3 {
			log.Fatalf("Expect at least 3 tokens on line.")
		}
		var options []string
		fieldName := ""
		fieldRegex := ""
		if len(parts) > 3 {
			// options are present
			options = removeDuplicateOptions(strings.Split(parts[1], ","))
			fieldName = parts[2]
			fieldRegex = strings.Join(parts[3:], " ")

		} else if len(parts) == 3 {

			fieldName = parts[1]
			fieldRegex = strings.Join(parts[2:], " ")
		}
		tValue := CreateNewTextFSMValue(TextFSMValue{Options: options, Name: fieldName, Regex: fieldRegex, LineNum: lineNum})
		*fields = append(*fields, tValue)
		return nil
	} else {
		return fmt.Errorf("\n In Line : %d / Value lines needs to be start with Value key word", lineNum)
	}
}

func removeDuplicateOptions(strings []string) []string {
	// Create a map to store unique strings
	uniqueStrings := make(map[string]bool)
	var result []string

	// Loop through the slice and add unique strings to the result
	for _, str := range strings {
		if _, found := uniqueStrings[str]; !found {
			uniqueStrings[str] = true
			result = append(result, str)
		} else {
			log.Println("You have dublicated Options in template!")
		}
	}
	return result
}
