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
}

const helptext = `
vf-operator sets up the VFs on selected devices
Options;
`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), helptext)
		flag.PrintDefaults()
	}

	input := &Input{
		configFile: flag.String("config", defaultConfigFile, "Path to config file"),
		port:       flag.Int("port", 5001, "gRPC port (tcp:port)"),
		gwPort:     flag.Int("gwPort", 15001, "API port (tcp:port)"),
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

	var wg sync.WaitGroup
	wg.Add(1)
	startGrpcServer(input, c)

	wg.Add(1)
	startGrpcGateway(input)

	wg.Wait()
}
