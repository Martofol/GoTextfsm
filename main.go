package main

import (
	"GoTextfsm/gotextfsm"
	"log"
)

func main() {
	cliFile := "./examples/cisco_ipv6_interface_example"
	templateFile := "./examples/cisco_ipv6_interface_template"

	tFSM := gotextfsm.NewTextFSM()
	err := tFSM.ParseTemplate(templateFile)
	if err != nil {
		log.Fatalln("TFSM Error", err)
	}
	data := gotextfsm.CreateParsedOutput(cliFile, &tFSM)
	log.Println("Parsed output :\n", data.Dict)
}
