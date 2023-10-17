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

package resource

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResourceServiceServer struct {
	network.UnimplementedResourceServiceServer
	resources map[string]resource
}

func NewResourceService(config []*network.ResourceConfig) *ResourceServiceServer {
	server := &ResourceServiceServer{
		resources: make(map[string]resource),
	}
	for _, c := range config {
		if _, err := server.CreateResourceConfig(context.TODO(), c); err != nil {
			fmt.Printf("Resource config %s not created: %v\n", c.GetName(), err)
		}
	}
	return server
}

func (s *ResourceServiceServer) getResource(resourceName string) *resource {
	if r, ok := s.resources[resourceName]; ok {
		return &r
	}
	return nil
}

func (s *ResourceServiceServer) CreateResourceConfig(_ context.Context, c *network.ResourceConfig) (*network.Resource, error) {
	if r := s.getResource(c.GetName()); r != nil {
		return nil, status.Errorf(codes.AlreadyExists, "resource id=%s already exists", c.GetName())
	}

	r := newResource(c)
	s.resources[r.config.GetName()] = *r
	return r.build(), nil
}

func (s *ResourceServiceServer) DeleteResourceConfig(_ context.Context, id *network.ResourceName) (*empty.Empty, error) {
	if r := s.getResource(id.GetName()); r != nil {
		delete(s.resources, r.config.GetName())
		return new(empty.Empty), nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *ResourceServiceServer) GetAllResourceConfigs(context.Context, *empty.Empty) (*network.ResourceConfigs, error) {
	c := &network.ResourceConfigs{}
	for _, r := range s.resources {
		c.ResourceConfigs = append(c.ResourceConfigs, r.config)
	}
	return c, nil
}

func (s *ResourceServiceServer) GetResourceConfig(_ context.Context, id *network.ResourceName) (*network.ResourceConfig, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return r.config, nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}

func (s *ResourceServiceServer) GetAllResources(context.Context, *empty.Empty) (*network.Resources, error) {
	resources := &network.Resources{}
	for _, r := range s.resources {
		resources.Resources = append(resources.Resources, r.build())
	}
	return resources, nil
}

func (s *ResourceServiceServer) GetResource(_ context.Context, id *network.ResourceName) (*network.Resource, error) {
	if r := s.getResource(id.GetName()); r != nil {
		return r.build(), nil
	}
	return nil, status.Errorf(codes.NotFound, "resource id=%s not found", id.GetName())
}
