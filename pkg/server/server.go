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

func (s *NetworkServiceServer) GetAllResources(empty *empty.Empty, stream network.NetworkService_GetAllResourcesServer) error {
	for _, r := range s.resources {
		devices := make([]*network.ResourceStatus, 0)
		for _, dev := range r.provider.GetDevices() {
			devices = append(devices, &network.ResourceStatus{
				Name:   dev.Name,
				Mac:    dev.MACAddress,
				Device: dev.PCIAddress,
				Vendor: dev.Vendor,
				Driver: dev.Driver,
				Vf:     utils.IsSriovVF(&dev.PCIAddress),
			})
			vfs, err := dev.GetVfs()
			if err == nil {
				for _, vf := range vfs {
					devices = append(devices, &network.ResourceStatus{
						Name:   vf.Name,
						Mac:    vf.MACAddress,
						Device: vf.PCIAddress,
						Vendor: vf.Vendor,
						Driver: vf.Driver,
						Vf:     true,
					})
				}
			}
		}
		res := &network.Resource{
			Name:   r.config.Name,
			NumVfs: uint32(r.config.NumVfs),
			NicSelector: &network.NicSelector{
				Vendors: r.config.GetVendors(),
				Drivers: r.config.GetDrivers(),
				Devices: r.config.GetDevices(),
				PfNames: r.config.GetPfNames(),
			},
			DeviceType: r.config.DeviceType,
			Status:     devices,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
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

func (s *NetworkServiceServer) configure(r *resource) {
	r.provider.Discover(&r.config)
	r.provider.Configure(&r.config)
}
