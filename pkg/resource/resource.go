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

package resource

import (
	"fmt"

	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/devices"
	"github.com/sriramy/vf-operator/pkg/utils"
)

type resource struct {
	config   *network.ResourceConfig
	provider *devices.NetDeviceProvider
}

func newResource(c *network.ResourceConfig) *resource {
	r := &resource{
		config: &network.ResourceConfig{
			Name:         c.GetName(),
			Mtu:          c.GetMtu(),
			NumVfs:       c.GetNumVfs(),
			NeedVhostNet: c.GetNeedVhostNet(),
			NicSelector: &network.NicSelector{
				Vendors: c.GetNicSelector().GetVendors(),
				Drivers: c.GetNicSelector().GetDrivers(),
				Devices: c.GetNicSelector().GetDevices(),
				PfNames: c.GetNicSelector().GetPfNames(),
			},
			DeviceType: c.GetDeviceType(),
		},
		provider: devices.NewNetDeviceProvider(),
	}
	if err := r.do(); err != nil {
		fmt.Printf("Error configuring resource id=%s, err:%v\n", r.config.GetName(), err)
	}

	return r
}

func (r *resource) build() *network.Resource {
	discoveredDevices := make([]string, 0)
	devices := make([]*network.ResourceStatus, 0)
	for _, dev := range r.provider.GetDevices() {
		vfDevices := make([]*network.VFResourceStatus, 0)
		for _, vf := range r.provider.GetVFDevices(dev) {
			vfDevices = append(vfDevices, &network.VFResourceStatus{
				Name:   vf.Name,
				Mac:    vf.MACAddress,
				Device: vf.PCIAddress,
				Vendor: vf.Vendor,
				Driver: vf.Driver,
			})
		}

		discoveredDevices = append(discoveredDevices, dev.PCIAddress)
		devices = append(devices, &network.ResourceStatus{
			Name:   dev.Name,
			Mtu:    utils.GetLinkMtu(dev.Name),
			NumVfs: utils.GetNumVfs(&dev.PCIAddress),
			Mac:    dev.MACAddress,
			Device: dev.PCIAddress,
			Vendor: dev.Vendor,
			Driver: dev.Driver,
			Vfs:    vfDevices,
		})
	}
	res := &network.Resource{
		Spec: &network.ResourceSpec{
			Name:         r.config.GetName(),
			Mtu:          r.config.GetMtu(),
			NumVfs:       r.config.GetNumVfs(),
			NeedVhostNet: r.config.GetNeedVhostNet(),
			Devices:      discoveredDevices,
			DeviceType:   r.config.GetDeviceType(),
		},
		Status: devices,
	}
	return res
}

func (r *resource) do() error {
	err := r.provider.Discover(r.config)
	if err != nil {
		return err
	}

	err = r.provider.Configure(r.config)
	if err != nil {
		return err
	}

	err = r.provider.Scan()
	if err != nil {
		return err
	}

	return nil
}
