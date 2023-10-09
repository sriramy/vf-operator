package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
)

func AddNameToNetworkAttachment(cniConfig map[string]interface{}, name string, resourceName string) (map[string]interface{}, error) {
	cniConfig["name"] = name
	cniConfig["resourceName"] = resourceName
	return cniConfig, nil
}

func AddDeviceIDToNetworkAttachment(cniConfig map[string]interface{}, pciAddress string) (map[string]interface{}, error) {
	pList, ok := cniConfig["plugins"]
	if !ok {
		if _, ok := cniConfig["type"]; ok {
			cniConfig["deviceID"] = pciAddress
			cniConfig["pciBusID"] = pciAddress
			return cniConfig, nil
		}
		return cniConfig, fmt.Errorf("AddDeviceIDToNetworkAttachment: plugins nor type key found")
	}

	pMap, ok := pList.([]interface{})
	if !ok {
		return cniConfig, fmt.Errorf("AddDeviceIDToNetworkAttachment: unable to typecast plugin list")
	}

	for idx, plugin := range pMap {
		currentPlugin, ok := plugin.(map[string]interface{})
		if !ok {
			return cniConfig, fmt.Errorf("AddDeviceIDToNetworkAttachment: unable to typecast plugin #%d", idx)
		}
		currentPlugin["deviceID"] = pciAddress
		currentPlugin["pciBusID"] = pciAddress
	}

	return cniConfig, nil
}

func AddNetworkAttachment(name string, config map[string]interface{}) error {
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	json, _ := json.MarshalIndent(config, "", " ")

	return os.WriteFile(file, json, 0o644)
}

func GetAllNetworkAttachments() []map[string]interface{} {
	cniConfigs := make([]map[string]interface{}, 0)
	names := make([]string, 0)
	err := filepath.Walk("/etc/cni/net.d", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			names = append(names, strings.TrimSuffix(info.Name(), filepath.Ext(path)))
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, name := range names {
		cniConfig, err := GetNetworkAttachment(name)
		if err == nil {
			cniConfigs = append(cniConfigs, cniConfig)
		}
	}
	return cniConfigs
}

func GetNetworkAttachment(name string) (map[string]interface{}, error) {
	var cniConfig map[string]interface{}
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	conf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(conf, &cniConfig)
	if err != nil {
		return nil, fmt.Errorf("addDeviceIDInConfList: failed to unmarshal inBytes: %v", err)
	}

	return cniConfig, nil
}

func RemoveNetworkAttachment(name string) error {
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	return os.Remove(file)
}

func IsAllocated(pciAddress string) network.VFStatus {
	for _, cniConfig := range GetAllNetworkAttachments() {
		pList, ok := cniConfig["plugins"]
		if !ok {
			if _, ok := cniConfig["type"]; ok {
				deviceID, ok := cniConfig["deviceID"].(string)
				if ok && deviceID == pciAddress {
					return network.VFStatus_USED
				}
			}
			continue
		}

		pMap, ok := pList.([]interface{})
		if !ok {
			continue
		}

		for _, plugin := range pMap {
			currentPlugin, ok := plugin.(map[string]interface{})
			if !ok {
				continue
			}

			deviceID, ok := currentPlugin["deviceID"].(string)
			if ok && deviceID == pciAddress {
				return network.VFStatus_USED
			}

		}

	}

	return network.VFStatus_FREE
}
