package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sriramy/vf-operator/pkg/config"
)

var defaultConfigFile = "/etc/cni/resource-pool.json"

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
	if configFile == nil {
		fmt.Println("No config file specified, mandatory argument")
		os.Exit(0)
	}

	c := config.GetResourceConfigList(*configFile)
	for _, r := range c.Resources {
		fmt.Printf("Configuration: NIC selector(Vendor: %s, DeviceID: %s, PfNames: %s, RootDevice: %s) NumVFs: %d, DeviceType: %s\n",
			r.GetVendor(), r.GetDeviceID(), r.GetPfNames(), r.GetRootDevices(), r.GetNumVfs(), r.GetDeviceType())
	}
}
