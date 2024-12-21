package gotextfsm

import (
	"bufio"
	"log"
	"os"
)

type Record struct {
	Name        []string            //equals to the field name these will be collum
	Value       map[string][]string //equals to extracted Values these will be rolls
	Options     map[string][]string //equals to the field Options
	MaxRowCount int
}

type re_Table struct {
	Collums []string          //colums of reTable
	Rows    map[string]string //Rows of reTable
}

var record Record

func ParseCLIOutput(filePath string, fields []TextFSMValue, startPatterns []StartPattern) Record {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CLI output file: %s", err)
	}
	defer file.Close()
	//Seting Up Record
	record = Record{}
	record.Value = make(map[string][]string)
	record.Options = make(map[string][]string)
	record.MaxRowCount = 0

	for _, tValue := range fields {
		record.Name = append(record.Name, tValue.Name)
		record.Options[tValue.Name] = tValue.Options
	}

	// Parse CLI output
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, pattern := range startPatterns {
			if pattern.Regex.MatchString(line) {
				match := pattern.Regex.FindStringSubmatch(line) //row value we want to save
				for i, name := range pattern.Regex.SubexpNames() {
					if i != 0 && name != "" {
						record.Set(name, pattern.Rules, match[i])
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading CLI output file: %s", err)
	}

	for _, collum := range record.Name {
		if len(record.Options[collum]) > 0 {
			for _, option := range record.Options[collum] {
				switch option {
				case "Filldown":
					for i := 0; i < record.MaxRowCount-1; i++ {
						record.Value[collum] = append(record.Value[collum], record.Value[collum][len(record.Value[collum])-1])
					}
				}
			}
		}
	}

	return record
}

func (r *Record) Set(name string, rules []string, value string) {
	if name == "Received_V4" {
		log.Println("Value of Received_V4 is :", value)
	}
	newValue := []string{value}
	for _, rule := range rules {
		switch rule {
		case "Next":
			continue
		case "Continue":
			//skip the pattern
			r.Value[name] = append(record.Value[name], "")
			return
		case "Record":
			//Do record functions
			newValue = append(r.Value[name], value)
		default:
			log.Fatalf("\n The rule : %s is not valid", rule)
		}
	}
	r.Value[name] = newValue
	if len(r.Value[name]) > r.MaxRowCount {
		r.MaxRowCount = len(r.Value[name])
	}
}
