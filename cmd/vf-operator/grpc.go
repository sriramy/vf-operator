package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	networkservice "github.com/sriramy/vf-operator/pkg/api/v1/gen/network"
	"github.com/sriramy/vf-operator/pkg/config"
	"github.com/sriramy/vf-operator/pkg/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startGrpcServer(i *Input, c config.ResourceConfigList) {
	defer i.wg.Done()

	serverEndpoint, err := net.Listen("tcp", fmt.Sprintf(":%d", *i.port))
	if err != nil {
		fmt.Printf("Listen failed: %v", err.Error())
		os.Exit(1)
	}
	defer serverEndpoint.Close()

	// start network service
	service := server.NewNetworkService(&c)
	service.Do()

	// start gRPC server
	grpcServer := grpc.NewServer()
	networkservice.RegisterNetworkServiceServer(grpcServer, service)
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

	// register gRPC gateway handler
	gwmux := runtime.NewServeMux()
	err = networkservice.RegisterNetworkServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// mount gRPC mux as root
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	// mount swagger-ui and swagger.json
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("/swagger-ui"))))
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "networkservice.swagger.json")
	})

	// start gRPC gateway
	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *i.gwPort),
		Handler: mux,
	}
	gwServer.ListenAndServe()
}
