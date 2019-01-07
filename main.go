package main

import (
	"flag"
	"os"
	"log"

	"github.com/alphagov/paas-go/broker"
	"github.com/alphagov/paas-drone-agent-broker/provider"
	"net"
	"fmt"
	"net/http"
	"code.cloudfoundry.org/lager"
)

var configFilePath string

func main() {
	flag.StringVar(&configFilePath, "config", "", "Location of the config file")
	flag.Parse()

	file, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Error opening config file %s: %s\n", configFilePath, err)
	}
	defer file.Close()

	config, err := broker.NewConfig(file)
	if err != nil {
		log.Fatalf("Error validating config file: %v\n", err)
	}

	droneAgentProvider, err := provider.NewDroneAgentProvider(config.Provider)
	if err != nil {
		log.Fatalf("Error creating Drone Agent Provider: %v\n", err)
	}

	logger := lager.NewLogger("drone-agent-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, config.API.LagerLogLevel))

	serviceBroker := broker.New(config, droneAgentProvider, logger)
	brokerAPI := broker.NewAPI(serviceBroker, logger, config)

	listener, err := net.Listen("tcp", ":"+config.API.Port)
	if err != nil {
		log.Fatalf("Error listening to port %s: %s", config.API.Port, err)
	}
	fmt.Println("Drone Agent Service Broker started on port " + config.API.Port + "...")
	http.Serve(listener, brokerAPI)
}
