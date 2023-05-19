package main

import (
	"context"
	"log"

	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/zencodinglab/xavier/gen/go/proto/xavier/v1/openai"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the OpenAI server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}()

	// Create a new OpenAI client.
	client := pb.NewOpenAIClient(conn)

	// Call the GetModels method.
	response, err := client.GetModels(context.Background(), &pb.GetModelsRequest{})
	if err != nil {
		log.Fatalf("failed to get models: %v", err)
	}

	// Process the response.
	for _, model := range response.Models {
		log.Printf("Model ID: %s", model.Id)
	}
}
