# vf-operator
[![Build check](https://github.com/sriramy/vf-operator/actions/workflows/build_check.yml/badge.svg)](https://github.com/sriramy/vf-operator/actions/workflows/build_check.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sriramy/vf-operator.svg)](https://pkg.go.dev/github.com/sriramy/vf-operator)

VF operator for podman

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

