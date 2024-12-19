package main

import (
	"GoTextfsm/gotextfsm"
	"fmt"
)

func main() {
	cliFile := "./examples/cisco_version_example"
	templateFile := "./examples/cisco_version_template"

	data := gotextfsm.Parser(templateFile, cliFile)

	// Display results
	fmt.Println("\nExtracted Records:")
	for _, record := range data {
		fmt.Printf("\n  %s: %s", record.Name, record.Value)
	}
}
