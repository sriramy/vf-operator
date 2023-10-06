package server

import (
	"encoding/json"
	"os"
	"path/filepath"

	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

const (
	cniVersion = "1.0.0"
	sriovType  = "sriov"
)

type PluginConfig struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	DeviceID string `json:"deviceID"`
	Vlan     uint32 `json:"vlan"`
	Mtu      uint32 `json:"mtu"`
}

type NetworkAttachmentConfig struct {
	CniVersion string         `json:"cniVersion"`
	Name       string         `json:"name"`
	Plugins    []PluginConfig `json:"plugins"`
}

func NewSriovNetworkAttachmentConfig(na *network.NetworkAttachment, pciAddress string) *NetworkAttachmentConfig {
	return &NetworkAttachmentConfig{
		CniVersion: cniVersion,
		Name:       na.GetName(),
		Plugins: []PluginConfig{
			{
				Name:     na.GetName(),
				Type:     sriovType,
				DeviceID: pciAddress,
				Vlan:     na.GetVlan(),
				Mtu:      na.GetMtu(),
			},
		},
	}
}

func AddNetworkAttachment(naConfig *NetworkAttachmentConfig) error {
	file := filepath.Join("/etc/cni/net.d", naConfig.Name, ".conflist")
	json, _ := json.MarshalIndent(naConfig, "", " ")

	return os.WriteFile(file, json, 0o644)
}

func RemoveNetworkAttachment(naConfig *NetworkAttachmentConfig) error {
	file := filepath.Join("/etc/cni/net.d", naConfig.Name, ".conflist")
	return os.Remove(file)
}
