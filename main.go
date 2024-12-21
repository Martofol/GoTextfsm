package main

import (
	"GoTextfsm/gotextfsm"
	"fmt"
	"strings"
)

func main() {
	cliFile := "./examples/f10_ip_bgp_summary_example"
	templateFile := "./examples/f10_ip_bgp_summary_template"

	data := gotextfsm.Parser(templateFile, cliFile)
	// Display results
	fmt.Println("\nExtracted Records:")
	for _, recordName := range data.Name {
		fmt.Printf("\n  %s: %s", recordName, strings.Join(data.Value[recordName], " , "))
	}
}
