package main

import (
	"GoTextfsm/gotextfsm"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	cliFilePath := "./examples/cisco_version_example"
	templateFilePath := "./examples/cisco_version_template"
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
	log.Println("Parsed Raw output :\n", data.Dict)

	client := redis.NewClient(&redis.Options{
		Addr:     "db:6379",
		Password: "",
		DB:       0,
	})

	err = client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to connect redis Error : %s", err.Error())
	}

	jsonString, err := json.MarshalIndent(data.Dict, "", "	")

	if err != nil {
		log.Fatalf("failed to marhal Error:%s\n", err.Error())
	}

	err = client.Set(context.Background(), "ParsedOutPut", jsonString, 0).Err()
	if err != nil {
		log.Fatalf("Failed to set value Error : %s\n", err.Error())
	}
	value, err := client.Get(context.Background(), "ParsedOutPut").Result()
	if err != nil {
		log.Fatalf("Failed to get value Error : %s\n", err)
	}
	fmt.Printf("Value retrieved from redis : %s\n", value)

	if err := client.Publish(context.Background(), "ParsedData", jsonString).Err(); err != nil {
		log.Fatalf("Cannot publish data to redis Error : %s\n", err)
	}
}
