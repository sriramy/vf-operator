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

package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/devices"
	"github.com/sriramy/vf-operator/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type resource struct {
	config   *network.ResourceConfig
	provider *devices.NetDeviceProvider
}

type NetworkServiceServer struct {
	network.UnimplementedNetworkServiceServer
	resources []resource
}

func (s *NetworkServiceServer) CreateResourceConfig(_ context.Context, c *network.ResourceConfig) (*network.Resource, error) {
	r := &resource{
		config: &network.ResourceConfig{
			Name:   c.GetName(),
			Mtu:    c.GetMtu(),
			NumVfs: c.GetNumVfs(),
			NicSelector: &network.NicSelector{
				Vendors: c.GetNicSelector().GetVendors(),
				Drivers: c.GetNicSelector().GetDrivers(),
				Devices: c.GetNicSelector().GetDevices(),
				PfNames: c.GetNicSelector().GetPfNames(),
			},
			DeviceType: c.GetDeviceType(),
		},
		provider: devices.NewNetDeviceProvider(),
	}

	s.resources = append(s.resources, *r)
	s.configure(r)
	return s.buildResource(r), nil
}

func (s *NetworkServiceServer) GetResourceConfig(_ context.Context, id *network.ResourceName) (*network.ResourceConfig, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return r.config, nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) DeleteResourceConfig(_ context.Context, id *network.ResourceName) (*empty.Empty, error) {
	for i, r := range s.resources {
		if r.config.GetName() == id.GetName() {
			s.resources = append(s.resources[:i], s.resources[i+1:]...)
			return new(empty.Empty), nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) GetAllResourceConfigs(context.Context, *empty.Empty) (*network.ResourceConfigs, error) {
	configs := make([]*network.ResourceConfig, 0)
	for _, r := range s.resources {
		configs = append(configs, r.config)
	}
	return &network.ResourceConfigs{Configs: configs}, nil
}

func (s *NetworkServiceServer) GetAllResources(context.Context, *empty.Empty) (*network.Resources, error) {
	resources := make([]*network.Resource, 0)
	for _, r := range s.resources {
		resources = append(resources, s.buildResource(&r))
	}

	return &network.Resources{Resources: resources}, nil
}

func (s *NetworkServiceServer) GetResource(_ context.Context, id *network.ResourceName) (*network.Resource, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return s.buildResource(r), nil
	}

	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) getResource(resourceName string) *resource {
	for _, r := range s.resources {
		if r.config.GetName() == resourceName {
			return &r
		}
	}
	return nil
}

func (s *NetworkServiceServer) buildResource(r *resource) *network.Resource {
	discoveredDevices := make([]string, 0)
	devices := make([]*network.ResourceStatus, 0)
	for _, dev := range r.provider.GetDevices() {
		vfDevices := make([]*network.VFResourceStatus, 0)
		for _, vf := range r.provider.GetVFDevices(dev) {
			vfDevices = append(vfDevices, &network.VFResourceStatus{
				Name:   vf.Name,
				Device: vf.PCIAddress,
				Vendor: vf.Vendor,
				Driver: vf.Driver,
				Status: IsAllocated(vf.PCIAddress),
			})
		}

		discoveredDevices = append(discoveredDevices, dev.PCIAddress)
		devices = append(devices, &network.ResourceStatus{
			Name:   dev.Name,
			Mtu:    utils.GetLinkMtu(dev.Name),
			NumVfs: utils.GetNumVfs(&dev.PCIAddress),
			Mac:    dev.MACAddress,
			Device: dev.PCIAddress,
			Vendor: dev.Vendor,
			Driver: dev.Driver,
			Vfs:    vfDevices,
		})
	}
	res := &network.Resource{
		Spec: &network.ResourceSpec{
			Name:    r.config.GetName(),
			Mtu:     r.config.GetMtu(),
			NumVfs:  r.config.GetNumVfs(),
			Devices: discoveredDevices,
		},
		Status: devices,
	}
	return res
}

func (s *NetworkServiceServer) GetAllNetworkAttachments(context.Context, *empty.Empty) (*network.NetworkAttachments, error) {
	nas := make([]*network.NetworkAttachment, 0)
	for _, na := range GetAllNetworkAttachments() {
		n, err := s.getNetworkAttachment(na)
		if err == nil {
			nas = append(nas, n)
		}
	}

	return &network.NetworkAttachments{Networkattachments: nas}, nil
}

func (s *NetworkServiceServer) CreateNetworkAttachment(_ context.Context, na *network.NetworkAttachment) (*empty.Empty, error) {
	cniConfig, _ := GetNetworkAttachment(na.GetName())
	if cniConfig != nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"network attachment name=%s already exists", na.GetName())
	}

	naConfig, err := AddNameToNetworkAttachment(na.GetConfig().AsMap(), na.GetName(), na.GetResourceName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot add name info: %v", err)
	}

	if na.GetResourceName() == "" {
		err = AddNetworkAttachment(na.GetName(), naConfig)
		if err != nil {
			RemoveNetworkAttachment(na.GetName())
			return new(empty.Empty), status.Errorf(codes.InvalidArgument, "Cannot add network attachment: %v", err)
		}
		return new(empty.Empty), nil
	}

	if r := s.getResource(na.GetResourceName()); r != nil {
		for _, dev := range r.provider.GetDevices() {
			for _, vf := range r.provider.GetVFDevices(dev) {
				if IsAllocated(vf.PCIAddress) == network.VFStatus_FREE {
					naConfig, err = AddDeviceIDToNetworkAttachment(naConfig, vf.PCIAddress)
					if err != nil {
						return nil, status.Errorf(codes.InvalidArgument, "Cannot add device info: %v", err)
					}
					err = AddNetworkAttachment(na.GetName(), naConfig)
					if err != nil {
						RemoveNetworkAttachment(na.GetName())
						return nil, status.Errorf(codes.InvalidArgument, "Cannot add network attachment: %v", err)
					}
					return new(empty.Empty), nil
				}
			}
		}

		return nil, status.Errorf(codes.ResourceExhausted, "resource id=%s has no free VFs", na.GetResourceName())

	}

	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", na.GetResourceName())
}

