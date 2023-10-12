# vf-operator
[![Build check](https://github.com/sriramy/vf-operator/actions/workflows/build_check.yml/badge.svg)](https://github.com/sriramy/vf-operator/actions/workflows/build_check.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sriramy/vf-operator.svg)](https://pkg.go.dev/github.com/sriramy/vf-operator)

VF operator for CNI only environments, like podman. If you are running with k8s, there are a number of options to configure and associate VFs to your container workloads.
For e.g. https://github.com/k8snetworkplumbingwg/sriov-network-operator 

Most of the configuration fields for the resources are inspired by sriov-network-operator

## Features
* Discover SR-IOV capable NICs based on the following selectors
  * Vendor names
  * PCI address list
  * Driver names
  * PF netdevice names
* Initialize SR-IOV capable NICs based on resource configuration
  * Set MTU
  * Set number of VFs
* Create network conflist for podman CNI backend based on network attachment configuration
  * For network attachments with resourceName specified, the resources discovered in above steps will be used to select a VF
  * Selected VF device ID is inserted into the input network attachment configuration

## gRPC server configuration API
[Protobuf format](https://github.com/sriramy/vf-operator/blob/main/docs/network/proto.md#networkservice-InitialConfiguration)

## Resource configuration example
To configure 4 VFs on a PF with name eth3, use the following configuration
```
curl -X 'POST' \
  'http://localhost:15001/api/v1/config/resources' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
		"name": "eth3-vfs",
		"mtu": 1500,
		"needVhostNet": false,
		"numVfs": 4,
		"nicSelector": {
				"drivers": ["igb"],
				"pfNames": ["eth3"]
		},
		"deviceType": "netdevice"
}
'
```

### Fetch all resource configurations, to verify that it is configured
```
curl -X 'GET' \
  'http://localhost:15001/api/v1/config/resources' \
  -H 'accept: application/json'
```
Example output:
```
{
  "resourceConfigs": [
    {
      "name": "eth3-vfs",
      "mtu": 1500,
      "numVfs": 4,
      "nicSelector": {
        "vendors": [],
        "drivers": [
          "igb"
        ],
        "devices": [],
        "pfNames": [
          "eth3"
        ]
      },
      "deviceType": "netdevice"
    }
  ]
}
```
## Fetching all discovered resources based on configuration
Use the API at /api/v1/resources
```
curl -X 'GET' \
  'http://localhost:15001/api/v1/resources' \
  -H 'accept: application/json' | jq

```
Example output
```
{
  "resources": [
    {
      "spec": {
        "name": "eth3-vfs",
        "mtu": 1500,
        "numVfs": 4,
        "devices": [
          "0000:02:00.0"
        ]
      },
      "status": [
        {
          "name": "eth3",
          "mtu": 1500,
          "numVfs": 4,
          "mac": "00:00:00:01:03:01",
          "vendor": "8086",
          "driver": "igb",
          "device": "0000:02:00.0",
          "vfs": [
            {
              "name": "eth5",
              "mac": "",
              "vendor": "8086",
              "driver": "igbvf",
              "device": "0000:02:10.0"
            },
            {
              "name": "eth6",
              "mac": "",
              "vendor": "8086",
              "driver": "igbvf",
              "device": "0000:02:10.2"
            },
            {
              "name": "eth7",
              "mac": "",
              "vendor": "8086",
              "driver": "igbvf",
              "device": "0000:02:10.4"
            },
            {
              "name": "eth8",
              "mac": "",
              "vendor": "8086",
              "driver": "igbvf",
              "device": "0000:02:10.6"
            }
          ]
        }
      ]
    }
  ]
}
```

## Network attachment example
Any valid network attachment configuration can be used. If resourceName field is defined, vf-operator will assign a valid VF device ID for the network attachment.
```
curl -X 'POST' \
  'http://localhost:15001/api/v1/config/networkattachments' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "f1-u",
	"resourceName": "eth3-vfs",
	"config": {
		"cniVersion": "1.0.0",
		"plugins": [
			{
				"type": "sriov"
			}
		]
	}
}
'
```
## Fetching all network attachments
```
curl -X 'GET' \
  'http://localhost:15001/api/v1/networkattachments' \
  -H 'accept: application/json'
```
Example output:
```
{
  "networkattachments": [
    {
      "name": "f1-u",
      "resourceName": "eth3-vfs",
      "config": {
        "cniVersion": "1.0.0",
        "plugins": [
          {
            "deviceID": "0000:02:10.0",
            "type": "sriov"
          }
        ]
      }
    }
  ]
}
```

## How these network attachments can be used with podman
### List all configured networks for use by podman
```
vm-001 ~ # podman network ls
NETWORK ID    NAME        DRIVER
770afe038c89  f1-u        sriov
2f259bab93aa  podman      bridge
```
### Run a container that uses the network, eth1 is the interface assigned by f1-u network attachment
```
podman run --name testvf --rm -dt --net podman --net f1-u library/alpine:latest
podman exec -ti testvf ip link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0@if18: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP 
    link/ether 86:4a:c2:a5:6a:76 brd ff:ff:ff:ff:ff:ff
11: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 6a:de:4f:de:9c:5a brd ff:ff:ff:ff:ff:ff
podman stop testvf
```
### Run a container with a specific interface name
```
podman run --name testvf --rm -dt --net podman --net f1-u:interface_name=f1-u library/alpine:latest
podman exec -ti testvf ip link
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0@if19: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue state UP 
    link/ether 02:ad:ab:51:b4:1c brd ff:ff:ff:ff:ff:ff
11: f1-u: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
    link/ether 6a:de:4f:de:9c:5a brd ff:ff:ff:ff:ff:ff
podman stop testvf
```
