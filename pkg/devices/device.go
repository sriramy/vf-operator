package devices

import (
	"fmt"

	"github.com/jaypipes/ghw"
	"github.com/sriramy/vf-operator/pkg/config"
	"github.com/sriramy/vf-operator/pkg/utils"
)

type NetDevice struct {
	Name       string
	MACAddress string
	PCIAddress string
	Vendor     string
	Driver     string

	device *ghw.PCIDevice
}

func (d *NetDevice) Configure(c *config.ResourceConfig) error {
	if utils.IsSriovPF(&d.device.Address) {
		err := utils.SetNumVfs(&d.device.Address, c.NumVfs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *NetDevice) GetVfs() ([]NetDevice, error) {
	devices := make([]NetDevice, 0)

	pci, err := ghw.PCI()
	if err != nil {
		return devices, fmt.Errorf("Couldn't get PCI info: %v", err)
	}

	numVfs := utils.GetNumVfs(&d.device.Address)
	for vfIndex := 0; vfIndex < numVfs; vfIndex++ {
		vfPciAddress, err := utils.GetVfPciAddressFromVFIndex(&d.device.Address, vfIndex)
		if err != nil {
			return devices, err
		}
		device := pci.GetDevice(vfPciAddress)
		nic, err := utils.GetVfNic(device)
		if err != nil {
			return devices, err
		}

		devices = append(devices, NetDevice{
			Name:       nic.Name,
			MACAddress: nic.MacAddress,
			PCIAddress: device.Address,
			Vendor:     device.Vendor.ID,
			Driver:     device.Driver,
			device:     device,
		})
	}

	return devices, nil
}
