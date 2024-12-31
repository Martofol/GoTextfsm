package gotextfsm

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type TextFSMError struct {
	msg string
}

func (e *TextFSMError) Error() string {
	return fmt.Sprintf("TextFSM Error: %s", e.msg)
}

type TextFSM struct {
	states       map[string][]TextFSMRule // A map of state names to their rules
	values       map[string]TextFSMValue  // List of values (variables) in FSM
	stateList    []string                 // List of states in order
	currentState string                   // The current state the FSM is in
}

func NewTextFSM() TextFSM {
	return TextFSM{
		states:       make(map[string][]TextFSMRule, 0),
		values:       make(map[string]TextFSMValue),
		stateList:    []string{},
		currentState: "Start", // Start in the Start state
	}
}

func (tFSM *TextFSM) ParseTemplate(templateFilePath string) error {
	file, err := os.Open(templateFilePath)
	if err != nil {
		log.Fatalf("Failed to open template file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNum := 0
	var stateName string
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") || line == "" {
			//we skip commet and empty lines
			continue
		}

		//Value parsing
		if strings.HasPrefix(line, "Value") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return &TextFSMError{msg: fmt.Sprintf("Line : %d / Expect at least 3 tokens on line", lineNum)}
			}
			var options []string
			fieldName := ""
			fieldRegexSTR := ""
			if len(parts) > 3 {
				// options are present
				options = strings.Split(parts[1], ",")
				fieldName = parts[2]
				fieldRegexSTR = strings.Join(parts[3:], " ")

			} else if len(parts) == 3 {

				fieldName = parts[1]
				fieldRegexSTR = strings.Join(parts[2:], " ")
			}
			if err != nil {
				return &TextFSMError{msg: fmt.Sprintf("Line : %d / Failed regex compile Error : %s", lineNum, err)}
			}
			tValue := NewTextFSMValue(TextFSMValue{Options: options, Name: fieldName, Regex: fieldRegexSTR, LineNum: lineNum})
			_, exist := tFSM.values[fieldName]
			if exist {
				return &TextFSMError{msg: fmt.Sprintf("Line : %d / the value name %s is already declared before", lineNum, fieldName)}
			} else {
				tFSM.values[fieldName] = tValue
			}
			continue
		}

		// Check if this is a state definition
		if !strings.HasPrefix(line, "^") {
			// This is the state name
			stateName = line
			if len(tFSM.states) != 0 {
				if _, exists := tFSM.states[stateName]; exists {
					return &TextFSMError{msg: fmt.Sprintf("Line %d : Duplicate state definition: %s", lineNum, stateName)}
				}
			}

			// Add state to FSM
			tFSM.states[stateName] = []TextFSMRule{}
			tFSM.stateList = append(tFSM.stateList, stateName)
			continue
		}

		// Parse rule if we are inside a state
		if len(stateName) > 0 && strings.HasPrefix(line, "^") {
			// Create new rule and add it to the current state
			rule := NewTextFSMRule(line, &tFSM.values)
			// log.Printf("\n\n Rule :\n  match : %s \n  regex : %s \n  Line_Op : %s \n  Record_Op : %s \n  newState : %s \n ", rule.match, rule.regex, rule.lineOp, rule.recordOp, rule.newState)
			tFSM.states[stateName] = append(tFSM.states[stateName], rule)
		}
	}
	// Validate FSM
	err = tFSM.validateFSM()
	if err != nil {
		log.Fatalln("FSM validation error:", err)
	}
	return nil
}

// validateFSM checks the validity of the FSM.
func (tFSM *TextFSM) validateFSM() error {
	// Must have a "Start" state
	if _, exists := tFSM.states["Start"]; !exists {
		return &TextFSMError{msg: "Missing state 'Start'"}
	}

	_, exsist := tFSM.states["End"]
	_, exsist_1 := tFSM.states["EOF"]
	if exsist || exsist_1 {
		// "End" and "EOF" states must be empty
		if len(tFSM.states["End"]) > 0 {
			return &TextFSMError{msg: "Non-empty 'End' state"}
		}
		if len(tFSM.states["EOF"]) > 0 {
			return &TextFSMError{msg: "Non-empty 'EOF' state"}
		}
	} else {
		tFSM.states["EOF"] = nil
		tFSM.stateList = append(tFSM.stateList, "EOF")
	}

	// Validate state transitions in each rule
	for stateName, rules := range tFSM.states {
		for _, rule := range rules {
			if rule.newState != "" && rule.newState != "End" && rule.newState != "EOF" {
				if _, exists := tFSM.states[rule.newState]; !exists {
					return &TextFSMError{
						msg: fmt.Sprintf("State '%s' referenced in '%s' does not exist", rule.newState, stateName),
					}
				}
			}
		}
	}

	return nil
}

func (tFSM *TextFSM) applyValue(matchedValues []string, matchedNames []string) ([]TextFSMValue, error) {
	fillUpValues := make([]TextFSMValue, 0)
	varMap := make(map[string]string, 0)
	for i, name := range matchedNames {
		if i != 0 && name != "" {
			varMap[name] = matchedValues[i]
		}
	}
	if len(varMap) > 0 {
		for key, value := range varMap {
			newValObj, exsist := tFSM.values[key]
			if !exsist {
				continue
			}
			if strings.Contains(newValObj.Regex, "(?P") {
				newValObj.AssignMapVar(varMap)
			} else {
				newValObj.AssignVar(value)
			}
			if Contains(&newValObj.Options, "Fillup") && newValObj.Value != nil {
				fillUpValues = append(fillUpValues, newValObj)
			}
			tFSM.values[key] = newValObj
		}
	}
	// tValue := tFSM.values[valueName]
	// tValue.AssignVar(matchedValue)
	// tFSM.values[valueName] = tValue
	// return tFSM.values[valueName], nil
	return fillUpValues, nil
}

// raiseError raises an error with the given message.
func (fsm *TextFSM) raiseError(msg string) {
	log.Fatalf("\n %s", msg)
}

func (tFSM *TextFSM) changeState(newState string) {
	tFSM.currentState = newState
}
