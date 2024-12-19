package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
)

type TextFSMValue struct {
	Options []string
	Name    string
	Regex   string
	LineNum int
}

const MAX_NAME_LEN = 48

func CreateNewTextFSMValue(textFSMValue TextFSMValue) TextFSMValue {
	textFSMValue.CheckFieldsValidness()
	return textFSMValue
}

func (t *TextFSMValue) CheckFieldsValidness() {
	var errList []error
	addError(&errList, t.hasValidName())
	addError(&errList, t.isValidRegex())
	addError(&errList, t.isValidOption())

	if len(errList) >= 1 {
		for _, err := range errList {
			fmt.Printf("\nError in Line : %d / %s ", t.LineNum, err)
		}
		log.Fatalf("\nError has been occured!")
	}
}

func (t *TextFSMValue) hasValidName() error {
	if len(t.Name) < MAX_NAME_LEN {
		return nil
	}
	return fmt.Errorf("the %s name is too long it needs to be under 48 letter", t.Name)
}

func (t *TextFSMValue) isValidRegex() error {
	// Regular expression to check if the string starts and ends with parentheses
	matched, _ := regexp.MatchString(`^\(.*\)$`, t.Regex)
	if matched {
		return nil
	} else {
		return fmt.Errorf("the %s Regex must be enclosed with parentheses", t.Regex)
	}
}

func (t *TextFSMValue) isValidOption() error {
	for _, option := range t.Options {
		switch option {
		case "Required", "Key", "List", "Filldown", "Fillup":
			continue
		}
		return fmt.Errorf("the %s option name is not valid", option)
	}
	return nil
}

func addError(errorList *[]error, err error) {
	if err != nil {
		*errorList = append(*errorList, err)
	}
}
