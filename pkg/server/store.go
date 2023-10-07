package server

import (
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

type NaEntry struct {
	pciAddress   string
	resourceName string
	naConfig     *NetworkAttachmentConfig
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

func Get(naConfigName string) *NaEntry {
	if na, ok := allocatedNetworkAttachments[naConfigName]; ok {
		return &na
	}

	return nil
}

func Store(naConfig *NetworkAttachmentConfig, pciAddress string, resourceName string) {
	allocatedNetworkAttachments[naConfig.Name] = NaEntry{
		pciAddress:   pciAddress,
		resourceName: resourceName,
		naConfig:     naConfig,
	}
}

func Erase(naConfigName string) {
	delete(allocatedNetworkAttachments, naConfigName)
}
