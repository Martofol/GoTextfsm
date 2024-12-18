package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// TemplateField represents a field and its corresponding regex.
type TemplateField struct {
	Options []string
	Name    string
	Regex   string
}

type Record struct {
	Name  string
	Value string
}

func main() {
	templateFile := "../examples/cisco_version_template"
	cliFile := "../examples/cisco_version_example"

	// Parse the template file
	templateFields, startPatterns := parseTemplateFile(templateFile)

	// Parse CLI output and match patterns
	data := parseCLIOutput(cliFile, templateFields, startPatterns)

	// Display results
	fmt.Println("\nExtracted Records:")
	for _, record := range data {
		fmt.Printf("  %s: %s\n", record.Name, record.Value)
	}
}

func parseTemplateFile(filePath string) ([]TemplateField, []string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open template file: %s", err)
	}
	defer file.Close()

	var fields []TemplateField
	var startPatterns []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "  ^") {
			startPatterns = append(startPatterns, line[2:])
			continue
		}

		if strings.HasPrefix(line, "Value") {
			parts := strings.Fields(line)
			if len(parts) < 3 {
				log.Fatalf("Expect at least 3 tokens on line.")
			}
			if len(parts) > 3 {

				// options are present
				options := removeDuplicateOptions(strings.Split(parts[1], ","))
				fieldName := parts[2]
				fieldRegex := strings.Join(parts[3:], " ")
				fields = append(fields, TemplateField{Name: fieldName, Options: options, Regex: fieldRegex})

			} else if len(parts) == 3 {

				fieldName := parts[1]
				fieldRegex := strings.Join(parts[2:], " ")
				fields = append(fields, TemplateField{Name: fieldName, Regex: fieldRegex})

			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading template file: %s", err)
	}

	return fields, startPatterns
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

func parseCLIOutput(filePath string, fields []TemplateField, patterns []string) []Record {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open CLI output file: %s", err)
	}
	defer file.Close()

	var records []Record

	//Debuging

	// for _, pattern := range patterns {
	// 	log.Println("Pattern :", pattern)
	// }

	// for _, field := range fields {
	// 	log.Println("\nField Name :", field.Name)
	// 	log.Println("Field Option length:", len(field.Options))
	// 	log.Println("Field Regex :", field.Regex)
	// }

	//--EndDebuging

	// Build regex for start patterns
	for _, pattern := range patterns {
		combinedRegex := pattern
		log.Println("Pattern :", pattern)
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

		// Parse CLI output
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				match := re.FindStringSubmatch(line)

				for i, name := range re.SubexpNames() {
					if i != 0 && name != "" {
						records = append(records, Record{Name: name, Value: match[i]})
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading CLI output file: %s", err)
		}
	}

	return records
}
