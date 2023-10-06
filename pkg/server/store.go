package server

import (
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

type naEntry struct {
	pciAddress string
	naConfig   *NetworkAttachmentConfig
}

var allocatedNetworkAttachments map[string]naEntry

func init() {
	allocatedNetworkAttachments = make(map[string]naEntry)
}

func IsAllocated(resourceName string) network.VFStatus {
	if _, ok := allocatedNetworkAttachments[resourceName]; ok {
		return network.VFStatus_USED
	}

	return network.VFStatus_FREE
}

func Get(resourceName string) *NetworkAttachmentConfig {
	if na, ok := allocatedNetworkAttachments[resourceName]; ok {
		return na.naConfig
	}

	return nil
}

func Store(naConfig *NetworkAttachmentConfig, pciAddress string) {
	allocatedNetworkAttachments[naConfig.Name] = naEntry{
		pciAddress: pciAddress,
		naConfig:   naConfig,
	}
}

func Erase(resourceName string) {
	delete(allocatedNetworkAttachments, resourceName)
}
