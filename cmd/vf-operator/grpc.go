package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sriramy/vf-operator/pkg/config"
	networkservice "github.com/sriramy/vf-operator/pkg/stubs/network"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGrpcServer(i *Input, c config.ResourceConfigList) *grpc.Server {
	defer i.wg.Done()

	serverEndpoint, err := net.Listen("tcp", fmt.Sprintf(":%d", *i.port))
	if err != nil {
		fmt.Printf("Listen failed: %v", err.Error())
		os.Exit(1)
	}
	defer serverEndpoint.Close()

	grpcServer := grpc.NewServer()
	networkservice.RegisterNetworkServiceServer(grpcServer, newServer(&c))
	grpcServer.Serve(serverEndpoint)

	return grpcServer
}

func startGrpcGateway(i *Input) *http.Server {
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

	gwmux := runtime.NewServeMux()

	err = networkservice.RegisterNetworkServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *i.gwPort),
		Handler: gwmux,
	}
	gwServer.ListenAndServe()

	return gwServer
}
