package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func BuildNetworkAttachmentConfig(cniConfig string, pciAddress string) (string, error) {
	var rawConfig map[string]interface{}
	var err error

	err = json.Unmarshal([]byte(cniConfig), &rawConfig)
	if err != nil {
		return "", fmt.Errorf("BuildNetworkAttachmentConfig: failed to unmarshal: %v", err)
	}

	pList, ok := rawConfig["plugins"]
	if !ok {
		return "", fmt.Errorf("BuildNetworkAttachmentConfig: unable to get plugin list")
	}

	pMap, ok := pList.([]interface{})
	if !ok {
		return "", fmt.Errorf("BuildNetworkAttachmentConfig: unable to typecast plugin list")
	}

	for idx, plugin := range pMap {
		currentPlugin, ok := plugin.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("BuildNetworkAttachmentConfig: unable to typecast plugin #%d", idx)
		}
		// Inject deviceID
		currentPlugin["deviceID"] = pciAddress
		currentPlugin["pciBusID"] = pciAddress
	}

	configBytes, err := json.Marshal(rawConfig)
	if err != nil {
		return "", fmt.Errorf("BuildNetworkAttachmentConfig: failed to marshal: %v", err)
	}

	return string(configBytes), nil
}

func AddNetworkAttachment(name string, config string) error {
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	json, _ := json.MarshalIndent(config, "", " ")

	return os.WriteFile(file, json, 0o644)
}

func RemoveNetworkAttachment(name string) error {
	file := filepath.Join("/etc/cni/net.d", name+".conflist")
	return os.Remove(file)
}
