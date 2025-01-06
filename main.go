package main

import (
	"GoTextfsm/gotextfsm"
	"log"
	"os"
)

func main() {
	cliFilePath := "./examples/unix_ifcfg_example"
	templateFilePath := "./examples/unix_ifcfg_template"
	file, err := os.Open(templateFilePath)
	if err != nil {
		log.Fatalf("Failed to open template file: %s", err)
	}
	tFSM := gotextfsm.NewTextFSM()
	err = tFSM.ParseTemplate(file)
	if err != nil {
		log.Fatalln("TFSM Error", err)
	}
	stringData, err := os.ReadFile(cliFilePath)
	if err != nil {
		log.Fatalf("Failed to read data file: %s", err)
	}
	data := gotextfsm.CreateParsedOutput(string(stringData), &tFSM)
	log.Println("Parsed output :\n", data.Dict)
}
