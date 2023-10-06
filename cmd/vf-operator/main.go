/*
 Copyright (c) 2023 Sriram Yagaraman

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"google.golang.org/protobuf/encoding/protojson"
)

var defaultConfigFile = "/etc/cni/vf-operator/resource-pool.json"

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
		os.Exit(1)
	}

	configJson, err := os.ReadFile(*input.configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := &network.ResourceConfigs{}
	err = protojson.Unmarshal(configJson, c)
	if err != nil {
		fmt.Printf("error reading configuration file: %v", err.Error())
	}

	input.wg.Add(1)
	go startGrpcServer(input, c)

	input.wg.Add(1)
	go startGrpcGateway(input)

	input.wg.Wait()
}
