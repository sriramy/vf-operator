# vf-operator
⚠ Experimental. Please take note that this is a pre-release version.

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
  * Create CDI spec for vhost-net and tun devices if needVhostNet is true
* Create network conflist for podman CNI backend based on network attachment configuration
  * For network attachments with resourceName specified, the resources discovered in above steps will be used to select a VF
  * Selected VF device ID is inserted into the input network attachment configuration

## Build
Run the following commands to build and install
```
make -j$(nproc)
DESTDIR="/path/to/install" PREFIX="/usr/local" make install 
```

Complete list of make targets
```
make help

 Make vf-operator, a network service that configures and discovers SR-IOV devices

 Targets;
  help                  This printout
  all (default)         Build gRPC stubs and the executable
  clean                 Remove built files
  dep                   Installs pre-requisites
  stubs                 Generates gRPC stubs and OpenAPIv2 specs
  test
  install               Installs the executable
  swagger_install       Installs static swagger UI

Binary:
  bin/vf-operator
```

## Test
Start the VF operator and test the exposed [APIs](docs/DESIGN.md)
```
vf-operator
```

Complete list of arguments
```
vf-operator --help

vf-operator discovers and configures SR-IOV capable NICs based on configured selectors and
creates network conflist for container runtimes to attach CNI networks to containers
Options;
  -config string
        Path to config file (default "/etc/cni/vf-operator/config.json")
  -gwPort int
        API port (tcp:port) (default 15001)
  -port int
        gRPC port (tcp:port) (default 5001)
```

## Basic design and API references
[Design](docs/DESIGN.md)

## Contributing
[How to contribute](docs/CONTRIBUTING.md)
