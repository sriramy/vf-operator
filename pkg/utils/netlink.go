package utils

import (
	"fmt"

	"github.com/jaypipes/ghw"
	nl "github.com/vishvananda/netlink"
)

func SetIfName(ifName string, newName string) error {
	link, err := nl.LinkByName(ifName)
	if err != nil {
		return err
	}
	return nl.LinkSetName(link, newName)
}

func GetVfNic(dev *ghw.PCIDevice) (*ghw.NIC, error) {
	net, err := ghw.Network()
	if err != nil {
		return nil, fmt.Errorf("Couldn't get NIC info: %v", err)
	}

	for _, nic := range net.NICs {
		if nic.PCIAddress == nil || *nic.PCIAddress != dev.Address {
			continue
		}
		return nic, nil
	}

	return nil, fmt.Errorf("No NIC found with PCI Address: %s", dev.Address)
}
