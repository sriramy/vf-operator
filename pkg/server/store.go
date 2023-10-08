package server

import (
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

type NaEntry struct {
	name         string
	resourceName string
	config       string
	pciAddress   string
}

var allocatedNetworkAttachments map[string]NaEntry

func init() {
	allocatedNetworkAttachments = make(map[string]NaEntry)
}

func GetAll() []*NaEntry {
	naEntries := make([]*NaEntry, 0)
	for _, na := range allocatedNetworkAttachments {
		naEntries = append(naEntries, &na)
	}
	return naEntries
}

func IsAllocated(pciAddress string) network.VFStatus {
	for _, na := range allocatedNetworkAttachments {
		if na.pciAddress == pciAddress {
			return network.VFStatus_USED
		}
	}

	return network.VFStatus_FREE
}

func Get(name string) *NaEntry {
	if na, ok := allocatedNetworkAttachments[name]; ok {
		return &na
	}

	return nil
}

func Store(name string, resourceName string, config string, pciAddress string) {
	allocatedNetworkAttachments[name] = NaEntry{
		name:         name,
		resourceName: resourceName,
		config:       config,
		pciAddress:   pciAddress,
	}
}

func Erase(name string) {
	delete(allocatedNetworkAttachments, name)
}
