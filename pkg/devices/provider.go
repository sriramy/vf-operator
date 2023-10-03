package devices

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/sriramy/vf-operator/pkg/config"
)

type NetDeviceProvider struct {
	devices []*NetDevice
}

// NewNetDeviceProvider DeviceProvider implementation from netDeviceProvider instance
func NewNetDeviceProvider() *NetDeviceProvider {
	return &NetDeviceProvider{
		devices: make([]*NetDevice, 0),
	}
}

func (p *NetDeviceProvider) GetDevices() []*NetDevice {
	return p.devices
}

func (p *NetDeviceProvider) Discover(c *config.ResourceConfig) error {
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
		if p.filter(c, device, nic.Name) {
			p.devices = append(p.devices, &NetDevice{
				Name:       nic.Name,
				MACAddress: nic.MacAddress,
				PCIAddress: device.Address,
				Vendor:     device.Vendor.ID,
				Driver:     device.Driver,
				device:     device,
			})
		}
	}

	return nil
}

func (p *NetDeviceProvider) filter(c *config.ResourceConfig, dev *ghw.PCIDevice, name string) bool {
	vendors := c.GetVendors()
	vendorMatch := (len(vendors) == 0)
	for _, v := range vendors {
		if v == dev.Vendor.ID {
			vendorMatch = true
			break
		}
	}

	pfNames := c.GetPfNames()
	pfNameMatch := (len(pfNames) == 0)
	for _, v := range pfNames {
		if v == name {
			pfNameMatch = true
			break
		}
	}

	drivers := c.GetDrivers()
	driverMatch := (len(drivers) == 0)
	for _, v := range drivers {
		if v == dev.Driver {
			driverMatch = true
			break
		}
	}

	devices := c.GetDevices()
	deviceMatch := (len(devices) == 0)
	for _, v := range devices {
		if v == dev.Address {
			deviceMatch = true
			break
		}
	}

	return vendorMatch && pfNameMatch && driverMatch && deviceMatch
}

func (p *NetDeviceProvider) Configure(c *config.ResourceConfig) error {
	for _, dev := range p.GetDevices() {
		err := dev.Configure(c)
		if err != nil {
			return err
		}
	}
	return nil
}
