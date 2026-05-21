package main

import (
	"context"
	"log"
	"time"

	pb "go-grpc-client/generated"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:9090",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewEmployeeClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	request := &pb.EmployeeRequest{
		Id: 1,
	}

	response, err := client.GetEmployee(ctx, request)

	if err != nil {
		log.Fatalf("Failed to make RPC: %v", err)
	}

	log.Println("EmpId: ", response.Id)
	log.Println("Name: ", response.Name)
	log.Println("Age: ", response.Age)
	log.Println("Salary: ", response.Salary)
}
