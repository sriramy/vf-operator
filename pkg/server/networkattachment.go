package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func BuildNetworkAttachmentConfig(cniConfig map[string]interface{}, pciAddress string) (map[string]interface{}, error) {
	pList, ok := cniConfig["plugins"]
	if !ok {
		if _, ok := cniConfig["type"]; ok {
			cniConfig["deviceID"] = pciAddress
			cniConfig["pciBusID"] = pciAddress
			return cniConfig, nil
		}
		return cniConfig, fmt.Errorf("BuildNetworkAttachmentConfig: plugins nor type key found")
	}

	pMap, ok := pList.([]interface{})
	if !ok {
		return cniConfig, fmt.Errorf("BuildNetworkAttachmentConfig: unable to typecast plugin list")
	}

	for idx, plugin := range pMap {
		currentPlugin, ok := plugin.(map[string]interface{})
		if !ok {
			return cniConfig, fmt.Errorf("BuildNetworkAttachmentConfig: unable to typecast plugin #%d", idx)
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

func RemoveNetworkAttachment(name string) error {
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	return os.Remove(file)
}
