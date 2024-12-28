package main

import (
	"GoTextfsm/gotextfsm"
	"log"
)

func main() {
	cliFile := "./examples/f10_ip_bgp_summary_example"
	templateFile := "./examples/f10_ip_bgp_summary_template"

	tFSM := gotextfsm.NewTextFSM()
	err := tFSM.ParseTemplate(templateFile)
	if err != nil {
		log.Fatalln("TFSM Error", err)
	}
	data := gotextfsm.CreateParsedOutput(cliFile, &tFSM)
	log.Println("Parsed output :\n", data.Dict)
}
