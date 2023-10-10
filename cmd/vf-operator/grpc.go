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

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	network "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/networkattachment"
	"github.com/sriramy/vf-operator/pkg/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGrpcServer(i *Input, c *network.InitialConfiguration) {
	defer i.wg.Done()

	serverEndpoint, err := net.Listen("tcp", fmt.Sprintf(":%d", *i.port))
	if err != nil {
		fmt.Printf("Listen failed: %v", err.Error())
		os.Exit(1)
	}
	defer serverEndpoint.Close()

	// start gRPC server
	grpcServer := grpc.NewServer()
	resourceServer := resource.NewResourceService(c.ResourceConfigs)
	networkAttachmentServer := networkattachment.NewNetworkAttachmentServer(resourceServer, c.Networkattachments)
	network.RegisterResourceServiceServer(grpcServer, resourceServer)
	network.RegisterNetworkAttachmentServiceServer(grpcServer, networkAttachmentServer)
	grpcServer.Serve(serverEndpoint)
}

func startGrpcGateway(i *Input) {
	defer i.wg.Done()

	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf(":%d", *i.port),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	// register gRPC gateway handlers
	gwmux := runtime.NewServeMux()
	err = network.RegisterResourceServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register resource service gateway:", err)
	}
	err = network.RegisterNetworkAttachmentServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register network attachment service gateway:", err)
	}

	// mount gRPC mux as root
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// mount swagger-ui as /sw/
	mux.Handle("/sw/", http.StripPrefix("/sw/", http.FileServer(http.Dir("/swagger-ui"))))

	// start gRPC gateway
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *i.gwPort),
		Handler: mux,
	}
	gwServer.ListenAndServe()
}
