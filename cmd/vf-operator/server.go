package main

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sriramy/vf-operator/pkg/config"
	network "github.com/sriramy/vf-operator/pkg/stubs/network"
)

type networkServiceServer struct {
	network.UnimplementedNetworkServiceServer
	config *config.ResourceConfigList
}

func (s *networkServiceServer) GetAllResources(empty *empty.Empty, stream network.NetworkService_GetAllResourcesServer) error {
	for _, r := range s.config.Resources {
		res := &network.Resource{
			Name:   r.Name,
			NumVfs: uint32(r.NumVfs),
			NicSelector: &network.NicSelector{
				Vendors: r.GetVendors(),
				Drivers: r.GetDrivers(),
				Devices: r.GetDevices(),
				PfNames: r.GetPfNames(),
			},
			DeviceType: r.DeviceType,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}

func newServer(c *config.ResourceConfigList) *networkServiceServer {
	return &networkServiceServer{config: c}
}
