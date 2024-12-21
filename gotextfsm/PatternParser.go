package gotextfsm

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type StartPattern struct {
	Regex *regexp.Regexp
	Rules []string
}

func CheckTheLineForPatterns(StartPatterns *[]StartPattern, fields []TextFSMValue, line string, lineNum int) error {
	if ParserStateController() == START && strings.HasPrefix(line, "  ^") {
		reRule := regexp.MustCompile(`^(.*)->\s*(.*)$`)
		match := reRule.FindStringSubmatch(line)
		patternRule := ""
		patternValue := ""
		if len(match) > 1 {
			//If line has rule like -> contunie
			patternRule = match[2]
			patternValue = match[1][3:]
		} else {
			patternValue = line[3:]
		}
		if patternRule == "" {
			patternRule = "Next" // Default value for patter rule
		}
		// Build regex for start patterns
		combinedRegex := patternValue
		// Replace placeholders (${FieldName}) with named regex groups (?P<FieldName>regex)
		for _, field := range fields {
			placeholder := fmt.Sprintf("${%s}", field.Name)
			namedGroup := fmt.Sprintf("(?P<%s>%s)", field.Name, field.Regex)
			combinedRegex = strings.ReplaceAll(combinedRegex, placeholder, namedGroup)
		}

		// Compile the resulting regex
		re, err := regexp.Compile(combinedRegex)
		if err != nil {
			log.Fatalf("Error compiling regex: %s", err)
		}
		*StartPatterns = append(*StartPatterns, StartPattern{Regex: re, Rules: strings.Split(patternRule, ".")})
		return nil
	} else {
		if ParserStateController() == START {
			ParserStateChanger(END)
			return nil
		}
		return fmt.Errorf("\n In Line : %d / Pattern lines needs to be start with (  ^)", lineNum)
	}
}
