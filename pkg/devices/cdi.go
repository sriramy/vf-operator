package devices

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/container-orchestrated-devices/container-device-interface/pkg/cdi"
	"github.com/container-orchestrated-devices/container-device-interface/specs-go"
)

const cdiDir = "/etc/cdi"

func init() {
	os.MkdirAll(cdiDir, os.ModePerm)
}

func generateVhostCDISpec(resourceName string) error {
	vendor := "vf-operator"
	class := "vhost"
	deviceNodes := make([]*specs.DeviceNode, 0)
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        "/dev/net/tun",
		HostPath:    "/dev/net/tun",
		Permissions: "rw",
	})
	deviceNodes = append(deviceNodes, &specs.DeviceNode{
		Path:        "/dev/vhost-net",
		HostPath:    "/dev/vhost-net",
		Permissions: "rw",
	})
	devices := make([]specs.Device, 0)
	devices = append(devices, specs.Device{
		Name:           "net",
		ContainerEdits: specs.ContainerEdits{DeviceNodes: deviceNodes},
	})
	spec := &specs.Spec{
		Kind:    vendor + "/" + class,
		Devices: devices,
	}
	version, err := cdi.MinimumRequiredVersion(spec)
	if err != nil {
		return err
	}
	spec.Version = version

	file := filepath.Join(cdiDir, resourceName+".json")
	json, _ := json.MarshalIndent(spec, "", " ")

	return os.WriteFile(file, json, 0o644)
}
