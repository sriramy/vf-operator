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

package networkattachment

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/resource"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NetworkAttachmentServiceServer struct {
	network.UnimplementedNetworkAttachmentServiceServer
	resourceService *resource.ResourceServiceServer
}

func NewNetworkAttachmentServer(r *resource.ResourceServiceServer, config []*network.NetworkAttachment) *NetworkAttachmentServiceServer {
	server := &NetworkAttachmentServiceServer{resourceService: r}
	for _, na := range config {
		if _, err := server.CreateNetworkAttachment(context.TODO(), na); err != nil {
			fmt.Print(err)
		}
	}
	return server
}

func (s *NetworkAttachmentServiceServer) CreateNetworkAttachment(_ context.Context, na *network.NetworkAttachment) (*empty.Empty, error) {
	n, _ := getNetworkAttachment(na.GetName())
	if n != nil {
		return nil, status.Errorf(codes.AlreadyExists,
			"network attachment id=%s already exists", na.GetName())
	}
	n = newNetworkAttachment(na.GetConfig().AsMap())
	n.addName(na.GetName(), na.GetResourceName())

	if na.GetResourceName() == "" {
		err := n.create(na.GetName())
		if err != nil {
			n.delete(na.GetName())
			return new(empty.Empty), status.Errorf(codes.Internal, "Cannot create network attachment: %v", err)
		}
		return new(empty.Empty), nil
	}

	r, err := s.resourceService.GetResource(context.TODO(), &network.ResourceName{Name: na.ResourceName})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "resource id=%s not found", na.GetResourceName())
	}

	for _, pf := range r.Status {
		for _, vf := range pf.Vfs {
			if !IsDeviceIDAllocated(vf.GetDevice()) {
				err := n.addDeviceID(vf.GetDevice())
				if err != nil {
					return nil, status.Errorf(codes.Internal, "Cannot add device info: %v", err)
				}
				err = n.create(na.GetName())
				if err != nil {
					n.delete(na.GetName())
					return nil, status.Errorf(codes.Internal, "Cannot create network attachment: %v", err)
				}
				return new(empty.Empty), nil
			}
		}
	}

	return nil, status.Errorf(codes.ResourceExhausted, "resource id=%s has no free VFs", na.GetResourceName())
}

func (s *NetworkAttachmentServiceServer) DeleteNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*empty.Empty, error) {
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

func (s *NetworkAttachmentServiceServer) GetAllNetworkAttachments(context.Context, *empty.Empty) (*network.NetworkAttachments, error) {
	nas := &network.NetworkAttachments{}
	for _, n := range getAllNetworkAttachments() {
		na, err := n.build()
		if err == nil {
			nas.Networkattachments = append(nas.Networkattachments, na)
		}
	}
	return nas, nil
}

func (s *NetworkAttachmentServiceServer) GetNetworkAttachment(_ context.Context, id *network.NetworkAttachmentName) (*network.NetworkAttachment, error) {
	n, err := getNetworkAttachment(id.GetName())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "network attachment id=%s not found", id.GetName())
	}
	return n.build()
}
