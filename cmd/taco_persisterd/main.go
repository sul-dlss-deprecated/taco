package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/generated/models"
)

func main() {
	mode := os.Getenv("TACO_ENV")
	if mode == "" {
		mode = "development"
	}

	configFile := fmt.Sprintf("../../config/%s.yaml", mode)
	rt, err := taco.NewRuntime(config.Init(configFile))
	if err != nil {
		log.Fatalln(err)
	}

	// TODO from env?
	shardID := "shardId-000000000000"
	iterator := rt.Stream().GetIterator(&shardID)

	messages, nextIterator := rt.Stream().GetRecords(iterator)

	log.Println("Here are our records: ")
	for i, msg := range messages {
		// TODO: make a custom message type
		message := &models.DepositNewResourceOKBody{}
		json.Unmarshal([]byte(msg), &message)
		log.Printf("Record %v: requestId: %s, resourceId: %s, state: %s",
			i, message.RequestID, message.ID, message.State)
		// resource = &persistence.Resource{ID: record.ID, Title: "None provided"}
		// rt.Repository().SaveItem(resource)
	}
	log.Printf("Next iterator: %s", *nextIterator)

}
