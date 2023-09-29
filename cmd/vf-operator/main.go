package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sriramy/vf-operator/pkg/config"
)

var defaultConfigFile = "/etc/cni/vf-operator.json"

const helptext = `
vf-operator sets up the VFs on selected devices
Options;
`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), helptext)
		flag.PrintDefaults()
	}

	configFile := flag.String("config", defaultConfigFile, "Path to config file")
	flag.Parse()
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	if configFile == nil {
		fmt.Println("No config file specified, mandatory argument")
		os.Exit(0)
	}

	c := config.GetConfig(*configFile)
	fmt.Printf("Configuration: NIC selector(Vendor: %s, DeviceID: %s, PfNames: %s, RootDevice: %s) NumVFs: %d, DeviceType: %s",
		c.NicSelector.Vendor, c.NicSelector.DeviceID, c.NicSelector.PfNames, c.NicSelector.RootDevices,
		c.NumVfs, c.DeviceType)
}
