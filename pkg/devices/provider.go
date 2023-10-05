package devices

import (
	"fmt"
	"log"

	"github.com/jaypipes/ghw"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

type NetDeviceProvider struct {
	devices []*NetDevice
	pci     *ghw.PCIInfo
	net     *ghw.NetworkInfo
}

// NewNetDeviceProvider DeviceProvider implementation from netDeviceProvider instance
func NewNetDeviceProvider() *NetDeviceProvider {
	pci, err := ghw.PCI()
	if err != nil {
		log.Fatalf("Couldn't get PCI info: %v", err)
	}

	net, err := ghw.Network()
	if err != nil {
		log.Fatalf("Couldn't get NIC info: %v", err)
	}

	return &NetDeviceProvider{
		devices: make([]*NetDevice, 0),
		pci:     pci,
		net:     net,
	}
}

func (p *NetDeviceProvider) GetDevices() []*NetDevice {
	return p.devices
}

func (p *NetDeviceProvider) Discover(c *network.ResourceConfig) error {
	for _, nic := range p.net.NICs {
		if nic.PCIAddress == nil {
			continue
		}

		device := p.pci.GetDevice(*nic.PCIAddress)
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

func (p *NetDeviceProvider) filter(c *network.ResourceConfig, dev *ghw.PCIDevice, name string) bool {
	vendors := c.GetNicSelector().GetVendors()
	vendorMatch := (len(vendors) == 0)
	for _, v := range vendors {
		if v == dev.Vendor.ID {
			vendorMatch = true
			break
		}
	}

	pfNames := c.GetNicSelector().GetPfNames()
	pfNameMatch := (len(pfNames) == 0)
	for _, v := range pfNames {
		if v == name {
			pfNameMatch = true
			break
		}
	}

	drivers := c.GetNicSelector().GetDrivers()
	driverMatch := (len(drivers) == 0)
	for _, v := range drivers {
		if v == dev.Driver {
			driverMatch = true
			break
		}
	}

	devices := c.GetNicSelector().GetDevices()
	deviceMatch := (len(devices) == 0)
	for _, v := range devices {
		if v == dev.Address {
			deviceMatch = true
			break
		}
	}

	return vendorMatch && pfNameMatch && driverMatch && deviceMatch
}

func (p *NetDeviceProvider) Configure(c *network.ResourceConfig) error {
	for _, dev := range p.GetDevices() {
		err := dev.configure(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *NetDeviceProvider) ProbeNics() error {
	pci, err := ghw.PCI()
	if err != nil {
		return fmt.Errorf("Couldn't get PCI info: %v", err)
	}

	net, err := ghw.Network()
	if err != nil {
		return fmt.Errorf("Couldn't get NIC info: %v", err)
	}

	p.pci = pci
	p.net = net
	return nil
}

func (p *NetDeviceProvider) GetVFDevices(dev *NetDevice) []*NetDevice {
	devices := make([]*NetDevice, 0)
	vfPCIs, err := dev.getVfPCIs()
	if err != nil {
		return devices
	}

	for _, vfPciAddress := range vfPCIs {
		vf, err := p.getNic(vfPciAddress)
		if err != nil {
			continue
		}
		device := p.pci.GetDevice(vfPciAddress)
		devices = append(devices, &NetDevice{
			Name:       vf.Name,
			MACAddress: vf.MacAddress,
			PCIAddress: vfPciAddress,
			Vendor:     device.Vendor.ID,
			Driver:     device.Driver,
			device:     device,
		})
	}

	return devices
}

func (p *NetDeviceProvider) getNic(pciAddress string) (*ghw.NIC, error) {
	stringMatch := func(a *string, b string) bool {
		return a != nil && *a == b
	}
	for _, nic := range p.net.NICs {
		if stringMatch(nic.PCIAddress, pciAddress) {
			return nic, nil
		}
	}

	return nil, fmt.Errorf("No NIC found with PCI Address: %s", pciAddress)
}
