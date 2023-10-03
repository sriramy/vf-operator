package devices

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/sriramy/vf-operator/pkg/config"
)

type NetDeviceProvider struct {
	config  *config.ResourceConfig
	devices []*NetDevice
}

type NetDevice struct {
	Name       string  `json:"name"`
	MACAddress string  `json:"mac_address"`
	PCIAddress *string `json:"pci_address,omitempty"`
	device     *ghw.PCIDevice
}

// NewNetDeviceProvider DeviceProvider implementation from netDeviceProvider instance
func NewNetDeviceProvider(c *config.ResourceConfig) *NetDeviceProvider {
	return &NetDeviceProvider{
		config:  c,
		devices: make([]*NetDevice, 0),
	}
}

func (p *NetDeviceProvider) GetDevices() []*NetDevice {
	return p.devices
}

func (p *NetDeviceProvider) Discover() error {
	pci, err := ghw.PCI()
	if err != nil {
		return fmt.Errorf("Couldn't get PCI info: %v", err)
	}

	net, err := ghw.Network()
	if err != nil {
		return fmt.Errorf("Couldn't get NIC info: %v", err)
	}

	for _, nic := range net.NICs {
		if nic.PCIAddress == nil {
			continue
		}

		device := pci.GetDevice(*nic.PCIAddress)
		if p.filter(device, nic.Name) {
			p.devices = append(p.devices, &NetDevice{
				Name:       nic.Name,
				MACAddress: nic.MacAddress,
				PCIAddress: nic.PCIAddress,
				device:     device,
			})
		}
	}

	return nil
}

func (p *NetDeviceProvider) filter(dev *ghw.PCIDevice, name string) bool {
	vendors := p.config.GetVendors()
	vendorMatch := (len(vendors) == 0)
	for _, v := range vendors {
		if v == dev.Vendor.ID {
			vendorMatch = true
			break
		}
	}

	pfNames := p.config.GetPfNames()
	pfNameMatch := (len(pfNames) == 0)
	for _, v := range pfNames {
		if v == name {
			pfNameMatch = true
			break
		}
	}

	drivers := p.config.GetDrivers()
	driverMatch := (len(drivers) == 0)
	for _, v := range drivers {
		if v == dev.Driver {
			driverMatch = true
			break
		}
	}

	devices := p.config.GetDevices()
	deviceMatch := (len(devices) == 0)
	for _, v := range devices {
		if v == dev.Address {
			deviceMatch = true
			break
		}
	}

	return vendorMatch && pfNameMatch && driverMatch && deviceMatch
}
