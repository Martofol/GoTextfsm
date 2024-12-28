package gotextfsm

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type ParserOutput struct {
	Dict             []map[string]interface{}
	lineNum          int
	currentStateName string
}

func CreateParsedOutput(dataFilePath string, tFSM *TextFSM) ParserOutput {
	parserOutput := ParserOutput{
		Dict:             make([]map[string]interface{}, 0),
		lineNum:          0,       //default number
		currentStateName: "Start", //default state name
	}
	parserOutput.ParseData(dataFilePath, tFSM)
	return parserOutput
}

func (output *ParserOutput) ParseData(dataFilePath string, tFSM *TextFSM) error {
	if tFSM.currentState != "Start" {
		return &TextFSMError{msg: fmt.Sprintf("\n Text Fsm needs to start with Start state current state : %s", tFSM.currentState)}
	}
	output.currentStateName = tFSM.currentState
	lineNum := 0
	file, err := os.Open(dataFilePath)
	if err != nil {
		return &TextFSMError{msg: fmt.Sprintf("Failed to open CLI output file: %s", err)}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		output.lineNum = lineNum
		output.checkLine(line, tFSM)
		if output.currentStateName == "EOf" || output.currentStateName == "End" {
			log.Println("EOF has been called")
			break
		}
	}
	output.appendRecord(*tFSM)
	return nil
}

func (output *ParserOutput) checkLine(line string, tFSM *TextFSM) {
	for _, rule := range tFSM.states[tFSM.currentState] {
		// log.Println("ruleRegex = ", rule.regex)
		if rule.regex.MatchString(line) {
			match := rule.regex.FindStringSubmatch(line)
			// log.Println("	RRegrex.SubexpNames = ", rule.regex.SubexpNames())
			for i, name := range rule.regex.SubexpNames() {
				if i != 0 && name != "" {
					tFSM.applyValue(name, match[i])
				}
			}
			output.HandleOperations(tFSM, &rule)
		}
	}
}

func (output *ParserOutput) HandleOperations(tFSM *TextFSM, rule *TextFSMRule) {

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
	case "Error":
		// Handle error, raise an exception
		tFSM.raiseError(fmt.Sprintf("Error raised in Template: %s", rule.match))
	}
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
		case SKIP_VALUE:
			newRow[name] = nil
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
