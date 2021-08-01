package main

import (
	delivery "github.com/gulfcoastdevops/buffertest/deliveries/timer"
	service "github.com/gulfcoastdevops/buffertest/one"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	// Create a channel to buffer messages sent from MQTT
	events := make(chan []byte)

	// Start processing messages from MQTT to channel
	err := delivery.Start(&wg, events, "alorica/http_response")
	if err != nil {
		log.Panic(err)
	}

	// Create a single instance of usecase to store messages
	usecase, err := service.Init(false)
	if err != nil {
		log.Panic(err)
	}

	// Create a formatter to read the event from the channel and
	// call the use case for each event
	formatter := &service.FormatAdapter{}
	err = formatter.Start(events, usecase)
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	err = usecase.CloseRepository()
	if err != nil {
		log.Panic(err)
	}
}