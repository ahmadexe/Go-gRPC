package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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

	chReq := make(chan *pb.UserRequest)
	chRes := make(chan *pb.UserResponse)

	// callToFetchStream(client)
	// callToFetchResponseViaStream(client, []string{"1", "2", "3", "4"})
	callToGetUsersBidirectionalStream(client, []string{"1", "2", "3", "4"}, chReq, chRes)

	for data := range chReq {
		fmt.Printf("Request: %v\n", data)
	}


	for data := range chRes {
		fmt.Printf("Response: %v\n", data)
	}
}

func callGetUser(client pb.UserServiceClient) {
	res, err := client.GetUser(context.TODO(), &pb.UserRequest{Id: "3"})
	if err != nil {
		log.Fatal(err)
	}

	user := data.User{Id: res.Id, Name: res.Name, Age: res.Age}

	log.Printf("User: %v", user)
}

func callToFetchStream(client pb.UserServiceClient) {
	stream, err := client.StreamAllUsers(context.TODO(), &pb.NoParam{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		user := data.User{Id: res.Id, Name: res.Name, Age: res.Age}
		log.Printf("User: %v", user)
	}

	log.Println("Stream completed")
}

func callToFetchResponseViaStream(client pb.UserServiceClient, ids []string) {
	stream, err := client.FetchStreamResponse(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for _, id := range ids {
		err := stream.Send(&pb.UserRequest{Id: id})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent request for user with id %v\n", id)
		time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatal(err)
	}

	for _, user := range res.Users {
		log.Printf("Recieved User: %v", user)
	}

	log.Println("Stream completed")
}

func callToGetUsersBidirectionalStream(client pb.UserServiceClient, ids []string, chReq chan *pb.UserRequest, chRes chan *pb.UserResponse) {
	stream, err := client.GetUsersBidirectionalStream(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for _, id := range ids {
			err := stream.Send(&pb.UserRequest{Id: id})
			if err != nil {
				log.Fatal(err)
			}

			chReq <- &pb.UserRequest{Id: id}
			time.Sleep(2 * time.Second)
		}
		close(chReq)
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(chRes)
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			chRes <- &pb.UserResponse{Id: res.Id, Name: res.Name, Age: res.Age}
		}
	}()
}
