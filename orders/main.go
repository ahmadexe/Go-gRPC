package main

import (
	"context"
	"log"

	"github.com/ahmadexe/go-grpc/data"
	pb "github.com/ahmadexe/go-grpc/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	con, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer con.Close()

	client := pb.NewUserServiceClient(con)

	res, err := client.GetUser(context.TODO(), &pb.UserRequest{Id: "3"})
	if err != nil {
		log.Fatal(err)
	}

	user := data.User{Id: res.Id, Name: res.Name, Age: res.Age}

	log.Printf("User: %v", user)
}
