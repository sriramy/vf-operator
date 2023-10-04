package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	sysBusPci   = "/sys/bus/pci/devices"
	sysClassNet = "/sys/class/net"
)

const (
	totalVfsFile = "sriov_totalvfs"
	numVfsFile   = "sriov_numvfs"
	physFnFile   = "physfn"
	virtFnPrefix = "virtfn"
)

func dirExists(pciAddress *string, file string) bool {
	dirPath := filepath.Join(sysBusPci, *pciAddress, file)
	info, err := os.Stat(dirPath)
	return err == nil && info.IsDir()
}

func fileExists(pciAddress *string, file string) bool {
	filePath := filepath.Join(sysBusPci, *pciAddress, file)
	info, err := os.Stat(filePath)
	return err == nil && !info.IsDir()
}

func readFile(pciAddress *string, file string) uint32 {
	filePath := filepath.Join(sysBusPci, *pciAddress, file)
	val, err := os.ReadFile(filePath)
	if err != nil {
		return 0
	}

	trimmedVal := bytes.TrimSpace(val)
	actualVal, err := strconv.Atoi(string(trimmedVal))
	if err != nil {
		return 0
	}

	return uint32(actualVal)
}

func writeFile(pciAddress *string, val uint32, file string) error {
	filePath := filepath.Join(sysBusPci, *pciAddress, file)
	bs := []byte(strconv.Itoa(int(val)))
	err := os.WriteFile(filePath, []byte("0"), os.ModeAppend)
	if err != nil {
		fmt.Printf("write(): fail to reset file %s", filePath)
		return err
	}
	err = os.WriteFile(filePath, bs, os.ModeAppend)
	if err != nil {
		fmt.Printf("write(): fail to set file %s with %d", filePath, val)
		return err
	}
	return nil
}

func IsSriovPF(pciAddress *string) bool {
	return fileExists(pciAddress, totalVfsFile)
}

func IsSriovVF(pciAddress *string) bool {
	return fileExists(pciAddress, physFnFile)
}

func GetTotalVfs(pciAddress *string) uint32 {
	return readFile(pciAddress, totalVfsFile)
}

func GetNumVfs(pciAddress *string) uint32 {
	return readFile(pciAddress, numVfsFile)
}

func SetNumVfs(pciAddress *string, numVfs uint32) error {
	totalVfs := GetTotalVfs(pciAddress)
	if numVfs > totalVfs {
		return fmt.Errorf("sriov_numvfs (%d) > sriov_total_vfs (%d)", numVfs, totalVfs)
	}

	if numVfs != GetNumVfs(pciAddress) {
		return writeFile(pciAddress, numVfs, numVfsFile)
	}

	return nil
}

func GetVfPciAddressFromVFIndex(pciAddress *string, vfIndex uint32) (string, error) {
	virtFn := fmt.Sprintf("%s%v", virtFnPrefix, vfIndex)
	virtFnLink := filepath.Join(sysBusPci, *pciAddress, virtFn)

	vfPciDir, err := os.Readlink(virtFnLink)
	if len(vfPciDir) <= 3 {
		return "", fmt.Errorf("Could not parse PCI Address for PF %s VF %s",
			*pciAddress, virtFn)
	}

	return vfPciDir[3:], err
}
