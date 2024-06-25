package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/ahmadexe/go-grpc/data"
	pb "github.com/ahmadexe/go-grpc/grpc"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UserServiceServer
}

var (
	usersSlice = []data.User{
		{Id: "1", Name: "Alice", Age: 25},
		{Id: "2", Name: "Bob", Age: 30},
		{Id: "3", Name: "Foo", Age: 35},
		{Id: "4", Name: "Bar", Age: 40},
	}
)

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

func (s *userServiceServer) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	for _, user := range usersSlice {
		if user.Id == in.Id {
			return &pb.UserResponse{Id: user.Id, Name: user.Name, Age: user.Age}, nil
		}
	}

	return nil, nil
}

func (s *userServiceServer) StreamAllUsers(in *pb.NoParam, stream pb.UserService_StreamAllUsersServer) error {
	for _, user := range usersSlice {
		if err := stream.Send(&pb.UserResponse{Id: user.Id, Name: user.Name, Age: user.Age}); err != nil {
			return err
		}
		time.Sleep(2 * time.Second)
	}

	return nil
}

func (s *userServiceServer) FetchStreamResponse(stream pb.UserService_FetchStreamResponseServer) error {
	var users []data.User

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			grpcUsers := make([]*pb.UserResponse, 0, len(users))
			for _, user := range users {
				grpcUsers = append(grpcUsers, &pb.UserResponse{Id: user.Id, Name: user.Name, Age: user.Age})
			}
			return stream.SendAndClose(&pb.UserList{Users: grpcUsers})
		}
		if err != nil {
			return err
		}

		for _, user := range usersSlice {
			if user.Id == req.Id {
				fmt.Printf("This is user %v\n", user)
				users = append(users, user)
			}
		}
		
	}
}