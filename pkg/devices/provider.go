package devices

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/k8snetworkplumbingwg/sriovnet"
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
	Vf         bool    `json:"vf"`
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
			vf := false
			_, err := sriovnet.GetPfPciFromVfPci(*nic.PCIAddress)
			if err == nil {
				vf = true
			}
			p.devices = append(p.devices, &NetDevice{
				Name:       nic.Name,
				MACAddress: nic.MacAddress,
				PCIAddress: nic.PCIAddress,
				Vf:         vf,
				device:     device,
			})
		}
	}

	return nil
}

func (p *NetDeviceProvider) filter(dev *ghw.PCIDevice, name string) bool {
	vendorMatch := (len(p.config.GetVendors()) == 0)
	for _, v := range p.config.GetVendors() {
		if v == dev.Vendor.ID {
			vendorMatch = true
			break
		}
	}

	pfNameMatch := (len(p.config.GetPfNames()) == 0)
	for _, v := range p.config.GetPfNames() {
		if v == name {
			pfNameMatch = true
			break
		}
	}

	driverMatch := (len(p.config.GetDrivers()) == 0)
	for _, v := range p.config.GetDrivers() {
		if v == dev.Driver {
			driverMatch = true
			break
		}
	}

	deviceMatch := (len(p.config.GetDevices()) == 0)
	for _, v := range p.config.GetDevices() {
		if v == dev.Address {
			deviceMatch = true
			break
		}
	}

	return vendorMatch && pfNameMatch && driverMatch && deviceMatch
}
