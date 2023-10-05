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
