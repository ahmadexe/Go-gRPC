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

	var i int
	fmt.Println("1. Get User\n2. Stream All Users\n3. Fetch Stream Response\n4. Get Users Bidirectional Stream")
	fmt.Scan(&i)

	switch i {
	case 1:
		var id int
		fmt.Println("Enter user id")
		fmt.Scan(&id)
		if id < 1 || id > 4 {
			log.Fatal("Invalid user id")
		}
		callGetUser(client, id)
	case 2:
		callToFetchStream(client)
	case 3:
		callToFetchResponseViaStream(client, []string{"1", "2", "3", "4"})
	case 4:
		callToGetUsersBidirectionalStream(client, []string{"1", "2", "3", "4"})
	}
}

func callGetUser(client pb.UserServiceClient, id int) {
	res, err := client.GetUser(context.TODO(), &pb.UserRequest{Id: fmt.Sprintf("%d", id)})
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

func callToGetUsersBidirectionalStream(client pb.UserServiceClient, ids []string) {
	stream, err := client.GetUsersBidirectionalStream(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Received User: %v\n", res)
		}
		close(waitc)
	}()

	for _, id := range ids {
		err := stream.Send(&pb.UserRequest{Id: id})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sent request for user with id %v\n", id)
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
}
