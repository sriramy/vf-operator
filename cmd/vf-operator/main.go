package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/sriramy/vf-operator/pkg/config"
)

var defaultConfigFile = "/etc/cni/resource-pool.json"

type Input struct {
	configFile *string
	port       *int
	gwPort     *int
	wg         *sync.WaitGroup
}

const helptext = `
vf-operator sets up the VFs on selected devices
Options;
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), helptext)
		flag.PrintDefaults()
	}

	input := &Input{
		configFile: flag.String("config", defaultConfigFile, "Path to config file"),
		port:       flag.Int("port", 5001, "gRPC port (tcp:port)"),
		gwPort:     flag.Int("gwPort", 15001, "API port (tcp:port)"),
		wg:         &sync.WaitGroup{},
	}

	flag.Parse()
	if input.configFile == nil {
		fmt.Println("No config file specified, mandatory argument")
		os.Exit(0)
	}

	c, err := config.GetResourceConfigList(*input.configFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input.wg.Add(1)
	go startGrpcServer(input, c)

	input.wg.Add(1)
	go startGrpcGateway(input)

	input.wg.Wait()
}
