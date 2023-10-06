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
