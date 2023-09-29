package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sriramy/vf-operator/pkg/config"
	"github.com/sriramy/vf-operator/pkg/devices"
	"github.com/sriramy/vf-operator/pkg/utils"
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
		provider := devices.NewNetDeviceProvider(&r)
		err := provider.Discover()
		if err != nil {
			log.Fatalf("Discover failed: %v", err)
		}

		fmt.Println("Discover results")
		for _, dev := range provider.GetDevices() {
			fmt.Println("==========")
			fmt.Printf("Name: %s\n", dev.Name)
			fmt.Printf("VF: %v\n", utils.IsSriovVF(dev.PCIAddress))
			fmt.Printf("MAC Address: %s\n", dev.MACAddress)
			fmt.Printf("PCI Address: %s\n", *dev.PCIAddress)
			fmt.Println("==========")
		}
	}
}
