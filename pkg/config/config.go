package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Mtu          int  `json:"mtu"`
	NeedVhostNet bool `json:"needVhostNet"`
	NumVfs       int  `json:"numVfs"`
	NicSelector  struct {
		Vendor      string `json:"vendor"`
		DeviceID    string `json:"deviceID"`
		PfNames     string `json:"pfNames"`
		RootDevices string `json:"rootDevices"`
	}
	DeviceType string `json:"deviceType"`
}

func GetConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func (c *Config) GetMtu() int {
	return c.Mtu
}

func (c *Config) GetVhostNet() bool {
	return c.NeedVhostNet
}

func (c *Config) GetVendor() string {
	return c.NicSelector.Vendor
}

func (c *Config) GetDeviceID() string {
	return c.NicSelector.DeviceID
}

func (c *Config) GetPfNames() string {
	return c.NicSelector.PfNames
}

func (c *Config) GetRootDevices() string {
	return c.NicSelector.RootDevices
}

func (c *Config) GetDeviceType() string {
	return c.DeviceType
}
