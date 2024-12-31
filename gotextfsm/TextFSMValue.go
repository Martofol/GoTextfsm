package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type TextFSMValue struct {
	Options       []string    // Options associated with the value (optional)
	Name          string      // Name of the variable
	Regex         string      // Compiled regex for matching input text
	Value         interface{} // The value extracted from the match
	FillDownValue interface{}
	LineNum       int
}

const MAX_NAME_LEN = 48

type ON_RECORD_TYPE int

const (
	SKIP_RECORD ON_RECORD_TYPE = iota
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
	matched, _ := regexp.MatchString(`^\(.*\)$`, t.Regex)
	if matched {
		//remove the most outer parantheses here
		t.Regex = t.Regex[1 : len(t.Regex)-1]
	} else {
		return fmt.Errorf("the %s Regex must be enclosed with parentheses", t.Regex)
	}
	if strings.Count(t.Regex, "(") != strings.Count(t.Regex, ")") {
		return fmt.Errorf(" %s Value regex '%s' must be contained within a '()' pair. ", t.Name, t.Regex)
	}

	return nil
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

func (t *TextFSMValue) AssignVar(matchedValue string) {
	var finalVal interface{} = matchedValue
	if Contains(&t.Options, "List") {
		if t.Value == nil {
			if Contains(&t.Options, "Filldown") && t.FillDownValue != nil {
				finalVal = append(t.FillDownValue.([]string), matchedValue)
			} else {
				finalVal = make([]string, 0)
				finalVal = append(finalVal.([]string), matchedValue)
			}
		} else {
			finalVal = append(t.Value.([]string), matchedValue)
		}

	}
	if Contains(&t.Options, "Filldown") {
		t.FillDownValue = finalVal
	}
	t.Value = finalVal
}

// var map = the values and corresponding names that matched from the input line
func (t *TextFSMValue) AssignMapVar(varMap map[string]string) {
	newVal := make(map[string]string)
	re := regexp.MustCompile(`\?P<([a-zA-Z0-9_]+)>`)
	matchedNames := re.FindStringSubmatch(t.Regex)
	if len(matchedNames) > 1 {
		for i := 1; i < len(matchedNames); i++ {
			newVal[matchedNames[i]] = varMap[matchedNames[i]]
		}
	} else {
		log.Fatalf("No name matches have been found in the nested values of the value : %s Make sure the nested value follows this rule : %s\n", t.Name, `\?P<([a-zA-Z0-9_]+)>`)
	}
	var finalVal interface{} = newVal
	if Contains(&t.Options, "List") {
		if t.Value == nil {
			if Contains(&t.Options, "Filldown") && t.FillDownValue != nil {
				finalVal = append(t.FillDownValue.([]map[string]string), newVal)
			} else {
				finalVal = make([]map[string]string, 0)
				finalVal = append(finalVal.([]map[string]string), newVal)
			}
		} else {
			finalVal = append(t.Value.([]map[string]string), newVal)
		}
	}
	if Contains(&t.Options, "Filldown") {
		t.FillDownValue = finalVal
	}
	t.Value = finalVal
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
	if (t.Value == nil || Contains(&t.Options, "Filldown")) && t.FillDownValue != nil {
		return t.FillDownValue
	}
	if t.Value != nil {
		return t.Value
	}
	return "Null"
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
