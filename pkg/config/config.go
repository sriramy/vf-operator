package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ResourceConfig struct {
	Mtu          int  `json:"mtu"`
	NeedVhostNet bool `json:"needVhostNet"`
	NumVfs       int  `json:"numVfs"`
	NicSelector  struct {
		Vendor      string `json:"vendor"`
		DeviceID    string `json:"deviceID"`
		PfNames     string `json:"pfNames"`
		RootDevices string `json:"rootDevices"`
	} `json:"nicSelector"`
	DeviceType string `json:"deviceType"`
}

type ResourceConfigList struct {
	Resources []ResourceConfig `json:"resources"`
}

func GetResourceConfigList(file string) ResourceConfigList {
	var config ResourceConfigList
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
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

func (c *ResourceConfig) GetVendor() string {
	return c.NicSelector.Vendor
}

func (c *ResourceConfig) GetDeviceID() string {
	return c.NicSelector.DeviceID
}

func (c *ResourceConfig) GetPfNames() string {
	return c.NicSelector.PfNames
}

func (c *ResourceConfig) GetRootDevices() string {
	return c.NicSelector.RootDevices
}

func (c *ResourceConfig) GetDeviceType() string {
	return c.DeviceType
}
