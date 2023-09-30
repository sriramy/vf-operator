package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var sysBusPci = "/sys/bus/pci/devices"

const (
	totalVfsFile = "sriov_totalvfs"
	numVfsFile   = "sriov_numvfs"
	physFnFile   = "physfn"
)

func fileExists(pciAddress *string, file string) bool {
	filePath := filepath.Join(sysBusPci, *pciAddress, file)
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
}

func readFile(pciAddress *string, file string) int {
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

	return actualVal
}

func writeFile(pciAddress *string, val int, file string) error {
	filePath := filepath.Join(sysBusPci, *pciAddress, file)
	bs := []byte(strconv.Itoa(val))
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

func GetTotalVfs(pciAddress *string) int {
	return readFile(pciAddress, totalVfsFile)
}

func GetNumVfs(pciAddress *string) int {
	return readFile(pciAddress, numVfsFile)
}

func SetNumVfs(pciAddress *string, numVfs int) error {
	totalVfs := GetTotalVfs(pciAddress)
	if numVfs > totalVfs {
		return fmt.Errorf("sriov_numvfs (%d) > sriov_total_vfs (%d)", numVfs, totalVfs)
	}
	return writeFile(pciAddress, numVfs, numVfsFile)
}
