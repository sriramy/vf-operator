/*
 Copyright (c) 2023 Sriram Yagaraman

 Permission is hereby granted, free of charge, to any person obtaining a copy of
 this software and associated documentation files (the "Software"), to deal in
 the Software without restriction, including without limitation the rights to
 use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 the Software, and to permit persons to whom the Software is furnished to do so,
 subject to the following conditions:

 The above copyright notice and this permission notice shall be included in all
 copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package devices

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jaypipes/ghw"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

const PCI_CLASS_NET = 0x02

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
	for _, device := range p.pci.ListDevices() {
		devClass, err := strconv.ParseInt(device.Class.ID, 16, 64)
		if err != nil {
			fmt.Printf("unable to get dev class %+v %q", device, err)
			continue
		}
		if devClass != PCI_CLASS_NET {
			continue
		}

		name := ""
		macAddress := ""
		nic, _ := p.getNic(device.Address)
		if nic != nil {
			name = nic.Name
			macAddress = nic.MacAddress
		}
		if p.filter(c, device, name) {
			p.devices = append(p.devices, &NetDevice{
				Name:       name,
				MACAddress: macAddress,
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
	if name != "" {
		for _, v := range pfNames {
			if v == name {
				pfNameMatch = true
				break
			}
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

func (p *NetDeviceProvider) Scan() error {
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
		device := p.pci.GetDevice(vfPciAddress)
		netdevice := &NetDevice{
			PCIAddress: vfPciAddress,
			Vendor:     device.Vendor.ID,
			Driver:     device.Driver,
			device:     device,
		}

		vf, _ := p.getNic(vfPciAddress)
		if vf != nil {
			netdevice.Name = vf.Name
			netdevice.MACAddress = vf.MacAddress
		}
		devices = append(devices, netdevice)
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
