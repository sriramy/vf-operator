package networkattachment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

var cniNetDir = "/etc/cni/net.d"

type networkattachment struct {
	config map[string]interface{}
}

func newNetworkAttachment(cniConfig map[string]interface{}) *networkattachment {
	os.MkdirAll(cniNetDir, os.ModePerm)
	return &networkattachment{config: cniConfig}
}

func getAllNetworkAttachments() []*networkattachment {
	cniConfigs := make([]*networkattachment, 0)
	err := filepath.Walk(cniNetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() {
			cniConfig, err := getNetworkAttachment(strings.TrimSuffix(info.Name(), filepath.Ext(path)))
			if err == nil {
				cniConfigs = append(cniConfigs, cniConfig)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	return cniConfigs
}

func getNetworkAttachment(name string) (*networkattachment, error) {
	var cniConfig map[string]interface{}
	file := filepath.Join(cniNetDir, name+".conflist")
	conf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(conf, &cniConfig)
	if err != nil {
		return nil, fmt.Errorf("addDeviceIDInConfList: failed to unmarshal inBytes: %v", err)
	}

	return newNetworkAttachment(cniConfig), nil
}

func IsDeviceIDAllocated(pciAddress string) bool {
	for _, n := range getAllNetworkAttachments() {
		deviceID, err := n.getDeviceID()
		if err == nil && deviceID == pciAddress {
			return true
		}
	}

	return false
}

func (n *networkattachment) addName(name string, resourceName string) {
	n.config["name"] = name
	n.config["resourceName"] = resourceName
}

func (n *networkattachment) addDeviceID(pciAddress string) error {
	pList, ok := n.config["plugins"]
	if !ok {
		if _, ok := n.config["type"]; ok {
			n.config["deviceID"] = pciAddress
			return nil
		}
		return fmt.Errorf("AddDeviceIDToNetworkAttachment: plugins nor type key found")
	}

	pMap, ok := pList.([]interface{})
	if !ok {
		return fmt.Errorf("AddDeviceIDToNetworkAttachment: unable to typecast plugin list")
	}

	for idx, plugin := range pMap {
		currentPlugin, ok := plugin.(map[string]interface{})
		if !ok {
			return fmt.Errorf("AddDeviceIDToNetworkAttachment: unable to typecast plugin #%d", idx)
		}
		currentPlugin["deviceID"] = pciAddress
	}

	return nil
}

func (n *networkattachment) getDeviceID() (string, error) {
	pList, ok := n.config["plugins"]
	if !ok {
		if _, ok := n.config["type"]; ok {
			deviceID, ok := n.config["deviceID"].(string)
			if ok {
				return deviceID, nil
			}
		}
	} else {
		pMap, ok := pList.([]interface{})
		if !ok {
			return "", fmt.Errorf("plugins key not found")
		}

		for _, plugin := range pMap {
			currentPlugin, ok := plugin.(map[string]interface{})
			if !ok {
				continue
			}

			deviceID, ok := currentPlugin["deviceID"].(string)
			if ok {
				return deviceID, nil
			}

		}
	}

	return "", fmt.Errorf("deviceID not found")
}

func (n *networkattachment) create(name string) error {
	file := filepath.Join(cniNetDir, name+".conflist")
	json, _ := json.MarshalIndent(n.config, "", " ")

	return os.WriteFile(file, json, 0o644)
}

func (n *networkattachment) delete(name string) error {
	file := filepath.Join(cniNetDir, name+".conflist")
	return os.Remove(file)
}

func (n *networkattachment) build() (*network.NetworkAttachment, error) {
	// obtain names from network attachment
	name, ok := n.config["name"].(string)
	if ok {
		delete(n.config, "name")
	}
	resourceName, ok := n.config["resourceName"].(string)
	if ok {
		delete(n.config, "resourceName")
	}

	config, err := structpb.NewStruct(n.config)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "network attachment config id=%s not valid", name)
	}
	return &network.NetworkAttachment{
		Name:         name,
		ResourceName: resourceName,
		Config:       config,
	}, nil
}
