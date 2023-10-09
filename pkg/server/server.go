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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NetworkServiceServer struct {
	network.UnimplementedNetworkServiceServer
	resources map[string]resource
}

func NewNetworkService(c *network.ResourceConfigs) *NetworkServiceServer {
	server := &NetworkServiceServer{}
	for _, c := range c.Configs {
		server.resources[c.GetName()] = *newResource(c)
	}
	return server
}

func (s *NetworkServiceServer) CreateResourceConfig(_ context.Context, c *network.ResourceConfig) (*network.Resource, error) {
	r := newResource(c)
	s.resources[r.config.GetName()] = *r
	return r.build(), nil
}

func (s *NetworkServiceServer) getResource(resourceName string) *resource {
	if r, ok := s.resources[resourceName]; ok {
		return &r
	}
	return nil
}

func (s *NetworkServiceServer) GetResourceConfig(_ context.Context, id *network.ResourceName) (*network.ResourceConfig, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return r.config, nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) DeleteResourceConfig(_ context.Context, id *network.ResourceName) (*empty.Empty, error) {
	if r := s.getResource(id.GetName()); r != nil {
		delete(s.resources, r.config.GetName())
		return new(empty.Empty), nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) GetAllResourceConfigs(context.Context, *empty.Empty) (*network.ResourceConfigs, error) {
	configs := &network.ResourceConfigs{}
	for _, r := range s.resources {
		configs.Configs = append(configs.Configs, r.config)
	}
	return configs, nil
}

func (s *NetworkServiceServer) GetAllResources(context.Context, *empty.Empty) (*network.Resources, error) {
	resources := &network.Resources{}
	for _, r := range s.resources {
		resources.Resources = append(resources.Resources, r.build())
	}
	return resources, nil
}

func (s *NetworkServiceServer) GetResource(_ context.Context, id *network.ResourceName) (*network.Resource, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return r.build(), nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *NetworkServiceServer) GetAllNetworkAttachments(context.Context, *empty.Empty) (*network.NetworkAttachments, error) {
	nas := &network.NetworkAttachments{}
	for _, n := range getAllNetworkAttachments() {
		na, err := n.build()
		if err == nil {
			nas.Networkattachments = append(nas.Networkattachments, na)
		}
	}
	return nas, nil
}

func (s *NetworkServiceServer) CreateNetworkAttachment(_ context.Context, na *network.NetworkAttachment) (*empty.Empty, error) {
	n, _ := getNetworkAttachment(na.GetName())
	if n != nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"network attachment name=%s already exists", na.GetName())
	}
	n = newNetworkAttachment(na.GetConfig().AsMap())
	n.addName(na.GetName(), na.GetResourceName())

	if na.GetResourceName() == "" {
		err := n.create(na.GetName())
		if err != nil {
			n.delete(na.GetName())
			return new(empty.Empty), status.Errorf(codes.InvalidArgument, "Cannot add network attachment: %v", err)
		}
		return new(empty.Empty), nil
	}

	if r := s.getResource(na.GetResourceName()); r != nil {
		for _, dev := range r.provider.GetDevices() {
			for _, vf := range r.provider.GetVFDevices(dev) {
				if isDeviceIDAllocated(vf.PCIAddress) == network.VFStatus_FREE {
					err := n.addDeviceID(vf.PCIAddress)
					if err != nil {
						return nil, status.Errorf(codes.InvalidArgument, "Cannot add device info: %v", err)
					}
					err = n.create(na.GetName())
					if err != nil {
						n.delete(na.GetName())
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
	n, err := getNetworkAttachment(id.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetName())
	}
	return n.build()
}

func (s *NetworkServiceServer) DeleteNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*empty.Empty, error) {
	n, err := getNetworkAttachment(id.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetName())
	}
	err = n.delete(id.GetName())
	if err != nil {
		return nil, err
	}
	return new(empty.Empty), nil
}
