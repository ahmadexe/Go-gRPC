package main

import (
	"log"
	"net"

	pb "github.com/ahmadexe/go-grpc/grpc"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UserServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", lis.Addr())

	// Create a gRPC server

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{})

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}