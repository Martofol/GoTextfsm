package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type ParserOutput struct {
	Dict             []map[string]interface{}
	lineNum          int
	currentStateName string
}

func CreateParsedOutput(stringData string, tFSM *TextFSM) ParserOutput {
	parserOutput := ParserOutput{
		Dict:             make([]map[string]interface{}, 0),
		lineNum:          0,       //default number
		currentStateName: "Start", //default state name
	}
	parserOutput.ParseData(stringData, tFSM)
	return parserOutput
}

func (output *ParserOutput) ParseData(stringData string, tFSM *TextFSM) error {
	if tFSM.currentState != "Start" {
		return &TextFSMError{msg: fmt.Sprintf("\n Text Fsm needs to start with Start state current state : %s", tFSM.currentState)}
	}
	output.currentStateName = tFSM.currentState
	lineNum := 0
	lines := strings.Split(stringData, "\n")
	for i, line := range lines {
		lineNum = i
		output.lineNum = lineNum
		output.checkLine(line, tFSM)
		if output.currentStateName == "EOf" || output.currentStateName == "End" {
			// log.Println("EOF has been called")
			break
		}
	}
	output.appendRecord(*tFSM)
	return nil
}

func (output *ParserOutput) checkLine(line string, tFSM *TextFSM) {
	for _, rule := range tFSM.states[tFSM.currentState] {
		ruleRegex := regexp.MustCompile(rule.regex)
		if ruleRegex.MatchString(line) {
			match := ruleRegex.FindStringSubmatch(line)
			fillupVals, err := tFSM.applyValue(match, ruleRegex.SubexpNames())
			if err != nil {
				log.Fatalf("\n Value field error : %s", err)
			}
			for _, val := range fillupVals {
				if output.Dict != nil {
					for i := len(output.Dict) - 1; i >= 0; i-- {
						if output.Dict[i][val.Name] == "Null" {
							output.Dict[i][val.Name] = val.Value
						} else {
							break
						}
					}
				}
			}
			next := output.HandleOperations(tFSM, &rule)
			if next {
				break
			}
		}
	}
}

func (output *ParserOutput) HandleOperations(tFSM *TextFSM, rule *TextFSMRule) bool {

	//Handle record operations
	switch rule.recordOp {
	case "NoRecord":
		//continue parsing
	case "Record":
		output.appendRecord(*tFSM)
	case "Clear":
		output.clearRecord(*tFSM)
	case "Clearall":
		output.clearAllRecord(*tFSM)
	}
	// Handle line operations
	switch rule.lineOp {
	case "Continue":
		// Do not change state, continue with the current line
	case "Next":
		// Transition to the next state (if specified)
		if rule.newState != "" {
			tFSM.changeState(rule.newState)
			output.currentStateName = tFSM.currentState
		}
		return true
	case "Error":
		// Handle error, raise an exception
		tFSM.raiseError(fmt.Sprintf("Error raised in Template: %s", rule.match))
	}
	return false
}

// means create the row with the values of tfsm
func (output *ParserOutput) appendRecord(tFSM TextFSM) {
	newRow := make(map[string]interface{})
	containsValue := false
	for name, value := range tFSM.values {
		ret := value.OnAppendRecord()
		switch ret {
		case SKIP_RECORD:
			output.clearRecord(tFSM)
			return
		case CONTINUE:
			newRow[name] = value.GetFinalValue()
			if !value.isEmpty() {
				containsValue = true
			}
		}
	}
	if containsValue {
		output.Dict = append(output.Dict, newRow)
	}
	output.clearRecord(tFSM)
}

func (output *ParserOutput) clearRecord(tFSM TextFSM) {
	for name, value := range tFSM.values {
		value.clearValue(false)
		tFSM.values[name] = value
	}
}

func (output *ParserOutput) clearAllRecord(tFSM TextFSM) {
	for name, value := range tFSM.values {
		value.clearValue(true)
		tFSM.values[name] = value
	}
}
