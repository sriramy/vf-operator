package devices

import (
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

func (d *NetDevice) configure(c *config.ResourceConfig) error {
	err := utils.SetLinkMtu(d.Name, c.Mtu)
	if err != nil {
		return err
	}

	if utils.IsSriovPF(&d.device.Address) {
		err := utils.SetNumVfs(&d.device.Address, c.NumVfs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *NetDevice) getVfPCIs() ([]string, error) {
	devices := make([]string, 0)

	numVfs := utils.GetNumVfs(&d.device.Address)
	for vfIndex := uint32(0); vfIndex < numVfs; vfIndex++ {
		vfPciAddress, err := utils.GetVfPciAddressFromVFIndex(&d.device.Address, vfIndex)
		if err != nil {
			return devices, err
		}
		devices = append(devices, vfPciAddress)
	}

	return devices, nil
}
