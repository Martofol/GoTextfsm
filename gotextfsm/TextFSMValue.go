package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
)

type TextFSMValue struct {
	Options       []string       // Options associated with the value (optional)
	Name          string         // Name of the variable
	Regex         *regexp.Regexp // Compiled regex for matching input text
	Value         interface{}    // The value extracted from the match
	FillDownValue interface{}
	LineNum       int
}

const MAX_NAME_LEN = 48

type ON_RECORD_TYPE int

const (
	SKIP_RECORD ON_RECORD_TYPE = iota
	SKIP_VALUE
	CONTINUE
)

func NewTextFSMValue(textFSMValue TextFSMValue) TextFSMValue {
	textFSMValue.CheckFieldsValidness()
	textFSMValue.Value = nil
	textFSMValue.FillDownValue = nil
	return textFSMValue
}

func (t *TextFSMValue) CheckFieldsValidness() {
	var errList []error
	addError(&errList, t.isValidName())
	addError(&errList, t.isValidRegex())
	addError(&errList, t.isValidOption())

	if len(errList) >= 1 {
		for _, err := range errList {
			fmt.Printf("\nError in Line : %d / %s ", t.LineNum, err)
		}
		log.Fatalf("\nError has been occured!")
	}
}

func (t *TextFSMValue) isValidName() error {
	if len(t.Name) > MAX_NAME_LEN {
		return fmt.Errorf("the %s name is too long it needs to be under 48 letter", t.Name)
	}
	return nil
}

func (t *TextFSMValue) isValidRegex() error {
	// Regular expression to check if the string starts and ends with parentheses
	reString := t.Regex.String()
	matched, _ := regexp.MatchString(`^\(.*\)$`, reString)
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

func (t *TextFSMValue) AssignVar(matchedValue interface{}) {
	var finalVal interface{}
	finalVal = matchedValue
	if Contains(&t.Options, "Filldown") {
		t.FillDownValue = finalVal
	}
	t.Value = finalVal
}

func (t *TextFSMValue) AppendValue(matchedValue string) {
	// value.Value = append(value.Value, matchedValue)
}

// this is a event called when we want to create new row
func (t *TextFSMValue) OnAppendRecord() ON_RECORD_TYPE {
	if Contains(&t.Options, "Required") {
		if t.isEmpty() {
			if Contains(&t.Options, "Filldown") {
				if t.FillDownValue == nil {
					return SKIP_RECORD
				} else {
					return CONTINUE
				}
			}
			return SKIP_RECORD
		}
	}
	return CONTINUE
}

func (t *TextFSMValue) clearValue(all bool) {
	if all && Contains(&t.Options, "Filldown") {
		t.FillDownValue = nil
	}
	t.Value = nil
}

func (t *TextFSMValue) GetFinalValue() interface{} {
	if t.Value == nil || Contains(&t.Options, "Filldown") {
		return t.FillDownValue
	}
	return t.Value
}

func (t *TextFSMValue) isEmpty() bool {
	if Contains(&t.Options, "Filldown") {
		return t.FillDownValue == nil
	}
	return t.Value == nil
}

func addError(errorList *[]error, err error) {
	if err != nil {
		*errorList = append(*errorList, err)
	}
}
