package utils

import nl "github.com/vishvananda/netlink"

func SetIfName(ifName string, newName string) error {
	link, err := nl.LinkByName(ifName)
	if err != nil {
		return err
	}
	return nl.LinkSetName(link, newName)
}
