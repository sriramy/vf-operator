package utils

import (
	"github.com/vishvananda/netlink"
	nl "github.com/vishvananda/netlink"
)

func GetLinkMtu(ifName string) uint32 {
	link, err := nl.LinkByName(ifName)
	if err != nil {
		return 0
	}

	dev, ok := link.(*netlink.Device)
	if !ok {
		return 0
	}

	return uint32(dev.MTU)
}

func SetLinkMtu(ifName string, mtu uint32) error {
	link, err := nl.LinkByName(ifName)
	if err != nil {
		return err
	}

	devMtu := 0
	dev, ok := link.(*netlink.Device)
	if ok {
		devMtu = dev.MTU
	}

	if devMtu != int(mtu) {
		err = nl.LinkSetMTU(link, int(mtu))
		if err != nil {
			return err
		}
	}

	return nil
}

func SetIfName(ifName string, newName string) error {
	link, err := nl.LinkByName(ifName)
	if err != nil {
		return err
	}
	return nl.LinkSetName(link, newName)
}
