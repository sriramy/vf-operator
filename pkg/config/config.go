package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ResourceConfig struct {
	Name         string `json:"name"`
	Mtu          int    `json:"mtu"`
	NeedVhostNet bool   `json:"needVhostNet"`
	NumVfs       int    `json:"numVfs"`
	NicSelector  struct {
		Vendors []string `json:"vendors,omitempty"`
		Drivers []string `json:"drivers,omitempty"`
		Devices []string `json:"devices,omitempty"`
		PfNames []string `json:"pfNames,omitempty"`
	} `json:"nicSelector"`
	DeviceType string `json:"deviceType"`
}

type ResourceConfigList struct {
	Resources []ResourceConfig `json:"resources"`
}

func GetResourceConfigList(file string) (ResourceConfigList, error) {
	var config ResourceConfigList
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("config file (%s) not found, skipping\n", file)
		return ResourceConfigList{}, nil
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return ResourceConfigList{}, err
	}
	return config, nil
}

func (c *ResourceConfig) GetName() string {
	return c.Name
}

func (c *ResourceConfig) GetMtu() int {
	return c.Mtu
}

func (c *ResourceConfig) GetVhostNet() bool {
	return c.NeedVhostNet
}

func (c *ResourceConfig) GetNumVfs() int {
	return c.NumVfs
}

func (c *ResourceConfig) GetVendors() []string {
	return c.NicSelector.Vendors
}

func (c *ResourceConfig) GetDrivers() []string {
	return c.NicSelector.Drivers
}

func (c *ResourceConfig) GetDevices() []string {
	return c.NicSelector.Devices
}

func (c *ResourceConfig) GetPfNames() []string {
	return c.NicSelector.PfNames
}

func (c *ResourceConfig) GetDeviceType() string {
	return c.DeviceType
}
