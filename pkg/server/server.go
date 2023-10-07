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
			Name:   &network.ResourceName{Id: c.GetName().GetId()},
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
	return s.getResource(r), nil
}

func (s *NetworkServiceServer) GetResourceConfig(_ context.Context, id *network.ResourceName) (*network.ResourceConfig, error) {
	for _, r := range s.resources {
		if r.config.GetName().GetId() == id.GetId() {
			return r.config, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetId())
}

func (s *NetworkServiceServer) DeleteResourceConfig(_ context.Context, id *network.ResourceName) (*empty.Empty, error) {
	for i, r := range s.resources {
		if r.config.GetName().GetId() == id.GetId() {
			s.resources = append(s.resources[:i], s.resources[i+1:]...)
			return new(empty.Empty), nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetId())
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
		resources = append(resources, s.getResource(&r))
	}

	return &network.Resources{Resources: resources}, nil
}

func (s *NetworkServiceServer) GetResource(_ context.Context, id *network.ResourceName) (*network.Resource, error) {
	for _, r := range s.resources {
		if r.config.GetName().GetId() == id.GetId() {
			return s.getResource(&r), nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetId())
}

func (s *NetworkServiceServer) getResource(r *resource) *network.Resource {
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
			Name:    &network.ResourceName{Id: r.config.GetName().GetId()},
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
	for _, na := range GetAll() {
		nas = append(nas, &network.NetworkAttachment{
			Name:         &network.NetworkAttachmentName{Id: na.naConfig.Name},
			ResourceName: &network.ResourceName{Id: na.resourceName},
			Mtu:          na.naConfig.Plugins[0].Mtu,
			Vlan:         na.naConfig.Plugins[0].Vlan,
		})
	}

	return &network.NetworkAttachments{Nas: nas}, nil
}

func (s *NetworkServiceServer) CreateNetworkAttachment(_ context.Context, na *network.NetworkAttachment) (*empty.Empty, error) {
	for _, r := range s.resources {
		if r.config.GetName().GetId() == na.GetResourceName().GetId() {
			for _, dev := range r.provider.GetDevices() {
				for _, vf := range r.provider.GetVFDevices(dev) {
					if IsAllocated(vf.PCIAddress) == network.VFStatus_FREE {
						naConfig := NewSriovNetworkAttachmentConfig(na, vf.PCIAddress)
						Store(naConfig, vf.PCIAddress, na.GetResourceName().GetId())
						err := AddNetworkAttachment(naConfig)
						if err != nil {
							Erase(naConfig.Name)
							return nil, status.Errorf(codes.Aborted, "Cannot add network attachment: %v", err)
						}
						return new(empty.Empty), nil
					}
				}
			}

			return nil, status.Errorf(codes.ResourceExhausted, "resource id=%s has no free VFs", na.GetResourceName().GetId())
		}
	}

	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", na.GetResourceName().GetId())
}

func (s *NetworkServiceServer) GetNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*network.NetworkAttachment, error) {
	if na := Get(id.GetId()); na != nil {
		return &network.NetworkAttachment{
			Name:         &network.NetworkAttachmentName{Id: na.naConfig.Name},
			ResourceName: &network.ResourceName{Id: na.resourceName},
			Mtu:          na.naConfig.Plugins[0].Mtu,
			Vlan:         na.naConfig.Plugins[0].Vlan,
		}, nil
	}

	return nil, status.Errorf(codes.ResourceExhausted, "network attachment id=%s not found", id.GetId())
}

func (s *NetworkServiceServer) DeleteNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*empty.Empty, error) {
	if na := Get(id.GetId()); na != nil {
		na := Get(id.GetId())
		RemoveNetworkAttachment(na.naConfig)
		Erase(id.GetId())
		return new(empty.Empty), nil
	}

	return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetId())
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
