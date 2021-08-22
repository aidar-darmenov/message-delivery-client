package main

import (
	"github.com/aidar-darmenov/message-delivery-client/config"
	"github.com/aidar-darmenov/message-delivery-client/service"
	"github.com/aidar-darmenov/message-delivery-client/webservice"
	"go.uber.org/zap"
	"log"
	"os"
	"strconv"
)

func main() {

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Incorrect command line arguments")
	}

	cfg := config.NewConfiguration("config/config.json", port)

	// Used uber zap logger for simple example. Now it writes in console
	// Usually, for this purposes we use logs sent to Kibana Elastic Search through Kafka
	var loggerConfig = zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}

	// Creating abstract service(business logic) layer
	s := service.NewService(cfg, logger)

	// Creating abstract webService(delivery) layer
	ws := webservice.NewWebService(s)
	go ws.Start()

	s.StartTcpClient()
}
