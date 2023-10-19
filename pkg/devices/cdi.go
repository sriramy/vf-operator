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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/container-orchestrated-devices/container-device-interface/pkg/cdi"
	"github.com/container-orchestrated-devices/container-device-interface/specs-go"
	"github.com/sriramy/vf-operator/pkg/utils"
)

var (
	cdiDir   = "/etc/cdi"
	tunDir   = "/dev/net/tun"
	vhostDir = "/dev/vhost-net"
	vfioDir  = "/dev/vfio/vfio"
)

func init() {
	os.MkdirAll(cdiDir, os.ModePerm)
}

func generateVhostCDISpec(resourceName string) error {
	vendor := "vf-operator"
	deviceNodes := make([]*specs.DeviceNode, 0)
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        tunDir,
		HostPath:    tunDir,
		Permissions: "rw",
	})
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        vhostDir,
		HostPath:    vhostDir,
		Permissions: "rw",
	})
	devices := make([]specs.Device, 0)
	devices = append(devices, specs.Device{
		Name:           "vhost",
		ContainerEdits: specs.ContainerEdits{DeviceNodes: deviceNodes},
	})
	spec := &specs.Spec{
		Kind:    vendor + "/" + resourceName,
		Devices: devices,
	}
	version, err := cdi.MinimumRequiredVersion(spec)
	if err != nil {
		return err
	}
	spec.Version = version

	file := filepath.Join(cdiDir, resourceName+"-vhost.json")
	json, _ := json.MarshalIndent(spec, "", " ")

	return os.WriteFile(file, json, 0o644)
}

func generateVfioCDISpec(resourceName string, pciAddress string) error {
	vendor := "vf-operator"
	deviceNodes := make([]*specs.DeviceNode, 0)
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        vfioDir,
		HostPath:    vfioDir,
		Permissions: "rw",
	})

	vfioDevHost, vfioDevContainer, err := utils.GetVFIODeviceFile(pciAddress)
	if err != nil {
		fmt.Printf("Cannot get vfio device file for: %s, %s\n", pciAddress, err.Error())
	} else {
		deviceNodes = append(deviceNodes, &specs.DeviceNode{
			Path:        vfioDevContainer,
			HostPath:    vfioDevHost,
			Permissions: "rw",
		})
	}
	devices := make([]specs.Device, 0)
	devices = append(devices, specs.Device{
		Name:           "vfio",
		ContainerEdits: specs.ContainerEdits{DeviceNodes: deviceNodes},
	})
	spec := &specs.Spec{
		Kind:    vendor + "/" + resourceName,
		Devices: devices,
	}
	version, err := cdi.MinimumRequiredVersion(spec)
	if err != nil {
		return err
	}
	spec.Version = version

	file := filepath.Join(cdiDir, resourceName+"-vfio.json")
	json, _ := json.MarshalIndent(spec, "", " ")

	return os.WriteFile(file, json, 0o644)
}
