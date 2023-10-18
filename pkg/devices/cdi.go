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
