package server

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sriramy/vf-operator/pkg/config"
	"github.com/sriramy/vf-operator/pkg/devices"
	network "github.com/sriramy/vf-operator/pkg/stubs/network"
	"github.com/sriramy/vf-operator/pkg/utils"
)

type resource struct {
	config   config.ResourceConfig
	provider *devices.NetDeviceProvider
}

type NetworkServiceServer struct {
	network.UnimplementedNetworkServiceServer
	resources []resource
}

func (s *NetworkServiceServer) CreateResourceConfig(c *network.ResourceConfig, stream network.NetworkService_CreateResourceConfigServer) error {
	r := &resource{
		config: config.ResourceConfig{
			Name:   c.GetName(),
			Mtu:    c.GetMtu(),
			NumVfs: c.GetNumVfs(),
			NicSelector: config.NicSelector{
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
	if err := stream.Send(s.getResource(r)); err != nil {
		return err
	}

	return nil
}

func (s *NetworkServiceServer) GetAllResourceConfigs(_ *empty.Empty, stream network.NetworkService_GetAllResourceConfigsServer) error {
	for _, r := range s.resources {
		res := &network.ResourceConfig{
			Name:   r.config.Name,
			Mtu:    r.config.Mtu,
			NumVfs: r.config.NumVfs,
			NicSelector: &network.NicSelector{
				Vendors: r.config.GetVendors(),
				Drivers: r.config.GetDrivers(),
				Devices: r.config.GetDevices(),
				PfNames: r.config.GetPfNames(),
			},
			DeviceType: r.config.DeviceType,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *NetworkServiceServer) GetAllResources(_ *empty.Empty, stream network.NetworkService_GetAllResourcesServer) error {
	for _, r := range s.resources {
		if err := stream.Send(s.getResource(&r)); err != nil {
			return err
		}
	}

	return nil
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
			Name:    r.config.Name,
			Mtu:     r.config.Mtu,
			NumVfs:  r.config.NumVfs,
			Devices: discoveredDevices,
		},
		Status: devices,
	}
	return res
}

func NewNetworkService(c *config.ResourceConfigList) *NetworkServiceServer {
	list := make([]resource, 0)
	for _, r := range c.Resources {
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
	err := r.provider.Discover(&r.config)
	if err != nil {
		return err
	}

	err = r.provider.Configure(&r.config)
	if err != nil {
		return err
	}

	err = r.provider.ProbeNics()
	if err != nil {
		return err
	}

	return nil
}
