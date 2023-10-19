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

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetVFIODeviceFile returns a vfio device files for vfio-pci bound PCI device's PCI address
func GetVFIODeviceFile(dev string) (devFileHost, devFileContainer string, err error) {
	// Get iommu group for this device
	devPath := filepath.Join(sysBusPciDevice, dev)
	_, err = os.Lstat(devPath)
	if err != nil {
		err = fmt.Errorf("GetVFIODeviceFile(): Could not get directory information for device: %s, Err: %v", dev, err)
		return devFileHost, devFileContainer, err
	}

	iommuDir := filepath.Join(devPath, "iommu_group")
	if err != nil {
		err = fmt.Errorf("GetVFIODeviceFile(): error reading iommuDir %v", err)
		return devFileHost, devFileContainer, err
	}

	dirInfo, err := os.Lstat(iommuDir)
	if err != nil {
		err = fmt.Errorf("GetVFIODeviceFile(): unable to find iommu_group %v", err)
		return devFileHost, devFileContainer, err
	}

	if dirInfo.Mode()&os.ModeSymlink == 0 {
		err = fmt.Errorf("GetVFIODeviceFile(): invalid symlink to iommu_group %v", err)
		return devFileHost, devFileContainer, err
	}

	linkName, err := filepath.EvalSymlinks(iommuDir)
	if err != nil {
		err = fmt.Errorf("GetVFIODeviceFile(): error reading symlink to iommu_group %v", err)
		return devFileHost, devFileContainer, err
	}
	devFileContainer = filepath.Join(devDir, "vfio", filepath.Base(linkName))
	devFileHost = devFileContainer

	// Get a file path to the iommu group name
	namePath := filepath.Join(linkName, "name")
	// Read the iommu group name
	// The name file will not exist on baremetal
	vfioName, errName := os.ReadFile(namePath)
	if errName == nil {
		vName := strings.TrimSpace(string(vfioName))

		// if the iommu group name == vfio-noiommu then we are in a VM, adjust path to vfio device
		if vName == "vfio-noiommu" {
			linkName = filepath.Join(filepath.Dir(linkName), "noiommu-"+filepath.Base(linkName))
			devFileHost = filepath.Join(devDir, "vfio", filepath.Base(linkName))
		}
	}

	return devFileHost, devFileContainer, err
}
