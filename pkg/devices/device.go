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
	"github.com/jaypipes/ghw"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/utils"
)

type NetDevice struct {
	Name       string
	MACAddress string
	PCIAddress string
	Vendor     string
	OrigDriver string
	Driver     string

	device *ghw.PCIDevice
}

func (d *NetDevice) configure(c *network.ResourceConfig) error {
	err := utils.SetLinkMtu(d.Name, c.GetMtu())
	if err != nil {
		return err
	}

	if utils.IsSriovPF(&d.device.Address) {
		if c.GetDeviceType() == DeviceTypeVfioPci && d.Driver != DeviceTypeVfioPci {
			err = utils.DriverOp(&d.device.Address, d.Driver, "unbind")
			if err != nil {
				return err
			}
			err = utils.DriverOp(&d.device.Address, c.GetDeviceType(), "bind")
			if err != nil {
				utils.DriverOp(&d.device.Address, d.Driver, "bind")
				return err
			}
			// store configured driver name
			d.OrigDriver = d.Driver
			d.Driver = c.GetDeviceType()
		}

		err = utils.SetNumVfs(&d.device.Address, c.GetNumVfs())
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
