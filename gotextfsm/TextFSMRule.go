package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// TextFSMRule represents a rule in the FSM, including a regex pattern, line and record operations, and state transition.
type TextFSMRule struct {
	match    string         // Regular expression to match against input lines
	regex    *regexp.Regexp // Compiled regular expression
	lineOp   string         // Line operation (Next, Continue, etc.)
	recordOp string         // Record operation (Record, Clear, etc.)
	newState string         // The next state to transition to
}

var LINE_OPs = []string{"Continue", "Next"}
var RECORD_OPs = []string{"NoRecord", "Record", "Clear", "Clearall", "Error"}

// NewTextFSMRule creates a new FSM rule with the given match string, line operation, and record operation.
func NewTextFSMRule(line string, fields *map[string]TextFSMValue) TextFSMRule {
	// Replace placeholders (${FieldName}) with named regex groups (?P<FieldName>regex)
	for name, value := range *fields {
		placeholder := fmt.Sprintf("${%s}", name)
		namedGroup := fmt.Sprintf("(?P<%s>%s)", name, value.Regex)
		line = strings.ReplaceAll(line, placeholder, namedGroup)
	}

	newTRule := TextFSMRule{
		lineOp:   "Next",     // default value for line_Op
		recordOp: "NoRecord", // default value for record_Op
		newState: "",         // default value for newState
		match:    line,
	}
	reRule := regexp.MustCompile(`^(.*)->\s*(.*)$`)
	match := reRule.FindStringSubmatch(line)

	if len(match) > 1 {
		//If line has rule like -> contunie
		err := newTRule.setOperations(match[2])
		if err != nil {
			log.Fatalln("Rule Operator Error : ", err)
		}
		re, err := regexp.Compile(match[1])
		if err != nil {
			log.Fatalln("Error :", err)
		}
		newTRule.regex = re
	} else {
		re, err := regexp.Compile(line)
		if err != nil {
			log.Fatalln("Error :", err)
		}
		newTRule.regex = re
	}
	return newTRule
}

func (tRule *TextFSMRule) setOperations(operationPart string) error {
	if operationPart == "" {
		tRule.newState = ""         // default value for newState
		tRule.lineOp = "Next"       // default value for line_Op
		tRule.recordOp = "NoRecord" // default value for record_Op
		return nil                  //we allow empty operator declaration
	}
	reRule_1 := regexp.MustCompile(`^(.*)\s(.*)$`) //is there a operation and newState declaration
	match_1 := reRule_1.FindStringSubmatch(operationPart)

	if len(match_1) > 1 {
		//new state has been found
		err, actionFound := tRule.setActions(match_1[1])
		if err != nil {
			return fmt.Errorf("\n Action syntax Error : %s ", err)
		}
		if !actionFound {
			return fmt.Errorf("\n Action syntax Error : Action %s doesn't exsist ", match_1[1])
		}
		tRule.newState = match_1[2]
		return nil
	}
	err, actionFound := tRule.setActions(operationPart)
	if err != nil {
		return fmt.Errorf("\n Action syntax Error : %s ", err)
	}
	if !actionFound {
		tRule.newState = operationPart
		//we check if the state exsist while trantitioning so no need to check here again
	}
	return nil
}

func (tRule *TextFSMRule) setActions(line string) (error, bool) {
	reRule_2 := regexp.MustCompile(`(\b\w+)\.(\w+\b)`) //is there a line operation like Next.Record
	match_2 := reRule_2.FindStringSubmatch(line)
	if len(match_2) > 1 {
		//there are multiple operations like Next.Record
		if !Contains(&LINE_OPs, match_2[1]) {
			return fmt.Errorf("\n There is no Line Action named %s exsist. Possible line actions :(%s) ", match_2[1], strings.Join(LINE_OPs, ",")), false
		}
		if !Contains(&RECORD_OPs, match_2[2]) {
			return fmt.Errorf("\n There is no Record Action named %s exsist. Possible record actions :(%s) ", match_2[2], strings.Join(RECORD_OPs, ",")), false
		}
		tRule.lineOp = match_2[1]
		tRule.recordOp = match_2[2]
		return nil, true
	}
	line = strings.TrimSpace(line)
	if Contains(&LINE_OPs, line) {
		tRule.lineOp = line
		return nil, true
	}
	if Contains(&RECORD_OPs, line) {
		tRule.recordOp = line
		return nil, true
	}

	return nil, false
}