func (s *NetworkServiceServer) GetNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*network.NetworkAttachment, error) {
	na, err := GetNetworkAttachment(id.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetName())
	}

	return s.getNetworkAttachment(na)
}

func (*NetworkServiceServer) getNetworkAttachment(na map[string]interface{}) (*network.NetworkAttachment, error) {
	// obtain names from network attachment
	name, ok := na["name"].(string)
	if ok {
		delete(na, "name")
	}
	resourceName, ok := na["resourceName"].(string)
	if ok {
		delete(na, "resourceName")
	}

	config, err := structpb.NewStruct(na)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "network attachment config id=%s not valid", name)
	}
	return &network.NetworkAttachment{
		Name:         name,
		ResourceName: resourceName,
		Config:       config,
	}, nil
}

func (s *NetworkServiceServer) DeleteNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*empty.Empty, error) {
	err := RemoveNetworkAttachment(id.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetName())
	}

	return new(empty.Empty), nil
}

func NewNetworkService(c *network.ResourceConfigs) *NetworkServiceServer {
	list := make([]resource, 0)
	for _, r := range c.Configs {
		list = append(list, resource{
			config:   r,
			provider: devices.NewNetDeviceProvider(),
		})
	}
	return &NetworkServiceServer{resources: list}
}

func (s *NetworkServiceServer) Do() {
	for _, r := range s.resources {
		s.configure(&r)
	}
}

func (s *NetworkServiceServer) configure(r *resource) error {
	err := r.provider.Discover(r.config)
	if err != nil {
		return err
	}

	err = r.provider.Configure(r.config)
	if err != nil {
		return err
	}

	err = r.provider.ProbeNics()
	if err != nil {
		return err
	}

	return nil
}
